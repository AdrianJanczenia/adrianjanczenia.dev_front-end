package renderer

import (
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

func TestPortfolioTemplatesRender(t *testing.T) {
	root := filepath.Join("..", "..", "web", "template")
	base := filepath.Join(root, "layout", "base.html")
	partials, err := filepath.Glob(filepath.Join(root, "partials", "*.html"))
	if err != nil {
		t.Fatal(err)
	}
	templates := map[string][]string{}
	for _, name := range []string{"index", "experience", "case", "explorer", "contact", "privacy"} {
		files := append([]string{base}, partials...)
		templates[name] = append(files, filepath.Join(root, name+".html"))
	}
	r, err := New(templates)
	if err != nil {
		t.Fatal(err)
	}
	content := &gateway_service.PageContent{}
	content.Meta.Title = "Portfolio test"
	content.Shared.Navigation.About = "About me"
	content.Pages.About.Hero.Title = "Backend test"
	content.Pages.Experience.Cases.Order = []string{"migration"}
	content.Pages.Experience.Cases.Items = map[string]gateway_service.Case{"migration": {Title: "Migration test"}}
	content.Pages.HowItWorks.ScenarioCopy = map[string]map[string]string{"content": {"injected": `</script><script>alert(1)</script>`}}
	content.UI.Errors = map[string]gateway_service.ErrorCopy{"error_message": {Message: "Temporary error"}}
	for _, lang := range []string{"pl", "en"} {
		for _, name := range []string{"index", "experience", "case", "explorer", "contact", "privacy"} {
			t.Run(lang+"/"+name, func(t *testing.T) {
				w := httptest.NewRecorder()
				r.Render(w, name, &data.TemplateData{Lang: lang, Page: name, CaseID: "migration", Path: "/", Content: content})
				if w.Code != 200 {
					t.Fatalf("status %d: %s", w.Code, w.Body.String())
				}
				body := w.Body.String()
				if !strings.Contains(body, `<html lang="`+lang+`"`) || !strings.Contains(body, `id="site-nav"`) || !strings.Contains(body, `id="cv-modal"`) {
					t.Fatal("shared layout missing")
				}
				if name == "explorer" && (!strings.Contains(body, `id="request-explorer"`) || !strings.Contains(body, `data-id="gateway"`)) {
					t.Fatal("explorer container or map missing")
				}
				if name == "explorer" && strings.Contains(body, `</script><script>alert(1)</script>`) {
					t.Fatal("explorer copy broke out of its JSON script element")
				}
				if name == "case" && !strings.Contains(body, "Migration test") {
					t.Fatal("case not rendered")
				}
			})
		}
	}
}
