package get_cv_token

import (
	"context"
	"errors"
	"testing"
)

type mockRequestCVTokenTask struct {
	executeFunc func(ctx context.Context, password, lang string) (string, error)
}

func (m *mockRequestCVTokenTask) Execute(ctx context.Context, password, lang string) (string, error) {
	return m.executeFunc(ctx, password, lang)
}

func TestProcess_GetCVToken(t *testing.T) {
	m := &mockRequestCVTokenTask{
		executeFunc: func(ctx context.Context, p, l string) (string, error) {
			if p == "valid" {
				return "tok", nil
			}
			return "", errors.New("invalid")
		},
	}
	p := NewProcess(m)

	t.Run("success", func(t *testing.T) {
		tok, err := p.Process(context.Background(), "valid", "pl")
		if err != nil || tok != "tok" {
			t.Errorf("expected tok, got %v, %v", tok, err)
		}
	})

	t.Run("error", func(t *testing.T) {
		_, err := p.Process(context.Background(), "wrong", "pl")
		if err == nil {
			t.Error("expected error")
		}
	})
}
