package download_cv

import (
	"io"
)

type Validator interface {
	Execute(token, lang string) error
}

type Streamer interface {
	Execute(token, lang string) (io.ReadCloser, string, error)
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

func (p *Process) Process(token, lang string) (io.ReadCloser, string, error) {
	if err := p.validator.Execute(token, lang); err != nil {
		return nil, "", err
	}
	return p.streamer.Execute(token, lang)
}
