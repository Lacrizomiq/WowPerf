package raiderio

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const apiURL = "https://raider.io/api/v1"

type RaiderIOClient struct {
	httpClient *http.Client
	baseURL    string
	limiter    *rate.Limiter
}

type RaiderIOService struct {
	Client *RaiderIOClient
}

func NewRaiderIOClient() (*RaiderIOClient, error) {
	client := &RaiderIOClient{
		httpClient: &http.Client{
			Timeout: 3 * time.Minute,
		},
		baseURL: apiURL,
		limiter: rate.NewLimiter(rate.Every(2*time.Second), 1), // 1 requête toutes les 2 secondes
	}
	return client, nil
}

// doRequestWithRetry performs a request with rate limiting and retries
func (c *RaiderIOClient) doRequestWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err = c.limiter.Wait(req.Context())
		if err != nil {
			return nil, err
		}

		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode != http.StatusGatewayTimeout {
			return resp, nil
		}

		if i < maxRetries-1 {
			waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
			time.Sleep(waitTime)
		}
	}
	return nil, err
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

	// Use the new doRequestWithRetry method
	log.Printf("Making GET request to: %s", url)

	resp, err := c.doRequestWithRetry(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	log.Printf("Received response with status code: %d", resp.StatusCode)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Verify the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode == http.StatusGatewayTimeout {
		return nil, fmt.Errorf("gateway timeout: the server took too long to respond")
	}

	// Parse the JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(body))
	}

	return result, nil
}
