// Package command wraps os/exec to easily execute the command command and get
// its output.
package command

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

type command struct {
	name string
	arg  []string
}

// New may panic if the uptime binary is not in the PATH.
func New(name string, args ...string) (*command, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, fmt.Errorf("unable to locate executable: %s", err)
	}
	return &command{path, args}, nil
}

// Run starts the process and waits for it to complete, then returns the
// standard output stream bytes.
func (c *command) Run(ctx context.Context, w io.Writer) error {
	cmd := exec.CommandContext(ctx, c.name, c.arg...)
	cmd.Stdout = w

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("os/exec: Cmd.Start: %v", err)
	}

	err := cmd.Wait()
	// cmd.Wait() does not propagate the context error (e.g. Canceled or
	// DeadlineExceeded), but instead returns an error indicating the process
	// was killed. In this case we care about the context error more.
	// See also: https://github.com/golang/go/issues/21880.
	if ctxerr := ctx.Err(); ctxerr != nil {
		return ctxerr
	}
	if err != nil {
		return fmt.Errorf("os/exec: Cmd.Wait: %v", err)
	}

	return nil
}
