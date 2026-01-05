package privacy_policy

import (
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
)

type FetchPrivacyContentTask interface {
	Execute(lang string) (*data.TemplateData, error)
}

type Process struct {
	fetchPrivacyContentTask FetchPrivacyContentTask
}

func NewProcess(fetchPrivacyContentTask FetchPrivacyContentTask) *Process {
	return &Process{
		fetchPrivacyContentTask: fetchPrivacyContentTask,
	}
}

func (p *Process) Process(lang string) (*data.TemplateData, error) {
	return p.fetchPrivacyContentTask.Execute(lang)
}
