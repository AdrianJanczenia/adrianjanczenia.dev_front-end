package privacy_policy

import (
	"log"
	"net/http"
	"strings"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
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
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	lang := "en"
	if strings.Contains(r.URL.Path, "polityka-prywatnosci") {
		lang = "pl"
	}
	if qLang := r.URL.Query().Get("lang"); qLang != "" {
		lang = qLang
	}

	templateData, err := h.processExecutor.Execute(lang)
	if err != nil {
		log.Printf("ERROR: failed to execute privacy policy process: %s", strings.ReplaceAll(err.Error(), "\n", " "))

		titles := map[string]string{"pl": "Wystąpił błąd", "en": "An error occurred"}
		messages := map[string]string{
			"pl": "Pracujemy nad rozwiązaniem problemu. Spróbuj ponownie później.",
			"en": "We are working on fixing the problem. Please try again later.",
		}

		errorData := struct {
			Lang    string
			Title   string
			Message string
		}{
			Lang:    lang,
			Title:   titles[lang],
			Message: messages[lang],
		}

		h.renderer.Render(w, "error", errorData)
		return
	}

	templateData.IsPrivacyPage = true
	h.renderer.Render(w, "privacy", templateData)
}
