package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
	 baseURL: baseURL,
	 httpClient: &http.Client{
	  Timeout: timeout,
	 },
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, headers map[string]string, body io.Reader, v any) error {
	url := c.baseURL + path
	
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
	 return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
	 req.Header.Set(key, value)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
	 return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
	 bodyBytes, _ := io.ReadAll(resp.Body)
	 return fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
	}
	
	if v != nil {
	 if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
	  return fmt.Errorf("failed to decode response: %w", err)
	 }
	}
	
	return nil
}
