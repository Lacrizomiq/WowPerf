package profile

import (
	"encoding/json"
	"fmt"
	"wowperf/internal/services/blizzard"
)

// Returns a summary of the items equipped by a character.
func GetCharacterEquipment(c *blizzard.Client, region, realmSlug, characterName, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf(apiURL+"/profile/wow/character/%s/%s/equipment", region, realmSlug, characterName)
	body, err := c.MakeRequest(endpoint, namespace, locale)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
