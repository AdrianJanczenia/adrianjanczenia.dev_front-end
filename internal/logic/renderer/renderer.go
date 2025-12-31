package renderer

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type Renderer interface {
	Render(w http.ResponseWriter, name string, data any)
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
			for i, v := range s {
				s[i] = `"` + v + `"`
			}
			return strings.Join(s, ", ")
		},
	}

	baseTmpl := filepath.Join(templatesPath, "layout", "base.html")
	indexTmpl := filepath.Join(templatesPath, "index.html")
	r.templates["index"] = template.Must(template.New("base.html").Funcs(funcs).ParseFiles(baseTmpl, indexTmpl))

	errorTmpl := filepath.Join(templatesPath, "error.html")
	r.templates["error"] = template.Must(template.New("error.html").Funcs(funcs).ParseFiles(errorTmpl))

	cvErrorTmpl := filepath.Join(templatesPath, "cv_error.html")
	r.templates["cv_error"] = template.Must(template.New("cv_error.html").Funcs(funcs).ParseFiles(cvErrorTmpl))

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
	if name == "error" {
		executeTmpl = "error.html"
	} else if name == "cv_error" {
		executeTmpl = "cv_error.html"
	}

	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, executeTmpl, data); err != nil {
		log.Printf("ERROR: could not execute template %s: %v", name, err)

		if name == "error" {
			http.Error(w, "Critical error rendering error page", http.StatusInternalServerError)
			return
		}

		errorData := struct {
			Lang    string
			Title   string
			Message string
		}{
			Lang:    "en",
			Title:   "Rendering Error",
			Message: "A critical error occurred while rendering the page.",
		}
		r.Render(w, "error", errorData)
		return
	}

	buf.WriteTo(w)
}
