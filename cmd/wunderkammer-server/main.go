package main

// THIS IS EARLY STAGES AND EVERYTHING IS IN FLUX

import (
	_ "github.com/mattn/go-sqlite3"
)

import (
	"context"
	"flag"
	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-wunderkammer/oembed"
	"github.com/aaronland/go-wunderkammer/www"
	"log"
	"net/http"
)

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080", "...")
	dsn := flag.String("database-dsn", "sql://sqlite3/oembed.db", "A valid wunderkammer database DSN string.")

	flag.Parse()

	ctx := context.Background()

	db, err := oembed.NewSQLOEmbedDatabase(ctx, *dsn)

	if err != nil {
		log.Fatalf("Failed to create database, %v", err)
	}

	defer db.Close()

	random_handler, err := www.NewRandomHandler(db)

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/random", random_handler)

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s", s.Address())
	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
