package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/peienxie/go-bank/api"
	db "github.com/peienxie/go-bank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@localhost:5432/go-bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	store := db.NewSQLStore(conn)
	server := api.NewServer(store)

	if err = server.Serve(":8080"); err != nil {
		log.Fatal(err)
	}
}
