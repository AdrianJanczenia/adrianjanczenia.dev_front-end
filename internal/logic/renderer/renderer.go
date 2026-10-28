package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
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

func New(templateMap map[string][]string) (Renderer, error) {
	r := &templateRenderer{
		templates: make(map[string]*template.Template),
	}

	funcs := template.FuncMap{
		"join": func(s []string) string {
			if len(s) == 0 {
				return ""
			}
			return "\"" + strings.Join(s, "\", \"") + "\""
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}

	for name, files := range templateMap {
		t, err := template.New("base.html").Funcs(funcs).ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		r.templates[name] = t
	}

	return r, nil
}

func (r *templateRenderer) renderBytes(name string, data any) ([]byte, error) {
	tmpl, ok := r.templates[name]
	if !ok {
		return nil, fmt.Errorf("template %s does not exist", name)
	}

	executeTmpl := "base"
	if name == "error" || name == "cv_error" {
		executeTmpl = name + ".html"
	}

	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, executeTmpl, data); err != nil {
		return nil, fmt.Errorf("execute template %s: %w", name, err)
	}
	return buf.Bytes(), nil
}

func (r *templateRenderer) Render(w http.ResponseWriter, name string, data any) {
	body, err := r.renderBytes(name, data)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(body)
}

func (r *templateRenderer) RenderError(w http.ResponseWriter, templateName string, appErr *errors.AppError, lang string, content *gateway_service.PageContent) {
	if appErr == nil {
		appErr = errors.ErrInternalServerError
	}
	if lang != "en" {
		lang = "pl"
	}
	fallback := errors.FallbackTranslations[lang]
	msg := ""
	if content != nil && content.Translations != nil && content.Translations[appErr.Slug] != "" {
		msg = content.Translations[appErr.Slug]
	} else {
		msg = fallback[appErr.Slug]
	}

	if msg == "" {
		msg = fallback["error_internal"]
	}

	titleKey := "error_title"
	if appErr.Slug == errors.ErrRouteNotFound.Slug {
		titleKey = "error_route_title"
	}
	title := fallback[titleKey]
	if content != nil && content.Translations != nil && content.Translations[titleKey] != "" {
		title = content.Translations[titleKey]
	}

	data := struct {
		Lang          string
		HTTPStatus    int
		ErrorTitle    string
		ErrorMessage  string
		HomeLabel     string
		ErrorHeader   string
		StatusKey     string
		MessageKey    string
		LanguageLabel string
		Content       *gateway_service.PageContent
	}{
		Lang:          lang,
		HTTPStatus:    appErr.HTTPStatus,
		ErrorTitle:    title,
		ErrorMessage:  msg,
		HomeLabel:     fallback["error_home"],
		ErrorHeader:   fallback["error_header"],
		StatusKey:     fallback["error_status_key"],
		MessageKey:    fallback["error_message_key"],
		LanguageLabel: fallback["error_language_label"],
		Content:       content,
	}

	body, err := r.renderBytes(templateName, data)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(appErr.HTTPStatus)
	_, _ = w.Write(body)
}
