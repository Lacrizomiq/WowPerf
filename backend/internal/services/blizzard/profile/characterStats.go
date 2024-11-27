package profile

import (
	"encoding/json"
	"fmt"
	"wowperf/internal/services/blizzard"
)

// GetCharacterStats returns the stats of a character.
func GetCharacterStats(s *blizzard.ProfileService, region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/statistics", region, realmSlug, characterName)
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
