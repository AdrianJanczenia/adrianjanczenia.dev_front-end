package main

import (
	"log"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/app"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/registry"
)

func main() {
	cfg, err := registry.LoadConfig()
	if err != nil {
		log.Fatalf("FATAL: could not load configuration: %v", err)
	}
	registry.Cfg = cfg

	application, err := app.Build(registry.Cfg)
	if err != nil {
		log.Fatalf("FATAL: could not build application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("FATAL: server failed to start: %v", err)
	}
}
