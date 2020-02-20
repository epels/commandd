// Package main facilitates end-to-end testing and exposes a flags for
// configuring the (concurrent) request count, the expected response and the
// target URL.
package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	count    int
	response string
	url      string
)

func init() {
	flag.IntVar(&count, "count", 50, "Number of requests to fire")
	flag.StringVar(&response, "response", "foo bar baz", "Expected response")
	flag.StringVar(&url, "url", "http://localhost:8081", "Target URL for end-to-end test")
	flag.Parse()
}

func main() {
	// Avoid unnecessary conversions between []byte and string.
	rb := []byte(response)

	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			c := http.Client{
				Timeout: 2 * time.Second,
			}
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				log.Fatalf("net/http: NewRequest: %s", err)
			}

			res, err := c.Do(req)
			if err != nil {
				log.Fatalf("net/http: Client.Do: %s", err)
			}
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Printf("%T.Close: %s", res.Body, err)
				}
			}()

			if res.StatusCode != http.StatusOK {
				log.Fatalf("Got %d, expected 200", res.StatusCode)
			}

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalf("io/ioutil: ReadAll: %s", err)
			}
			if !bytes.Equal(b, rb) {
				log.Fatalf("Got %s, expected %s", b, rb)
			}

			wg.Done()
		}()
	}
	wg.Wait()
	log.Printf("OK")
}
