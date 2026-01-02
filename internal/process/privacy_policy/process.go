package privacy_policy

import (
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
)

type Task interface {
	Run(lang string) (*data.TemplateData, error)
}

type Process struct {
	fetchTask Task
}

func NewProcess(fetchTask Task) *Process {
	return &Process{
		fetchTask: fetchTask,
	}
}

func (p *Process) Execute(lang string) (*data.TemplateData, error) {
	return p.fetchTask.Run(lang)
}
