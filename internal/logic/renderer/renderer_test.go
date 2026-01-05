package renderer

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
)

func setupTestTemplates(t *testing.T) (string, map[string][]string) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "layout"), 0755)

	dummyBase := `{{define "base"}}Base{{template "content" .}}{{end}}`
	dummyContent := `{{define "content"}}Content{{end}}`
	dummyError := `{{.ErrorMessage}}`

	basePath := filepath.Join(tmpDir, "layout", "base.html")
	indexPath := filepath.Join(tmpDir, "index.html")
	errorPath := filepath.Join(tmpDir, "error.html")

	os.WriteFile(basePath, []byte(dummyBase), 0644)
	os.WriteFile(indexPath, []byte(dummyContent), 0644)
	os.WriteFile(errorPath, []byte(dummyError), 0644)

	templateMap := map[string][]string{
		"index": {basePath, indexPath},
		"error": {errorPath},
	}

	return tmpDir, templateMap
}

func TestTemplateRenderer_Render(t *testing.T) {
	_, templateMap := setupTestTemplates(t)
	r, err := New(templateMap)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

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
}

func TestTemplateRenderer_RenderError(t *testing.T) {
	_, templateMap := setupTestTemplates(t)
	r, _ := New(templateMap)

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
}
