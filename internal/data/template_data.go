package data

import "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/adrianjanczenia.dev_content-service"

type TemplateData struct {
	Lang    string
	Content adrianjanczenia_dev_content_service.PageContent
}
