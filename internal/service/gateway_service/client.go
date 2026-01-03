package gateway_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

func (c *Client) GetPageContent(lang string) (*PageContent, error) {
	url := fmt.Sprintf("%s/api/v1/content?lang=%s", c.baseURL, lang)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "PortfolioFrontend/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned non-200 status: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("unexpected content type: got %s, want application/json", contentType)
	}

	var pageContent PageContent
	if err := json.NewDecoder(resp.Body).Decode(&pageContent); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return &pageContent, nil
}

func (c *Client) RequestCVToken(password, lang string) (string, error) {
	reqBody, err := json.Marshal(map[string]string{
		"password": password,
		"lang":     lang,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/v1/cv-request", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status_%d", resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}

func (c *Client) DownloadCVStream(token, lang string) (io.ReadCloser, string, int, error) {
	params := url.Values{}
	params.Add("token", token)
	params.Add("lang", lang)

	fullURL := fmt.Sprintf("%s/api/v1/download/cv?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(fullURL)
	if err != nil {
		return nil, "", http.StatusInternalServerError, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, "", resp.StatusCode, fmt.Errorf("status_%d", resp.StatusCode)
	}

	return resp.Body, resp.Header.Get("Content-Type"), resp.StatusCode, nil
}
