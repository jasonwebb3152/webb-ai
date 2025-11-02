package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/jasonwebb3152/webb-ai/api/gapi"
	"github.com/jasonwebb3152/webb-ai/api/pb"
	db "github.com/jasonwebb3152/webb-ai/db/sqlc"

	"github.com/golang-migrate/migrate/v4"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	DATABASE_URL        = "postgresql://root:secret@localhost:5432/webb_ai?sslmode=disable"
	HTTP_SERVER_ADDRESS = "0.0.0.0:8080"
	MIGRATION_URL       = "file://db/migration"
	SYMMETRIC_KEY       = "12345678901234567890123456789012"
)

func main() {
	ctx := context.Background()

	// Connect to DB
	pool, err := pgxpool.New(ctx, DATABASE_URL)
	if err != nil {
		log.Fatalf("Failed to connext to db: %s", err)
	}
	log.Println("Successfully connected to DB!")

	store := db.NewStore(pool)

	// Migrate DB to newest schema
	runDbMigration()

	// Start HTTP Gateway Server
	runHttpGatewayServer(store)
}

func runDbMigration() {
	m, err := migrate.New(MIGRATION_URL, DATABASE_URL)
	if err != nil {
		log.Fatalf("Failed to instantiate migration: %s", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migratrion: %s", err)
	}

	log.Printf("Successfully ran migration!")
}

func runHttpGatewayServer(store db.Store) {
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server, err := gapi.NewServer(store, SYMMETRIC_KEY)
	if err != nil {
		log.Fatalf("Failed to create server (token maker): %s", err)
	}

	err = pb.RegisterWebbAiHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatalf("Failed to register grpcMux with server: %s", err)
	}

	// Create server that reroutes to grpc server
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", HTTP_SERVER_ADDRESS)
	if err != nil {
		log.Fatalf("Failed to create listener: %s", err)
	}

	log.Println("Starting HTTP Server!")
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalf("Failed to start http server: %s", err)
	}
}
