package get_cv_link

type CVTaskExecutor interface {
	Execute(password, lang string) (string, error)
}

type Process struct {
	taskExecutor CVTaskExecutor
}

func NewProcess(taskExecutor CVTaskExecutor) *Process {
	return &Process{
		taskExecutor: taskExecutor,
	}
}

func (p *Process) Execute(password, lang string) (string, error) {
	return p.taskExecutor.Execute(password, lang)
}
