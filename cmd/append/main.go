package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/aaronland/go-wunderkammer/oembed"
	"github.com/sfomuseum/go-flags/multi"	
	"github.com/tidwall/pretty"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	format_json := flag.Bool("format", false, "Emit results as formatted JSON.")
	as_json := flag.Bool("json", false, "Emit results as a JSON array.")

	to_stdout := flag.Bool("stdout", true, "Emit to STDOUT")
	to_devnull := flag.Bool("null", false, "Emit to /dev/null")

	timings := flag.Bool("timings", false, "Log timings (time to wait to process, time to complete processing")

	fragment := flag.String("fragment", "", "A valid URI fragment to append to an OEmbed URL.")

	var parameters multi.KeyValue
	flag.Var(&parameters, "parameter", "A valid URI parameter (key=value) to append an OEmbed URL.")

	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

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

	mu := new(sync.RWMutex)

	t0 := time.Now()

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

		u, err := url.Parse(rec.URL)

		if err != nil {
			log.Fatalf("Failed to parse URL %s, %v", rec.URL, err)
		}

		u.Fragment = *fragment

		if len(parameters) > 0 {

			q := u.Query()

			for _, param := range parameters {
				q.Set(param.Key, param.Value)
			}

			u.RawQuery = q.Encode()
		}
		
		rec.URL = u.String()

		body, err = json.Marshal(rec)

		if err != nil {
			log.Fatalf("Failed to marshal record, %v", err)
		}

		if *format_json {
			body = pretty.Pretty(body)
		}

		new_count := atomic.AddInt32(&count, 1)

		mu.Lock()

		if *as_json && new_count > 1 {
			wr.Write([]byte(","))
		}

		wr.Write(body)
		wr.Write([]byte("\n"))

		mu.Unlock()
	}

	if *as_json {
		wr.Write([]byte("]"))
	}

	if *timings {
		log.Printf("Time to process %d records, %v\n", count, time.Since(t0))
	}
}
