package main

import (
	"context"
	"errors"
	"github.com/auronvila/simple-bank/api"
	db "github.com/auronvila/simple-bank/db/sqlc"
	_ "github.com/auronvila/simple-bank/doc/statik"
	"github.com/auronvila/simple-bank/gapi"
	"github.com/auronvila/simple-bank/mail"
	accountPb "github.com/auronvila/simple-bank/pb/account"
	userPb "github.com/auronvila/simple-bank/pb/user"
	"github.com/auronvila/simple-bank/util"
	"github.com/auronvila/simple-bank/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.DbSource)
	prettyOutput := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err != nil {
		prettyOutput.Fatal().Err(err).Msg("cannot connect to the db")
	}

	// run db migrations
	runDbMigration(config.MigrationUrl, config.DbSource)
	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{Addr: config.RedisAddress}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, config, redisOpt, store)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor)
	runGrpcServer(ctx, waitGroup, config, store, taskDistributor)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runDbMigration(migrationUrl string, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot initialize migration!!: ")
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("error during migration process!!: ")
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.SmtpSenderName, config.SmtpUsername, config.SmtpPass)

	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to run task processor")
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown task processor")

		taskProcessor.ShutDown()
		log.Info().Msg("task processor is stoped")
		return nil
	})
}

func runGrpcServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating the server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	userPb.RegisterUsersServer(grpcServer, server)
	accountPb.RegisterAccountsServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start the gRPC server at %s", listener.Addr().String())
		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Msg("grpc server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown grpc server")

		grpcServer.GracefulStop()
		log.Info().Msg("grpc server is stoped")
		return nil
	})
}

func runGatewayServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating the server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption, runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		switch key {
		case "user-agent", "x-real-ip", "x-forwarded-for":
			return key, true
		default:
			return runtime.DefaultHeaderMatcher(key)
		}
	}))

	err = userPb.RegisterUsersHandlerServer(ctx, grpcMux, server)
	err = accountPb.RegisterAccountsHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server: ")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFs, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create statik file system")
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	httpServer := &http.Server{
		Handler: gapi.HttpLogger(mux),
		Addr:    config.HttpServerAddress,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("started the HTTP gateway server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("HTTP gateway server failed to serve: ")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("gracefully shut down http gateway")

		err = httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shut down http gateway server")
		}

		log.Info().Msg("http gateway server is stoped")
		return nil
	})

}

// ? deprecated: used before to spin up a gin server but now grpc gateway is used
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating the server")
	}
	err = server.Start(config.HttpServerAddress)

	if err != nil {
		log.Fatal().Err(err)
	}
}
