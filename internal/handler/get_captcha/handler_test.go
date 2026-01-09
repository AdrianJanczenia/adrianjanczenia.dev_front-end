package get_captcha

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockGetCaptchaProcess struct {
	processFunc func(ctx context.Context, seed, signature, nonce string) (map[string]string, error)
}

func (m *mockGetCaptchaProcess) Process(ctx context.Context, seed, signature, nonce string) (map[string]string, error) {
	return m.processFunc(ctx, seed, signature, nonce)
}

func TestHandler_GetCaptcha(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		body        interface{}
		processFunc func(context.Context, string, string, string) (map[string]string, error)
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
				"seed":      "s",
				"signature": "sig",
				"nonce":     "n",
			},
			processFunc: func(ctx context.Context, s, sig, n string) (map[string]string, error) {
				return nil, errors.New("fail")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   map[string]string{"error": "error_cv_server"},
		},
		{
			name:   "success",
			method: http.MethodPost,
			body: map[string]string{
				"seed":      "s",
				"signature": "sig",
				"nonce":     "n",
			},
			processFunc: func(ctx context.Context, s, sig, n string) (map[string]string, error) {
				return map[string]string{"captchaId": "1", "captchaImg": "data"}, nil
			},
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"captchaId": "1", "captchaImg": "data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockGetCaptchaProcess{processFunc: tt.processFunc}
			h := NewHandler(m)

			var reqBody []byte
			if s, ok := tt.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/api/captcha", bytes.NewBuffer(reqBody))
			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Handle() status = %v, want %v", rr.Code, tt.wantStatus)
			}

			var gotBody map[string]string
			json.Unmarshal(rr.Body.Bytes(), &gotBody)

			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("Handle() body = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}
