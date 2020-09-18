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
	"net/url"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	format_json := flag.Bool("format", false, "Emit results as formatted JSON.")
	as_json := flag.Bool("json", false, "Emit results as a JSON array.")

	to_stdout := flag.Bool("stdout", true, "Emit to STDOUT")
	to_devnull := flag.Bool("null", false, "Emit to /dev/null")

	workers := flag.Int("workers", runtime.NumCPU(), "The number of concurrent workers to append data URLs with")
	timings := flag.Bool("timings", false, "Log timings (time to wait to process, time to complete processing")

	fragment := flag.String("fragment", "", "A valid URI fragment to append to an OEmbed URL.")

	// handle URL parameters here...

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

	throttle := make(chan bool, *workers)

	for i := 0; i < *workers; i++ {
		throttle <- true
	}

	mu := new(sync.RWMutex)
	wg := new(sync.WaitGroup)

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

		t1 := time.Now()

		<-throttle

		if *timings {
			log.Printf("Time to wait to process %s, %v\n", rec.URL, time.Since(t1))
		}

		wg.Add(1)

		go func(rec *oembed.Photo) {

			u, err := url.Parse(rec.URL)

			if err != nil {
				log.Fatalf("Failed to parse URL %s, %v", rec.URL, err)
			}

			u.Fragment = *fragment

			rec.URL = u.String()

			body, err := json.Marshal(rec)

			if err != nil {
				log.Fatalf("Failed to marshal record, %v", err)
			}

			if *format_json {
				body = pretty.Pretty(body)
			}

			new_count := atomic.AddInt32(&count, 1)

			mu.Lock()
			defer mu.Unlock()

			if *as_json && new_count > 1 {
				wr.Write([]byte(","))
			}

			wr.Write(body)
			wr.Write([]byte("\n"))

		}(rec)

	}

	if *as_json {
		wr.Write([]byte("]"))
	}

	wg.Wait()

	if *timings {
		log.Printf("Time to process %d records, %v\n", count, time.Since(t0))
	}
}
