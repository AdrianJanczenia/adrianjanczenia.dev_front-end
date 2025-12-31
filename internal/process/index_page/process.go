package index_page

import (
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type ContentFetcher interface {
	Fetch(lang string) (*gateway_service.PageContent, error)
}

type Process struct {
	contentFetcher ContentFetcher
}

func NewProcess(contentFetcher ContentFetcher) *Process {
	return &Process{
		contentFetcher: contentFetcher,
	}
}

func (p *Process) Execute(lang string) (*data.TemplateData, error) {
	content, err := p.contentFetcher.Fetch(lang)
	if err != nil {
		return nil, err
	}

	templateData := &data.TemplateData{
		Lang:    lang,
		Content: content,
	}

	return templateData, nil
}
