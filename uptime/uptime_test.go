package uptime

import (
	"context"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		cmd := command{
			name: "echo",
			arg:  []string{"-n", "foo bar"},
		}

		b, err := cmd.Run(context.TODO())
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
			arg:  []string{"3"},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
		defer cancel()

		if _, err := cmd.Run(ctx); err == nil {
			t.Fatalf("Got nil, expected error")
		}
	})
}
