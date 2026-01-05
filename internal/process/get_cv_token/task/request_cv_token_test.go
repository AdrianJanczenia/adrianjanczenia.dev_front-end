package task

import (
	"errors"
	"testing"
)

type mockGatewayClient struct {
	requestCVTokenFunc func(password, lang string) (string, error)
}

func (m *mockGatewayClient) RequestCVToken(password, lang string) (string, error) {
	return m.requestCVTokenFunc(password, lang)
}

func TestRequestCVTokenTask_Execute(t *testing.T) {
	tests := []struct {
		name               string
		requestCVTokenFunc func(string, string) (string, error)
		wantErr            bool
	}{
		{
			name:               "success",
			requestCVTokenFunc: func(p, l string) (string, error) { return "token", nil },
			wantErr:            false,
		},
		{
			name:               "error",
			requestCVTokenFunc: func(p, l string) (string, error) { return "", errors.New("fail") },
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockGatewayClient{requestCVTokenFunc: tt.requestCVTokenFunc}
			task := NewRequestCVTokenTask(m)
			_, err := task.Execute("pass", "pl")
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
