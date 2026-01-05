package privacy_policy

import (
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
)

type mockFetchPrivacyContentTask struct {
	executeFunc func(lang string) (*data.TemplateData, error)
}

func (m *mockFetchPrivacyContentTask) Execute(lang string) (*data.TemplateData, error) {
	return m.executeFunc(lang)
}

func TestProcess_PrivacyPolicy(t *testing.T) {
	m := &mockFetchPrivacyContentTask{executeFunc: func(l string) (*data.TemplateData, error) {
		return &data.TemplateData{Lang: l}, nil
	}}
	p := NewProcess(m)
	res, _ := p.Process("pl")
	if res.Lang != "pl" {
		t.Errorf("expected pl, got %s", res.Lang)
	}
}
