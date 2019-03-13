package mock

import "context"

type Runner struct {
	RunFunc func(context.Context) ([]byte, error)
}

func (r *Runner) Run(ctx context.Context) ([]byte, error) {
	return r.RunFunc(ctx)
}
