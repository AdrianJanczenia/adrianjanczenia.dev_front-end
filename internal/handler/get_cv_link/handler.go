package cv

import (
	"encoding/json"
	"net/http"
)

type Executor interface {
	Execute(password, lang string) (string, error)
}

type Handler struct {
	processExecutor Executor
}

func NewHandler(processExecutor Executor) *Handler {
	return &Handler{
		processExecutor: processExecutor,
	}
}

func (h *Handler) HandleCVRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	password := r.URL.Query().Get("password")
	lang := r.URL.Query().Get("lang")

	if password == "" || lang == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := h.processExecutor.Execute(password, lang)
	if err != nil {
		if err.Error() == "status_401" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}
