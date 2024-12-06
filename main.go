package main

import (
	"database/sql"
	"github.com/auronvila/simple-bank/api"
	simplebank "github.com/auronvila/simple-bank/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:1@localhost:5432/simple-bank?sslmode=disable"
	address  = "0.0.0.0:3002"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}

	store := simplebank.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(address)

	if err != nil {
		log.Fatal(err)
	}
}
