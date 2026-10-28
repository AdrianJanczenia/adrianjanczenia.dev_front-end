package app

import (
	"context"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	handlerDownloadCV "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/download_cv"
	handlerGetCaptcha "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/get_captcha"
	handlerGetCVToken "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/get_cv_token"
	handlerGetPow "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/get_pow"
	handlerIndexPage "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/index_page"
	handlerPrivacyPolicy "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/privacy_policy"
	handlerVerifyCaptcha "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/handler/verify_captcha"
	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/renderer"
	processDownloadCV "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/download_cv"
	taskDownloadCVLink "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/download_cv/task"
	processGetCaptcha "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/get_captcha"
	taskGetCaptcha "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/get_captcha/task"
	processGetCVToken "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/get_cv_token"
	taskRequestCVToken "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/get_cv_token/task"
	processGetPow "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/get_pow"
	taskGetPow "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/get_pow/task"
	processIndexPage "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/index_page"
	taskIndexPage "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/index_page/task"
	processPrivacyPolicy "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/privacy_policy"
	taskPrivacyPolicy "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/privacy_policy/task"
	processVerifyCaptcha "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/verify_captcha"
	taskVerifyCaptcha "github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/process/verify_captcha/task"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/registry"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/service/gateway_service"
)

type App struct {
	server *http.Server
}

func Build(cfg *registry.Config) (*App, error) {
	if cfg == nil || cfg.Server.Port == "" || cfg.Api.GatewayURL == "" || cfg.Templates.Path == "" {
		return nil, fmt.Errorf("invalid front-end configuration")
	}
	layoutBase := filepath.Join(cfg.Templates.Path, "layout", "base.html")
	partials, err := filepath.Glob(filepath.Join(cfg.Templates.Path, "partials", "*.html"))
	if err != nil {
		return nil, fmt.Errorf("find template partials: %w", err)
	}

	templateMap := map[string][]string{
		"index":    append(append([]string{layoutBase}, partials...), filepath.Join(cfg.Templates.Path, "index.html")),
		"privacy":  append(append([]string{layoutBase}, partials...), filepath.Join(cfg.Templates.Path, "privacy.html")),
		"error":    {filepath.Join(cfg.Templates.Path, "error.html")},
		"cv_error": {filepath.Join(cfg.Templates.Path, "cv_error.html")},
	}

	pageRenderer, err := renderer.New(templateMap)
	if err != nil {
		return nil, err
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = 15 * time.Second
	httpClient := &http.Client{Transport: transport}
	gatewayService := gateway_service.NewClient(httpClient, cfg.Api.GatewayURL)

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

	getPowTask := taskGetPow.NewFetchPowTask(gatewayService)
	getPowProcess := processGetPow.NewProcess(getPowTask)
	getPowHandler := handlerGetPow.NewHandler(getPowProcess)

	getCaptchaTask := taskGetCaptcha.NewFetchCaptchaTask(gatewayService)
	getCaptchaProcess := processGetCaptcha.NewProcess(getCaptchaTask)
	getCaptchaHandler := handlerGetCaptcha.NewHandler(getCaptchaProcess)

	verifyCaptchaTask := taskVerifyCaptcha.NewVerifyCaptchaTask(gatewayService)
	verifyCaptchaProcess := processVerifyCaptcha.NewProcess(verifyCaptchaTask)
	verifyCaptchaHandler := handlerVerifyCaptcha.NewHandler(verifyCaptchaProcess)

	mux := http.NewServeMux()
	renderNotFoundPage := func(w http.ResponseWriter, r *http.Request) {
		lang := "pl"
		if r.URL.Query().Get("lang") == "en" {
			lang = "en"
		}
		pageRenderer.RenderError(w, "error", appErrors.ErrRouteNotFound, lang, nil)
	}
	renderUnknownAPI := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Vary", "Accept")
		if prefersHTML(r.Header.Values("Accept")) {
			renderNotFoundPage(w, r)
			return
		}
		appErrors.WriteJSON(w, appErrors.ErrRouteNotFound)
	}
	staticFs := http.FileServer(http.Dir(filepath.Join(filepath.Dir(cfg.Templates.Path), "static")))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	mux.HandleFunc("/{$}", indexPageHandler.Handle)
	mux.HandleFunc("/privacy-policy", privacyHandler.Handle)
	mux.HandleFunc("/polityka-prywatnosci", privacyHandler.Handle)
	mux.HandleFunc("/api/cv-token", getCVTokenHandler.Handle)
	mux.HandleFunc("/api/download/cv", downloadCVHandler.Handle)
	mux.HandleFunc("/api/pow", getPowHandler.Handle)
	mux.HandleFunc("/api/captcha", getCaptchaHandler.Handle)
	mux.HandleFunc("/api/captcha-verify", verifyCaptchaHandler.Handle)
	mux.HandleFunc("/api", renderUnknownAPI)
	mux.HandleFunc("/api/", renderUnknownAPI)
	mux.HandleFunc("/", renderNotFoundPage)

	serverAddr := ":" + cfg.Server.Port
	server := &http.Server{
		Addr: serverAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
			mux.ServeHTTP(w, r)
		}),
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    32 << 10,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       120 * time.Second,
	}

	return &App{server: server}, nil
}

func prefersHTML(acceptValues []string) bool {
	type preference struct {
		quality     float64
		specificity int
	}
	html := preference{specificity: -1}
	json := preference{specificity: -1}
	apply := func(current *preference, quality float64, specificity int) {
		if specificity > current.specificity || (specificity == current.specificity && quality > current.quality) {
			*current = preference{quality: quality, specificity: specificity}
		}
	}
	for _, value := range acceptValues {
		for _, entry := range strings.Split(value, ",") {
			mediaType, params, err := mime.ParseMediaType(strings.TrimSpace(entry))
			if err != nil {
				continue
			}
			quality := 1.0
			if raw, ok := params["q"]; ok {
				quality, err = strconv.ParseFloat(raw, 64)
				if err != nil || !(quality >= 0 && quality <= 1) {
					continue
				}
			}
			switch mediaType {
			case "text/html":
				apply(&html, quality, 2)
			case "text/*":
				apply(&html, quality, 1)
			case "application/json":
				apply(&json, quality, 2)
			case "application/*":
				apply(&json, quality, 1)
			case "*/*":
				apply(&html, quality, 0)
				apply(&json, quality, 0)
			}
		}
	}
	return html.quality > json.quality
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
