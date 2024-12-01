package simplebank

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:1@localhost:5432/simple-bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}

	testQueries = New(conn)

	m.Run()
}
