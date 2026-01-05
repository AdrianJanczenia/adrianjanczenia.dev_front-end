package task

import (
	"context"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type GatewayClient interface {
	GetPageContent(ctx context.Context, lang string) (*gateway_service.PageContent, error)
}

type ContentFetcherTask struct {
	gatewayClient GatewayClient
}

func NewContentFetcherTask(gatewayClient GatewayClient) *ContentFetcherTask {
	return &ContentFetcherTask{
		gatewayClient: gatewayClient,
	}
}

func (t *ContentFetcherTask) Execute(ctx context.Context, lang string) (*gateway_service.PageContent, error) {
	return t.gatewayClient.GetPageContent(ctx, lang)
}
