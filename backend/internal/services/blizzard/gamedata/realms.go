package gamedata

import (
	"fmt"
	"wowperf/internal/services/blizzard"
)

// GetRealmsIndex retrieves an index of realms
func GetRealmsIndex(s *blizzard.GameDataService, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/realm/index", region)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetConnectedRealmIndex retrieves an index of connected realms
func GetConnectedRealmIndex(s *blizzard.GameDataService, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/connected-realm/index", region)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}
