package gamedata

import (
	"fmt"
	"wowperf/internal/services/blizzard"
)

// GetMythicKeystoneIndex retrieves an index of mythic keystones
func GetMythicKeystoneIndex(s *blizzard.GameDataService, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/mythic-keystone/index", region)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetMythicKeystoneDungeonsIndex retrieves an index of mythic keystone dungeons
func GetMythicKeystoneDungeonsIndex(s *blizzard.GameDataService, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/mythic-keystone/dungeon/index", region)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetMythicKeystone retrieves a mythic keystone by ID
func GetMythicKeystoneByID(s *blizzard.GameDataService, mythicKeystoneID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/mythic-keystone/dungeon/%d", region, mythicKeystoneID)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetMythicKeystonePeriodsIndex retrieves an index of mythic keystone periods
func GetMythicKeystonePeriodsIndex(s *blizzard.GameDataService, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/mythic-keystone/period/index", region)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetMythicKeystone retrieves a mythic keystone periiodby periodID
func GetMythicKeystonePeriodByID(s *blizzard.GameDataService, periodID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/mythic-keystone/period/%d", region, periodID)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetMythicKeystoneDungeons retrieves a mythic keystone dungeon by dungeonID
func GetMythicKeystoneSeasonsIndex(s *blizzard.GameDataService, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/mythic-keystone/season/index", region)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}

// GetMythicKeystoneSeasonByID retrieves a mythic keystone season by seasonID
func GetMythicKeystoneSeasonByID(s *blizzard.GameDataService, seasonID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/mythic-keystone/season/%d", region, seasonID)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}
