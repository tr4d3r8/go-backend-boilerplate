package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"gitlab.com/tr4d3r/backend-master-golang/api"
	db "gitlab.com/tr4d3r/backend-master-golang/db/sqlc"
	"gitlab.com/tr4d3r/backend-master-golang/util"
)

func main(){
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