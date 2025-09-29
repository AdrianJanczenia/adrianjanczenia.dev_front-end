package index_page

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
)

type mockExecutor struct {
	retData *data.TemplateData
	retErr  error
}

func (m *mockExecutor) Execute(lang string) (*data.TemplateData, error) {
	return m.retData, m.retErr
}

type mockRenderer struct {
	calledWithTmplName string
	calledWithData     any
}

func (m *mockRenderer) Render(w http.ResponseWriter, name string, data any) {
	m.calledWithTmplName = name
	m.calledWithData = data
}

func TestHandleIndexPage(t *testing.T) {
	testCases := []struct {
		name                   string
		url                    string
		method                 string
		mockExecutorResult     *data.TemplateData
		mockExecutorError      error
		expectedStatusCode     int
		expectedTmplName       string
		expectedLangInTmplData string
	}{
		{
			name:                   "Success case - Polish",
			url:                    "/?lang=pl",
			method:                 http.MethodGet,
			mockExecutorResult:     &data.TemplateData{Lang: "pl"},
			mockExecutorError:      nil,
			expectedStatusCode:     http.StatusOK,
			expectedTmplName:       "index",
			expectedLangInTmplData: "pl",
		},
		{
			name:                   "Success case - English",
			url:                    "/?lang=en",
			method:                 http.MethodGet,
			mockExecutorResult:     &data.TemplateData{Lang: "en"},
			mockExecutorError:      nil,
			expectedStatusCode:     http.StatusOK,
			expectedTmplName:       "index",
			expectedLangInTmplData: "en",
		},
		{
			name:               "Failure case - Process returns error",
			url:                "/",
			method:             http.MethodGet,
			mockExecutorResult: nil,
			mockExecutorError:  errors.New("something went wrong"),
			expectedStatusCode: http.StatusOK,
			expectedTmplName:   "error",
		},
		{
			name:               "Failure case - Wrong HTTP method",
			url:                "/",
			method:             http.MethodPost,
			mockExecutorResult: nil,
			mockExecutorError:  nil,
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedTmplName:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockExec := &mockExecutor{
				retData: tc.mockExecutorResult,
				retErr:  tc.mockExecutorError,
			}
			mockRend := &mockRenderer{}
			handler := NewHandler(mockExec, mockRend)

			req := httptest.NewRequest(tc.method, tc.url, nil)
			rr := httptest.NewRecorder()

			handler.HandleIndexPage(rr, req)

			if rr.Code != tc.expectedStatusCode {
				t.Errorf("expected status code %d, but got %d", tc.expectedStatusCode, rr.Code)
			}

			if mockRend.calledWithTmplName != tc.expectedTmplName {
				t.Errorf("expected renderer to be called with template '%s', but got '%s'",
					tc.expectedTmplName, mockRend.calledWithTmplName)
			}

			if tc.expectedLangInTmplData != "" {
				if templateData, ok := mockRend.calledWithData.(*data.TemplateData); ok {
					if templateData.Lang != tc.expectedLangInTmplData {
						t.Errorf("expected language in template data to be '%s', but got '%s'",
							tc.expectedLangInTmplData, templateData.Lang)
					}
				} else {
					t.Errorf("data passed to renderer has unexpected type")
				}
			}
		})
	}
}
