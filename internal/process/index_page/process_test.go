package index_page

import (
	"context"
	"errors"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type mockContentFetcherTask struct {
	executeFunc func(ctx context.Context, lang string) (*gateway_service.PageContent, error)
}

func (m *mockContentFetcherTask) Execute(ctx context.Context, lang string) (*gateway_service.PageContent, error) {
	return m.executeFunc(ctx, lang)
}

func TestProcess_IndexPage(t *testing.T) {
	tests := []struct {
		name        string
		executeFunc func(context.Context, string) (*gateway_service.PageContent, error)
		wantErr     bool
	}{
		{
			name: "success",
			executeFunc: func(ctx context.Context, l string) (*gateway_service.PageContent, error) {
				return &gateway_service.PageContent{}, nil
			},
			wantErr: false,
		},
		{
			name: "task error",
			executeFunc: func(ctx context.Context, l string) (*gateway_service.PageContent, error) {
				return nil, errors.New("fail")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockContentFetcherTask{executeFunc: tt.executeFunc}
			p := NewProcess(m)
			res, err := p.Process(context.Background(), "pl")
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && res.Lang != "pl" {
				t.Errorf("Process() lang = %s, want pl", res.Lang)
			}
		})
	}
}
