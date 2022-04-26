package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

const (
	dbDriver = "mysql"
	dbSource = "root:root@tcp(localhost:3306)/simple_bank_test?parseTime=true"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
