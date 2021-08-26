package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/tr4d3r8/go-backend-boilerplate/api"
	db "github.com/tr4d3r8/go-backend-boilerplate/db/sqlc"
	"github.com/tr4d3r8/go-backend-boilerplate/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configurations", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
