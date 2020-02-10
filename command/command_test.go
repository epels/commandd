package command

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		cmd := command{
			name: "echo",
			arg:  []string{"-n", "foo bar"},
		}

		b, err := cmd.Run(context.Background())
		if err != nil {
			t.Fatalf("Unexpected error running process: %v", err)
		}

		if s := string(b); s != "foo bar" {
			t.Errorf("Got %q, expected foo bar", s)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		cmd := command{
			name: "sleep",
			arg:  []string{"1"},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		if _, err := cmd.Run(ctx); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Got %T (%s), expected context.DeadlineExceeded", err, err)
		}
	})
}
