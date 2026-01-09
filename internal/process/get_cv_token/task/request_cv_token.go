package task

import (
	"context"
)

type GatewayClient interface {
	RequestCVToken(ctx context.Context, password, lang, captchaID string) (string, error)
}

type RequestCVTokenTask struct {
	gatewayClient GatewayClient
}

func NewRequestCVTokenTask(gatewayClient GatewayClient) *RequestCVTokenTask {
	return &RequestCVTokenTask{
		gatewayClient: gatewayClient,
	}
}

func (t *RequestCVTokenTask) Execute(ctx context.Context, password, lang, captchaID string) (string, error) {
	return t.gatewayClient.RequestCVToken(ctx, password, lang, captchaID)
}
