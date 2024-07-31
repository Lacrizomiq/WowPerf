package blizzard

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	authURL = "https://oauth.battle.net/token"
	apiURL  = "https://%s.api.blizzard.com"
)

type Client struct {
	httpClient *http.Client
	region     string
	token      *oauth2.Token
}

func NewClient() (*Client, error) {
	clientID := os.Getenv("BLIZZARD_CLIENT_ID")
	clientSecret := os.Getenv("BLIZZARD_CLIENT_SECRET")
	region := os.Getenv("BLIZZARD_REGION")

	if clientID == "" || clientSecret == "" || region == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     authURL,
	}

	client := &Client{
		httpClient: &http.Client{},
		region:     region,
	}

	if err := client.refreshToken(config); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) refreshToken(config *clientcredentials.Config) error {
	token, err := config.Token(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	c.token = token
	return nil
}

func (c *Client) makeRequest(endpoint, namespace, locale string) ([]byte, error) {
	if c.token.Expiry.Before(time.Now()) {
		if err := c.refreshToken(&clientcredentials.Config{
			ClientID:     os.Getenv("BLIZZARD_CLIENT_ID"),
			ClientSecret: os.Getenv("BLIZZARD_CLIENT_SECRET"),
			TokenURL:     authURL,
		}); err != nil {
			return nil, err
		}
	}

	params := url.Values{}
	params.Add("namespace", namespace)
	if locale != "" {
		params.Add("locale", locale)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", endpoint, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func (c *Client) GetCharacterProfile(region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s", region, realmSlug, characterName)
	body, err := c.makeRequest(endpoint, namespace, locale)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func (c *Client) GetCharacterMythicKeystoneProfile(region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/mythic-keystone-profile", region, realmSlug, characterName)
	body, err := c.makeRequest(endpoint, namespace, locale)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func (c *Client) GetCharacterEquipment(region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/equipment", region, realmSlug, characterName)
	body, err := c.makeRequest(endpoint, namespace, locale)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
