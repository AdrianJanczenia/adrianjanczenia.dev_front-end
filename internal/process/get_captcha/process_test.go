package get_captcha

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type mockFetchCaptchaTask struct {
	executeFunc func(ctx context.Context, seed, signature, nonce string) (map[string]string, error)
}

func (m *mockFetchCaptchaTask) Execute(ctx context.Context, seed, signature, nonce string) (map[string]string, error) {
	return m.executeFunc(ctx, seed, signature, nonce)
}

func TestProcess_GetCaptcha(t *testing.T) {
	tests := []struct {
		name        string
		executeFunc func(context.Context, string, string, string) (map[string]string, error)
		want        map[string]string
		wantErr     bool
	}{
		{
			name: "success",
			executeFunc: func(ctx context.Context, s, sig, n string) (map[string]string, error) {
				return map[string]string{"captchaId": "123", "captchaImg": "base64data"}, nil
			},
			want:    map[string]string{"captchaId": "123", "captchaImg": "base64data"},
			wantErr: false,
		},
		{
			name: "error",
			executeFunc: func(ctx context.Context, s, sig, n string) (map[string]string, error) {
				return nil, errors.New("task error")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockFetchCaptchaTask{executeFunc: tt.executeFunc}
			p := NewProcess(m)
			got, err := p.Process(context.Background(), "seed", "sig", "nonce")
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Process() got = %v, want %v", got, tt.want)
			}
		})
	}
}
