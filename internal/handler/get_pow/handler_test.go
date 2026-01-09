package get_pow

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockGetPowProcess struct {
	processFunc func(ctx context.Context) (map[string]string, error)
}

func (m *mockGetPowProcess) Process(ctx context.Context) (map[string]string, error) {
	return m.processFunc(ctx)
}

func TestHandler_GetPow(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		processFunc func(context.Context) (map[string]string, error)
		wantStatus  int
		wantBody    map[string]string
	}{
		{
			name:       "method not allowed",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   map[string]string{"error": "error_message"},
		},
		{
			name:   "process error",
			method: http.MethodGet,
			processFunc: func(ctx context.Context) (map[string]string, error) {
				return nil, errors.New("internal failure")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   map[string]string{"error": "error_cv_server"},
		},
		{
			name:   "success",
			method: http.MethodGet,
			processFunc: func(ctx context.Context) (map[string]string, error) {
				return map[string]string{"seed": "s1", "signature": "sig1"}, nil
			},
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"seed": "s1", "signature": "sig1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mockGetPowProcess{processFunc: tt.processFunc}
			h := NewHandler(m)

			req := httptest.NewRequest(tt.method, "/api/pow", nil)
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
