package download_cv

import (
	"io"
)

type Validator interface {
	Execute(token, lang string) error
}

type Streamer interface {
	Execute(token, lang string) (io.ReadCloser, string, int, error)
}

type Process struct {
	validator Validator
	streamer  Streamer
}

func NewProcess(v Validator, s Streamer) *Process {
	return &Process{
		validator: v,
		streamer:  s,
	}
}

func (p *Process) Execute(token, lang string) (io.ReadCloser, string, int, error) {
	if err := p.validator.Execute(token, lang); err != nil {
		return nil, "", 400, err
	}
	return p.streamer.Execute(token, lang)
}
