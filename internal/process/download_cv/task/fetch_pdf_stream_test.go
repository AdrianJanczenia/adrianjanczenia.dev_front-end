package task

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type mockGatewayClient struct {
	downloadCVStreamFunc func(token, lang string) (io.ReadCloser, string, error)
}

func (m *mockGatewayClient) DownloadCVStream(token, lang string) (io.ReadCloser, string, error) {
	return m.downloadCVStreamFunc(token, lang)
}

func TestFetchPDFStreamTask_Execute(t *testing.T) {
	m := &mockGatewayClient{
		downloadCVStreamFunc: func(t, l string) (io.ReadCloser, string, error) {
			if t == "valid" {
				return io.NopCloser(strings.NewReader("pdf")), "application/pdf", nil
			}
			return nil, "", errors.New("fail")
		},
	}
	task := NewFetchPDFStreamTask(m)

	t.Run("success", func(t *testing.T) {
		stream, contentType, err := task.Execute("valid", "pl")
		if err != nil || contentType != "application/pdf" {
			t.Errorf("unexpected results: %v, %s", err, contentType)
		}
		stream.Close()
	})

	t.Run("error", func(t *testing.T) {
		_, _, err := task.Execute("invalid", "pl")
		if err == nil {
			t.Error("expected error")
		}
	})
}
