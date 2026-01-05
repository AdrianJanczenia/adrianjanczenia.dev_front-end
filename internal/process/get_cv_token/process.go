package get_cv_token

type RequestCVTokenTask interface {
	Execute(password, lang string) (string, error)
}

type Process struct {
	cvTokenTask RequestCVTokenTask
}

func NewProcess(task RequestCVTokenTask) *Process {
	return &Process{
		cvTokenTask: task,
	}
}

func (p *Process) Process(password, lang string) (string, error) {
	return p.cvTokenTask.Execute(password, lang)
}
