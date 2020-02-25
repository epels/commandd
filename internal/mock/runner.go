package mock

import (
	"context"
	"io"
)

type Runner struct {
	RunFunc func(context.Context, io.Writer) error
}

func (r *Runner) Run(ctx context.Context, w io.Writer) error {
	return r.RunFunc(ctx, w)
}
