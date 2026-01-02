package task

import (
	"fmt"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type ContentProvider interface {
	GetPageContent(lang string) (*gateway_service.PageContent, error)
}

type FetchPrivacyContentTask struct {
	gatewayClient ContentProvider
}

func NewFetchPrivacyContentTask(gatewayClient ContentProvider) *FetchPrivacyContentTask {
	return &FetchPrivacyContentTask{
		gatewayClient: gatewayClient,
	}
}

func (t *FetchPrivacyContentTask) Run(lang string) (*data.TemplateData, error) {
	content, err := t.gatewayClient.GetPageContent(lang)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content from gateway: %w", err)
	}

	return &data.TemplateData{
		Lang:    lang,
		Content: content,
	}, nil
}
