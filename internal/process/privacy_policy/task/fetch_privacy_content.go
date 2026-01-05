package task

import (
	"context"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type ContentService interface {
	GetPageContent(ctx context.Context, lang string) (*gateway_service.PageContent, error)
}

type FetchPrivacyContentTask struct {
	contentService ContentService
}

func NewFetchPrivacyContentTask(contentService ContentService) *FetchPrivacyContentTask {
	return &FetchPrivacyContentTask{
		contentService: contentService,
	}
}

func (t *FetchPrivacyContentTask) Execute(ctx context.Context, lang string) (*data.TemplateData, error) {
	content, err := t.contentService.GetPageContent(ctx, lang)
	if err != nil {
		return nil, err
	}

	return &data.TemplateData{
		Lang:    lang,
		Content: content,
	}, nil
}
