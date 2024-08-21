package gamedata

import (
	"fmt"
	"wowperf/internal/services/blizzard"
)

// GetMythicKeystoneAffixIndex retrieves an index of mythic keystone affixes
func GetMythicKeystoneAffixIndex(s *blizzard.GameDataService, region, locale string) (map[string]interface{}, error) {
	namespace := fmt.Sprintf("static-%s", region)
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/keystone-affix/index", region)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetMythicKeystoneAffix retrieves a mythic keystone affix by ID
func GetMythicKeystoneAffixByID(s *blizzard.GameDataService, affixID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/keystone-affix/%d", region, affixID)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetMythicKeystoneAffixMedia retrieves the media assets for a mythic keystone affix
func GetMythicKeystoneAffixMedia(s *blizzard.GameDataService, affixID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/media/keystone-affix/%d", region, affixID)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}
