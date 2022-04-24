package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/peienxie/go-bank/api"
	"github.com/peienxie/go-bank/config"
	db "github.com/peienxie/go-bank/db/sqlc"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	store := db.NewSQLStore(conn)
	server := api.NewServer(store)

	if err = server.Serve(config.ServerAddress); err != nil {
		log.Fatal(err)
	}
}
