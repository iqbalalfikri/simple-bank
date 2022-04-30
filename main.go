package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iqbalalfikri/simple-bank/api"
	db "github.com/iqbalalfikri/simple-bank/db/sqlc"
	"github.com/iqbalalfikri/simple-bank/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".", "dev")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DatabaseConfig.Driver, config.DatabaseConfig.Source)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerConfig.Address)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
