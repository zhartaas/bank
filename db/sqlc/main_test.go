package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	testQueries = New(conn)

	os.Exit(m.Run())
}
