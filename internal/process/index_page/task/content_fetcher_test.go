package task

import (
	"errors"
	"reflect"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/adrianjanczenia.dev_content-service"
)

type mockContentProvider struct {
	retContent *adrianjanczenia_dev_content_service.PageContent
	retErr     error
}

func (m *mockContentProvider) GetPageContent(lang string) (*adrianjanczenia_dev_content_service.PageContent, error) {
	return m.retContent, m.retErr
}

func TestContentFetcherTask_Fetch(t *testing.T) {
	mockContent := &adrianjanczenia_dev_content_service.PageContent{
		Profile: adrianjanczenia_dev_content_service.Profile{Name: "Test Name"},
	}
	mockError := errors.New("service failed")

	testCases := []struct {
		name               string
		lang               string
		mockProviderResult *adrianjanczenia_dev_content_service.PageContent
		mockProviderError  error
		expectedResult     *adrianjanczenia_dev_content_service.PageContent
		expectedError      error
	}{
		{
			name:               "Success case",
			lang:               "pl",
			mockProviderResult: mockContent,
			mockProviderError:  nil,
			expectedResult:     mockContent,
			expectedError:      nil,
		},
		{
			name:               "Failure case - provider returns error",
			lang:               "en",
			mockProviderResult: nil,
			mockProviderError:  mockError,
			expectedResult:     nil,
			expectedError:      mockError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockProvider := &mockContentProvider{
				retContent: tc.mockProviderResult,
				retErr:     tc.mockProviderError,
			}
			task := NewContentFetcherTask(mockProvider)

			result, err := task.Fetch(tc.lang)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error '%v', but got '%v'", tc.expectedError, err)
			}

			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Errorf("expected result '%+v', but got '%+v'", tc.expectedResult, result)
			}
		})
	}
}
