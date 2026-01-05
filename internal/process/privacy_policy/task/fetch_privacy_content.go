package task

import (
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type GatewayClient interface {
	GetPageContent(lang string) (*gateway_service.PageContent, error)
}

type FetchPrivacyContentTask struct {
	gatewayClient GatewayClient
}

func NewFetchPrivacyContentTask(gatewayClient GatewayClient) *FetchPrivacyContentTask {
	return &FetchPrivacyContentTask{
		gatewayClient: gatewayClient,
	}
}

func (t *FetchPrivacyContentTask) Execute(lang string) (*data.TemplateData, error) {
	content, err := t.gatewayClient.GetPageContent(lang)
	if err != nil {
		return nil, err
	}

	return &data.TemplateData{
		Lang:    lang,
		Content: content,
	}, nil
}
