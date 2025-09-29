package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlerIndexPage "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/index_page"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/renderer"
	processIndexPage "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/index_page"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/index_page/task"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/registry"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/adrianjanczenia.dev_content-service"
)

type App struct {
	server *http.Server
}

func Build(cfg *registry.Config) (*App, error) {
	pageRenderer := renderer.New(cfg.Templates.Path)
	httpClient := &http.Client{Timeout: 10 * time.Second}

	contentService := adrianjanczenia_dev_content_service.NewClient(httpClient, cfg.Api.BaseURL)

	contentFetcherTask := task.NewContentFetcherTask(contentService)

	indexProcess := processIndexPage.NewProcess(contentFetcherTask)

	pageHandler := handlerIndexPage.NewHandler(indexProcess, pageRenderer)

	mux := http.NewServeMux()
	staticFs := http.FileServer(http.Dir("./internal/web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	mux.HandleFunc("/", pageHandler.HandleIndexPage)

	serverAddr := ":" + cfg.Server.Port
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	return &App{server: server}, nil
}

func (a *App) Run() error {
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("INFO: starting server on %s", a.server.Addr)
		serverErrors <- a.server.ListenAndServe()
	}()

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return err
	case sig := <-shutdownChannel:
		log.Printf("INFO: received signal %s. starting graceful shutdown...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := a.server.Shutdown(ctx); err != nil {
			log.Printf("ERROR: could not gracefully shutdown the server: %v", err)
			return a.server.Close()
		}
		log.Println("INFO: server shutdown complete.")
	}

	return nil
}
