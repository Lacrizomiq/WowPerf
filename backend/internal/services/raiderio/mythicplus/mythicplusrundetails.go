package raiderioMythicPlus

import (
	"strconv"
	"wowperf/internal/services/raiderio"
)

func GetMythicPlusRunsDetails(s *raiderio.RaiderIOService, season string, id int) (map[string]interface{}, error) {
	params := map[string]string{
		"season": season,
		"id":     strconv.Itoa(id),
	}

	endpoint := "/mythic-plus/run-details"
	return s.Client.Get(endpoint, params)
}
