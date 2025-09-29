package index_page

import (
	"errors"
	"reflect"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/data"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/adrianjanczenia.dev_content-service"
)

type mockContentFetcher struct {
	retContent *adrianjanczenia_dev_content_service.PageContent
	retErr     error
}

func (m *mockContentFetcher) Fetch(lang string) (*adrianjanczenia_dev_content_service.PageContent, error) {
	return m.retContent, m.retErr
}

func TestProcess_Execute(t *testing.T) {
	mockContent := &adrianjanczenia_dev_content_service.PageContent{
		Profile: adrianjanczenia_dev_content_service.Profile{Name: "Test Name"},
	}
	mockError := errors.New("task failed")

	testCases := []struct {
		name              string
		lang              string
		mockFetcherResult *adrianjanczenia_dev_content_service.PageContent
		mockFetcherError  error
		expectedResult    *data.TemplateData
		expectedError     error
	}{
		{
			name:              "Success case",
			lang:              "pl",
			mockFetcherResult: mockContent,
			mockFetcherError:  nil,
			expectedResult: &data.TemplateData{
				Lang:    "pl",
				Content: *mockContent,
			},
			expectedError: nil,
		},
		{
			name:              "Failure case - task returns error",
			lang:              "en",
			mockFetcherResult: nil,
			mockFetcherError:  mockError,
			expectedResult:    nil,
			expectedError:     mockError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockFetcher := &mockContentFetcher{
				retContent: tc.mockFetcherResult,
				retErr:     tc.mockFetcherError,
			}
			process := NewProcess(mockFetcher)

			result, err := process.Execute(tc.lang)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error '%v', but got '%v'", tc.expectedError, err)
			}

			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Errorf("expected result '%+v', but got '%+v'", tc.expectedResult, result)
			}
		})
	}
}
