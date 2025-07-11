package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/auronvila/simple-bank/api"
	db "github.com/auronvila/simple-bank/db/sqlc"
	_ "github.com/auronvila/simple-bank/doc/statik"
	"github.com/auronvila/simple-bank/gapi"
	"github.com/auronvila/simple-bank/pb"
	"github.com/auronvila/simple-bank/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}

	// run db migrations
	runDbMigration(config.MigrationUrl, config.DbSource)
	store := db.NewStore(conn)

	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runDbMigration(migrationUrl string, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal("cannot initialize migration!!: ", err)
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("error during migration process!!: ", err)
	}

	log.Println("db migrated successfully")
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("error creating the server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start the gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start the gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("error creating the server", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
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
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal("cannot register handler server: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFs, err := fs.New()
	if err != nil {
		log.Fatal("Cannot create statik file system %v", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("started the HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start the HTTP gateway server: ", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("error creating the server", err)
	}
	err = server.Start(config.HttpServerAddress)

	if err != nil {
		log.Fatal(err)
	}
}
