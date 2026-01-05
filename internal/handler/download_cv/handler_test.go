package cv_download

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type mockDownloadCVProcess struct {
	processFunc func(ctx context.Context, token, lang string) (io.ReadCloser, string, error)
}

func (m *mockDownloadCVProcess) Process(ctx context.Context, token, lang string) (io.ReadCloser, string, error) {
	return m.processFunc(ctx, token, lang)
}

type mockContentProvider struct {
	getPageContentFunc func(ctx context.Context, lang string) (*gateway_service.PageContent, error)
}

func (m *mockContentProvider) GetPageContent(ctx context.Context, lang string) (*gateway_service.PageContent, error) {
	return m.getPageContentFunc(ctx, lang)
}

type mockRenderer struct {
	renderErrorFunc func(w http.ResponseWriter, templateName string, appErr *appErrors.AppError, lang string, content *gateway_service.PageContent)
}

func (m *mockRenderer) Render(w http.ResponseWriter, name string, data any) {}
func (m *mockRenderer) RenderError(w http.ResponseWriter, templateName string, appErr *appErrors.AppError, lang string, content *gateway_service.PageContent) {
	m.renderErrorFunc(w, templateName, appErr, lang, content)
}

func TestHandler_DownloadCV(t *testing.T) {
	t.Run("missing params calls render error", func(t *testing.T) {
		var errorCalled bool
		r := &mockRenderer{renderErrorFunc: func(w http.ResponseWriter, tn string, ae *appErrors.AppError, l string, c *gateway_service.PageContent) {
			errorCalled = true
		}}
		cp := &mockContentProvider{getPageContentFunc: func(ctx context.Context, l string) (*gateway_service.PageContent, error) { return nil, nil }}
		h := NewHandler(nil, cp, r)
		req := httptest.NewRequest(http.MethodGet, "/cv", nil)
		h.Handle(httptest.NewRecorder(), req)
		if !errorCalled {
			t.Error("expected RenderError to be called")
		}
	})

	t.Run("successful download", func(t *testing.T) {
		p := &mockDownloadCVProcess{processFunc: func(ctx context.Context, t, l string) (io.ReadCloser, string, error) {
			return io.NopCloser(strings.NewReader("pdf data")), "application/pdf", nil
		}}
		h := NewHandler(p, nil, nil)
		req := httptest.NewRequest(http.MethodGet, "/cv?token=tok&lang=pl", nil)
		w := httptest.NewRecorder()
		h.Handle(w, req)
		if w.Header().Get("Content-Type") != "application/pdf" {
			t.Errorf("expected application/pdf, got %s", w.Header().Get("Content-Type"))
		}
		if w.Body.String() != "pdf data" {
			t.Errorf("expected pdf data, got %s", w.Body.String())
		}
	})
}
