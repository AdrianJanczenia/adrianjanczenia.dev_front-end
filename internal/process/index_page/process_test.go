package index_page

import (
	"errors"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type mockContentFetcherTask struct {
	executeFunc func(lang string) (*gateway_service.PageContent, error)
}

func (m *mockContentFetcherTask) Execute(lang string) (*gateway_service.PageContent, error) {
	return m.executeFunc(lang)
}

func TestProcess_IndexPage(t *testing.T) {
	tests := []struct {
		name        string
		executeFunc func(string) (*gateway_service.PageContent, error)
		wantErr     bool
	}{
		{
			name: "success",
			executeFunc: func(l string) (*gateway_service.PageContent, error) {
				return &gateway_service.PageContent{}, nil
			},
			wantErr: false,
		},
		{
			name: "task error",
			executeFunc: func(l string) (*gateway_service.PageContent, error) {
				return nil, errors.New("fail")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockContentFetcherTask{executeFunc: tt.executeFunc}
			p := NewProcess(m)
			res, err := p.Process("pl")
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && res.Lang != "pl" {
				t.Errorf("Process() lang = %s, want pl", res.Lang)
			}
		})
	}
}
