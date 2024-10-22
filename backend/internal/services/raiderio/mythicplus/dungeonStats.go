package raiderioMythicPlus

import (
	"fmt"
	"log"
	"sort"
	"sync"
	models "wowperf/internal/models/raiderio/mythicrundetails"
	"wowperf/internal/services/raiderio"
)

// DungeonStats is a struct that contains the statistics for a dungeon
type DungeonStats struct {
	DungeonName string
	RoleStats   map[string]map[string]int
	SpecStats   map[string]map[string]int
	LevelStats  map[int]int
	TeamComp    models.TeamCompMap
}

// GetDungeonStats retrieves the dungeon statistics for a given dungeon slug, specially the class composition of the runs
func GetDungeonStats(s *raiderio.RaiderIOService, season, region, dungeonSlug string) (*DungeonStats, error) {
	stats := &DungeonStats{
		DungeonName: dungeonSlug,
		RoleStats:   make(map[string]map[string]int),
		SpecStats:   make(map[string]map[string]int),
		LevelStats:  make(map[int]int),
		TeamComp:    make(models.TeamCompMap),
	}

	// Create a wait group and a channel for errors
	var wg sync.WaitGroup
	errors := make(chan error, 2)

	// Iterate over the pages of rankings
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

			// Iterate over the rankings
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

					// Create a composition for the team
					comp := models.TeamComposition{}
					dpsMembers := []models.TeamMember{}
					for _, member := range roster {
						memberMap := member.(map[string]interface{})
						role := memberMap["role"].(string)
						character := memberMap["character"].(map[string]interface{})
						class := character["class"].(map[string]interface{})["name"].(string)
						spec := character["spec"].(map[string]interface{})["name"].(string)

						member := models.TeamMember{
							Class: class,
							Spec:  spec,
						}

						// Add the member to the composition
						switch role {
						case "tank":
							comp.Tank = member
						case "healer":
							comp.Healer = member
						case "dps":
							dpsMembers = append(dpsMembers, member)
						}
					}

					// Sort the dps members by class and spec
					sort.Slice(dpsMembers, func(i, j int) bool {
						if dpsMembers[i].Class == dpsMembers[j].Class {
							return dpsMembers[i].Spec < dpsMembers[j].Spec
						}
						return dpsMembers[i].Class < dpsMembers[j].Class
					})

					// Add the dps members to the composition
					if len(dpsMembers) >= 1 {
						comp.Dps1 = dpsMembers[0]
					}
					if len(dpsMembers) >= 2 {
						comp.Dps2 = dpsMembers[1]
					}
					if len(dpsMembers) >= 3 {
						comp.Dps3 = dpsMembers[2]
					}

					// Create a key for the team composition
					key := fmt.Sprintf("%s_%s_%s_%s_%s_%s_%s_%s_%s_%s",
						comp.Tank.Class, comp.Tank.Spec,
						comp.Healer.Class, comp.Healer.Spec,
						comp.Dps1.Class, comp.Dps1.Spec,
						comp.Dps2.Class, comp.Dps2.Spec,
						comp.Dps3.Class, comp.Dps3.Spec)
					if teamComp, exists := stats.TeamComp[key]; exists {
						teamComp.Count++
						stats.TeamComp[key] = teamComp
					} else {
						stats.TeamComp[key] = models.TeamCompStats{Count: 1, Composition: comp}
					}

					// Add the role to the role stats
					if stats.RoleStats[role] == nil {
						stats.RoleStats[role] = make(map[string]int)
					}
					stats.RoleStats[role][class]++

					// Add the spec to the spec stats
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

// GetAllDungeonStats retrieves the dungeon statistics for a given season and region for a list of dungeon slugs
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
