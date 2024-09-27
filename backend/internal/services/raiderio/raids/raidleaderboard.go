package raiderioRaid

import (
	"wowperf/internal/services/raiderio"
)

// GetRaidLeaderboard gets the raid leaderboard for a given raid
func GetRaidLeaderboard(s *raiderio.RaiderIOService, raid, difficulty, region string) (map[string]interface{}, error) {
	params := map[string]string{
		"raid":       raid,
		"difficulty": difficulty,
		"region":     region,
	}

	endpoint := "/raiding/progression"
	return s.Client.Get(endpoint, params)
}
