package task

import (
	"context"
)

type GatewayClient interface {
	GetCaptcha(ctx context.Context, seed, signature, nonce string) (map[string]string, error)
}

type FetchCaptchaTask struct {
	client GatewayClient
}

func NewFetchCaptchaTask(c GatewayClient) *FetchCaptchaTask {
	return &FetchCaptchaTask{client: c}
}

func (t *FetchCaptchaTask) Execute(ctx context.Context, seed, signature, nonce string) (map[string]string, error) {
	return t.client.GetCaptcha(ctx, seed, signature, nonce)
}
