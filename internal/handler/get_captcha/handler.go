package get_captcha

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
)

type GetCaptchaProcess interface {
	Process(ctx context.Context, seed, signature, nonce string) (map[string]string, error)
}

type Handler struct {
	process GetCaptchaProcess
}

func NewHandler(p GetCaptchaProcess) *Handler {
	return &Handler{process: p}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.WriteJSON(w, errors.ErrMethodNotAllowed)
		return
	}

	var req struct {
		Seed      string `json:"seed"`
		Signature string `json:"signature"`
		Nonce     string `json:"nonce"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	data, err := h.process.Process(r.Context(), req.Seed, req.Signature, req.Nonce)
	if err != nil {
		errors.WriteJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
