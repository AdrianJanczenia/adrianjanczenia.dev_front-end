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
	url := fmt.Sprintf("%s/api/v1/content?lang=%s", c.baseURL, lang)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
		var errorResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		if errorResp.Error != "" {
			return nil, errors.FromSlug(errorResp.Error)
		}
		return nil, errors.FromHTTPStatus(resp.StatusCode)
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

func (c *Client) RequestCVToken(ctx context.Context, password, lang string) (string, error) {
	reqBody, err := json.Marshal(map[string]string{
		"password": password,
		"lang":     lang,
	})
	if err != nil {
		return "", errors.ErrInternalServerError
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
		var errorResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		if errorResp.Error != "" {
			return "", errors.FromSlug(errorResp.Error)
		}
		return "", errors.FromHTTPStatus(resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", errors.ErrInternalServerError
	}

	return result.Token, nil
}

func (c *Client) DownloadCVStream(ctx context.Context, token, lang string) (io.ReadCloser, string, error) {
	params := url.Values{}
	params.Add("token", token)
	params.Add("lang", lang)

	fullURL := fmt.Sprintf("%s/api/v1/download/cv?%s", c.baseURL, params.Encode())

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, "", errors.ErrInternalServerError
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", errors.ErrServiceUnavailable
	}

	if resp.StatusCode != http.StatusOK {
		defer cancel()

		var errorResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		resp.Body.Close()

		var appErr *errors.AppError
		if errorResp.Error != "" {
			appErr = errors.FromSlug(errorResp.Error)
		} else {
			appErr = errors.FromHTTPStatus(resp.StatusCode)
		}

		return nil, "", appErr
	}

	return resp.Body, resp.Header.Get("Content-Type"), nil
}
