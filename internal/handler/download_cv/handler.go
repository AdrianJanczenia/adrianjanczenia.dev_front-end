package cv_download

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/renderer"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type DownloadCVProcess interface {
	Process(ctx context.Context, token, lang string) (io.ReadCloser, string, error)
}

type ContentProvider interface {
	GetPageContent(ctx context.Context, lang string) (*gateway_service.PageContent, error)
}

type Handler struct {
	process       DownloadCVProcess
	contentClient ContentProvider
	renderer      renderer.Renderer
}

func NewHandler(process DownloadCVProcess, contentClient ContentProvider, renderer renderer.Renderer) *Handler {
	return &Handler{
		process:       process,
		contentClient: contentClient,
		renderer:      renderer,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		appErrors.WriteJSON(w, appErrors.ErrMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	lang := r.URL.Query().Get("lang")

	if token == "" || lang == "" {
		h.renderErrorPage(r.Context(), w, lang, appErrors.ErrInvalidInput)
		return
	}

	stream, contentType, err := h.process.Process(r.Context(), token, lang)
	if err != nil {
		var appErr *appErrors.AppError
		if !errors.As(err, &appErr) {
			appErr = appErrors.ErrCVExpired
		}
		h.renderErrorPage(r.Context(), w, lang, appErr)
		return
	}
	defer stream.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "inline; filename=\"cv_adrian_janczenia.pdf\"")

	_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})
	if _, err := io.Copy(w, stream); err != nil {
		log.Printf("ERROR: CV stream interrupted: %v", err)
	}
}

func (h *Handler) renderErrorPage(ctx context.Context, w http.ResponseWriter, lang string, appErr *appErrors.AppError) {
	content, _ := h.contentClient.GetPageContent(ctx, lang)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.renderer.RenderError(w, "cv_error", appErr, lang, content)
}
