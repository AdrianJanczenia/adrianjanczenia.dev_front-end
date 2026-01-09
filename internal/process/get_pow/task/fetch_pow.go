package task

import (
	"context"
)

type GatewayClient interface {
	GetPow(ctx context.Context) (map[string]string, error)
}

type FetchPowTask struct {
	client GatewayClient
}

func NewFetchPowTask(c GatewayClient) *FetchPowTask {
	return &FetchPowTask{client: c}
}

func (t *FetchPowTask) Execute(ctx context.Context) (map[string]string, error) {
	return t.client.GetPow(ctx)
}
