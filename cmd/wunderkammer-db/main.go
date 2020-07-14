package main

import (
	_ "github.com/mattn/go-sqlite3"
)

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/aaronland/go-wunderkammer/oembed"
	"io"
	"log"
	"os"
)

func main() {

	dsn := flag.String("database-dsn", "sql://sqlite3/oembed.db", "...")
	populate_data_url := flag.Bool("populate-data-url", false, "")

	content_aware_resize := flag.Bool("content-aware-resize", false, "...")
	content_aware_height := flag.Int("content-aware-height", 0, "...")
	content_aware_width := flag.Int("content-aware-width", 0, "...")

	halftone := flag.Bool("halftone", false, "...")
	resize := flag.Bool("resize", false, "...")
	resize_max_dimension := flag.Int("resize-max-dimension", 0, "...")

	flag.Parse()

	if *content_aware_resize {

		if *content_aware_height == 0 {
			log.Fatalf("Missing -content-aware-height value")
		}

		if *content_aware_width == 0 {
			log.Fatalf("Missing -content-aware-width value")
		}
	}

	if *resize {

		if *resize_max_dimension == 0 {
			log.Fatalf("Missing -resize-max-dimension value")
		}
	}

	ctx := context.Background()

	db, err := oembed.NewSQLOEmbedDatabase(ctx, *dsn)

	if err != nil {
		log.Fatalf("Failed to create database, %v", err)
	}

	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	for {

		select {
		case <-ctx.Done():
			break
		default:
			// pass
		}

		body, err := reader.ReadBytes('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Failed to read bytes, %v", err)
		}

		body = bytes.TrimSpace(body)

		var rec *oembed.Photo

		err = json.Unmarshal(body, &rec)

		if err != nil {
			log.Fatalf("Failed to unmarshal OEmbed record, %v", err)
		}

		if *populate_data_url && rec.DataURL == "" {

			opts := &oembed.DataURLOptions{
				ContentAwareResize: *content_aware_resize,
				ContentAwareWidth:  *content_aware_width,
				ContentAwareHeight: *content_aware_height,
				Resize:             *resize,
				ResizeMaxDimension: *resize_max_dimension,
				Halftone:           *halftone,
			}

			data_url, err := oembed.DataURL(ctx, rec.URL, opts)

			if err != nil {
				log.Fatalf("Failed to populate data URL for '%s', %v", rec.URL, err)
			}

			rec.DataURL = data_url
		}

		err = db.AddOEmbed(ctx, rec)

		if err != nil {
			log.Fatalf("Failed to add OEmbed record %s to database, %v", rec.URL, err)
		}
	}

}
