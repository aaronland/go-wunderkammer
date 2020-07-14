package main

import (
	_ "github.com/mattn/go-sqlite3"
)

import (
	"context"
	"database/sql"
	"encoding/base64"
	"flag"
	"github.com/aaronland/go-http-server"
	"log"
	"net/http"
	"strings"
)

func NewImageHandler(db *sql.DB) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		path := req.URL.Path
		path = strings.TrimLeft(path, "/")

		q := "SELECT body FROM images WHERE id = ?"
		row := db.QueryRowContext(ctx, q, path)

		var body string
		err := row.Scan(&body)

		if err != nil {
			log.Println(err)
			return
		}

		raw, err := base64.StdEncoding.DecodeString(body)

		if err != nil {
			log.Println(err)
			return
		}

		rsp.Header().Set("Content-type", "image/png")
		rsp.Write(raw)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080", "...")
	sqlite_dsn := flag.String("sqlite-dsn", ":memory:", "...")

	flag.Parse()

	ctx := context.Background()

	db, err := sql.Open("sqlite3", *sqlite_dsn)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	image_handler, err := NewImageHandler(db)

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", image_handler)

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
