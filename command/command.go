// Package command wraps os/exec to easily execute the command command and get
// its output.
package command

import (
	"context"
	"fmt"
	"io/ioutil"
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
func (c *command) Run(ctx context.Context) ([]byte, error) {
	cmd := exec.CommandContext(ctx, c.name, c.arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("os/exec: Cmd.StdoutPipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("os/exec: Cmd.Start: %v", err)
	}

	// @todo: Consider having Run accept an io.Writer and pipe the command's
	//        output to it, so there is no need to read all output into memory
	//        at once.
	b, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("io/ioutil: ReadAll: %v", err)
	}

	err = cmd.Wait()
	// cmd.Wait() does not propagate the context error (e.g. Canceled or
	// DeadlineExceeded), but instead returns an error indicating the process
	// was killed. In this case we care about the context error more.
	// See also: https://github.com/golang/go/issues/21880.
	if cerr := ctx.Err(); cerr != nil {
		return nil, cerr
	}
	if err != nil {
		return nil, fmt.Errorf("os/exec: Cmd.Wait: %v", err)
	}

	return b, nil
}
