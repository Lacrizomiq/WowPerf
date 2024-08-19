package seeds

import (
	"fmt"
	"log"
	"time"
	models "wowperf/internal/models/mythicdj"
	"wowperf/internal/services/blizzard"
	"wowperf/internal/services/blizzard/gamedata"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SeedMythicDjs seeds mythic dungeons and their related data
func SeedMythicDjs(db *gorm.DB, gameDataService *blizzard.GameDataService) error {
	log.Println("Seeding mythic dungeons")

	// Get all mythic dungeons from the Blizzard API
	dungeonIndex, err := gamedata.GetMythicKeystoneDungeonsIndex(gameDataService, "eu", "dynamic-eu", "en_GB")
	if err != nil {
		log.Fatalf("Failed to get mythic dungeons index: %v", err)
	}

	dungeons := dungeonIndex["dungeons"].([]interface{})

	for _, dungeon := range dungeons {
		dungeonData := dungeon.(map[string]interface{})
		id := int(dungeonData["id"].(float64))

		// Get the dungeon details from the Blizzard API
		dungeonDetails, err := gamedata.GetMythicKeystoneByID(gameDataService, id, "eu", "dynamic-eu", "en_GB")
		if err != nil {
			log.Fatalf("Failed to get mythic dungeon details: %v", err)
		}

		mythicDj := models.MythicDj{
			ID:   uint(id),
			Name: dungeonDetails["name"].(string),
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "icon_url"}),
		}).Create(&mythicDj).Error; err != nil {
			return err
		}

		for _, upgrade := range dungeonDetails["keystone_upgrades"].([]interface{}) {
			upgradeData := upgrade.(map[string]interface{})
			keystoneUpgrade := models.KeystoneUpgrade{
				MythicDjID:         mythicDj.ID,
				QualifyingDuration: int(upgradeData["qualifying_duration"].(float64)),
				UpgradeLevel:       int(upgradeData["upgrade_level"].(float64)),
			}
			if err := db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "mythic_dj_id"}, {Name: "upgrade_level"}},
				DoUpdates: clause.AssignmentColumns([]string{"qualifying_duration"}),
			}).Create(&keystoneUpgrade).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// SeedAffixes seeds mythic dungeon affixes and their related data
func SeedAffixes(db *gorm.DB, gameDataService *blizzard.GameDataService) error {
	log.Println("Seeding mythic dungeon affixes")

	affixIndex, err := gamedata.GetMythicKeystoneAffixIndex(gameDataService, "eu", "en_GB")
	if err != nil {
		log.Printf("Error getting mythic dungeon affix index: %v", err)
		return err
	}

	affixes, ok := affixIndex["affixes"].([]interface{})
	if !ok {
		log.Println("No affixes found or unexpected data structure")
		return nil
	}

	for _, affix := range affixes {
		affixData := affix.(map[string]interface{})
		id := int(affixData["id"].(float64))

		// Get the affix details from the Blizzard API
		affixDetails, err := gamedata.GetMythicKeystoneAffixByID(gameDataService, id, "eu", "static-eu", "en_GB")
		if err != nil {
			log.Printf("Failed to get mythic affix details for ID %d: %v", id, err)
			continue
		}

		affixMedia, err := gamedata.GetMythicKeystoneAffixMedia(gameDataService, id, "eu", "static-eu", "en_GB")
		if err != nil {
			log.Printf("Failed to get mythic affix media for ID %d: %v", id, err)
			continue
		}

		iconURL := ""
		if assets, ok := affixMedia["assets"].([]interface{}); ok && len(assets) > 0 {
			if asset, ok := assets[0].(map[string]interface{}); ok {
				if value, ok := asset["value"].(string); ok {
					iconURL = value
				}
			}
		}

		affix := models.Affix{
			AffixID: uint(id),
			Name:    affixDetails["name"].(string),
			IconURL: iconURL,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "affix_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "icon_url"}),
		}).Create(&affix).Error; err != nil {
			log.Printf("Error creating/updating affix with ID %d: %v", id, err)
			continue
		}
	}

	return nil
}

// SeedSeasons seeds mythic dungeon seasons and their related data
func SeedSeasons(db *gorm.DB, gameDataService *blizzard.GameDataService) error {
	log.Println("Seeding mythic dungeon seasons")

	seasonIndex, err := gamedata.GetMythicKeystoneSeasonsIndex(gameDataService, "eu", "dynamic-eu", "en_GB")
	if err != nil {
		log.Printf("Failed to get mythic dungeon seasons index: %v", err)
		return err
	}

	seasons, ok := seasonIndex["seasons"].([]interface{})
	if !ok {
		log.Println("No seasons found or unexpected data structure")
		return nil
	}

	for _, season := range seasons {
		seasonData, ok := season.(map[string]interface{})
		if !ok {
			log.Println("Unexpected season data structure")
			continue
		}

		id, ok := seasonData["id"].(float64)
		if !ok {
			log.Println("Unexpected season ID type")
			continue
		}

		// Get the season details from the Blizzard API
		seasonDetails, err := gamedata.GetMythicKeystoneSeasonByID(gameDataService, int(id), "eu", "dynamic-eu", "en_GB")
		if err != nil {
			log.Printf("Failed to get mythic season details for ID %d: %v", int(id), err)
			continue
		}

		seasonName, ok := seasonDetails["season_name"].(string)
		if !ok {
			seasonName = fmt.Sprintf("Season %d", int(id))
		}

		season := models.Season{
			SeasonID: uint(id),
			Name:     seasonName,
		}

		if err := db.Create(&season).Error; err != nil {
			log.Printf("Error creating season with ID %d: %v", int(id), err)
			continue
		}

		if err := SeedsPeriodsForSeason(db, gameDataService, int(id)); err != nil {
			log.Printf("Error seeding periods for season %d: %v", int(id), err)
			continue
		}
	}
	return nil
}

// SeedsPeriodsForSeason seeds mythic dungeon season periods and their related data
func SeedsPeriodsForSeason(db *gorm.DB, gameDataService *blizzard.GameDataService, seasonID int) error {
	log.Println("Seeding mythic dungeon season periods")

	// Get the periods index
	periodsIndex, err := gamedata.GetMythicKeystonePeriodsIndex(gameDataService, "eu", "dynamic-eu", "en_GB")
	if err != nil {
		return fmt.Errorf("failed to get periods index: %v", err)
	}

	periods, ok := periodsIndex["periods"].([]interface{})
	if !ok {
		return fmt.Errorf("unexpected periods data structure in index")
	}

	for _, periodData := range periods {
		period, ok := periodData.(map[string]interface{})
		if !ok {
			log.Println("Unexpected period data structure in index")
			continue
		}

		periodID, ok := period["id"].(float64)
		if !ok {
			log.Println("Unexpected period ID type in index")
			continue
		}

		// Get period details
		periodDetails, err := gamedata.GetMythicKeystonePeriodByID(gameDataService, int(periodID), "eu", "dynamic-eu", "en_GB")
		if err != nil {
			log.Printf("Failed to get period details for ID %d: %v", int(periodID), err)
			continue
		}

		startTime, ok := periodDetails["start_timestamp"].(float64)
		if !ok {
			log.Printf("Unexpected start_timestamp type for period %d", int(periodID))
			continue
		}

		endTime, ok := periodDetails["end_timestamp"].(float64)
		if !ok {
			log.Printf("Unexpected end_timestamp type for period %d", int(periodID))
			continue
		}

		// Check if this period belongs to the current season
		if periodSeasonID, ok := periodDetails["season"].(map[string]interface{})["id"].(float64); !ok || int(periodSeasonID) != seasonID {
			continue // Skip this period if it doesn't belong to the current season
		}

		newPeriod := models.Period{
			PeriodID:  uint(periodID),
			SeasonID:  uint(seasonID),
			StartTime: time.Unix(int64(startTime/1000), 0).Format(time.RFC3339),
			EndTime:   time.Unix(int64(endTime/1000), 0).Format(time.RFC3339),
		}

		if err := db.Create(&newPeriod).Error; err != nil {
			log.Printf("Error creating period with ID %d: %v", int(periodID), err)
			continue
		}
	}

	return nil
}

func SeedAll(db *gorm.DB) error {
	blizzardService, err := blizzard.NewService()
	if err != nil {
		return err
	}

	if err := SeedMythicDjs(db, blizzardService.GameData); err != nil {
		return err
	}

	if err := SeedAffixes(db, blizzardService.GameData); err != nil {
		return err
	}

	if err := SeedSeasons(db, blizzardService.GameData); err != nil {
		return err
	}

	return nil
}
