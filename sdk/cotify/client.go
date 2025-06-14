package cotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client represents a cotify API client
type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient creates a new cotify client
func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{
		baseURL: baseURL,
		client:  httpClient,
	}
}

// StoreRequest represents the request structure for storing an item
type StoreRequest struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Metadata string `json:"metadata"`
}

// StoreResponse represents the response structure for storing an item
type StoreResponse struct {
	Item struct {
		ID       uint   `json:"id"`
		Url      string `json:"url"`
		Title    string `json:"title"`
		Type     string `json:"type"`
		Metadata string `json:"metadata"`
	} `json:"item"`
	IsNew bool `json:"is_new"`
}

// Store stores a new item in the cotify service
func (c *Client) Store(req *StoreRequest) (*StoreResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/items", c.baseURL)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response StoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
