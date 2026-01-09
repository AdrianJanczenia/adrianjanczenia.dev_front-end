package verify_captcha

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
)

type VerifyCaptchaProcess interface {
	Process(ctx context.Context, captchaID, captchaValue string) (string, error)
}

type Handler struct {
	process VerifyCaptchaProcess
}

func NewHandler(p VerifyCaptchaProcess) *Handler {
	return &Handler{process: p}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.WriteJSON(w, errors.ErrMethodNotAllowed)
		return
	}

	var req struct {
		CaptchaID    string `json:"captchaId"`
		CaptchaValue string `json:"captchaValue"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	captchaID, err := h.process.Process(r.Context(), req.CaptchaID, req.CaptchaValue)
	if err != nil {
		errors.WriteJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"captchaId": captchaID})
}
