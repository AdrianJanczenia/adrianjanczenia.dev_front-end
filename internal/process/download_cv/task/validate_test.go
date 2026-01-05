package task

import (
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
	"github.com/google/uuid"
)

func TestValidateTask_Execute(t *testing.T) {
	task := NewValidateLinkTask()

	t.Run("valid", func(t *testing.T) {
		token, _ := uuid.NewUUID()
		if err := task.Execute(token.String(), "pl"); err != nil {
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
