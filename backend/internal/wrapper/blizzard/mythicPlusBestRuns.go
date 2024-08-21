package wrapper

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
	mythicplus "wowperf/internal/models/mythicplus"

	"gorm.io/gorm"
)

// TransformMythicPlusBestRuns transforms the best runs from the Mythic+ API into a struct that is easier to use than the Blizzard API response
func TransformMythicPlusBestRuns(data map[string]interface{}, db *gorm.DB) ([]mythicplus.MythicPlusRun, error) {
	bestRuns, ok := data["best_runs"].([]interface{})
	if !ok {
		log.Println("Error: best runs not found or not a slice")
		return nil, fmt.Errorf("best runs not found or not a slice")
	}

	var wg sync.WaitGroup
	runChan := make(chan mythicplus.MythicPlusRun, len(bestRuns))
	errChan := make(chan error, len(bestRuns))

	for i, run := range bestRuns {
		wg.Add(1)
		go func(i int, run interface{}) {
			defer wg.Done()
			log.Printf("Processing run %d", i)
			mythicRun, err := processMythicPlusRun(run, db)
			if err != nil {
				log.Printf("Error processing run %d: %v", i, err)
				errChan <- err
				return
			}
			runChan <- mythicRun
			log.Printf("Finished processing run %d", i)
		}(i, run)
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
		err := <-errChan
		log.Printf("Error in TransformMythicPlusBestRuns: %v", err)
		return nil, err
	}

	log.Printf("Successfully transformed %d runs", len(results))
	return results, nil
}

// processMythicPlusRun transforms a single Mythic+ run from the Blizzard API into a struct.
func processMythicPlusRun(runData interface{}, db *gorm.DB) (mythicplus.MythicPlusRun, error) {
	log.Println("Starting processMythicPlusRun")

	runMap, ok := runData.(map[string]interface{})
	if !ok {
		log.Println("Error: invalid run data structure")
		return mythicplus.MythicPlusRun{}, fmt.Errorf("invalid run data structure")
	}

	dungeonMap := runMap["dungeon"].(map[string]interface{})
	challengeModeID := uint(dungeonMap["id"].(float64))
	log.Printf("Processing dungeon with ChallengeModeID: %d", challengeModeID)

	var dungeon mythicplus.Dungeon
	err := db.Preload("KeyStoneUpgrades").Preload("Seasons").Where("challenge_mode_id = ?", challengeModeID).First(&dungeon).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Warning: Dungeon with ChallengeModeID %d not found in database", challengeModeID)
			dungeon = mythicplus.Dungeon{
				ChallengeModeID: &challengeModeID,
				Name:            dungeonMap["name"].(string),
				ShortName:       dungeonMap["name"].(string),
			}
		} else {
			log.Printf("Error fetching dungeon: %v", err)
			return mythicplus.MythicPlusRun{}, fmt.Errorf("error fetching dungeon: %v", err)
		}
	}

	log.Printf("Dungeon found: %+v", dungeon)

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
		log.Printf("Error getting affixes: %v", err)
		return mythicplus.MythicPlusRun{}, fmt.Errorf("error getting affixes: %v", err)
	}

	members, err := getMembers(runMap["members"].([]interface{}))
	if err != nil {
		log.Printf("Error getting members: %v", err)
		return mythicplus.MythicPlusRun{}, fmt.Errorf("error getting members: %v", err)
	}

	var seasonID uint
	var season mythicplus.Season
	if len(dungeon.Seasons) > 0 {
		seasonID = dungeon.Seasons[0].ID
		season = dungeon.Seasons[0]
	} else {
		log.Printf("Warning: No seasons associated with dungeon %d", dungeon.ID)
	}

	return mythicplus.MythicPlusRun{
		CompletedTimestamp:    completedTimestamp,
		DungeonID:             dungeon.ID,
		Dungeon:               dungeon,
		ShortName:             dungeon.ShortName,
		Duration:              duration,
		IsCompletedWithinTime: runMap["is_completed_within_time"].(bool),
		KeyStoneUpgrades:      KeystoneUpgrades,
		KeystoneLevel:         keystoneLevel,
		MythicRating:          score,
		SeasonID:              seasonID,
		Season:                season,
		Affixes:               affixes,
		Members:               members,
	}, nil
}

// getAffixes returns a slice of affixes from a slice of affix IDs
func getAffixes(affixesData []interface{}, db *gorm.DB) ([]mythicplus.Affix, error) {
	log.Printf("Starting getAffixes with %d affixes", len(affixesData))

	var affixes []mythicplus.Affix

	var wg sync.WaitGroup
	affixChan := make(chan mythicplus.Affix, len(affixesData))
	errorChan := make(chan error, len(affixesData))

	for i, affixData := range affixesData {
		wg.Add(1)
		go func(i int, affixData interface{}) {
			defer wg.Done()
			log.Printf("Processing affix %d: %+v", i, affixData)

			affixMap, ok := affixData.(map[string]interface{})
			if !ok {
				errorChan <- fmt.Errorf("affix data is not a map for affix %d", i)
				return
			}

			affixIDFloat, ok := affixMap["id"].(float64)
			if !ok {
				errorChan <- fmt.Errorf("invalid affix ID type for affix %d", i)
				return
			}

			affixID := uint(affixIDFloat)

			var affix mythicplus.Affix
			if err := db.First(&affix, affixID).Error; err != nil {
				log.Printf("Error fetching affix %d with ID %d: %v", i, affixID, err)
				errorChan <- fmt.Errorf("error fetching affix %d: %v", i, err)
				return
			}
			log.Printf("Successfully fetched affix %d: ID %d, Name: %s", i, affix.ID, affix.Name)
			affixChan <- affix
		}(i, affixData)
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
		err := <-errorChan
		log.Printf("Error in getAffixes: %v", err)
		return nil, err
	}

	log.Printf("Successfully processed %d affixes", len(affixes))
	return affixes, nil
}

// getMembers returns a slice of members from a slice of member data
func getMembers(membersData []interface{}) ([]mythicplus.MythicPlusRunMember, error) {
	var members []mythicplus.MythicPlusRunMember
	for i, memberData := range membersData {
		log.Printf("Processing member %d: %+v", i, memberData)

		memberMap, ok := memberData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid member data structure for member %d", i)
		}

		character, ok := memberMap["character"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid character data structure for member %d", i)
		}

		realm, ok := character["realm"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid realm data structure for member %d", i)
		}

		race, ok := memberMap["race"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid race data structure for member %d", i)
		}

		spec, ok := memberMap["specialization"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid specialization data structure for member %d", i)
		}

		member := mythicplus.MythicPlusRunMember{
			CharacterID:       uint(character["id"].(float64)),
			CharacterName:     getStringValue(character, "name"),
			RealmID:           uint(realm["id"].(float64)),
			RealmName:         getStringValue(realm, "name"),
			RealmSlug:         getStringValue(realm, "slug"),
			EquippedItemLevel: int(memberMap["equipped_item_level"].(float64)),
			RaceID:            uint(race["id"].(float64)),
			RaceName:          getStringValue(race, "name"),
			SpecializationID:  uint(spec["id"].(float64)),
			Specialization:    getStringValue(spec, "name"),
		}
		members = append(members, member)
	}
	return members, nil
}

func getStringValue(data map[string]interface{}, key string) string {
	if value, ok := data[key]; ok {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return ""
}
