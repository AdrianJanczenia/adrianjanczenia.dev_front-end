package renderer

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type Renderer interface {
	Render(w http.ResponseWriter, name string, data any)
	RenderError(w http.ResponseWriter, templateName string, appErr *errors.AppError, lang string, content *gateway_service.PageContent)
}

type templateRenderer struct {
	templates map[string]*template.Template
}

func New(templatesPath string) Renderer {
	r := &templateRenderer{
		templates: make(map[string]*template.Template),
	}

	funcs := template.FuncMap{
		"join": func(s []string) string {
			return strings.Join(s, ", ")
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}

	layoutBase := filepath.Join(templatesPath, "layout", "base.html")
	partials, _ := filepath.Glob(filepath.Join(templatesPath, "partials", "*.html"))

	contentPages := []string{"index", "privacy"}
	for _, page := range contentPages {
		files := append([]string{layoutBase}, partials...)
		files = append(files, filepath.Join(templatesPath, page+".html"))
		r.templates[page] = template.Must(template.New("base.html").Funcs(funcs).ParseFiles(files...))
	}

	errorPages := []string{"error", "cv_error"}
	for _, page := range errorPages {
		r.templates[page] = template.Must(template.New(page + ".html").Funcs(funcs).ParseFiles(filepath.Join(templatesPath, page+".html")))
	}

	return r
}

func (r *templateRenderer) Render(w http.ResponseWriter, name string, data any) {
	tmpl, ok := r.templates[name]
	if !ok {
		log.Printf("ERROR: template %s does not exist", name)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	executeTmpl := "base"
	if name == "error" || name == "cv_error" {
		executeTmpl = name + ".html"
	}

	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, executeTmpl, data); err != nil {
		log.Printf("ERROR: could not execute template %s: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

func (r *templateRenderer) RenderError(w http.ResponseWriter, templateName string, appErr *errors.AppError, lang string, content *gateway_service.PageContent) {
	msg := ""
	if content != nil && content.Translations != nil && content.Translations[appErr.Slug] != "" {
		msg = content.Translations[appErr.Slug]
	} else if trans, ok := errors.FallbackTranslations[lang]; ok {
		msg = trans[appErr.Slug]
	}

	if msg == "" {
		msg = errors.FallbackTranslations[lang]["error_cv_server"]
	}

	title := ""
	if content != nil && content.Translations != nil && content.Translations["error_title"] != "" {
		title = content.Translations["error_title"]
	} else if lang == "pl" {
		title = "Wystąpił błąd"
	} else {
		title = "An error occurred"
	}

	data := struct {
		Lang         string
		HTTPStatus   int
		ErrorTitle   string
		ErrorMessage string
		Content      *gateway_service.PageContent
	}{
		Lang:         lang,
		HTTPStatus:   appErr.HTTPStatus,
		ErrorTitle:   title,
		ErrorMessage: msg,
		Content:      content,
	}

	w.WriteHeader(appErr.HTTPStatus)
	r.Render(w, templateName, data)
}
