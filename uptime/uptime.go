// Package uptime wraps os/exec to easily execute the uptime command and get
// its output.
package uptime

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
func New() *command {
	path, err := exec.LookPath("uptime")
	if err != nil {
		panic("Unable to locate uptime binary: " + err.Error())
	}
	return &command{path, nil}
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

	b, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("ioutil: ReadAll: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("os/exec: Cmd.Wait: %v", err)
	}

	return b, nil
}
