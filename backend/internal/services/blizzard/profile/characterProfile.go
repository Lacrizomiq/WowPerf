package profile

import (
	"encoding/json"
	"fmt"
	"wowperf/internal/services/blizzard"
)

const apiURL = "https://%s.api.blizzard.com"

// Returns a profile summary for a character.
func GetCharacterProfile(s *blizzard.ProfileService, region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s", region, realmSlug, characterName)
	body, err := s.Client.MakeRequest(endpoint, namespace, locale)
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
func GetCharacterMedia(s *blizzard.ProfileService, region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/character-media", region, realmSlug, characterName)
	body, err := s.Client.MakeRequest(endpoint, namespace, locale)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
