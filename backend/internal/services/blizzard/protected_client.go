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

	"wowperf/internal/services/blizzard/auth"
)

// ProtectedClient is a client for the Blizzard API that is protected by OAuth.
type ProtectedClient struct {
	httpClient    *http.Client
	region        string
	battleNetAuth *auth.BattleNetAuthService
}

// NewProtectedClient creates a new ProtectedClient.
func NewProtectedClient(region string, battleNetAuth *auth.BattleNetAuthService) *ProtectedClient {
	return &ProtectedClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		region:        region,
		battleNetAuth: battleNetAuth,
	}
}

// MakeProtectedRequest makes a request to the Blizzard API with OAuth protection.
func (c *ProtectedClient) MakeProtectedRequest(ctx context.Context, userID uint, endpoint, namespace, locale string) ([]byte, error) {
	// Get the user token
	token, err := c.battleNetAuth.GetUserToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user token: %w", err)
	}
	log.Printf("Token: %s", token.AccessToken)

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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
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
