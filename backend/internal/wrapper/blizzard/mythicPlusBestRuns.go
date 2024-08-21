package wrapper

import (
	"fmt"
	"sync"
	"time"
	mythicplus "wowperf/internal/models/mythicplus"

	"gorm.io/gorm"
)

func TransformMythicPlusBestRuns(data map[string]interface{}, db *gorm.DB) ([]mythicplus.MythicPlusRun, error) {
	bestRuns, ok := data["best_runs"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("best runs not found or not a slice")
	}

	var wg sync.WaitGroup
	runChan := make(chan mythicplus.MythicPlusRun, len(bestRuns))
	errChan := make(chan error, len(bestRuns))

	for _, run := range bestRuns {
		wg.Add(1)
		go func(run interface{}) {
			defer wg.Done()
			mythicRun, err := processMythicPlusRun(run, db)
			if err != nil {
				errChan <- err
				return
			}
			runChan <- mythicRun
		}(run)
	}

	go func() {
		wg.Wait()
		close(runChan)
		close(errChan)
	}()

	var results []mythicplus.MythicPlusRun
	for run := range runChan {
		results = append(results, run)
	}

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return results, nil
}

// processMythicPlusRun transforms a single Mythic+ run from the Blizzard API into a struct.
func processMythicPlusRun(runData interface{}, db *gorm.DB) (mythicplus.MythicPlusRun, error) {
	runMap, ok := runData.(map[string]interface{})
	if !ok {
		return mythicplus.MythicPlusRun{}, fmt.Errorf("invalid run data structure")
	}

	dungeonMap := runMap["dungeon"].(map[string]interface{})
	dungeonID := uint(dungeonMap["id"].(float64))

	var dungeon mythicplus.Dungeon
	if err := db.Preload("KeyStoneUpgrades").First(&dungeon, dungeonID).Error; err != nil {
		return mythicplus.MythicPlusRun{}, fmt.Errorf("error fetching dungeon: %v", err)
	}

	completedTimestamp := time.Unix(int64(runMap["completed_timestamp"].(float64))/1000, 0)
	duration := int64(runMap["duration"].(float64))

	keystoneLevel := int(runMap["keystone_level"].(float64))

	var KeystoneUpgrades int
	for _, upgrade := range dungeon.KeyStoneUpgrades {
		if duration <= int64(upgrade.QualifyingDuration) {
			KeystoneUpgrades = upgrade.UpgradeLevel
		} else {
			break
		}
	}

	mythicRating := runMap["mythic_rating"].(map[string]interface{})
	score := mythicRating["rating"].(float64)

	affixes, err := getAffixes(runMap["keystone_affixes"].([]interface{}), db)
	if err != nil {
		return mythicplus.MythicPlusRun{}, fmt.Errorf("error getting affixes: %v", err)
	}

	members, err := getMembers(runMap["members"].([]interface{}))
	if err != nil {
		return mythicplus.MythicPlusRun{}, fmt.Errorf("error getting members: %v", err)
	}

	return mythicplus.MythicPlusRun{
		CompletedTimestamp:    completedTimestamp,
		DungeonID:             dungeonID,
		Dungeon:               dungeon,
		ShortName:             dungeon.ShortName,
		Duration:              duration,
		IsCompletedWithinTime: runMap["is_completed_within_time"].(bool),
		KeyStoneUpgrades:      KeystoneUpgrades,
		KeystoneLevel:         keystoneLevel,
		MythicRating:          score,
		SeasonID:              dungeon.Seasons[0].ID,
		Season:                dungeon.Seasons[0],
		Affixes:               affixes,
		Members:               members,
	}, nil
}

// getAffixes returns a slice of affixes from a slice of affix IDs
func getAffixes(affixesData []interface{}, db *gorm.DB) ([]mythicplus.Affix, error) {
	var affixes []mythicplus.Affix

	var wg sync.WaitGroup
	affixChan := make(chan mythicplus.Affix, len(affixesData))
	errorChan := make(chan error, len(affixesData))

	for _, affixID := range affixesData {
		wg.Add(1)
		go func(affixID interface{}) {
			defer wg.Done()
			affixIDInt, ok := affixID.(float64)
			if !ok {
				errorChan <- fmt.Errorf("invalid affix ID type")
				return
			}

			var affix mythicplus.Affix
			if err := db.First(&affix, uint(affixIDInt)).Error; err != nil {
				errorChan <- fmt.Errorf("error fetching affix: %v", err)
				return
			}
			affixChan <- affix
		}(affixID)
	}

	go func() {
		wg.Wait()
		close(affixChan)
		close(errorChan)
	}()

	for affix := range affixChan {
		affixes = append(affixes, affix)
	}

	if len(errorChan) > 0 {
		return nil, <-errorChan
	}

	return affixes, nil
}

// getMembers returns a slice of members from a slice of member data
func getMembers(membersData []interface{}) ([]mythicplus.MythicPlusRunMember, error) {
	var members []mythicplus.MythicPlusRunMember
	for _, memberData := range membersData {
		memberMap := memberData.(map[string]interface{})
		character := memberMap["character"].(map[string]interface{})
		realm := character["realm"].(map[string]interface{})
		spec := memberMap["specialization"].(map[string]interface{})
		race := memberMap["race"].(map[string]interface{})

		member := mythicplus.MythicPlusRunMember{
			CharacterID:       uint(character["id"].(float64)),
			CharacterName:     character["name"].(string),
			RealmID:           uint(realm["id"].(float64)),
			RealmName:         realm["name"].(string),
			RealmSlug:         realm["slug"].(string),
			EquippedItemLevel: int(memberMap["equipped_item_level"].(float64)),
			RaceID:            uint(race["id"].(float64)),
			RaceName:          race["name"].(string),
			SpecializationID:  uint(spec["id"].(float64)),
			Specialization:    spec["name"].(string),
		}
		members = append(members, member)
	}
	return members, nil
}
