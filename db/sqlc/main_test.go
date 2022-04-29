package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbDriver = "mysql"
	dbSource = "root:root@tcp(localhost:3306)/simple_bank_test?parseTime=true"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	testQueries = New(testDB)
	tx := &sql.Tx{}
	testQueries.WithTx(tx)
	os.Exit(m.Run())
}
