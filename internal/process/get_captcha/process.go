package get_captcha

import (
	"context"
)

type FetchCaptchaTask interface {
	Execute(ctx context.Context, seed, signature, nonce string) (map[string]string, error)
}

type Process struct {
	fetchCaptchaTask FetchCaptchaTask
}

func NewProcess(t FetchCaptchaTask) *Process {
	return &Process{fetchCaptchaTask: t}
}

func (p *Process) Process(ctx context.Context, seed, signature, nonce string) (map[string]string, error) {
	return p.fetchCaptchaTask.Execute(ctx, seed, signature, nonce)
}
