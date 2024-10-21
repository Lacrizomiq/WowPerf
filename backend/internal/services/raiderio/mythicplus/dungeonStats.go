package raiderioMythicPlus

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"wowperf/internal/services/raiderio"
)

type DungeonStats struct {
	DungeonName string
	RoleStats   map[string]map[string]int
	SpecStats   map[string]map[string]int
	LevelStats  map[int]int
	TeamComp    map[string]int
}

// GetDungeonStats retrieves the dungeon statistics for a given dungeon slug, specially the class composition of the runs
func GetDungeonStats(s *raiderio.RaiderIOService, season, region, dungeonSlug string) (*DungeonStats, error) {
	stats := &DungeonStats{
		DungeonName: dungeonSlug,
		RoleStats:   make(map[string]map[string]int),
		SpecStats:   make(map[string]map[string]int),
		LevelStats:  make(map[int]int),
		TeamComp:    make(map[string]int),
	}

	var wg sync.WaitGroup
	errors := make(chan error, 2)

	for page := 0; page <= 2; page++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()

			log.Printf("Requesting page %d for %s %s", page, season, region)

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

				runDetails, ok := run["run"].(map[string]interface{})
				if !ok {
					continue
				}

				mythicLevel, ok := runDetails["mythic_level"].(float64)
				if ok {
					stats.LevelStats[int(mythicLevel)]++
				}

				roster, ok := runDetails["roster"].([]interface{})
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

					spec, ok := character["spec"].(map[string]interface{})
					if !ok {
						continue
					}

					specName, ok := spec["name"].(string)
					if !ok {
						continue
					}

					var teamComp strings.Builder
					for _, member := range roster {
						memberMap, ok := member.(map[string]interface{})
						if !ok {
							continue
						}

						role, _ := memberMap["role"].(string)
						character, _ := memberMap["character"].(map[string]interface{})
						class, _ := character["class"].(map[string]interface{})["name"].(string)
						spec, _ := character["spec"].(map[string]interface{})["name"].(string)
						teamComp.WriteString(fmt.Sprintf("%s-%s-%s, ", role, class, spec))
					}
					compString := strings.TrimRight(teamComp.String(), ", ")
					stats.TeamComp[compString]++

					if stats.RoleStats[role] == nil {
						stats.RoleStats[role] = make(map[string]int)
					}
					stats.RoleStats[role][class]++

					if stats.SpecStats[class] == nil {
						stats.SpecStats[class] = make(map[string]int)
					}
					stats.SpecStats[class][specName]++
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
