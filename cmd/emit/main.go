package main

import (
	_ "github.com/mattn/go-sqlite3"
)

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aaronland/go-json-query"
	"github.com/aaronland/go-wunderkammer/oembed"
	"github.com/tidwall/pretty"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {

	dsn := flag.String("database-dsn", "sql://sqlite3/oembed.db", "A valid wunderkammer database DSN string.")

	to_stdout := flag.Bool("stdout", true, "Emit to STDOUT")
	to_devnull := flag.Bool("null", false, "Emit to /dev/null")

	as_json := flag.Bool("json", false, "Emit a JSON list.")
	format_json := flag.Bool("format-json", false, "Format JSON output for each record.")

	var queries query.QueryFlags
	flag.Var(&queries, "query", "One or more {PATH}={REGEXP} parameters for filtering records.")

	valid_modes := strings.Join([]string{query.QUERYSET_MODE_ALL, query.QUERYSET_MODE_ANY}, ", ")
	desc_modes := fmt.Sprintf("Specify how query filtering should be evaluated. Valid modes are: %s", valid_modes)

	query_mode := flag.String("query-mode", query.QUERYSET_MODE_ALL, desc_modes)

	flag.Parse()

	ctx := context.Background()

	db, err := oembed.NewSQLOEmbedDatabase(ctx, *dsn)

	if err != nil {
		log.Fatalf("Failed to create database, %v", err)
	}

	defer db.Close()

	var qs *query.QuerySet

	if len(queries) > 0 {

		qs = &query.QuerySet{
			Queries: queries,
			Mode:    *query_mode,
		}
	}

	writers := make([]io.Writer, 0)

	if *to_stdout {
		writers = append(writers, os.Stdout)
	}

	if *to_devnull {
		writers = append(writers, ioutil.Discard)
	}

	if len(writers) == 0 {
		log.Fatal("Nothing to write to.")
	}

	wr := io.MultiWriter(writers...)

	count := int32(0)

	mu := new(sync.RWMutex)

	if *as_json {
		wr.Write([]byte("["))
	}

	cb := func(ctx context.Context, ph *oembed.Photo) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		body, err := json.Marshal(ph)

		if err != nil {
			return err
		}

		if qs != nil {

			matches, err := query.Matches(ctx, qs, body)

			if err != nil {
				return err
			}

			if !matches {
				return nil
			}
		}

		if *format_json {
			body = pretty.Pretty(body)
		}

		body = bytes.TrimSpace(body)

		new_count := atomic.AddInt32(&count, 1)

		mu.Lock()
		defer mu.Unlock()

		if *as_json && new_count > 1 {
			wr.Write([]byte(","))
		}

		wr.Write(body)
		wr.Write([]byte("\n"))

		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = db.GetOEmbedWithCallback(ctx, cb)

	if err != nil {
		log.Fatal(err)
	}

	if *as_json {
		wr.Write([]byte("]"))
	}

}
