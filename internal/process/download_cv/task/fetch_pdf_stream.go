package task

import (
	"context"
	"io"
)

type GatewayClient interface {
	DownloadCVStream(ctx context.Context, token, lang string) (io.ReadCloser, string, error)
}

type FetchPDFStreamTask struct {
	gatewayClient GatewayClient
}

func NewFetchPDFStreamTask(gatewayClient GatewayClient) *FetchPDFStreamTask {
	return &FetchPDFStreamTask{
		gatewayClient: gatewayClient,
	}
}

func (t *FetchPDFStreamTask) Execute(ctx context.Context, token, lang string) (io.ReadCloser, string, error) {
	return t.gatewayClient.DownloadCVStream(ctx, token, lang)
}
