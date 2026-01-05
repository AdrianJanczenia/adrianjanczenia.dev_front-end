package privacy_policy

import (
	"context"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
)

type FetchPrivacyContentTask interface {
	Execute(ctx context.Context, lang string) (*data.TemplateData, error)
}

type Process struct {
	fetchPrivacyContentTask FetchPrivacyContentTask
}

func NewProcess(fetchPrivacyContentTask FetchPrivacyContentTask) *Process {
	return &Process{
		fetchPrivacyContentTask: fetchPrivacyContentTask,
	}
}

func (p *Process) Process(ctx context.Context, lang string) (*data.TemplateData, error) {
	return p.fetchPrivacyContentTask.Execute(ctx, lang)
}
