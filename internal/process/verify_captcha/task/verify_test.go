package task

import (
	"context"
	"errors"
	"testing"
)

type mockGatewayClient struct {
	verifyCaptchaFunc func(ctx context.Context, captchaId, captchaValue string) (string, error)
}

func (m *mockGatewayClient) VerifyCaptcha(ctx context.Context, captchaId, captchaValue string) (string, error) {
	return m.verifyCaptchaFunc(ctx, captchaId, captchaValue)
}

func TestVerifyTask_Execute(t *testing.T) {
	tests := []struct {
		name              string
		verifyCaptchaFunc func(context.Context, string, string) (string, error)
		want              string
		wantErr           bool
	}{
		{
			name: "success",
			verifyCaptchaFunc: func(ctx context.Context, id, val string) (string, error) {
				return "verified-id", nil
			},
			want:    "verified-id",
			wantErr: false,
		},
		{
			name: "error",
			verifyCaptchaFunc: func(ctx context.Context, id, val string) (string, error) {
				return "", errors.New("verification failed")
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockGatewayClient{verifyCaptchaFunc: tt.verifyCaptchaFunc}
			task := NewVerifyCaptchaTask(m)
			got, err := task.Execute(context.Background(), "id", "value")
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
