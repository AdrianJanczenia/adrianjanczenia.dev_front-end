package task

import (
	"context"
)

type GatewayClient interface {
	VerifyCaptcha(ctx context.Context, captchaId, captchaValue string) (string, error)
}

type VerifyTask struct {
	client GatewayClient
}

func NewVerifyCaptchaTask(c GatewayClient) *VerifyTask {
	return &VerifyTask{client: c}
}

func (t *VerifyTask) Execute(ctx context.Context, captchaID, captchaValue string) (string, error) {
	return t.client.VerifyCaptcha(ctx, captchaID, captchaValue)
}
