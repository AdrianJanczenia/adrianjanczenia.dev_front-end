package task

import (
	"context"
	"errors"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type mockGatewayClient struct {
	getPageContentFunc func(ctx context.Context, lang string) (*gateway_service.PageContent, error)
}

func (m *mockGatewayClient) GetPageContent(ctx context.Context, lang string) (*gateway_service.PageContent, error) {
	return m.getPageContentFunc(ctx, lang)
}

func TestFetchPrivacyContentTask_Execute(t *testing.T) {
	m := &mockGatewayClient{getPageContentFunc: func(ctx context.Context, l string) (*gateway_service.PageContent, error) {
		if l == "pl" {
			return &gateway_service.PageContent{}, nil
		}
		return nil, errors.New("fail")
	}}
	task := NewFetchPrivacyContentTask(m)

	if _, err := task.Execute(context.Background(), "pl"); err != nil {
		t.Errorf("expected success, got %v", err)
	}
	if _, err := task.Execute(context.Background(), "en"); err == nil {
		t.Error("expected error")
	}
}
