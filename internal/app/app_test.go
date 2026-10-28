package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/registry"
)

func TestUnknownRoutes(t *testing.T) {
	cfg := &registry.Config{}
	cfg.Server.Port = "8080"
	cfg.Api.GatewayURL = "http://127.0.0.1:1"
	cfg.Templates.Path = filepath.Join("..", "web", "template")

	application, err := Build(cfg)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name        string
		path        string
		accept      string
		contentType string
		body        string
		htmlLang    string
	}{
		{"polish page", "/missing-page", "", "text/html", "Nie znaleziono strony", "pl"},
		{"english page", "/missing?lang=en", "", "text/html", "Page not found", "en"},
		{"invalid CV download path", "/api/download/cv/extra?token=example&lang=pl", "text/html,application/xhtml+xml,*/*;q=0.8", "text/html", "Nie znaleziono strony", "pl"},
		{"misspelled CV download path", "/api/download/cva?token=example&lang=pl", "text/html,application/xhtml+xml,*/*;q=0.8", "text/html", "Nie znaleziono strony", "pl"},
		{"misspelled CV download path in English", "/api/download/cva?token=example&lang=en", "text/html,application/xhtml+xml,*/*;q=0.8", "text/html", "Page not found", "en"},
		{"unknown browser API route", "/api/unknown?lang=en", "text/html", "text/html", "Page not found", "en"},
		{"unknown API route", "/api/unknown", "application/json", "application/json", "error_route_not_found", ""},
		{"API root", "/api", "", "application/json", "error_route_not_found", ""},
		{"generic API client", "/api/download/cv/extra", "*/*", "application/json", "error_route_not_found", ""},
		{"HTML rejected", "/api/unknown", "text/html;q=0,*/*;q=1", "application/json", "error_route_not_found", ""},
		{"JSON preferred", "/api/unknown", "text/html;q=0.2,application/json;q=1", "application/json", "error_route_not_found", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.accept != "" {
				request.Header.Set("Accept", tt.accept)
			}
			application.server.Handler.ServeHTTP(w, request)
			if w.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404", w.Code)
			}
			if !strings.HasPrefix(w.Header().Get("Content-Type"), tt.contentType) {
				t.Errorf("Content-Type = %q, want prefix %q", w.Header().Get("Content-Type"), tt.contentType)
			}
			if !strings.Contains(w.Body.String(), tt.body) {
				t.Errorf("body does not contain %q", tt.body)
			}
			if tt.contentType == "text/html" {
				for _, marker := range []string{"error-centered-body", "styled-block", "?lang=pl", "?lang=en"} {
					if !strings.Contains(w.Body.String(), marker) {
						t.Errorf("HTML 404 is missing %q", marker)
					}
				}
				if !strings.Contains(w.Body.String(), `<html lang="`+tt.htmlLang+`">`) {
					t.Errorf("HTML 404 has the wrong document language, want %q", tt.htmlLang)
				}
			}
			if strings.HasPrefix(tt.path, "/api") {
				if got := w.Header().Get("Cache-Control"); got != "no-store" {
					t.Errorf("Cache-Control = %q, want no-store", got)
				}
				if got := w.Header().Get("Referrer-Policy"); got != "no-referrer" {
					t.Errorf("Referrer-Policy = %q, want no-referrer", got)
				}
				if got := w.Header().Get("Vary"); got != "Accept" {
					t.Errorf("Vary = %q, want Accept", got)
				}
				if strings.Contains(tt.path, "token=example") && strings.Contains(w.Body.String(), "token=example") {
					t.Error("404 response reflects the CV token")
				}
			}
			if tt.contentType == "application/json" {
				var response map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response["error"] != tt.body {
					t.Errorf("invalid API error: response=%v err=%v", response, err)
				}
			}
		})
	}

	w := httptest.NewRecorder()
	application.server.Handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/static/fonts/Inter-latin.woff2", nil))
	if w.Code != http.StatusOK || !strings.HasPrefix(w.Body.String(), "wOF2") {
		t.Fatalf("local font response: status=%d, magic=%q", w.Code, w.Body.String()[:min(4, w.Body.Len())])
	}
	if got := w.Header().Get("Content-Type"); got != "font/woff2" {
		t.Errorf("font Content-Type = %q, want font/woff2", got)
	}
}
