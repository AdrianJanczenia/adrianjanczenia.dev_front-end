package verify_captcha

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
)

type mockVerifyCaptchaProcess struct {
	processFunc func(ctx context.Context, captchaID, captchaValue string) (string, error)
}

func (m *mockVerifyCaptchaProcess) Process(ctx context.Context, captchaID, captchaValue string) (string, error) {
	return m.processFunc(ctx, captchaID, captchaValue)
}

func TestHandler_VerifyCaptcha(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		body        interface{}
		processFunc func(context.Context, string, string) (string, error)
		wantStatus  int
		wantBody    map[string]string
	}{
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   map[string]string{"error": "error_message"},
		},
		{
			name:       "invalid input",
			method:     http.MethodPost,
			body:       "invalid-json",
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]string{"error": "error_message"},
		},
		{
			name:   "process error",
			method: http.MethodPost,
			body: map[string]string{
				"captchaId":    "cid",
				"captchaValue": "val",
			},
			processFunc: func(ctx context.Context, id, val string) (string, error) {
				return "", errors.ErrCaptchaInvalid
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]string{"error": "error_captcha_invalid"},
		},
		{
			name:   "success",
			method: http.MethodPost,
			body: map[string]string{
				"captchaId":    "cid",
				"captchaValue": "val",
			},
			processFunc: func(ctx context.Context, id, val string) (string, error) {
				return "ok-id", nil
			},
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"captchaId": "ok-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockVerifyCaptchaProcess{processFunc: tt.processFunc}
			h := NewHandler(m)

			var reqBody []byte
			if s, ok := tt.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/api/captcha-verify", bytes.NewBuffer(reqBody))
			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Handle() status = %v, want %v", rr.Code, tt.wantStatus)
			}

			var gotBody map[string]string
			json.Unmarshal(rr.Body.Bytes(), &gotBody)

			for k, v := range tt.wantBody {
				if gotBody[k] != v {
					t.Errorf("Handle() body[%s] = %v, want %v", k, gotBody[k], v)
				}
			}
		})
	}
}
