package command

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Run("OK with args", func(t *testing.T) {
		cmd, err := New("echo", "-n", "foo", "bar", "baz 42")
		if err != nil {
			t.Fatalf("New: %s", err)
		}

		b, err := cmd.Run(context.Background())
		if err != nil {
			t.Fatalf("Unexpected error running process: %v", err)
		}

		if s := string(b); s != "foo bar baz 42" {
			t.Errorf("Got %q, expected foo bar", s)
		}
	})

	t.Run("OK without args", func(t *testing.T) {
		cmd, err := New("echo")
		if err != nil {
			t.Fatalf("New: %s", err)
		}

		b, err := cmd.Run(context.Background())
		if err != nil {
			t.Fatalf("Unexpected error running process: %v", err)
		}

		if s := string(b); s != "\n" {
			t.Errorf("Got %q, expected %q", s, "\n")
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		cmd, err := New("sleep", "1")
		if err != nil {
			t.Fatalf("New: %s", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		if _, err := cmd.Run(ctx); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Got %T (%s), expected context.DeadlineExceeded", err, err)
		}
	})
}
