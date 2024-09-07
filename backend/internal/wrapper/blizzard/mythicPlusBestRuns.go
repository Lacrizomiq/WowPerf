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

	characterName := character["name"].(string)
	realm, ok := character["realm"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("realm information not found")
	}
	realmSlug := realm["slug"].(string)

	// Extract mythic rating
	mythicRating, ok := data["mythic_rating"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("mythic rating information not found")
	}
	overallRating := mythicRating["rating"].(float64)
	color, ok := mythicRating["color"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("color information not found")
	}
	r := int(color["r"].(float64))
	g := int(color["g"].(float64))
	b := int(color["b"].(float64))
	colorHex := fmt.Sprintf("#%02X%02X%02X", r, g, b)

	// Extract season information
	season, ok := data["season"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("season information not found")
	}
	seasonID := uint(season["id"].(float64))

	var wg sync.WaitGroup
	runChan := make(chan mythicplus.MythicPlusRun, len(bestRuns))
	errChan := make(chan error, len(bestRuns))

	for i, run := range bestRuns {
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
		SeasonID:               seasonID,
		OverallMythicRating:    overallRating,
		OverallMythicRatingHex: colorHex,
		BestRuns:               results,
	}

	log.Printf("Successfully transformed %d runs", len(results))
	return seasonInfo, nil
}

// processMythicPlusRun transforms a single Mythic+ run from the Blizzard API into a struct.
func processMythicPlusRun(runData interface{}, db *gorm.DB, seasonSlug string) (mythicplus.MythicPlusRun, error) {
	log.Println("Starting processMythicPlusRun")

	runMap, ok := runData.(map[string]interface{})
	if !ok {
		log.Println("Error: invalid run data structure")
		return mythicplus.MythicPlusRun{}, fmt.Errorf("invalid run data structure")
	}

	dungeonMap := runMap["dungeon"].(map[string]interface{})
	challengeModeID := uint(dungeonMap["id"].(float64))
	log.Printf("Processing dungeon with ChallengeModeID: %d", challengeModeID)

	blizzardSeasonID, exists := SeasonIDMapping[seasonSlug]
	if !exists {
		return mythicplus.MythicPlusRun{}, fmt.Errorf("unknown season slug: %s", seasonSlug)
	}

	var dungeon mythicplus.Dungeon
	err := db.Preload("Seasons", "slug = ?", seasonSlug).Where("challenge_mode_id = ?", challengeModeID).First(&dungeon).Error
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
	log.Printf("KeyStoneUpgrades for dungeon %d: %+v", dungeon.ID, dungeon.KeyStoneUpgrades)

	completedTimestamp := time.Unix(int64(runMap["completed_timestamp"].(float64))/1000, 0)
	duration := int64(runMap["duration"].(float64))

	keystoneLevel := int(runMap["keystone_level"].(float64))

	var keyStoneUpgrades []mythicplus.KeyStoneUpgrade
	err = db.Where("challenge_mode_id = ?", challengeModeID).Order("qualifying_duration DESC").Find(&keyStoneUpgrades).Error
	if err != nil {
		log.Printf("Error fetching KeyStoneUpgrades: %v", err)
	}

	log.Printf("KeyStoneUpgrades for ChallengeModeID %d: %+v", challengeModeID, keyStoneUpgrades)

	var KeystoneUpgrades int
	if len(keyStoneUpgrades) > 0 {
		for _, upgrade := range keyStoneUpgrades {
			if duration <= int64(upgrade.QualifyingDuration) {
				KeystoneUpgrades = upgrade.UpgradeLevel
			} else {
				break
			}
		}
		log.Printf("Calculated KeystoneUpgrades: %d for duration %d", KeystoneUpgrades, duration)
	} else {
		log.Printf("Warning: No keystone upgrades found for ChallengeModeID %d", challengeModeID)
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
		KeyStoneUpgrades:      KeystoneUpgrades,
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
