package index_page

import (
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type ContentFetcherTask interface {
	Execute(lang string) (*gateway_service.PageContent, error)
}

type Process struct {
	contentFetcher ContentFetcherTask
}

func NewProcess(contentFetcher ContentFetcherTask) *Process {
	return &Process{
		contentFetcher: contentFetcher,
	}
}

func (p *Process) Process(lang string) (*data.TemplateData, error) {
	content, err := p.contentFetcher.Execute(lang)
	if err != nil {
		return nil, err
	}

	templateData := &data.TemplateData{
		Lang:    lang,
		Content: content,
	}

	return templateData, nil
}
