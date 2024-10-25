// package warcraftlogs/client.go
package warcraftlogs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	authURL      = "https://www.warcraftlogs.com/oauth/token"
	clientAPIURL = "https://www.warcraftlogs.com/api/v2/client"
	userAPIURL   = "https://www.warcraftlogs.com/api/v2/user"
)

type Client struct {
	httpClient *http.Client
	token      *oauth2.Token
	isPublic   bool
}

// NewClient creates a new Warcraft Logs API client for public (client credentials) access
func NewClient() (*Client, error) {
	clientID := os.Getenv("WARCRAFTLOGS_CLIENT_ID")
	clientSecret := os.Getenv("WARCRAFTLOGS_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing required environment variables WARCRAFTLOGS_CLIENT_ID or WARCRAFTLOGS_CLIENT_SECRET")
	}

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     authURL,
	}

	client := &Client{
		httpClient: &http.Client{},
		isPublic:   true,
	}

	if err := client.refreshToken(config); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) refreshToken(config *clientcredentials.Config) error {
	log.Println("Refreshing Warcraft Logs token...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := config.Token(ctx)
	if err != nil {
		log.Printf("Failed to get Warcraft Logs token: %v", err)
		return fmt.Errorf("failed to get token: %w", err)
	}

	c.token = token
	log.Printf("Warcraft Logs token refreshed successfully. Expires at: %v", token.Expiry)
	return nil
}

// MakeGraphQLRequest makes a GraphQL request to the Warcraft Logs API
func (c *Client) MakeGraphQLRequest(query string, variables map[string]interface{}) ([]byte, error) {
	if c.token.Expiry.Before(time.Now()) {
		if err := c.refreshToken(&clientcredentials.Config{
			ClientID:     os.Getenv("WARCRAFTLOGS_CLIENT_ID"),
			ClientSecret: os.Getenv("WARCRAFTLOGS_CLIENT_SECRET"),
			TokenURL:     authURL,
		}); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
	}

	apiURL := clientAPIURL
	if !c.isPublic {
		apiURL = userAPIURL
	}

	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	return body, nil
}
