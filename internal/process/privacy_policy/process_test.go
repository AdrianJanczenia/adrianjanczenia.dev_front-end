package privacy_policy

import (
	"context"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
)

type mockFetchPrivacyContentTask struct {
	executeFunc func(ctx context.Context, lang string) (*data.TemplateData, error)
}

func (m *mockFetchPrivacyContentTask) Execute(ctx context.Context, lang string) (*data.TemplateData, error) {
	return m.executeFunc(ctx, lang)
}

func TestProcess_PrivacyPolicy(t *testing.T) {
	m := &mockFetchPrivacyContentTask{executeFunc: func(ctx context.Context, l string) (*data.TemplateData, error) {
		return &data.TemplateData{Lang: l}, nil
	}}
	p := NewProcess(m)
	res, _ := p.Process(context.Background(), "pl")
	if res.Lang != "pl" {
		t.Errorf("expected pl, got %s", res.Lang)
	}
}
