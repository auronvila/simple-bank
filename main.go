package main

import (
	"context"
	"database/sql"
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
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	prettyOutput := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err != nil {
		prettyOutput.Fatal().Err(err).Msg("cannot connect to the db")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// run db migrations
	runDbMigration(config.MigrationUrl, config.DbSource)
	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{Addr: config.RedisAddress}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(config, redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
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

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.SmtpSenderName, config.SmtpUsername, config.SmtpPass)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to run task processor")
	}
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
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

	log.Info().Msgf("start the gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start the gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener: ")
	}

	log.Info().Msgf("started the HTTP gateway server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start the HTTP gateway server: ")
	}
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
