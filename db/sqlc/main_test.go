package db

import (
	"database/sql"
	_ "github.com/auronvila/simple-bank/doc/statik"
	"github.com/auronvila/simple-bank/util"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"testing"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var config util.Config
	var err error

	if err := godotenv.Load("../../app.env"); err != nil {
		log.Println("Warning: Could not load .env file:", err)
	}

	config, err = util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDb, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}

	testQueries = New(testDb)

	m.Run()
}
