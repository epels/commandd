package handler

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/epels/commandd/internal/mock"
)

func TestServeHTTP(t *testing.T) {
	defaultTimeout := 10 * time.Second // Arbitrary.
	noopLogger := log.New(ioutil.Discard, "", 0)

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/run", nil)

		r := mock.Runner{
			RunFunc: func(ctx context.Context) ([]byte, error) {
				return []byte("hello world"), nil
			},
		}
		h := New(noopLogger, &r, defaultTimeout)

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Got %d, expected 200", rec.Code)
		}

		if s := rec.Body.String(); s != "hello world" {
			t.Errorf("Got %q, expected hello world", s)
		}
	})

	t.Run("Runner error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/run", nil)

		var logBuf bytes.Buffer
		r := mock.Runner{
			RunFunc: func(ctx context.Context) ([]byte, error) {
				return nil, errors.New("some-error")
			},
		}
		h := New(log.New(&logBuf, "", log.LstdFlags), &r, defaultTimeout)

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Got %d, expected 500", rec.Code)
		}

		// Assert it writes a log on error.
		if s := logBuf.String(); !strings.Contains(s, "some-error") {
			t.Fatalf("Got %q, expected to contain some-error", s)
		}
	})

	t.Run("Runner canceled", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/run", nil)

		r := mock.Runner{
			RunFunc: func(ctx context.Context) ([]byte, error) {
				return nil, context.Canceled
			},
		}
		h := New(noopLogger, &r, defaultTimeout)

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("Got %d, expected 503", rec.Code)
		}
	})

	t.Run("Runner timeout", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/run", nil)

		r := mock.Runner{
			RunFunc: func(ctx context.Context) ([]byte, error) {
				return nil, context.DeadlineExceeded
			},
		}
		h := New(noopLogger, &r, defaultTimeout)

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusRequestTimeout {
			t.Errorf("Got %d, expected 408", rec.Code)
		}
	})
}
