package get_pow

import (
	"context"
)

type FetchPowTask interface {
	Execute(ctx context.Context) (map[string]string, error)
}

type Process struct {
	fetchPowTask FetchPowTask
}

func NewProcess(t FetchPowTask) *Process {
	return &Process{fetchPowTask: t}
}

func (p *Process) Process(ctx context.Context) (map[string]string, error) {
	return p.fetchPowTask.Execute(ctx)
}
