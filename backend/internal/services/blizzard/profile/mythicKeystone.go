package profile

import (
	"encoding/json"
	"fmt"
	"wowperf/internal/services/blizzard"
)

// Returns a summary of the items equipped by a character.
func GetCharacterMythicKeystoneProfile(s *blizzard.ProfileService, region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/mythic-keystone-profile", region, realmSlug, characterName)
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

// Returns the Mythic Keystone season details for a character.
// Returns a 404 Not Found for characters that have not yet completed a Mythic Keystone dungeon for the specified season.
func GetCharacterMythicKeystoneSeasonDetails(s *blizzard.ProfileService, region, realmSlug, characterName, seasonId, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/mythic-keystone-profile/season/%s", region, realmSlug, characterName, seasonId)
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
