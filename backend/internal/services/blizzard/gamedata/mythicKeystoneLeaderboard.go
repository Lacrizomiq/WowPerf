package gamedata

import (
	"fmt"
	"wowperf/internal/services/blizzard"
)

// GetMythicKeystoneLeaderboardIndex retrieves an index of mythic keystone leaderboards for a connected realm.
func GetMythicKeystoneLeaderboardIndex(s *blizzard.GameDataService, connectedRealmID int, region, namespace, locale string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/connected-realm/%d/mythic-leaderboard/index", region, connectedRealmID)
	return s.Client.MakeRequest(endpoint, namespace, locale)
}
