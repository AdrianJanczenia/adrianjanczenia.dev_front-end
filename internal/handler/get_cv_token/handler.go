package get_cv_token

import (
	"encoding/json"
	"net/http"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
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
		errors.WriteJSON(w, errors.ErrMethodNotAllowed)
		return
	}

	var input struct {
		Password string `json:"password"`
		Lang     string `json:"lang"`
		FullName string `json:"fullName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	if input.FullName != "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": "token"})
		return
	}

	if input.Password == "" || input.Lang == "" {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	token, err := h.processExecutor.Execute(input.Password, input.Lang)
	if err != nil {
		errors.WriteJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
