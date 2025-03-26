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
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	warcraftlogsTypes "wowperf/internal/services/warcraftlogs/types"
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

// GraphQLResponse is the response from the Warcraft Logs API
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

// GraphQLError is the error from the Warcraft Logs API
type GraphQLError struct {
	Message    string          `json:"message"`
	Path       []string        `json:"path,omitempty"`
	Extensions json.RawMessage `json:"extensions,omitempty"`
}

// NewClient creates a new Warcraft Logs API client for public (client credentials) access
func NewClient() (*Client, error) {
	clientID := os.Getenv("WARCRAFTLOGS_CLIENT_ID")
	clientSecret := os.Getenv("WARCRAFTLOGS_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, &warcraftlogsTypes.WarcraftLogsError{
			Type:      warcraftlogsTypes.ErrorTypeValidation,
			Message:   "missing required environment variables",
			Retryable: false,
		}
	}

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     authURL,
	}

	client := &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 60, // Increased to 60 seconds
		},
		isPublic: true,
	}

	if err := client.refreshToken(config); err != nil {
		return nil, err
	}

	return client, nil
}

// refreshToken obtains or refreshes the OAuth token
func (c *Client) refreshToken(config *clientcredentials.Config) error {
	log.Println("[DEBUG] Refreshing WarcraftLogs token...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	token, err := config.Token(ctx)
	if err != nil {
		return &warcraftlogsTypes.WarcraftLogsError{
			Type:      warcraftlogsTypes.ErrorTypeAPI,
			Message:   "failed to obtain OAuth token",
			Cause:     err,
			Retryable: true,
		}
	}

	c.token = token
	log.Printf("[DEBUG] WarcraftLogs token refreshed successfully. Expires at: %v", token.Expiry)
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
		return nil, &warcraftlogsTypes.WarcraftLogsError{
			Type:      warcraftlogsTypes.ErrorTypeNetwork,
			Message:   fmt.Sprintf("failed to send request: %v", err),
			Cause:     err,
			Retryable: true,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 500 {
		return nil, warcraftlogsTypes.NewAPIError(resp.StatusCode, nil)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, warcraftlogsTypes.NewAPIError(resp.StatusCode, nil)
	}

	var graphQLResp GraphQLResponse
	if err := json.Unmarshal(body, &graphQLResp); err != nil {
		return nil, fmt.Errorf("failed to parse GraphQL response: %w", err)
	}

	if len(graphQLResp.Errors) > 0 {
		for _, gqlErr := range graphQLResp.Errors {
			if isRateLimitError(gqlErr) {
				return nil, &warcraftlogsTypes.WarcraftLogsError{
					Type:      warcraftlogsTypes.ErrorTypeRateLimit,
					Message:   gqlErr.Message,
					Retryable: true,
				}
			}
		}
		return nil, warcraftlogsTypes.NewAPIError(resp.StatusCode, fmt.Errorf(graphQLResp.Errors[0].Message))
	}

	return graphQLResp.Data, nil
}

// isRateLimitError checks if a GraphQL error is a rate limit error
func isRateLimitError(err GraphQLError) bool {
	return strings.Contains(err.Message, "Rate limit exceeded")
}
