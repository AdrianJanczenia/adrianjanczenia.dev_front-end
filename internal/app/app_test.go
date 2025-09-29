package app

import (
	"testing"
	"time"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/registry"
)

func TestBuild(t *testing.T) {
	t.Run("it builds the application without errors", func(t *testing.T) {
		cfg := &registry.Config{}
		cfg.Server.Port = "8081"
		cfg.Server.ReadTimeout = 5 * time.Second
		cfg.Server.WriteTimeout = 10 * time.Second
		cfg.Api.BaseURL = "http://localhost:9999"
		cfg.Templates.Path = "../../internal/web/template"

		app, err := Build(cfg)

		if err != nil {
			t.Fatalf("expected no error from Build, but got %v", err)
		}

		if app == nil {
			t.Fatal("expected app not to be nil")
		}

		if app.server == nil {
			t.Fatal("expected app.server not to be nil")
		}

		if app.server.Addr != ":8081" {
			t.Errorf("expected server address to be ':8081', but got '%s'", app.server.Addr)
		}
	})
}
