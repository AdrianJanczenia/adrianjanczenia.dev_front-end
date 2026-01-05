package renderer

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

func setupTestTemplates(t *testing.T) string {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "layout"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "partials"), 0755)

	dummyBase := `{{define "base"}}Base{{template "content" .}}{{end}}`
	dummyContent := `{{define "content"}}Content{{end}}`
	dummyError := `{{.ErrorMessage}}`

	os.WriteFile(filepath.Join(tmpDir, "layout", "base.html"), []byte(dummyBase), 0644)
	os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(dummyContent), 0644)
	os.WriteFile(filepath.Join(tmpDir, "privacy.html"), []byte(dummyContent), 0644)
	os.WriteFile(filepath.Join(tmpDir, "error.html"), []byte(dummyError), 0644)
	os.WriteFile(filepath.Join(tmpDir, "cv_error.html"), []byte(dummyError), 0644)

	return tmpDir
}

func TestTemplateRenderer_Render(t *testing.T) {
	tmpDir := setupTestTemplates(t)
	r := New(tmpDir)

	t.Run("renders index template", func(t *testing.T) {
		w := httptest.NewRecorder()
		r.Render(w, "index", nil)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if w.Body.String() != "BaseContent" {
			t.Errorf("expected BaseContent, got %s", w.Body.String())
		}
	})

	t.Run("renders error template directly", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := struct{ ErrorMessage string }{ErrorMessage: "CustomError"}
		r.Render(w, "error", data)
		if w.Body.String() != "CustomError" {
			t.Errorf("expected CustomError, got %s", w.Body.String())
		}
	})

	t.Run("handles non-existent template", func(t *testing.T) {
		w := httptest.NewRecorder()
		r.Render(w, "missing", nil)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", w.Code)
		}
	})
}

func TestTemplateRenderer_RenderError(t *testing.T) {
	tmpDir := setupTestTemplates(t)
	r := New(tmpDir)

	t.Run("renders error with fallback translation", func(t *testing.T) {
		w := httptest.NewRecorder()
		r.RenderError(w, "error", errors.ErrInvalidPassword, "pl", nil)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
		expectedMsg := errors.FallbackTranslations["pl"]["error_cv_auth"]
		if w.Body.String() != expectedMsg {
			t.Errorf("expected %s, got %s", expectedMsg, w.Body.String())
		}
	})

	t.Run("renders error with content translations", func(t *testing.T) {
		w := httptest.NewRecorder()
		content := &gateway_service.PageContent{
			Translations: map[string]string{
				"error_cv_auth": "Błędne hasło z API",
			},
		}
		r.RenderError(w, "error", errors.ErrInvalidPassword, "pl", content)
		if w.Body.String() != "Błędne hasło z API" {
			t.Errorf("expected translated message, got %s", w.Body.String())
		}
	})
}
