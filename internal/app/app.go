package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlerDownloadCV "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/download_cv"
	handlerGetCVToken "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/get_cv_token"
	handlerIndexPage "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/index_page"
	handlerPrivacyPolicy "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/privacy_policy"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/renderer"
	processDownloadCV "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/download_cv"
	taskDownloadCVLink "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/download_cv/task"
	processGetCVToken "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/get_cv_token"
	taskRequestCVToken "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/get_cv_token/task"
	processIndexPage "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/index_page"
	taskIndexPage "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/index_page/task"
	processPrivacyPolicy "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/privacy_policy"
	taskPrivacyPolicy "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/privacy_policy/task"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/registry"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type App struct {
	server *http.Server
}

func Build(cfg *registry.Config) (*App, error) {
	pageRenderer := renderer.New(cfg.Templates.Path)
	httpClient := &http.Client{Timeout: 10 * time.Second}

	gatewayService := gateway_service.NewClient(httpClient, cfg.Api.BaseURL)

	contentFetcherTask := taskIndexPage.NewContentFetcherTask(gatewayService)
	indexPageProcess := processIndexPage.NewProcess(contentFetcherTask)
	indexPageHandler := handlerIndexPage.NewHandler(indexPageProcess, pageRenderer)

	privacyFetcherTask := taskPrivacyPolicy.NewFetchPrivacyContentTask(gatewayService)
	privacyProcess := processPrivacyPolicy.NewProcess(privacyFetcherTask)
	privacyHandler := handlerPrivacyPolicy.NewHandler(privacyProcess, pageRenderer)

	requestCVTokenTask := taskRequestCVToken.NewRequestCVTokenTask(gatewayService)
	getCVTokenProcess := processGetCVToken.NewProcess(requestCVTokenTask)
	getCVTokenHandler := handlerGetCVToken.NewHandler(getCVTokenProcess)

	validateCVLinkTask := taskDownloadCVLink.NewValidateLinkTask()
	streamCVLinkTask := taskDownloadCVLink.NewFetchPDFStreamTask(gatewayService)
	downloadCVProcess := processDownloadCV.NewProcess(validateCVLinkTask, streamCVLinkTask)
	downloadCVHandler := handlerDownloadCV.NewHandler(downloadCVProcess, gatewayService, pageRenderer)

	mux := http.NewServeMux()
	staticFs := http.FileServer(http.Dir("./internal/web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	mux.HandleFunc("/", indexPageHandler.HandleIndexPage)
	mux.HandleFunc("/privacy-policy", privacyHandler.HandlePrivacyPage)
	mux.HandleFunc("/polityka-prywatnosci", privacyHandler.HandlePrivacyPage)
	mux.HandleFunc("/api/cv-token", getCVTokenHandler.HandleCVRequest)
	mux.HandleFunc("/api/download/cv", downloadCVHandler.HandleDownload)

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
