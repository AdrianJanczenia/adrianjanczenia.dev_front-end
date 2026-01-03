package get_cv_token

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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Password string `json:"password"`
		Lang     string `json:"lang"`
		FullName string `json:"fullName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if input.FullName != "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": "token"})
		return
	}

	if input.Password == "" || input.Lang == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.processExecutor.Execute(input.Password, input.Lang)
	if err != nil {
		if err.Error() == "status_401" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
