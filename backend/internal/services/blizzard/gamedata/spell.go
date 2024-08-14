package gamedata

import (
	"fmt"
	"wowperf/internal/services/blizzard"
)

// GetSpellMedia retrieves the media assets for a spell
func GetSpellMedia(s *blizzard.GameDataService, spellId int, region, namespace, locale string) (map[string]interface{}, error) {
	baseURL := fmt.Sprintf("https://%s.api.blizzard.com", region)
	if region == "cn" {
		baseURL = "https://gateway.battlenet.com.cn"
	}

	endpoint := fmt.Sprintf("%s/data/wow/media/spell/%d", baseURL, spellId)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}
