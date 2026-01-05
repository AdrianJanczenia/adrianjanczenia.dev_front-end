package privacy_policy

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type mockPrivacyPolicyProcess struct {
	processFunc func(lang string) (*data.TemplateData, error)
}

func (m *mockPrivacyPolicyProcess) Process(lang string) (*data.TemplateData, error) {
	return m.processFunc(lang)
}

type mockRenderer struct {
	renderFunc      func(w http.ResponseWriter, name string, data any)
	renderErrorFunc func(w http.ResponseWriter, templateName string, appErr *appErrors.AppError, lang string, content *gateway_service.PageContent)
}

func (m *mockRenderer) Render(w http.ResponseWriter, name string, data any) {
	m.renderFunc(w, name, data)
}

func (m *mockRenderer) RenderError(w http.ResponseWriter, templateName string, appErr *appErrors.AppError, lang string, content *gateway_service.PageContent) {
	m.renderErrorFunc(w, templateName, appErr, lang, content)
}

func TestHandler_PrivacyPolicy(t *testing.T) {
	t.Run("successful render", func(t *testing.T) {
		var renderCalled bool
		r := &mockRenderer{renderFunc: func(w http.ResponseWriter, n string, d any) {
			renderCalled = true
			td := d.(*data.TemplateData)
			if !td.IsPrivacyPage {
				t.Error("expected IsPrivacyPage to be true")
			}
		}}
		p := &mockPrivacyPolicyProcess{processFunc: func(l string) (*data.TemplateData, error) {
			return &data.TemplateData{}, nil
		}}
		h := NewHandler(p, r)
		req := httptest.NewRequest(http.MethodGet, "/privacy-policy", nil)
		h.Handle(httptest.NewRecorder(), req)
		if !renderCalled {
			t.Error("render not called")
		}
	})

	t.Run("process error", func(t *testing.T) {
		var errorCalled bool
		r := &mockRenderer{renderErrorFunc: func(w http.ResponseWriter, tn string, ae *appErrors.AppError, l string, c *gateway_service.PageContent) {
			errorCalled = true
		}}
		p := &mockPrivacyPolicyProcess{processFunc: func(l string) (*data.TemplateData, error) {
			return nil, errors.New("fail")
		}}
		h := NewHandler(p, r)
		req := httptest.NewRequest(http.MethodGet, "/privacy-policy", nil)
		h.Handle(httptest.NewRecorder(), req)
		if !errorCalled {
			t.Error("error render not called")
		}
	})
}
