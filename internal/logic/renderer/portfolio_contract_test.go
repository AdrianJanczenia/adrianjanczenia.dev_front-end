package renderer

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

func TestBundledContentRendersEveryPage(t *testing.T) {
	root := filepath.Join("..", "..", "web", "template")
	partials, err := filepath.Glob(filepath.Join(root, "partials", "*.html"))
	if err != nil {
		t.Fatal(err)
	}
	templates := map[string][]string{}
	for _, name := range []string{"index", "experience", "case", "explorer", "contact", "privacy"} {
		files := append([]string{filepath.Join(root, "layout", "base.html")}, partials...)
		templates[name] = append(files, filepath.Join(root, name+".html"))
	}
	pageRenderer, err := New(templates)
	if err != nil {
		t.Fatal(err)
	}
	explorerScript, err := os.ReadFile(filepath.Join("..", "..", "web", "static", "js", "explorer.js"))
	if err != nil {
		t.Fatal(err)
	}
	scriptHash := sha256.Sum256(explorerScript)
	explorerScriptURL := `/static/js/explorer.js?v=` + hex.EncodeToString(scriptHash[:8])
	iconSprite, err := os.ReadFile(filepath.Join("..", "..", "web", "static", "img", "technology-icons.svg"))
	if err != nil {
		t.Fatal(err)
	}
	for _, lang := range []string{"pl", "en"} {
		path := filepath.Join("..", "..", "..", "adrianjanczenia.dev_content-service", "content", lang+".json")
		raw, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			t.Skip("sibling content-service repository not present")
		}
		if err != nil {
			t.Fatal(err)
		}
		var content gateway_service.PageContent
		if err := json.Unmarshal(raw, &content); err != nil {
			t.Fatal(err)
		}
		if content.SchemaVersion != 2 || content.Meta.Locale != lang {
			t.Fatalf("%s content version or locale mismatch", lang)
		}
		for _, groupID := range content.Pages.Experience.TechnologyGroups.Order {
			group := content.Pages.Experience.TechnologyGroups.Items[groupID]
			for _, item := range group.Items {
				if item.ID == "" || item.Label == "" || !strings.Contains(string(iconSprite), `id="`+item.ID+`"`) {
					t.Errorf("%s technology %q has no label or matching icon", lang, item.ID)
				}
			}
		}
		for _, tc := range []struct{ name, page, caseID, path, marker string }{
			{"index", "about", "", "/", content.Pages.About.Hero.Title},
			{"experience", "experience", "", "/experience", content.Pages.Experience.Intro.Title},
			{"case", "experience", "migration", "/experience/monolith-migration", content.Pages.Experience.Cases.Items["migration"].Title},
			{"case", "experience", "checkout", "/experience/orders-and-payments", content.Pages.Experience.Cases.Items["checkout"].Title},
			{"case", "experience", "reservation", "/experience/reservation-waitlist", content.Pages.Experience.Cases.Items["reservation"].Title},
			{"explorer", "how_it_works", "", "/how-this-site-works", "data-id=\"gateway\""},
			{"contact", "contact", "", "/contact", content.Pages.Contact.Intro.Title},
			{"privacy", "privacy", "", "/privacy-policy", content.Pages.Privacy.Title},
		} {
			t.Run(lang+"/"+tc.path, func(t *testing.T) {
				w := httptest.NewRecorder()
				pageRenderer.Render(w, tc.name, &data.TemplateData{Lang: lang, Page: tc.page, CaseID: tc.caseID, Path: tc.path, Content: &content})
				if w.Code != 200 || !strings.Contains(w.Body.String(), tc.marker) {
					t.Fatalf("render failed: status=%d marker=%q", w.Code, tc.marker)
				}
				if strings.Contains(w.Body.String(), "Miinto") {
					t.Fatal("employer name leaked into public HTML")
				}
				if tc.name == "experience" || tc.name == "case" {
					if strings.Contains(w.Body.String(), "CASE STUDY /") {
						t.Fatal("case category has an unwanted prefix")
					}
					for _, caseID := range content.Pages.Experience.Cases.Order {
						if tc.name == "case" && caseID != tc.caseID {
							continue
						}
						category := content.Pages.Experience.Cases.Items[caseID].Category
						markup := `<p class="eyebrow">` + category + `</p>`
						if tc.name == "case" {
							markup = `<p class="eyebrow case-kicker">` + category + `</p>`
						}
						if category == "" || category == caseID || !strings.Contains(w.Body.String(), markup) {
							t.Fatalf("%s category is missing or uses the technical ID", caseID)
						}
					}
				}
				if tc.name == "experience" && !strings.Contains(w.Body.String(), "/static/img/technology-icons.svg?v=") {
					t.Fatal("versioned technology icons are missing")
				}
				if !strings.Contains(w.Body.String(), `src="/static/js/cv.js?v=`) {
					t.Fatal("CV script is missing or unversioned")
				}
				const cvCopyPrefix = `<script type="application/json" id="cv-error-copy">`
				cvCopyStart := strings.Index(w.Body.String(), cvCopyPrefix)
				if cvCopyStart < 0 {
					t.Fatal("CV error translations are missing")
				}
				cvCopyJSON := w.Body.String()[cvCopyStart+len(cvCopyPrefix):]
				cvCopyJSON = cvCopyJSON[:strings.Index(cvCopyJSON, "</script>")]
				var cvErrorCopy map[string]gateway_service.ErrorCopy
				if err := json.Unmarshal([]byte(cvCopyJSON), &cvErrorCopy); err != nil || len(cvErrorCopy) == 0 {
					t.Fatalf("invalid CV error translations: %v", err)
				}
				if tc.name == "index" {
					expectedAlt := `alt="Portret Adriana Janczeni"`
					if lang == "en" {
						expectedAlt = `alt="Portrait of Adrian Janczenia"`
					}
					if !strings.Contains(w.Body.String(), expectedAlt) {
						t.Fatalf("portrait alternative text does not match %s", lang)
					}
				}
				if tc.name == "explorer" {
					if !strings.Contains(w.Body.String(), `src="`+explorerScriptURL+`"`) {
						t.Fatalf("explorer script is not versioned from its current contents")
					}
					start := strings.Index(w.Body.String(), `<script type="application/json" id="explorer-copy">`)
					if start < 0 {
						t.Fatal("explorer copy JSON missing")
					}
					jsonText := w.Body.String()[start+len(`<script type="application/json" id="explorer-copy">`):]
					jsonText = jsonText[:strings.Index(jsonText, "</script>")]
					var copy gateway_service.HowItWorksPage
					if err := json.Unmarshal([]byte(jsonText), &copy); err != nil || len(copy.ScenarioCopy) == 0 || len(copy.UICopy) == 0 {
						t.Fatalf("invalid explorer copy: %v", err)
					}
				}
			})
		}
	}
}

func TestExplorerCopyReferencesExist(t *testing.T) {
	js, err := os.ReadFile(filepath.Join("..", "..", "web", "static", "js", "explorer.js"))
	if err != nil {
		t.Fatal(err)
	}
	references := regexp.MustCompile(`\bscenarioText\('([a-z]+)\.([a-z0-9_]+)'\)`).FindAllStringSubmatch(string(js), -1)
	uiReferences := regexp.MustCompile(`\bui\.([A-Za-z][A-Za-z0-9]*)`).FindAllStringSubmatch(string(js), -1)
	markup, err := os.ReadFile(filepath.Join("..", "..", "web", "template", "explorer.html"))
	if err != nil {
		t.Fatal(err)
	}
	uiReferences = append(uiReferences, regexp.MustCompile(`data-i="([A-Za-z][A-Za-z0-9]*)"`).FindAllStringSubmatch(string(markup), -1)...)
	if len(references) == 0 {
		t.Fatal("explorer has no scenario copy references")
	}
	for _, lang := range []string{"pl", "en"} {
		path := filepath.Join("..", "..", "..", "adrianjanczenia.dev_content-service", "content", lang+".json")
		raw, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			t.Skip("sibling content-service repository not present")
		}
		if err != nil {
			t.Fatal(err)
		}
		var content gateway_service.PageContent
		if err := json.Unmarshal(raw, &content); err != nil {
			t.Fatal(err)
		}
		for _, reference := range references {
			if content.Pages.HowItWorks.ScenarioCopy[reference[1]][reference[2]] == "" {
				t.Errorf("%s is missing explorer copy %s.%s", lang, reference[1], reference[2])
			}
		}
		for _, reference := range uiReferences {
			if reference[1] == "title" || reference[1] == "lead" {
				continue
			}
			if content.Pages.HowItWorks.UICopy[reference[1]] == "" {
				t.Errorf("%s is missing explorer UI copy %s", lang, reference[1])
			}
		}
	}
}
