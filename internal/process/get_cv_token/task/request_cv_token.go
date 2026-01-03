package task

type GatewayClient interface {
	RequestCVToken(password, lang string) (string, error)
}

type RequestCVTokenTask struct {
	client GatewayClient
}

func NewRequestCVTokenTask(client GatewayClient) *RequestCVTokenTask {
	return &RequestCVTokenTask{
		client: client,
	}
}

func (t *RequestCVTokenTask) Execute(password, lang string) (string, error) {
	return t.client.RequestCVToken(password, lang)
}
