package task

import (
	"regexp"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
)

type ValidateLinkTask struct{}

func NewValidateLinkTask() *ValidateLinkTask {
	return &ValidateLinkTask{}
}

func (t *ValidateLinkTask) Execute(token, lang string) error {
	if lang != "pl" && lang != "en" {
		return errors.ErrInvalidInput
	}

	match, _ := regexp.MatchString("^[a-fA-F0-9-]{36}$", token)
	if !match {
		return errors.ErrInvalidInput
	}

	return nil
}
