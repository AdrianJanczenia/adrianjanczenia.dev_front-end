package task

import (
	"errors"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type mockGatewayClient struct {
	getPageContentFunc func(lang string) (*gateway_service.PageContent, error)
}

func (m *mockGatewayClient) GetPageContent(lang string) (*gateway_service.PageContent, error) {
	return m.getPageContentFunc(lang)
}

func TestFetchPrivacyContentTask_Execute(t *testing.T) {
	m := &mockGatewayClient{getPageContentFunc: func(l string) (*gateway_service.PageContent, error) {
		if l == "pl" {
			return &gateway_service.PageContent{}, nil
		}
		return nil, errors.New("fail")
	}}
	task := NewFetchPrivacyContentTask(m)

	if _, err := task.Execute("pl"); err != nil {
		t.Errorf("expected success, got %v", err)
	}
	if _, err := task.Execute("en"); err == nil {
		t.Error("expected error")
	}
}
