package task

import (
	"errors"
	"regexp"
)

type ValidateLinkTask struct{}

func NewValidateLinkTask() *ValidateLinkTask {
	return &ValidateLinkTask{}
}

func (t *ValidateLinkTask) Execute(token, lang string) error {
	if lang != "pl" && lang != "en" {
		return errors.New("invalid_lang")
	}

	match, _ := regexp.MatchString("^[a-fA-F0-9-]{36}$", token)
	if !match {
		return errors.New("invalid_token_format")
	}

	return nil
}
