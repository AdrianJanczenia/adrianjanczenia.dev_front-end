package renderer

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func setupTestTemplates(t *testing.T) string {
	tmpDir := t.TempDir()
	layoutDir := filepath.Join(tmpDir, "layout")
	err := os.MkdirAll(layoutDir, 0755)
	if err != nil {
		return ""
	}

	baseContent := `{{define "base"}}{{template "content" .}}{{end}}`
	err = os.WriteFile(filepath.Join(layoutDir, "base.html"), []byte(baseContent), 0644)
	if err != nil {
		return ""
	}

	indexContent := `{{define "content"}}<title>{{.Content.Meta.Title}}</title>{{end}}`
	err = os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(indexContent), 0644)
	if err != nil {
		return ""
	}

	errorContent := `{{define "error.html"}}Error Page{{end}}`
	err = os.WriteFile(filepath.Join(tmpDir, "error.html"), []byte(errorContent), 0644)
	if err != nil {
		return ""
	}

	return tmpDir
}

func TestNew(t *testing.T) {
	t.Run("it builds renderer with valid templates", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("New() panicked, but it was not expected to")
			}
		}()
		templateDir := setupTestTemplates(t)
		New(templateDir)
	})

	t.Run("it panics when a template is missing", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("New() did not panic, but it was expected to")
			}
		}()
		tmpDir := t.TempDir()
		New(tmpDir)
	})
}

func TestRenderer_Render(t *testing.T) {
	templateDir := setupTestTemplates(t)
	renderer := New(templateDir)

	t.Run("it renders a valid template", func(t *testing.T) {
		rr := httptest.NewRecorder()
		data := struct {
			Content struct{ Meta struct{ Title string } }
		}{Content: struct{ Meta struct{ Title string } }{Meta: struct{ Title string }{Title: "Test Title"}}}

		renderer.Render(rr, "index", data)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, rr.Code)
		}
		if !bytes.Contains(rr.Body.Bytes(), []byte("<title>Test Title</title>")) {
			t.Errorf("rendered HTML does not contain the expected content")
		}
	})
}
