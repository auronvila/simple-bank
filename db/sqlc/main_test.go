package db

import (
	"database/sql"
	"github.com/auronvila/simple-bank/util"
	_ "github.com/lib/pq"
	"log"
	"testing"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var config util.Config
	var err error

	config, err = util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config on test main")
	}

	testDb, err = sql.Open(config.DBDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to the db", err)
	}

	testQueries = New(testDb)

	m.Run()
}
