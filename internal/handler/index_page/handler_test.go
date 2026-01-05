package index_page

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type mockIndexPageProcess struct {
	processFunc func(ctx context.Context, lang string) (*data.TemplateData, error)
}

func (m *mockIndexPageProcess) Process(ctx context.Context, lang string) (*data.TemplateData, error) {
	return m.processFunc(ctx, lang)
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

func TestHandler_IndexPage(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		url             string
		processFunc     func(context.Context, string) (*data.TemplateData, error)
		wantRenderName  string
		wantErrorCalled bool
	}{
		{
			name:   "successful handle",
			method: http.MethodGet,
			url:    "/",
			processFunc: func(ctx context.Context, l string) (*data.TemplateData, error) {
				return &data.TemplateData{Lang: l}, nil
			},
			wantRenderName:  "index",
			wantErrorCalled: false,
		},
		{
			name:            "method not allowed",
			method:          http.MethodPost,
			url:             "/",
			wantErrorCalled: true,
		},
		{
			name:   "process error",
			method: http.MethodGet,
			url:    "/",
			processFunc: func(ctx context.Context, l string) (*data.TemplateData, error) {
				return nil, errors.New("fail")
			},
			wantErrorCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var renderName string
			var errorCalled bool

			mRenderer := &mockRenderer{
				renderFunc: func(w http.ResponseWriter, name string, data any) { renderName = name },
				renderErrorFunc: func(w http.ResponseWriter, templateName string, appErr *appErrors.AppError, lang string, content *gateway_service.PageContent) {
					errorCalled = true
				},
			}

			h := NewHandler(&mockIndexPageProcess{processFunc: tt.processFunc}, mRenderer)
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			h.Handle(w, req)

			if tt.wantErrorCalled != errorCalled {
				t.Errorf("Handle() errorCalled = %v, want %v", errorCalled, tt.wantErrorCalled)
			}
			if tt.wantRenderName != "" && renderName != tt.wantRenderName {
				t.Errorf("Handle() renderName = %s, want %s", renderName, tt.wantRenderName)
			}
		})
	}
}
