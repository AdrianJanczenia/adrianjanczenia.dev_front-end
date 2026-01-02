package data

import "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"

type TemplateData struct {
	Lang          string
	IsPrivacyPage bool
	Content       *gateway_service.PageContent
}
