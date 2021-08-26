package db

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/tr4d3r8/go-backend-boilerplate/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("could not load config", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	m.Run()
}
