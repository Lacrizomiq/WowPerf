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

type GameDataClient struct {
	httpClient *http.Client
	region     string
	token      *oauth2.Token
}

// NewGameDataClient creates a new Blizzard Game Data API client
func NewGameDataClient() (*GameDataClient, error) {
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

	client := &GameDataClient{
		httpClient: &http.Client{},
		region:     region,
	}

	if err := client.refreshToken(config); err != nil {
		return nil, err
	}

	return client, nil
}

// refreshToken refreshes the token using the provided config
func (c *GameDataClient) refreshToken(config *clientcredentials.Config) error {

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

// makeRequest makes a request to the Blizzard Game Data API
func (c *GameDataClient) makeRequest(endpoint, namespace, locale string) (map[string]interface{}, error) {
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

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetItemMedia retrieves the media assets for an item
func (c *GameDataClient) GetItemMedia(itemID int, region, namespace, locale string) (map[string]interface{}, error) {
	baseURL := fmt.Sprintf("https://%s.api.blizzard.com", region)
	if region == "cn" {
		baseURL = "https://gateway.battlenet.com.cn"
	}

	endpoint := fmt.Sprintf("%s/data/wow/media/item/%d", baseURL, itemID)
	return c.makeRequest(endpoint, namespace, locale)
}

// GetSpellMedia retrieves the media assets for a spell
func (c *GameDataClient) GetSpellMedia(spellId int, region, namespace, locale string) (map[string]interface{}, error) {
	baseURL := fmt.Sprintf("https://%s.api.blizzard.com", region)
	if region == "cn" {
		baseURL = "https://gateway.battlenet.com.cn"
	}

	endpoint := fmt.Sprintf("%s/data/wow/media/spell/%d", baseURL, spellId)
	return c.makeRequest(endpoint, namespace, locale)
}

// GetPlayableSpecializationIndex retrieves an index of playable specializations
func (c *GameDataClient) GetPlayableSpecializationIndex(region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/playable-specialization/index", region)
	return c.makeRequest(endpoint, namespace, locale)
}

// GetPlayableSpecialization retrieves a playable specialization
func (c *GameDataClient) GetPlayableSpecialization(specID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/playable-specialization/%d", region, specID)
	return c.makeRequest(endpoint, namespace, locale)
}

// GetPlayableSpecializationMedia retrieves the media assets for a playable specialization
func (c *GameDataClient) GetPlayableSpecializationMedia(specID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/media/playable-specialization/%d", region, specID)
	return c.makeRequest(endpoint, namespace, locale)
}

// GetTalentTreeIndex retrieves an index of talent trees
func (c *GameDataClient) GetTalentTreeIndex(region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/talent-tree/index", region)
	return c.makeRequest(endpoint, namespace, locale)
}

// GetTalentTree retrieves a talent tree by spec ID
func (c *GameDataClient) GetTalentTree(talentTreeID, specID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/talent-tree/%d/playable-specialization/%d", region, talentTreeID, specID)
	return c.makeRequest(endpoint, namespace, locale)
}

// GetTalentTreeNodes retrieves the nodes of a talent tree as well as links to associated playable specializations given a talent tree id
func (c *GameDataClient) GetTalentTreeNodes(talentTreeID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/talent-tree/%d", region, talentTreeID)
	return c.makeRequest(endpoint, namespace, locale)
}

// GetTalentIndex retrieves an index of talents
func (c *GameDataClient) GetTalentIndex(region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/talent/index", region)
	return c.makeRequest(endpoint, namespace, locale)
}

// GetTalentByID retrieves a talent by ID
func (c *GameDataClient) GetTalentByID(talentID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/talent/%d", region, talentID)
	return c.makeRequest(endpoint, namespace, locale)
}

// GetPlayableClassIndex retrieves an index of playable classes
func (c *GameDataClient) GetPlayableClassIndex(region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/playable-class/index", region)
	return c.makeRequest(endpoint, namespace, locale)
}
