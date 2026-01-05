package privacy_policy

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/renderer"
)

type Executor interface {
	Execute(lang string) (*data.TemplateData, error)
}

type Handler struct {
	processExecutor Executor
	renderer        renderer.Renderer
}

func NewHandler(processExecutor Executor, renderer renderer.Renderer) *Handler {
	return &Handler{
		processExecutor: processExecutor,
		renderer:        renderer,
	}
}

func (h *Handler) HandlePrivacyPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.renderer.RenderError(w, "error", appErrors.ErrMethodNotAllowed, "pl", nil)
		return
	}

	lang := "pl"
	if strings.Contains(r.URL.Path, "privacy-policy") {
		lang = "en"
	}
	if qLang := r.URL.Query().Get("lang"); qLang != "" {
		lang = qLang
	}

	templateData, err := h.processExecutor.Execute(lang)
	if err != nil {
		log.Printf("ERROR: failed to execute privacy policy process: %s", strings.ReplaceAll(err.Error(), "\n", " "))

		var appErr *appErrors.AppError
		if !errors.As(err, &appErr) {
			appErr = appErrors.ErrServiceUnavailable
		}

		h.renderer.RenderError(w, "error", appErr, lang, nil)
		return
	}

	templateData.IsPrivacyPage = true
	h.renderer.Render(w, "privacy", templateData)
}
