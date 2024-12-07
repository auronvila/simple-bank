package main

import (
	"database/sql"
	"github.com/auronvila/simple-bank/api"
	simplebank "github.com/auronvila/simple-bank/db/sqlc"
	"github.com/auronvila/simple-bank/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}
	conn, err := sql.Open(config.DBDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}

	store := simplebank.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal(err)
	}
}
