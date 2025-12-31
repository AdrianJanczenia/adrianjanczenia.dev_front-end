package index_page

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

func (h *Handler) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	lang := getLanguage(r)

	templateData, err := h.processExecutor.Execute(lang)
	if err != nil {
		log.Printf("ERROR: failed to execute index process: %s", strings.ReplaceAll(err.Error(), "\n", " "))

		errorData := struct {
			Lang    string
			Title   string
			Message string
		}{
			Lang:    lang,
			Title:   "Internal Server Error",
			Message: "Something went wrong. Please try again later.",
		}

		if templateData != nil && templateData.Content != nil && templateData.Content.Translations["error_title"] != "" {
			errorData.Title = templateData.Content.Translations["error_title"]
			errorData.Message = templateData.Content.Translations["error_message"]
		}

		h.renderer.Render(w, "error", errorData)
		return
	}

	h.renderer.Render(w, "index", templateData)
}

func getLanguage(r *http.Request) string {
	if r.URL.Query().Get("lang") == "en" {
		return "en"
	}
	return "pl"
}
