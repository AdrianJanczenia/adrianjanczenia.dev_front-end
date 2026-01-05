package get_cv_token

import "context"

type RequestCVTokenTask interface {
	Execute(ctx context.Context, password, lang string) (string, error)
}

type Process struct {
	cvTokenTask RequestCVTokenTask
}

func NewProcess(task RequestCVTokenTask) *Process {
	return &Process{
		cvTokenTask: task,
	}
}

func (p *Process) Process(ctx context.Context, password, lang string) (string, error) {
	return p.cvTokenTask.Execute(ctx, password, lang)
}
