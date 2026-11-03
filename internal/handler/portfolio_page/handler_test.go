package portfolio_page

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type fakeProcess struct {
	lang  string
	calls int
}

func (p *fakeProcess) Process(_ context.Context, lang string) (*data.TemplateData, error) {
	p.lang = lang
	p.calls++
	return &data.TemplateData{Lang: lang, Content: &gateway_service.PageContent{}}, nil
}

type fakeRenderer struct {
	name string
	data *data.TemplateData
	err  *appErrors.AppError
}

func (r *fakeRenderer) Render(_ http.ResponseWriter, name string, value any) {
	r.name = name
	r.data = value.(*data.TemplateData)
}
func (r *fakeRenderer) RenderError(_ http.ResponseWriter, _ string, err *appErrors.AppError, _ string, _ *gateway_service.PageContent) {
	r.err = err
}

func TestPageHandlerPreservesPageAndLanguage(t *testing.T) {
	p, r := &fakeProcess{}, &fakeRenderer{}
	h := NewHandler(p, r)
	req := httptest.NewRequest(http.MethodGet, "/experience/monolith-migration?lang=en", nil)
	h.Handle("case", "experience", "migration")(httptest.NewRecorder(), req)
	if p.lang != "en" || r.name != "case" || r.data.Page != "experience" || r.data.CaseID != "migration" || r.data.Path != "/experience/monolith-migration" {
		t.Fatalf("unexpected render: lang=%q template=%q data=%+v", p.lang, r.name, r.data)
	}
}

func TestPageHandlerRejectsNonGET(t *testing.T) {
	p, r := &fakeProcess{}, &fakeRenderer{}
	w := httptest.NewRecorder()
	NewHandler(p, r).Handle("contact", "contact", "")(w, httptest.NewRequest(http.MethodPost, "/contact", nil))
	if p.calls != 0 || r.err != appErrors.ErrMethodNotAllowed || w.Header().Get("Allow") != "GET" {
		t.Fatal("non-GET was not rejected")
	}
}
