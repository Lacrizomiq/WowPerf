package raiderio

import (
	"fmt"
	"time"
	"wowperf/internal/services/raiderio"
	raiderioRaid "wowperf/internal/services/raiderio/raids"
	"wowperf/pkg/cache"
)

func StartRaidLeaderboardCacheUpdater(service *raiderio.RaiderIOService) {
	updateFunc := func() error {
		raids := []string{"nerubar-palace"}
		difficulties := []string{"normal", "heroic", "mythic"}
		regions := []string{"world", "us", "eu", "tw", "kr", "cn"}

		for _, raid := range raids {
			for _, difficulty := range difficulties {
				for _, region := range regions {
					leaderboard, err := raiderioRaid.GetRaidLeaderboard(service, raid, difficulty, region, 100, 0)
					if err != nil {
						return fmt.Errorf("error updating cache for %s %s %s: %v", raid, difficulty, region, err)
					}
					cacheKey := fmt.Sprintf("raid_leaderboard_%s_%s_%s_%d_%d", raid, difficulty, region, 100, 0)
					cache.Set(cacheKey, leaderboard, 1*time.Hour)
				}
			}
		}
		return nil
	}

	cache.StartPeriodicUpdate("raid_leaderboard_nerubar_palace_normal_world", updateFunc, 1*time.Hour)
}
