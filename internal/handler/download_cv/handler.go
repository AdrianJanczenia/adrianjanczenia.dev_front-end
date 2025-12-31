package cv_download

import (
	"io"
	"net/http"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/renderer"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type Executor interface {
	Execute(token, lang string) (io.ReadCloser, string, int, error)
}

type ContentProvider interface {
	GetPageContent(lang string) (*gateway_service.PageContent, error)
}

type Handler struct {
	processExecutor Executor
	contentClient   ContentProvider
	renderer        renderer.Renderer
}

func NewHandler(e Executor, c ContentProvider, r renderer.Renderer) *Handler {
	return &Handler{
		processExecutor: e,
		contentClient:   c,
		renderer:        r,
	}
}

func (h *Handler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	lang := r.URL.Query().Get("lang")

	if token == "" || lang == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stream, contentType, status, err := h.processExecutor.Execute(token, lang)
	if err != nil {
		h.renderErrorPage(w, lang, status)
		return
	}
	defer stream.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "inline; filename=\"cv.pdf\"")

	io.Copy(w, stream)
}

func (h *Handler) renderErrorPage(w http.ResponseWriter, lang string, status int) {
	content, err := h.contentClient.GetPageContent(lang)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	errorKey := "error_cv_server"
	if status == http.StatusGone || status == http.StatusNotFound {
		errorKey = "error_cv_expired"
	}

	data := struct {
		Lang         string
		Content      *gateway_service.PageContent
		ErrorMessage string
	}{
		Lang:         lang,
		Content:      content,
		ErrorMessage: content.Translations[errorKey],
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	h.renderer.Render(w, "cv_error", data)
}
