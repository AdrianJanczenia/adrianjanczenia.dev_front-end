package task

import "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"

type GatewayClient interface {
	GetPageContent(lang string) (*gateway_service.PageContent, error)
}

type ContentFetcherTask struct {
	gatewayClient GatewayClient
}

func NewContentFetcherTask(gatewayClient GatewayClient) *ContentFetcherTask {
	return &ContentFetcherTask{
		gatewayClient: gatewayClient,
	}
}

func (t *ContentFetcherTask) Execute(lang string) (*gateway_service.PageContent, error) {
	return t.gatewayClient.GetPageContent(lang)
}
