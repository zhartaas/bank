package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testConn *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	testConn, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal(err)
	}
	defer testConn.Close()

	testQueries = New(testConn)

	os.Exit(m.Run())
}
