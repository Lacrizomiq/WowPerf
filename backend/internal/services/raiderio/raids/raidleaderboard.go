package raiderioRaid

import (
	"strconv"
	"wowperf/internal/services/raiderio"
)

// GetRaidLeaderboard gets the raid leaderboard for a given raid
func GetRaidLeaderboard(s *raiderio.RaiderIOService, raid, difficulty, region string, limit, page int) (map[string]interface{}, error) {
	params := map[string]string{
		"raid":       raid,
		"difficulty": difficulty,
		"region":     region,
		"limit":      strconv.Itoa(limit),
		"page":       strconv.Itoa(page),
	}

	endpoint := "/raiding/raid-rankings"
	return s.Client.Get(endpoint, params)
}
