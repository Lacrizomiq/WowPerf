package wrapper

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
	mythicplus "wowperf/internal/models/mythicplus"

	"gorm.io/gorm"
)

var SeasonIDMapping = map[string]int{
	"season-tww-1": 13,
	"season-df-4":  12,
	"season-df-3":  11,
	"season-df-2":  10,
	"season-df-1":  9,
}

var SeasonSlugMapping = make(map[int]string)

func init() {
	for slug, id := range SeasonIDMapping {
		SeasonSlugMapping[id] = slug
	}
}

// TransformMythicPlusBestRuns transforms the best runs from the Mythic+ API into a struct that is easier to use than the Blizzard API response
func TransformMythicPlusBestRuns(data map[string]interface{}, db *gorm.DB, seasonSlug string) (*mythicplus.MythicPlusSeasonInfo, error) {
	bestRuns, ok := data["best_runs"].([]interface{})
	if !ok {
		log.Println("Error: best runs not found or not a slice")
		return nil, fmt.Errorf("best runs not found or not a slice")
	}

	// Extract character information
	character, ok := data["character"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("character information not found")
	}

	realm := character["realm"].(map[string]interface{})
	characterName := character["name"].(string)
	realmSlug := realm["slug"].(string)

	// Get season ID
	seasonID, exists := SeasonIDMapping[seasonSlug]
	if !exists {
		return nil, fmt.Errorf("unknown season slug: %s", seasonSlug)
	}

	// Get overall mythic rating
	mythicRating, ok := data["mythic_rating"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("mythic rating not found")
	}

	overallRating := mythicRating["rating"].(float64)
	color := mythicRating["color"].(map[string]interface{})
	colorHex := fmt.Sprintf("#%02x%02x%02x", int(color["r"].(float64)), int(color["g"].(float64)), int(color["b"].(float64)))

	// Create a map to store the highest level run for each dungeon
	dungeonBestRuns := make(map[int]interface{})

	// First pass: organize runs by dungeon and keep only the highest level
	for _, run := range bestRuns {
		runMap, ok := run.(map[string]interface{})
		if !ok {
			continue
		}

		dungeon, ok := runMap["dungeon"].(map[string]interface{})
		if !ok {
			continue
		}

		dungeonID := int(dungeon["id"].(float64))
		keystoneLevel := int(runMap["keystone_level"].(float64))

		if existingRun, exists := dungeonBestRuns[dungeonID]; exists {
			existingRunMap := existingRun.(map[string]interface{})
			existingLevel := int(existingRunMap["keystone_level"].(float64))

			if keystoneLevel > existingLevel {
				dungeonBestRuns[dungeonID] = runMap
			}
		} else {
			dungeonBestRuns[dungeonID] = runMap
		}
	}

	// Convert map back to slice for processing
	filteredBestRuns := make([]interface{}, 0, len(dungeonBestRuns))
	for _, run := range dungeonBestRuns {
		filteredBestRuns = append(filteredBestRuns, run)
	}

	// Process the filtered runs concurrently
	var wg sync.WaitGroup
	runChan := make(chan mythicplus.MythicPlusRun, len(filteredBestRuns))
	errChan := make(chan error, len(filteredBestRuns))

	for i, run := range filteredBestRuns {
		wg.Add(1)
		go func(i int, run interface{}) {
			defer wg.Done()
			mythicRun, err := processMythicPlusRun(run, db, seasonSlug)
			if err != nil {
				log.Printf("Error processing run %d: %v", i, err)
				errChan <- err
				return
			}
			runChan <- mythicRun
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

	seasonInfo := &mythicplus.MythicPlusSeasonInfo{
		CharacterName:          characterName,
		RealmSlug:              realmSlug,
		SeasonID:               uint(seasonID),
		OverallMythicRating:    overallRating,
		OverallMythicRatingHex: colorHex,
		BestRuns:               results,
	}

	return seasonInfo, nil
}

// processMythicPlusRun transforms a single Mythic+ run from the Blizzard API into a struct.
func processMythicPlusRun(runData interface{}, db *gorm.DB, seasonSlug string) (mythicplus.MythicPlusRun, error) {

	runMap, ok := runData.(map[string]interface{})
	if !ok {
		log.Println("Error: invalid run data structure")
		return mythicplus.MythicPlusRun{}, fmt.Errorf("invalid run data structure")
	}

	dungeonMap := runMap["dungeon"].(map[string]interface{})
	challengeModeID := uint(dungeonMap["id"].(float64))

	blizzardSeasonID, exists := SeasonIDMapping[seasonSlug]
	if !exists {
		return mythicplus.MythicPlusRun{}, fmt.Errorf("unknown season slug: %s", seasonSlug)
	}

	var dungeon mythicplus.Dungeon
	err := db.Preload("Seasons", "slug = ?", seasonSlug).
		Preload("KeyStoneUpgrades").
		Where("challenge_mode_id = ?", challengeModeID).
		First(&dungeon).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Warning: Dungeon with ChallengeModeID %d not found in database", challengeModeID)
			dungeon = mythicplus.Dungeon{
				ChallengeModeID: challengeModeID,
				Name:            dungeonMap["name"].(string),
				ShortName:       dungeonMap["name"].(string),
			}
		} else {
			log.Printf("Error fetching dungeon: %v", err)
			return mythicplus.MythicPlusRun{}, fmt.Errorf("error fetching dungeon: %v", err)
		}
	}

	var keyStoneUpgrades []mythicplus.KeyStoneUpgrade
	err = db.Where("challenge_mode_id = ?", dungeon.ChallengeModeID).Find(&keyStoneUpgrades).Error
	if err != nil {
		log.Printf("Error loading KeyStoneUpgrades for ChallengeModeID %d: %v", dungeon.ChallengeModeID, err)
	} else {
		log.Printf("Loaded KeyStoneUpgrades for ChallengeModeID %d: %v", dungeon.ChallengeModeID, keyStoneUpgrades)
	}

	completedTimestamp := time.Unix(int64(runMap["completed_timestamp"].(float64))/1000, 0)
	duration := int64(runMap["duration"].(float64))

	keystoneLevel := int(runMap["keystone_level"].(float64))

	keystoneUpgrades := 0
	if len(dungeon.KeyStoneUpgrades) > 0 {
		// Sort KeyStoneUpgrades by UpgradeLevel descending
		sort.Slice(dungeon.KeyStoneUpgrades, func(i, j int) bool {
			return dungeon.KeyStoneUpgrades[i].UpgradeLevel > dungeon.KeyStoneUpgrades[j].UpgradeLevel
		})

		for _, upgrade := range dungeon.KeyStoneUpgrades {
			if duration <= upgrade.QualifyingDuration {
				keystoneUpgrades = upgrade.UpgradeLevel
				break // Found the highest possible upgrade level
			}
		}
		log.Printf("Calculated KeystoneUpgrades: %d for duration %d", keystoneUpgrades, duration)
	} else {
		log.Printf("Warning: No keystone upgrades found for ChallengeModeID %d", dungeon.ChallengeModeID)
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

	var season mythicplus.Season
	if len(dungeon.Seasons) > 0 {
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
		KeyStoneUpgrades:      keystoneUpgrades,
		KeystoneLevel:         keystoneLevel,
		MythicRating:          score,
		SeasonID:              uint(blizzardSeasonID),
		Season:                season,
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

	for i, affixData := range affixesData {
		wg.Add(1)
		go func(i int, affixData interface{}) {
			defer wg.Done()

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

	return affixes, nil
}

// getMembers returns a slice of members from a slice of member data
func getMembers(membersData []interface{}) ([]mythicplus.MythicPlusRunMember, error) {
	var members []mythicplus.MythicPlusRunMember
	for i, memberData := range membersData {

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
