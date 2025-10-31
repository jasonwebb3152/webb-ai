package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const DATABASE_URL = "postgresql://root:secret@localhost:5432/webb_ai?sslmode=disable"

var (
	testStore Store
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, DATABASE_URL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testStore = NewStore(connPool)
	os.Exit(m.Run()) // Start running unit test
}
