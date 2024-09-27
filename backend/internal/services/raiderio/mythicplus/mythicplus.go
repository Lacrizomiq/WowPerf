package raiderioMythicPlus

import (
	"strconv"
	"wowperf/internal/services/raiderio"
)

func GetMythicPlusBestRuns(s *raiderio.RaiderIOService, season, region, dungeon string, page int) (map[string]interface{}, error) {
	params := map[string]string{
		"season":  season,
		"region":  region,
		"dungeon": dungeon,
		"page":    strconv.Itoa(page),
	}

	endpoint := "/mythic-plus/runs"
	return s.Client.Get(endpoint, params)
}
