package task

type GatewayClient interface {
	RequestCV(password, lang string) (string, error)
}

type RequestCVLinkTask struct {
	client GatewayClient
}

func NewRequestCVLinkTask(client GatewayClient) *RequestCVLinkTask {
	return &RequestCVLinkTask{
		client: client,
	}
}

func (t *RequestCVLinkTask) Execute(password, lang string) (string, error) {
	return t.client.RequestCV(password, lang)
}
