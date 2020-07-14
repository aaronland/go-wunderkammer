package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/aaronland/go-wunderkammer/oembed"
	"github.com/tidwall/pretty"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
)

func main() {

	content_aware_resize := flag.Bool("content-aware-resize", false, "...")
	content_aware_height := flag.Int("content-aware-height", 0, "...")
	content_aware_width := flag.Int("content-aware-width", 0, "...")

	halftone := flag.Bool("halftone", false, "...")
	resize := flag.Bool("resize", false, "...")
	resize_max_dimension := flag.Int("resize-max-dimension", 0, "...")

	overwrite := flag.Bool("overwrite", false, "...")

	format_json := flag.Bool("format", false, "...")
	as_json := flag.Bool("json", false, "...")

	to_stdout := flag.Bool("stdout", true, "Emit to STDOUT")
	to_devnull := flag.Bool("null", false, "Emit to /dev/null")

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

	reader := bufio.NewReader(os.Stdin)

	count := int32(0)

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

		if rec.DataURL == "" || *overwrite {

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

		body, err = json.Marshal(rec)

		if err != nil {
			log.Fatalf("Failed to marshal record, %v", err)
		}

		if *format_json {
			body = pretty.Pretty(body)
		}

		new_count := atomic.AddInt32(&count, 1)

		if *as_json && new_count > 1 {
			wr.Write([]byte(","))
		}

		wr.Write(body)
		wr.Write([]byte("\n"))

	}

	if *as_json {
		wr.Write([]byte("]"))
	}

}
