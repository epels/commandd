package command

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Run("OK with args", func(t *testing.T) {
		cmd, err := New("echo", "-n", "foo", "bar", "baz 42")
		if err != nil {
			t.Fatalf("New: %s", err)
		}

		var sb strings.Builder
		if err = cmd.Run(context.Background(), &sb); err != nil {
			t.Fatalf("Unexpected error running process: %v", err)
		}

		if s := sb.String(); s != "foo bar baz 42" {
			t.Errorf("Got %q, expected foo bar", s)
		}
	})

	t.Run("OK without args", func(t *testing.T) {
		cmd, err := New("echo")
		if err != nil {
			t.Fatalf("New: %s", err)
		}

		var sb strings.Builder
		if err = cmd.Run(context.Background(), &sb); err != nil {
			t.Fatalf("Unexpected error running process: %v", err)
		}

		if s := sb.String(); s != "\n" {
			t.Errorf("Got %q, expected %q", s, "\n")
		}
	})

	t.Run("Cancel", func(t *testing.T) {
		cmd, err := New("sleep", "1")
		if err != nil {
			t.Fatalf("New: %s", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(25 * time.Millisecond)
			cancel()
		}()

		var sb strings.Builder
		if err := cmd.Run(ctx, &sb); !errors.Is(err, context.Canceled) {
			t.Fatalf("Got %T (%s), expected context.DeadlineExceeded", err, err)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		cmd, err := New("sleep", "1")
		if err != nil {
			t.Fatalf("New: %s", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		var sb strings.Builder
		if err := cmd.Run(ctx, &sb); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Got %T (%s), expected context.DeadlineExceeded", err, err)
		}
	})
}
