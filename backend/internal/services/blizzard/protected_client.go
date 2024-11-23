// Package protected_client provides a client for the Blizzard API that is protected by OAuth.
package blizzard

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

// ProtectedClient is a client for the Blizzard API that is protected by OAuth.
type ProtectedClient struct {
	httpClient  *http.Client
	region      string
	oauthConfig *oauth2.Config
}

// NewProtectedClient creates a new ProtectedClient.
func NewProtectedClient(region string, oauthConfig *oauth2.Config) *ProtectedClient {
	return &ProtectedClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		region:      region,
		oauthConfig: oauthConfig,
	}
}

// MakeProtectedRequest makes a request to the Blizzard API with OAuth protection.
func (c *ProtectedClient) MakeProtectedRequest(ctx context.Context, userToken *oauth2.Token, endpoint, namespace, locale string) ([]byte, error) {
	// Verify if the user token is valid
	if userToken.Expiry.Before(time.Now()) {
		log.Println("User token is expired, refreshing")
		return nil, fmt.Errorf("user token is expired")
	}

	// Create the URL with the params
	params := url.Values{}
	params.Add("namespace", namespace)
	if locale != "" {
		params.Add("locale", locale)
	}

	apiEndpoint := fmt.Sprintf("https://%s.api.blizzard.com%s", c.region, endpoint)
	if len(params) > 0 {
		apiEndpoint = fmt.Sprintf("%s?%s", apiEndpoint, params.Encode())
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add the Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken.AccessToken))
	req.Header.Set("Battlenet-Namespace", namespace)
	req.Header.Set("Accept", "application/json")

	log.Printf("Making protected request to: %s", apiEndpoint)
	log.Printf("Headers: %s", logSafeHeaders(req.Header))

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Verify the response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("API request failed. Status: %d, Body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	return body, nil
}
