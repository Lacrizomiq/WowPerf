package raiderio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const apiURL = "https://raider.io/api/v1"

type RaiderIOClient struct {
	httpClient *http.Client
	baseURL    string
}

type RaiderIOService struct {
	Client *RaiderIOClient
}

func NewRaiderIOClient() (*RaiderIOClient, error) {
	client := &RaiderIOClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: apiURL,
	}
	return client, nil
}

// Get makes a GET request to the Raider.io API
func (c *RaiderIOClient) Get(endpoint string, params map[string]string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	// Add query parameters to the URL if they are provided
	if len(params) > 0 {
		url += "?"
		for key, value := range params {
			url += fmt.Sprintf("%s=%s&", key, value)
		}
		url = url[:len(url)-1] // Remove the last '&'
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "wowperf/1.0")

	// Create a new GET request
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Verify the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(body))
	}

	return result, nil
}
