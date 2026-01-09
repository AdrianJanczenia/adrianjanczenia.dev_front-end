package task

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type mockGatewayClient struct {
	getPowFunc func(ctx context.Context) (map[string]string, error)
}

func (m *mockGatewayClient) GetPow(ctx context.Context) (map[string]string, error) {
	return m.getPowFunc(ctx)
}

func TestFetchPowTask_Execute(t *testing.T) {
	tests := []struct {
		name       string
		getPowFunc func(context.Context) (map[string]string, error)
		want       map[string]string
		wantErr    bool
	}{
		{
			name: "success",
			getPowFunc: func(ctx context.Context) (map[string]string, error) {
				return map[string]string{"seed": "test-seed", "signature": "test-sig"}, nil
			},
			want:    map[string]string{"seed": "test-seed", "signature": "test-sig"},
			wantErr: false,
		},
		{
			name: "error",
			getPowFunc: func(ctx context.Context) (map[string]string, error) {
				return nil, errors.New("gateway error")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockGatewayClient{getPowFunc: tt.getPowFunc}
			task := NewFetchPowTask(m)
			got, err := task.Execute(context.Background())
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
