package verify_captcha

import (
	"context"
	"errors"
	"testing"
)

type mockVerifyTask struct {
	executeFunc func(ctx context.Context, captchaID, captchaValue string) (string, error)
}

func (m *mockVerifyTask) Execute(ctx context.Context, captchaID, captchaValue string) (string, error) {
	return m.executeFunc(ctx, captchaID, captchaValue)
}

func TestProcess_VerifyCaptcha(t *testing.T) {
	tests := []struct {
		name        string
		executeFunc func(context.Context, string, string) (string, error)
		want        string
		wantErr     bool
	}{
		{
			name: "success",
			executeFunc: func(ctx context.Context, id, val string) (string, error) {
				return "ok-id", nil
			},
			want:    "ok-id",
			wantErr: false,
		},
		{
			name: "error",
			executeFunc: func(ctx context.Context, id, val string) (string, error) {
				return "", errors.New("task failure")
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockVerifyTask{executeFunc: tt.executeFunc}
			p := NewProcess(m)
			got, err := p.Process(context.Background(), "id", "val")
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Process() got = %v, want %v", got, tt.want)
			}
		})
	}
}
