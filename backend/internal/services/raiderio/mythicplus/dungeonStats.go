package raiderioMythicPlus

import (
	"fmt"
	"sync"
	"wowperf/internal/services/raiderio"
)

type DungeonStats struct {
	DungeonName string
	RoleStats   map[string]map[string]int
}

// GetDungeonStats retrieves the dungeon statistics for a given dungeon slug, specially the class composition of the runs
func GetDungeonStats(s *raiderio.RaiderIOService, season, region, dungeonSlug string) (*DungeonStats, error) {
	stats := &DungeonStats{
		DungeonName: dungeonSlug,
		RoleStats:   make(map[string]map[string]int),
	}

	var wg sync.WaitGroup
	errors := make(chan error, 2)

	for page := 0; page <= 1; page++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			runs, err := GetMythicPlusBestRuns(s, season, region, dungeonSlug, page)
			if err != nil {
				errors <- err
				return
			}

			rankings, ok := runs["rankings"].([]interface{})
			if !ok {
				errors <- fmt.Errorf("unexpected format for rankings")
				return
			}

			for _, ranking := range rankings {
				run, ok := ranking.(map[string]interface{})
				if !ok {
					continue
				}

				roster, ok := run["run"].(map[string]interface{})["roster"].([]interface{})
				if !ok {
					continue
				}

				for _, member := range roster {
					memberMap, ok := member.(map[string]interface{})
					if !ok {
						continue
					}

					role, ok := memberMap["role"].(string)
					if !ok {
						continue
					}

					character, ok := memberMap["character"].(map[string]interface{})
					if !ok {
						continue
					}

					class, ok := character["class"].(map[string]interface{})["name"].(string)
					if !ok {
						continue
					}

					if stats.RoleStats[role] == nil {
						stats.RoleStats[role] = make(map[string]int)
					}
					stats.RoleStats[role][class]++
				}
			}
		}(page)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return stats, nil
}

func GetAllDungeonStats(s *raiderio.RaiderIOService, season, region string, dungeonSlugs []string) ([]*DungeonStats, error) {
	var wg sync.WaitGroup
	stats := make([]*DungeonStats, len(dungeonSlugs))
	errors := make([]error, len(dungeonSlugs))

	for i, slug := range dungeonSlugs {
		wg.Add(1)
		go func(i int, slug string) {
			defer wg.Done()
			stat, err := GetDungeonStats(s, season, region, slug)
			stats[i] = stat
			errors[i] = err
		}(i, slug)
	}

	wg.Wait()

	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return stats, nil
}
