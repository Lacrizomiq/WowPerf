package blizzard

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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

// NewClient creates a new Blizzard API client
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

// refreshToken refreshes the token using the provided config
func (c *Client) refreshToken(config *clientcredentials.Config) error {

	log.Println("Refreshing token...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := config.Token(ctx)
	if err != nil {
		log.Printf("Failed to get token: %v", err)
		return fmt.Errorf("failed to get token: %w", err)
	}

	c.token = token
	log.Printf("Token refreshed successfully. Expires at: %v", token.Expiry)
	if scopes, ok := token.Extra("scope").(string); ok {
		log.Printf("Token scopes: %s", scopes)
		if !strings.Contains(scopes, "wow.profile") {
			return fmt.Errorf("token does not have the required 'wow.profile' scope")
		}
	} else {
		log.Println("Unable to retrieve token scopes")
	}
	return nil
}

// makeRequest makes a request to the Blizzard API
func (c *Client) makeRequest(endpoint, namespace, locale string) ([]byte, error) {
	if c.token.Expiry.Before(time.Now()) {
		if err := c.refreshToken(&clientcredentials.Config{
			ClientID:     os.Getenv("BLIZZARD_CLIENT_ID"),
			ClientSecret: os.Getenv("BLIZZARD_CLIENT_SECRET"),
			TokenURL:     authURL,
		}); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
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
	req.Header.Set("Battlenet-Namespace", namespace)
	req.Header.Set("Accept", "application/json")

	log.Printf("Request URL: %s", req.URL.String())
	log.Printf("Request headers: %s", logSafeHeaders(req.Header))

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
		log.Printf("API request failed. Status: %d, Body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("API request failed with status code: %d, Body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// logSafeHeaders logs the headers safely
func logSafeHeaders(headers http.Header) string {
	safeHeaders := make(http.Header)
	for k, v := range headers {
		if k != "Authorization" {
			safeHeaders[k] = v
		}
	}
	return fmt.Sprintf("%v", safeHeaders)
}

// Returns a profile summary for a character.
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

// Returns the Mythic Keystone profile index for a character.
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

// Returns the Mythic Keystone season details for a character.
// Returns a 404 Not Found for characters that have not yet completed a Mythic Keystone dungeon for the specified season.
func (c *Client) GetCharacterMythicKeystoneSeasonDetails(region, realmSlug, characterName, seasonId, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/mythic-keystone-profile/season/%s", region, realmSlug, characterName, seasonId)
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

// Returns a summary of the items equipped by a character.
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

// Returns a summary of a character's specializations.
func (c *Client) GetCharacterSpecializations(region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/specializations", region, realmSlug, characterName)
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

// Returns a summary of the media assets available for a character (such as an avatar render).
func (c *Client) GetCharacterMedia(region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/character-media", region, realmSlug, characterName)
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
