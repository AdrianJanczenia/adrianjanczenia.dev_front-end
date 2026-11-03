package renderer

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	assetURLs, err := versionedStaticAssets(templateMap)
	if err != nil {
		return nil, err
	}

	funcs := template.FuncMap{
		"staticAssetURL": func(path string) string {
			if versioned, ok := assetURLs[path]; ok {
				return versioned
			}
			return path
		},
		"jsonData": func(v any) template.JS {
			encoded, err := json.Marshal(v)
			if err != nil {
				return template.JS("null")
			}
			return template.JS(encoded)
		},
		"casePath": func(id string) string {
			switch id {
			case "migration":
				return "/experience/monolith-migration"
			case "checkout":
				return "/experience/orders-and-payments"
			case "reservation":
				return "/experience/reservation-waitlist"
			default:
				return "/experience"
			}
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

func versionedStaticAssets(templateMap map[string][]string) (map[string]string, error) {
	urls := make(map[string]string)
	for _, files := range templateMap {
		for _, file := range files {
			if filepath.Base(file) != "base.html" || filepath.Base(filepath.Dir(file)) != "layout" {
				continue
			}

			staticDir := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(file))), "static")
			info, err := os.Stat(staticDir)
			if stderrors.Is(err, os.ErrNotExist) {
				return urls, nil
			}
			if err != nil {
				return nil, fmt.Errorf("find static assets: %w", err)
			}
			if !info.IsDir() {
				return nil, fmt.Errorf("static asset path is not a directory: %s", staticDir)
			}

			for _, path := range []string{
				"/static/css/portfolio.css",
				"/static/img/technology-icons.svg",
				"/static/js/site.js",
				"/static/js/cv.js",
				"/static/js/explorer.js",
			} {
				assetPath := filepath.Join(staticDir, filepath.FromSlash(path[len("/static/"):]))
				asset, err := os.ReadFile(assetPath)
				if err != nil {
					return nil, fmt.Errorf("read static asset %s: %w", path, err)
				}
				sum := sha256.Sum256(asset)
				urls[path] = path + "?v=" + hex.EncodeToString(sum[:8])
			}
			return urls, nil
		}
	}
	return urls, nil
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

	if content != nil && content.UI.Errors[appErr.Slug].Message != "" {
		msg = content.UI.Errors[appErr.Slug].Message
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
	if content != nil && content.UI.Errors[appErr.Slug].Title != "" {
		title = content.UI.Errors[appErr.Slug].Title
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
