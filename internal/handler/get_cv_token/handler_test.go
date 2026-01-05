package get_cv_token

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockGetCVTokenProcess struct {
	processFunc func(ctx context.Context, password, lang string) (string, error)
}

func (m *mockGetCVTokenProcess) Process(ctx context.Context, password, lang string) (string, error) {
	return m.processFunc(ctx, password, lang)
}

func TestHandler_GetCVToken(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		body        any
		processFunc func(context.Context, string, string) (string, error)
		wantStatus  int
		wantToken   string
	}{
		{
			name:   "success",
			method: http.MethodPost,
			body:   input{Password: "123", Lang: "pl"},
			processFunc: func(ctx context.Context, p, l string) (string, error) {
				return "tok123", nil
			},
			wantStatus: http.StatusOK,
			wantToken:  "tok123",
		},
		{
			name:       "honeypot triggered",
			method:     http.MethodPost,
			body:       input{FullName: "bot"},
			wantStatus: http.StatusOK,
			wantToken:  "token",
		},
		{
			name:       "invalid input",
			method:     http.MethodPost,
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "process error",
			method: http.MethodPost,
			body:   input{Password: "wrong", Lang: "en"},
			processFunc: func(ctx context.Context, p, l string) (string, error) {
				return "", errors.New("fail")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(&mockGetCVTokenProcess{processFunc: tt.processFunc})
			var body []byte
			if s, ok := tt.body.(string); ok {
				body = []byte(s)
			} else {
				body, _ = json.Marshal(tt.body)
			}
			req := httptest.NewRequest(tt.method, "/", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			h.Handle(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Handle() status = %v, wantStatus %v", w.Code, tt.wantStatus)
			}
			if tt.wantToken != "" {
				var resp map[string]string
				json.Unmarshal(w.Body.Bytes(), &resp)
				if resp["token"] != tt.wantToken {
					t.Errorf("Handle() token = %s, want %s", resp["token"], tt.wantToken)
				}
			}
		})
	}
}
