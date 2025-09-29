package task

import "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/adrianjanczenia.dev_content-service"

type ContentProvider interface {
	GetPageContent(lang string) (*adrianjanczenia_dev_content_service.PageContent, error)
}

// ContentFetcherTask is a task for fetching page content.
type ContentFetcherTask struct {
	contentProvider ContentProvider
}

func NewContentFetcherTask(contentProvider ContentProvider) *ContentFetcherTask {
	return &ContentFetcherTask{
		contentProvider: contentProvider,
	}
}

func (t *ContentFetcherTask) Fetch(lang string) (*adrianjanczenia_dev_content_service.PageContent, error) {
	return t.contentProvider.GetPageContent(lang)
}
