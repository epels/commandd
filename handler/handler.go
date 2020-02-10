// Package handler exposes a net/http Handler. It invokes anything runnable and
// writes the output to the response stream.
package handler

import (
	"context"
	"log"
	"net/http"
	"time"
)

type handler struct {
	errLog  *log.Logger
	r       runner
	timeout time.Duration
}

type runner interface {
	Run(context.Context) ([]byte, error)
}

// New gets a new handler that writes errors to the logger and invokes r as its
// data source.
func New(errLog *log.Logger, r runner, timeout time.Duration) *handler {
	return &handler{
		errLog:  errLog,
		r:       r,
		timeout: timeout,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	b, err := h.r.Run(ctx)
	if err == context.DeadlineExceeded {
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}
	if err != nil {
		h.errLog.Printf("Unexpected error running: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		h.errLog.Printf("Unexpected error writing: %v", err)
	}
}
