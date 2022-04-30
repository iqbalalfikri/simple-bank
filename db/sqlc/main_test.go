package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iqbalalfikri/simple-bank/util"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..", "test")
	testDB, err = sql.Open(config.DatabaseConfig.Driver, config.DatabaseConfig.Source)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	testQueries = New(testDB)
	tx := &sql.Tx{}
	testQueries.WithTx(tx)
	os.Exit(m.Run())
}
