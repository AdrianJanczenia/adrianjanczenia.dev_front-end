package get_cv_token

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
)

type input struct {
	Password  string `json:"password"`
	Lang      string `json:"lang"`
	CaptchaID string `json:"captchaId"`
	FullName  string `json:"fullName"`
}

type GetCVTokenProcess interface {
	Process(ctx context.Context, password, lang, captchaID string) (string, error)
}

type Handler struct {
	process GetCVTokenProcess
}

func NewHandler(process GetCVTokenProcess) *Handler {
	return &Handler{
		process: process,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.WriteJSON(w, errors.ErrMethodNotAllowed)
		return
	}

	var i input

	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	if i.FullName != "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": "token"})
		return
	}

	if i.Password == "" || i.Lang == "" || i.CaptchaID == "" {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	token, err := h.process.Process(r.Context(), i.Password, i.Lang, i.CaptchaID)
	if err != nil {
		errors.WriteJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
