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

type PrivacyPolicyProcess interface {
	Process(lang string) (*data.TemplateData, error)
}

type Handler struct {
	process  PrivacyPolicyProcess
	renderer renderer.Renderer
}

func NewHandler(process PrivacyPolicyProcess, renderer renderer.Renderer) *Handler {
	return &Handler{
		process:  process,
		renderer: renderer,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
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

	templateData, err := h.process.Process(lang)
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
