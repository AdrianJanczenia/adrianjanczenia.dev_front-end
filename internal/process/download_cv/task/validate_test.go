package task

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
)

func TestValidateTask_Execute(t *testing.T) {
	task := NewValidateLinkTask()

	t.Run("valid", func(t *testing.T) {
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

		b := make([]byte, 32)
		for i := range b {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				t.Errorf("rand error %v", err)
			}
			b[i] = charset[num.Int64()]
		}

		token := string(b)
		if err := task.Execute(token, "pl"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("missing token", func(t *testing.T) {
		if err := task.Execute("", "pl"); err != errors.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("missing lang", func(t *testing.T) {
		if err := task.Execute("tok", ""); err != errors.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}
