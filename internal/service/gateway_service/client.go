package gateway_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_front-end/internal/logic/errors"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *Client) GetPageContent(ctx context.Context, lang string) (*PageContent, error) {
	params := url.Values{}
	params.Set("lang", lang)
	fullURL := fmt.Sprintf("%s/api/v1/content?%s", c.baseURL, params.Encode())

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	req.Header.Set("User-Agent", "PortfolioFrontend/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.ErrServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, decodeErrorResponse(resp)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, errors.ErrInternalServerError
	}

	var pageContent PageContent
	if err := json.NewDecoder(resp.Body).Decode(&pageContent); err != nil {
		return nil, errors.ErrInternalServerError
	}

	return &pageContent, nil
}

func (c *Client) RequestCVToken(ctx context.Context, password, lang, captchaID string) (string, error) {
	reqBody, err := json.Marshal(map[string]string{
		"password":  password,
		"lang":      lang,
		"captchaId": captchaID,
	})
	if err != nil {
		return "", errors.ErrInternalServerError
	}

	ctx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/cv-request", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", errors.ErrInternalServerError
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", errors.ErrServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", decodeErrorResponse(resp)
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", errors.ErrInternalServerError
	}
	if result.Token == "" {
		return "", errors.ErrServiceUnavailable
	}

	return result.Token, nil
}

func (c *Client) DownloadCVStream(ctx context.Context, token, lang string) (io.ReadCloser, string, error) {
	params := url.Values{}
	params.Add("token", token)
	params.Add("lang", lang)

	fullURL := fmt.Sprintf("%s/api/v1/download/cv?%s", c.baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, "", errors.ErrInternalServerError
	}

	streamClient := *c.httpClient
	streamClient.Timeout = 0
	resp, err := streamClient.Do(req)
	if err != nil {
		return nil, "", errors.ErrServiceUnavailable
	}

	if resp.StatusCode != http.StatusOK {
		appErr := decodeErrorResponse(resp)
		resp.Body.Close()
		return nil, "", appErr
	}

	return resp.Body, resp.Header.Get("Content-Type"), nil
}

func (c *Client) GetPow(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v1/pow", nil)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.ErrServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, decodeErrorResponse(resp)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.ErrInternalServerError
	}

	return result, nil
}

func (c *Client) GetCaptcha(ctx context.Context, seed, signature, nonce string) (map[string]string, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"seed":      seed,
		"signature": signature,
		"nonce":     nonce,
	})

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/captcha", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.ErrServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, decodeErrorResponse(resp)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.ErrInternalServerError
	}

	return result, nil
}

func (c *Client) VerifyCaptcha(ctx context.Context, captchaID, captchaValue string) (string, error) {
	reqBody, _ := json.Marshal(map[string]string{
		"captchaId":    captchaID,
		"captchaValue": captchaValue,
	})

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/captcha-verify", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", errors.ErrInternalServerError
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", errors.ErrServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", decodeErrorResponse(resp)
	}

	var result struct {
		CaptchaID string `json:"captchaId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", errors.ErrInternalServerError
	}

	return result.CaptchaID, nil
}

func decodeErrorResponse(resp *http.Response) *errors.AppError {
	data, err := io.ReadAll(io.LimitReader(resp.Body, 4097))
	if err != nil || len(data) > 4096 {
		return errors.ErrBadGateway
	}

	var payload struct {
		Error string `json:"error"`
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return errors.ErrBadGateway
	}

	return errors.FromUpstream(resp.StatusCode, payload.Error)
}
