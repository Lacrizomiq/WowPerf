package warcraftlogs

import (
	"context"
	"fmt"
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

// NewClient creates a new WarcraftLogs client for public (client credentials) access.
func NewClient() (*Client, error) {
	clientID := os.Getenv("WARCRAFTLOGS_CLIENT_ID")
	clientSecret := os.Getenv("WARCRAFTLOGS_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing required environment variables")
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

// refreshToken refreshes the token using the provided config.
func (c *Client) refreshToken(config *clientcredentials.Config) error {
	log.Println("Refreshing WarcraftLogs token...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := config.Token(ctx)
	if err != nil {
		log.Printf("Failed to refresh WarcraftLogs token: %v", err)
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	c.token = token
	log.Println("WarcraftLogs token refreshed successfully")
	return nil
}
