package gamedata

import (
	"fmt"
	"wowperf/internal/services/blizzard"
)

// GetItemMedia retrieves the media assets for an item
func GetItemMedia(c *blizzard.GameDataClient, itemID int, region, namespace, locale string) (map[string]interface{}, error) {
	baseURL := fmt.Sprintf("https://%s.api.blizzard.com", region)
	if region == "cn" {
		baseURL = "https://gateway.battlenet.com.cn"
	}

	endpoint := fmt.Sprintf("%s/data/wow/media/item/%d", baseURL, itemID)
	return c.MakeRequest(endpoint, namespace, locale)
}
