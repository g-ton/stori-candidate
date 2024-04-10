package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/g-ton/stori-candidate/internal/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// ../.. means to go to the root of the project
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load configuration:", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
