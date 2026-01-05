package task

import (
	"io"
)

type GatewayClient interface {
	DownloadCVStream(token, lang string) (io.ReadCloser, string, error)
}

type FetchPDFStreamTask struct {
	client GatewayClient
}

func NewFetchPDFStreamTask(client GatewayClient) *FetchPDFStreamTask {
	return &FetchPDFStreamTask{
		client: client,
	}
}

func (t *FetchPDFStreamTask) Execute(token, lang string) (io.ReadCloser, string, error) {
	return t.client.DownloadCVStream(token, lang)
}
