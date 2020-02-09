// Package main is the executable for commandd. It spins up a HTTP server.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/epels/commandd/handler"
	"github.com/epels/commandd/uptime"
)

var (
	errLog  = log.New(os.Stderr, "[ERROR]: ", log.LstdFlags|log.Lshortfile)
	infoLog = log.New(os.Stdout, "[INFO]: ", log.LstdFlags|log.Lshortfile)

	addr    = flag.String("addr", ":8080", "Address to listen on")
	pattern = flag.String("pattern", "/uptime", "Pattern to respond to. Set to / for any path")
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle(*pattern, handler.New(errLog, uptime.New()))

	s := &http.Server{
		Addr:    *addr,
		Handler: mux,

		IdleTimeout:  60 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		infoLog.Printf("Starting server on %q", *addr)
		errChan <- s.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		errLog.Printf("net/http: Server.ListenAndServe: %v", err)
	case sig := <-sigChan:
		infoLog.Printf("Exiting with signal: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		errLog.Printf("net/http: Server.Shutdown: %v", err)
	}
}
