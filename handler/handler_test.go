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

	"github.com/epels/uptimed/internal/mock"
)

func TestServeHTTP(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/uptime", nil)

		r := mock.Runner{
			RunFunc: func(ctx context.Context) ([]byte, error) {
				return []byte("hello world"), nil
			},
		}
		h := New(log.New(ioutil.Discard, "", log.LstdFlags), &r)

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
		req := httptest.NewRequest(http.MethodGet, "/uptime", nil)

		var logBuf bytes.Buffer
		r := mock.Runner{
			RunFunc: func(ctx context.Context) ([]byte, error) {
				return nil, errors.New("some-error")
			},
		}
		h := New(log.New(&logBuf, "", log.LstdFlags), &r)

		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Got %d, expected 500", rec.Code)
		}

		// Assert it writes a log on error.
		if s := logBuf.String(); !strings.Contains(s, "some-error") {
			t.Fatalf("Got %q, expected to contain some-error", s)
		}
	})
}
