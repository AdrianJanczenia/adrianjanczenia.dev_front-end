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

func TestContentFetcherTask_Execute(t *testing.T) {
	tests := []struct {
		name               string
		getPageContentFunc func(context.Context, string) (*gateway_service.PageContent, error)
		wantErr            bool
	}{
		{
			name: "success",
			getPageContentFunc: func(ctx context.Context, l string) (*gateway_service.PageContent, error) {
				return &gateway_service.PageContent{}, nil
			},
			wantErr: false,
		},
		{
			name: "error",
			getPageContentFunc: func(ctx context.Context, l string) (*gateway_service.PageContent, error) {
				return nil, errors.New("fail")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockGatewayClient{getPageContentFunc: tt.getPageContentFunc}
			task := NewContentFetcherTask(m)
			_, err := task.Execute(context.Background(), "pl")
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
