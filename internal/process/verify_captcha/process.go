package verify_captcha

import (
	"context"
)

type VerifyTaskTask interface {
	Execute(ctx context.Context, captchaID, captchaValue string) (string, error)
}

type Process struct {
	verifyTask VerifyTaskTask
}

func NewProcess(t VerifyTaskTask) *Process {
	return &Process{verifyTask: t}
}

func (p *Process) Process(ctx context.Context, captchaID, captchaValue string) (string, error) {
	return p.verifyTask.Execute(ctx, captchaID, captchaValue)
}
