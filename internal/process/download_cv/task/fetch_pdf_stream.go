package task

import (
	"io"
)

type GatewayStreamClient interface {
	DownloadCVStream(token, lang string) (io.ReadCloser, string, int, error)
}

type FetchPDFStreamTask struct {
	client GatewayStreamClient
}

func NewFetchPDFStreamTask(client GatewayStreamClient) *FetchPDFStreamTask {
	return &FetchPDFStreamTask{
		client: client,
	}
}

func (t *FetchPDFStreamTask) Execute(token, lang string) (io.ReadCloser, string, int, error) {
	return t.client.DownloadCVStream(token, lang)
}
