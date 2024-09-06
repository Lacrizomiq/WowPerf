package profile

import (
	"encoding/json"
	"fmt"
	"wowperf/internal/services/blizzard"
)

// GetCharacterEncounterSummary returns the encounter summary for a character.
func GetCharacterEncounterSummary(s *blizzard.ProfileService, region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/encounters", region, realmSlug, characterName)
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

// GetCharacterDungeonEncounters returns the dungeon encounters for a character.
func GetCharacterDungeonEncounters(s *blizzard.ProfileService, region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/encounters/dungeons", region, realmSlug, characterName)
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

// GetCharacterRaidEncounters returns the raid encounters for a character.
func GetCharacterRaidEncounters(s *blizzard.ProfileService, region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/encounters/raids", region, realmSlug, characterName)
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
