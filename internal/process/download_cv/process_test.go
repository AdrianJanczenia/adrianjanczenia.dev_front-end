package download_cv

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

type mockValidator struct {
	executeFunc func(token, lang string) error
}

func (m *mockValidator) Execute(token, lang string) error {
	return m.executeFunc(token, lang)
}

type mockStreamer struct {
	executeFunc func(ctx context.Context, token, lang string) (io.ReadCloser, string, error)
}

func (m *mockStreamer) Execute(ctx context.Context, token, lang string) (io.ReadCloser, string, error) {
	return m.executeFunc(ctx, token, lang)
}

func TestProcess_DownloadCV(t *testing.T) {
	t.Run("validation fails", func(t *testing.T) {
		v := &mockValidator{executeFunc: func(t, l string) error { return errors.New("val fail") }}
		s := &mockStreamer{}
		p := NewProcess(v, s)
		_, _, err := p.Process(context.Background(), "t", "l")
		if err == nil || err.Error() != "val fail" {
			t.Errorf("expected validation error, got %v", err)
		}
	})

	t.Run("streaming succeeds", func(t *testing.T) {
		v := &mockValidator{executeFunc: func(t, l string) error { return nil }}
		s := &mockStreamer{executeFunc: func(ctx context.Context, t, l string) (io.ReadCloser, string, error) {
			return io.NopCloser(strings.NewReader("data")), "type", nil
		}}
		p := NewProcess(v, s)
		stream, _, err := p.Process(context.Background(), "t", "l")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		stream.Close()
	})
}
