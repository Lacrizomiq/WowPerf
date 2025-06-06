package raiderioMythicPlus

import (
	"fmt"
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

	var wg sync.WaitGroup
	errors := make(chan error, 2)

	for page := 0; page <= 2; page++ {
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
				processRun(ranking, stats)
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

// processRun processes a run and updates the stats
// processRun processes a single run and updates the stats
func processRun(ranking interface{}, stats *DungeonStats) {
	run, ok := ranking.(map[string]interface{})
	if !ok {
		return
	}

	runDetails, ok := run["run"].(map[string]interface{})
	if !ok {
		return
	}

	// Process mythic level
	if mythicLevel, ok := runDetails["mythic_level"].(float64); ok {
		stats.LevelStats[int(mythicLevel)]++
	}

	// Process roster
	roster, ok := runDetails["roster"].([]interface{})
	if !ok {
		return
	}

	var comp models.TeamComposition
	var dpsMembers []models.TeamMember

	for _, member := range roster {
		memberMap, ok := member.(map[string]interface{})
		if !ok {
			continue
		}

		role, class, spec := extractMemberInfo(memberMap)
		if role == "" || class == "" || spec == "" {
			continue
		}

		member := models.TeamMember{
			Class: class,
			Spec:  spec,
		}

		// Update role and spec stats
		updateStats(stats, role, class, spec)

		// Build team composition
		switch role {
		case "tank":
			comp.Tank = member
		case "healer":
			comp.Healer = member
		case "dps":
			dpsMembers = append(dpsMembers, member)
		}
	}

	// Sort and assign DPS
	processDPSMembers(&comp, dpsMembers)

	// Create and store team composition
	key := createTeamCompKey(comp)
	if teamComp, exists := stats.TeamComp[key]; exists {
		teamComp.Count++
		stats.TeamComp[key] = teamComp
	} else {
		stats.TeamComp[key] = models.TeamCompStats{
			Count:       1,
			Composition: comp,
		}
	}
}

// extractMemberInfo extracts the role, class, and spec from a member map
func extractMemberInfo(memberMap map[string]interface{}) (role, class, spec string) {
	role, _ = memberMap["role"].(string)

	character, ok := memberMap["character"].(map[string]interface{})
	if !ok {
		return "", "", ""
	}

	classInfo, ok := character["class"].(map[string]interface{})
	if !ok {
		return role, "", ""
	}
	class, _ = classInfo["name"].(string)

	specInfo, ok := character["spec"].(map[string]interface{})
	if !ok {
		return role, class, ""
	}
	spec, _ = specInfo["name"].(string)

	return role, class, spec
}

// updateStats updates the role and spec stats
func updateStats(stats *DungeonStats, role, class, spec string) {
	if stats.RoleStats[role] == nil {
		stats.RoleStats[role] = make(map[string]int)
	}
	stats.RoleStats[role][class]++

	if stats.SpecStats[class] == nil {
		stats.SpecStats[class] = make(map[string]int)
	}
	stats.SpecStats[class][spec]++
}

// processDPSMembers processes the dps members and updates the team composition
func processDPSMembers(comp *models.TeamComposition, dpsMembers []models.TeamMember) {
	sort.Slice(dpsMembers, func(i, j int) bool {
		if dpsMembers[i].Class == dpsMembers[j].Class {
			return dpsMembers[i].Spec < dpsMembers[j].Spec
		}
		return dpsMembers[i].Class < dpsMembers[j].Class
	})

	if len(dpsMembers) >= 1 {
		comp.Dps1 = dpsMembers[0]
	}
	if len(dpsMembers) >= 2 {
		comp.Dps2 = dpsMembers[1]
	}
	if len(dpsMembers) >= 3 {
		comp.Dps3 = dpsMembers[2]
	}
}

// createTeamCompKey creates a key for the team composition
func createTeamCompKey(comp models.TeamComposition) string {
	return fmt.Sprintf("%s_%s_%s_%s_%s_%s_%s_%s_%s_%s",
		comp.Tank.Class, comp.Tank.Spec,
		comp.Healer.Class, comp.Healer.Spec,
		comp.Dps1.Class, comp.Dps1.Spec,
		comp.Dps2.Class, comp.Dps2.Spec,
		comp.Dps3.Class, comp.Dps3.Spec)
}
