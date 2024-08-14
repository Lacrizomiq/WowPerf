package wrapper

import (
	"fmt"
	"strings"
	"time"
	"wowperf/internal/models"
)

func TransformMythicPlusBestRuns(data map[string]interface{}) ([]models.MythicPlusRun, error) {
	var bestRuns []models.MythicPlusRun

	rawBestRuns, ok := data["best_runs"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("best runs not found")
	}

	for _, run := range rawBestRuns {
		runMap, ok := run.(map[string]interface{})
		if !ok {
			continue
		}

		dungeon, _ := runMap["dungeon"].(map[string]interface{})
		dungeonName, _ := dungeon["name"].(string)

		completedTimestamp, _ := runMap["completed_timestamp"].(float64)
		completedAt := time.Unix(int64(completedTimestamp), 0).Format("2006-01-02 15:04:05")

		mythicRating, _ := runMap["mythic_rating"].(map[string]interface{})
		score, _ := mythicRating["rating"].(float64)

		bestRun := models.MythicPlusRun{
			Dungeon:             dungeonName,
			ShortName:           getShortName(dungeonName),
			MythicLevel:         int(runMap["keystone_level"].(float64)),
			CompletedAt:         completedAt,
			ClearTimeMS:         int(runMap["duration"].(float64)),
			NumKeystoneUpgrades: getKeystoneUpgrades(runMap),
			Score:               score,
			URL:                 "",
			Affixes:             getAffixes(runMap),
		}

		bestRuns = append(bestRuns, bestRun)
	}

	return bestRuns, nil
}

// getShortName returns a short name for a dungeon name
func getShortName(dungeonName string) string {
	return strings.ReplaceAll(dungeonName, " ", "_")
}

// getKeystoneUpgrades returns the number of keystone upgrades for a run
func getKeystoneUpgrades(runMap map[string]interface{}) int {
	keystoneUpgrades, _ := runMap["keystone_upgrades"].([]interface{})
	return len(keystoneUpgrades)
}

// getAffixes returns the affixes for a run
func getAffixes(runMap map[string]interface{}) []models.Affix {
	affixes := []models.Affix{}

	for _, affix := range runMap["affixes"].([]interface{}) {
		affixMap, ok := affix.(map[string]interface{})
		if !ok {
			continue
		}

		affixID, _ := affixMap["affix_id"].(float64)
		affixName, _ := affixMap["affix_name"].(string)

		affixes = append(affixes, models.Affix{
			ID:   int(affixID),
			Name: affixName,
		})
	}

	return affixes
}
