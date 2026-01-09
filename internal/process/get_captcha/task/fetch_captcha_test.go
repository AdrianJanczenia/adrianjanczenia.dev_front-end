package task

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type mockGatewayClient struct {
	getCaptchaFunc func(ctx context.Context, seed, signature, nonce string) (map[string]string, error)
}

func (m *mockGatewayClient) GetCaptcha(ctx context.Context, seed, signature, nonce string) (map[string]string, error) {
	return m.getCaptchaFunc(ctx, seed, signature, nonce)
}

func TestFetchCaptchaTask_Execute(t *testing.T) {
	tests := []struct {
		name           string
		getCaptchaFunc func(context.Context, string, string, string) (map[string]string, error)
		want           map[string]string
		wantErr        bool
	}{
		{
			name: "success",
			getCaptchaFunc: func(ctx context.Context, s, sig, n string) (map[string]string, error) {
				return map[string]string{"captchaId": "cid1", "captchaImg": "img1"}, nil
			},
			want:    map[string]string{"captchaId": "cid1", "captchaImg": "img1"},
			wantErr: false,
		},
		{
			name: "error",
			getCaptchaFunc: func(ctx context.Context, s, sig, n string) (map[string]string, error) {
				return nil, errors.New("gateway error")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockGatewayClient{getCaptchaFunc: tt.getCaptchaFunc}
			task := NewFetchCaptchaTask(m)
			got, err := task.Execute(context.Background(), "s", "sig", "n")
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
