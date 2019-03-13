// Package handler exposes a net/http Handler. It invokes anything runnable and
// writes the output to the response stream.
package handler

import (
	"context"
	"log"
	"net/http"
)

type handler struct {
	errLog *log.Logger
	r      runner
}

type runner interface {
	Run(context.Context) ([]byte, error)
}

// New gets a new handler that writes errors to the logger and invokes r as its
// data source.
func New(errLog *log.Logger, r runner) *handler {
	return &handler{errLog, r}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := h.r.Run(r.Context())
	if err != nil {
		h.errLog.Printf("Unexpected error running: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		h.errLog.Printf("Unexpected error writing: %v", err)
	}
}
