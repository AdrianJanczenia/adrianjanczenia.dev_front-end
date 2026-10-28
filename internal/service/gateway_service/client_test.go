package gateway_service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
)

func TestClient_GetPageContent(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		wantErr        error
	}{
		{
			name: "successful response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"meta":{"title":"Test"}}`))
			},
			wantErr: nil,
		},
		{
			name: "invalid content type",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte(`ok`))
			},
			wantErr: errors.ErrInternalServerError,
		},
		{
			name: "missing error slug",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr: errors.ErrBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer ts.Close()

			c := NewClient(&http.Client{}, ts.URL)
			res, err := c.GetPageContent(context.Background(), "pl")

			if err != tt.wantErr {
				t.Errorf("GetPageContent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && res.Meta.Title != "Test" {
				t.Errorf("GetPageContent() title = %s, want Test", res.Meta.Title)
			}
		})
	}
}

func TestClient_RequestCVToken(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		wantToken      string
		wantErr        error
	}{
		{
			name: "successful token",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"token":"secret"}`))
			},
			wantToken: "secret",
			wantErr:   nil,
		},
		{
			name: "invalid password",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"error_cv_auth"}`))
			},
			wantToken: "",
			wantErr:   errors.ErrInvalidPassword,
		},
		{
			name: "missing token",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{}`))
			},
			wantToken: "",
			wantErr:   errors.ErrServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer ts.Close()

			c := NewClient(&http.Client{}, ts.URL)
			token, err := c.RequestCVToken(context.Background(), "pass", "pl", "123")

			if err != tt.wantErr {
				t.Errorf("RequestCVToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if token != tt.wantToken {
				t.Errorf("RequestCVToken() got = %s, want %s", token, tt.wantToken)
			}
		})
	}
}

func TestClient_DownloadCVStream(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("token") == "valid" {
			w.Header().Set("Content-Type", "application/pdf")
			w.Write([]byte("%PDF"))
			return
		}
		w.WriteHeader(http.StatusGone)
		w.Write([]byte(`{"error":"error_cv_expired"}`))
	}))
	defer ts.Close()

	c := NewClient(&http.Client{}, ts.URL)

	t.Run("successful stream", func(t *testing.T) {
		body, contentType, err := c.DownloadCVStream(context.Background(), "valid", "pl")
		if err != nil {
			t.Fatalf("DownloadCVStream() unexpected error: %v", err)
		}
		defer body.Close()
		content, _ := io.ReadAll(body)
		if string(content) != "%PDF" || contentType != "application/pdf" {
			t.Errorf("DownloadCVStream() got invalid content or type")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		_, _, err := c.DownloadCVStream(context.Background(), "expired", "pl")
		if err != errors.ErrCVExpired {
			t.Errorf("DownloadCVStream() expected ErrCVExpired, got %v", err)
		}
	})
}

func TestDownloadCVStreamContinuesAfterClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write([]byte("%PDF"))
		w.(http.Flusher).Flush()
		time.Sleep(50 * time.Millisecond)
		_, _ = w.Write([]byte("-complete"))
	}))
	defer server.Close()
	client := NewClient(&http.Client{Timeout: 10 * time.Millisecond}, server.URL)
	body, _, err := client.DownloadCVStream(context.Background(), "token", "pl")
	if err != nil {
		t.Fatal(err)
	}
	defer body.Close()
	content, err := io.ReadAll(body)
	if err != nil || string(content) != "%PDF-complete" {
		t.Fatalf("stream = %q, %v", content, err)
	}
}
