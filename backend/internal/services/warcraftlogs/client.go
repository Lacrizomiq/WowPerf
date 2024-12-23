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
	"strconv"
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
			Timeout: time.Second * 30,
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
	// Check if the token is expired and refresh it if necessary
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

	// Prepare the request body
	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, &warcraftlogsTypes.WarcraftLogsError{
			Type:      warcraftlogsTypes.ErrorTypeValidation,
			Message:   "failed to marshal request body",
			Cause:     err,
			Retryable: false,
		}
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, &warcraftlogsTypes.WarcraftLogsError{
			Type:      warcraftlogsTypes.ErrorTypeAPI,
			Message:   "failed to create request",
			Cause:     err,
			Retryable: false,
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &warcraftlogsTypes.WarcraftLogsError{
			Type:      warcraftlogsTypes.ErrorTypeAPI,
			Message:   "failed to send request",
			Cause:     err,
			Retryable: false,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &warcraftlogsTypes.WarcraftLogsError{
			Type:      warcraftlogsTypes.ErrorTypeAPI,
			Message:   "failed to read response body",
			Cause:     err,
			Retryable: false,
		}
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			retryAfter := time.Second * 5 // default
			if s := resp.Header.Get("Retry-After"); s != "" {
				if seconds, err := strconv.Atoi(s); err == nil {
					retryAfter = time.Duration(seconds) * time.Second
				}
			}
			return nil, warcraftlogsTypes.NewRateLimitError(&warcraftlogsTypes.RateLimitInfo{
				ResetIn: retryAfter,
			}, nil)
		case http.StatusUnauthorized:
			return nil, &warcraftlogsTypes.WarcraftLogsError{
				Type:      warcraftlogsTypes.ErrorTypeAPI,
				Message:   "unauthorized: invalid or expired token",
				Retryable: true,
			}
		default:
			return nil, warcraftlogsTypes.NewAPIError(resp.StatusCode, fmt.Errorf("response: %s", body))
		}
	}

	// Parse the response body into a GraphQLResponse
	var graphQLResponse GraphQLResponse
	if err := json.Unmarshal(body, &graphQLResponse); err != nil {
		return nil, &warcraftlogsTypes.WarcraftLogsError{
			Type:      warcraftlogsTypes.ErrorTypeAPI,
			Message:   "failed to parse GraphQL response",
			Cause:     err,
			Retryable: false,
		}
	}

	// Handle GraphQL errors
	if len(graphQLResponse.Errors) > 0 {
		// Check for specific GraphQL error types
		for _, gqlErr := range graphQLResponse.Errors {
			if isRateLimitError(gqlErr) {
				return nil, warcraftlogsTypes.NewRateLimitError(nil, fmt.Errorf(gqlErr.Message))
			}
		}

		// Generic GraphQL error
		return nil, &warcraftlogsTypes.WarcraftLogsError{
			Type:      warcraftlogsTypes.ErrorTypeAPI,
			Message:   fmt.Sprintf("GraphQL error: %s", graphQLResponse.Errors[0].Message),
			Retryable: false,
		}
	}

	return graphQLResponse.Data, nil
}

// isRateLimitError checks if a GraphQL error is a rate limit error
func isRateLimitError(err GraphQLError) bool {
	return strings.Contains(err.Message, "Rate limit exceeded")
}

// containsRateLimitKeywords checks if the error message contains rate limit related keywords
func containsRateLimitKeywords(message string) bool {
	rateLimitKeywords := []string{
		"rate limit",
		"too many requests",
		"quota exceeded",
	}

	messageLower := strings.ToLower(message)
	for _, keyword := range rateLimitKeywords {
		if strings.Contains(messageLower, keyword) {
			return true
		}
	}
	return false
}
