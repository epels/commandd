// Package handler exposes a net/http Handler. It invokes anything runnable and
// writes the output to the response stream.
package handler

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
)

type handler struct {
	errLog  *log.Logger
	r       runner
	timeout time.Duration
}

type runner interface {
	Run(context.Context, io.Writer) error
}

// New gets a new handler that writes errors to the logger and invokes r as its
// data source.
func New(errLog *log.Logger, r runner, timeout time.Duration) http.Handler {
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		// Log an error, but proceed anyway - not being able to expose metrics
		// is not critical.
		errLog.Printf("Error registering server metric views: %s", err)
	}
	inner := &handler{
		errLog:  errLog,
		r:       r,
		timeout: timeout,
	}
	return &ochttp.Handler{Handler: inner}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	err := h.r.Run(ctx, w)
	switch {
	case errors.Is(err, context.Canceled):
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	case errors.Is(err, context.DeadlineExceeded):
		w.WriteHeader(http.StatusRequestTimeout)
		return
	case err != nil:
		h.errLog.Printf("%T.Run: %s", h.r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
