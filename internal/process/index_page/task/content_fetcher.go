package task

import "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"

type ContentProvider interface {
	GetPageContent(lang string) (*gateway_service.PageContent, error)
}

type ContentFetcherTask struct {
	contentProvider ContentProvider
}

func NewContentFetcherTask(contentProvider ContentProvider) *ContentFetcherTask {
	return &ContentFetcherTask{
		contentProvider: contentProvider,
	}
}

func (t *ContentFetcherTask) Fetch(lang string) (*gateway_service.PageContent, error) {
	return t.contentProvider.GetPageContent(lang)
}
