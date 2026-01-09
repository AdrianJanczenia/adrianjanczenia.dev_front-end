package get_pow

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type mockFetchPowTask struct {
	executeFunc func(ctx context.Context) (map[string]string, error)
}

func (m *mockFetchPowTask) Execute(ctx context.Context) (map[string]string, error) {
	return m.executeFunc(ctx)
}

func TestProcess_GetPow(t *testing.T) {
	tests := []struct {
		name        string
		executeFunc func(context.Context) (map[string]string, error)
		want        map[string]string
		wantErr     bool
	}{
		{
			name: "success",
			executeFunc: func(ctx context.Context) (map[string]string, error) {
				return map[string]string{"seed": "s1", "signature": "sig1"}, nil
			},
			want:    map[string]string{"seed": "s1", "signature": "sig1"},
			wantErr: false,
		},
		{
			name: "error",
			executeFunc: func(ctx context.Context) (map[string]string, error) {
				return nil, errors.New("task failure")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockFetchPowTask{executeFunc: tt.executeFunc}
			p := NewProcess(m)
			got, err := p.Process(context.Background())
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
