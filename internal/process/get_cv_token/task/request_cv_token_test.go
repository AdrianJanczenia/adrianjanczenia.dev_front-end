package task

import (
	"context"
	"errors"
	"testing"
)

type mockGatewayClient struct {
	requestCVTokenFunc func(ctx context.Context, password, lang, captchaID string) (string, error)
}

func (m *mockGatewayClient) RequestCVToken(ctx context.Context, password, lang, captchaID string) (string, error) {
	return m.requestCVTokenFunc(ctx, password, lang, captchaID)
}

func TestRequestCVTokenTask_Execute(t *testing.T) {
	tests := []struct {
		name               string
		requestCVTokenFunc func(context.Context, string, string, string) (string, error)
		wantErr            bool
	}{
		{
			name:               "success",
			requestCVTokenFunc: func(ctx context.Context, p, l, c string) (string, error) { return "token", nil },
			wantErr:            false,
		},
		{
			name:               "error",
			requestCVTokenFunc: func(ctx context.Context, p, l, c string) (string, error) { return "", errors.New("fail") },
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockGatewayClient{requestCVTokenFunc: tt.requestCVTokenFunc}
			task := NewRequestCVTokenTask(m)
			_, err := task.Execute(context.Background(), "pass", "pl", "123")
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
