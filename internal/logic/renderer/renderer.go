package renderer

import (
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
	}

	if err := tmpl.ExecuteTemplate(w, executeTmpl, data); err != nil {
		if name == "error" {
			log.Printf("CRITICAL: failed to execute error template: %v", err)
			http.Error(w, "a critical error occurred while rendering the error page.", http.StatusInternalServerError)
			return
		}
		log.Printf("ERROR: could not execute template %s: %v", name, err)
		r.Render(w, "error", data)
	}
}
