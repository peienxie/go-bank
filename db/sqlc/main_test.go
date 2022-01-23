package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open("postgres", "postgresql://root:root@localhost:5432/go-bank?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
