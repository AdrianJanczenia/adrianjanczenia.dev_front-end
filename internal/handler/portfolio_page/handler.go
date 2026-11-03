package portfolio_page

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/renderer"
)

type PageProcess interface {
	Process(context.Context, string) (*data.TemplateData, error)
}

type Handler struct {
	process  PageProcess
	renderer renderer.Renderer
}

func NewHandler(process PageProcess, pageRenderer renderer.Renderer) *Handler {
	return &Handler{process: process, renderer: pageRenderer}
}

func (h *Handler) Handle(templateName, page, caseID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := "pl"

		if r.URL.Query().Get("lang") == "en" {
			lang = "en"
		}

		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			h.renderer.RenderError(w, "error", appErrors.ErrMethodNotAllowed, lang, nil)

			return
		}

		content, err := h.process.Process(r.Context(), lang)
		if err != nil {
			log.Printf("ERROR: failed to load portfolio page: %s", strings.ReplaceAll(err.Error(), "\n", " "))

			var appErr *appErrors.AppError
			if !errors.As(err, &appErr) {
				appErr = appErrors.ErrServiceUnavailable
			}

			h.renderer.RenderError(w, "error", appErr, lang, nil)

			return
		}

		content.Page = page
		content.CaseID = caseID
		content.Path = r.URL.Path

		h.renderer.Render(w, templateName, content)
	}
}
