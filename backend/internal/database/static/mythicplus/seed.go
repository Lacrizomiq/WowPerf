package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	mythicplus "wowperf/internal/models/mythicplus"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeasonData struct {
	Seasons []struct {
		Slug      string `json:"slug"`
		Name      string `json:"name"`
		ShortName string `json:"short_name"`

		SeasonalAffix *struct{} `json:"seasonal_affix"`
		Starts        struct {
			US string `json:"us"`
			EU string `json:"eu"`
			TW string `json:"tw"`
			KR string `json:"kr"`
			CN string `json:"cn"`
		} `json:"starts"`
		Ends struct {
			US string `json:"us"`
			EU string `json:"eu"`
			TW string `json:"tw"`
			KR string `json:"kr"`
			CN string `json:"cn"`
		} `json:"ends"`
		Dungeons []struct {
			ID               uint   `json:"id"`
			ChallengeModeID  uint   `json:"challenge_mode_id"`
			EncounterID      uint   `json:"encounter_id"`
			Slug             string `json:"slug"`
			Name             string `json:"name"`
			ShortName        string `json:"short_name"`
			MediaURL         string `json:"mediaURL"`
			Icon             string `json:"icon"`
			KeystoneUpgrades []struct {
				QualifyingDuration int `json:"qualifying_duration"`
				UpgradeLevel       int `json:"upgrade_level"`
			} `json:"keystone_upgrades"`
		} `json:"dungeons"`
	} `json:"seasons"`
}

type AffixData struct {
	AffixDetails []struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Icon        string `json:"icon"`
		Description string `json:"description"`
		WowheadURL  string `json:"wowhead_url"`
	} `json:"affix_details"`
}

// SeedSeasons seeds the Mythic+ seasons from a JSON file
func SeedSeasons(db *gorm.DB, filePath string) error {
	var seasonData SeasonData
	seasonFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	if err := json.Unmarshal(seasonFile, &seasonData); err != nil {
		return fmt.Errorf("error unmarshaling file %s: %v", filePath, err)
	}

	// First, create all dungeons
	for _, s := range seasonData.Seasons {
		for _, d := range s.Dungeons {
			dungeon := mythicplus.Dungeon{
				ID:              d.ID,
				ChallengeModeID: d.ChallengeModeID,
				Slug:            d.Slug,
				Name:            d.Name,
				ShortName:       d.ShortName,
				MediaURL:        d.MediaURL,
				Icon:            &d.Icon,
				// EncounterID is a pointer, will be nil by default
			}

			// Only set EncounterID if it's greater than 0
			if d.EncounterID > 0 {
				encounterID := d.EncounterID
				dungeon.EncounterID = &encounterID
			}

			// First try to find existing dungeon
			var existingDungeon mythicplus.Dungeon
			result := db.First(&existingDungeon, d.ID)

			if result.Error != nil {
				// Dungeon doesn't exist, create it
				if err := db.Create(&dungeon).Error; err != nil {
					return fmt.Errorf("error creating dungeon %s: %v", d.Name, err)
				}
			} else {
				// Dungeon exists, update it
				existingDungeon.ChallengeModeID = d.ChallengeModeID
				existingDungeon.Slug = d.Slug
				existingDungeon.Name = d.Name
				existingDungeon.ShortName = d.ShortName
				existingDungeon.MediaURL = d.MediaURL
				existingDungeon.Icon = &d.Icon

				// Only update EncounterID if the new value is greater than 0
				if d.EncounterID > 0 {
					encounterID := d.EncounterID
					existingDungeon.EncounterID = &encounterID
				}

				if err := db.Save(&existingDungeon).Error; err != nil {
					return fmt.Errorf("error updating dungeon %s: %v", d.Name, err)
				}
			}

			// Handle KeyStoneUpgrades
			if err := db.Where("challenge_mode_id = ?", d.ChallengeModeID).Delete(&mythicplus.KeyStoneUpgrade{}).Error; err != nil {
				return fmt.Errorf("error deleting old keystone upgrades for dungeon %s: %v", d.Name, err)
			}

			for _, ku := range d.KeystoneUpgrades {
				keyStoneUpgrade := mythicplus.KeyStoneUpgrade{
					ChallengeModeID:    d.ChallengeModeID,
					QualifyingDuration: int64(ku.QualifyingDuration),
					UpgradeLevel:       ku.UpgradeLevel,
				}

				if err := db.Create(&keyStoneUpgrade).Error; err != nil {
					return fmt.Errorf("error creating keystone upgrade for dungeon %s: %v", d.Name, err)
				}
			}
		}
	}

	// Then create seasons and associations
	for _, s := range seasonData.Seasons {
		season := mythicplus.Season{
			Slug:      s.Slug,
			Name:      s.Name,
			ShortName: s.ShortName,
			StartsUS:  parseTime(s.Starts.US),
			StartsEU:  parseTime(s.Starts.EU),
			StartsTW:  parseTime(s.Starts.TW),
			StartsKR:  parseTime(s.Starts.KR),
			StartsCN:  parseTime(s.Starts.CN),
			EndsUS:    parseTime(s.Ends.US),
			EndsEU:    parseTime(s.Ends.EU),
			EndsTW:    parseTime(s.Ends.TW),
			EndsKR:    parseTime(s.Ends.KR),
			EndsCN:    parseTime(s.Ends.CN),
		}

		// Create or update season
		if err := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&season).Error; err != nil {
			return fmt.Errorf("error creating/updating season %s: %v", s.Name, err)
		}

		// Create associations
		for _, d := range s.Dungeons {
			var dungeon mythicplus.Dungeon
			if err := db.First(&dungeon, d.ID).Error; err != nil {
				return fmt.Errorf("error finding dungeon %d: %v", d.ID, err)
			}

			if err := db.Model(&season).Association("Dungeons").Append(&dungeon); err != nil {
				return fmt.Errorf("error associating dungeon %s with season %s: %v", d.Name, s.Name, err)
			}
		}
	}

	return nil
}

// UpdateSeasons updates existing seasons and adds new ones from a JSON file
// This function checks for existing records and updates them, or creates new ones if needed
func UpdateSeasons(db *gorm.DB, filePath string) error {
	var seasonData SeasonData
	seasonFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	if err := json.Unmarshal(seasonFile, &seasonData); err != nil {
		return fmt.Errorf("error unmarshaling file %s: %v", filePath, err)
	}

	for _, s := range seasonData.Seasons {
		// Check if season exists
		var existingSeason mythicplus.Season
		result := db.Where("slug = ?", s.Slug).First(&existingSeason)

		season := mythicplus.Season{
			Slug:      s.Slug,
			Name:      s.Name,
			ShortName: s.ShortName,
			StartsUS:  parseTime(s.Starts.US),
			StartsEU:  parseTime(s.Starts.EU),
			StartsTW:  parseTime(s.Starts.TW),
			StartsKR:  parseTime(s.Starts.KR),
			StartsCN:  parseTime(s.Starts.CN),
			EndsUS:    parseTime(s.Ends.US),
			EndsEU:    parseTime(s.Ends.EU),
			EndsTW:    parseTime(s.Ends.TW),
			EndsKR:    parseTime(s.Ends.KR),
			EndsCN:    parseTime(s.Ends.CN),
		}

		if result.Error != nil {
			// Season doesn't exist, create it
			if err := db.Create(&season).Error; err != nil {
				return fmt.Errorf("error creating season %s: %v", s.Name, err)
			}
			log.Printf("Added new season: %s", s.Name)
		} else {
			// Season exists, update it
			existingSeason.Name = s.Name
			existingSeason.ShortName = s.ShortName
			existingSeason.StartsUS = parseTime(s.Starts.US)
			existingSeason.StartsEU = parseTime(s.Starts.EU)
			existingSeason.StartsTW = parseTime(s.Starts.TW)
			existingSeason.StartsKR = parseTime(s.Starts.KR)
			existingSeason.StartsCN = parseTime(s.Starts.CN)
			existingSeason.EndsUS = parseTime(s.Ends.US)
			existingSeason.EndsEU = parseTime(s.Ends.EU)
			existingSeason.EndsTW = parseTime(s.Ends.TW)
			existingSeason.EndsKR = parseTime(s.Ends.KR)
			existingSeason.EndsCN = parseTime(s.Ends.CN)

			if err := db.Save(&existingSeason).Error; err != nil {
				return fmt.Errorf("error updating season %s: %v", s.Name, err)
			}
			log.Printf("Updated existing season: %s", s.Name)

			// Use the existing season for dungeons
			season = existingSeason
		}

		// Process dungeons for this season
		for _, d := range s.Dungeons {
			// Check if dungeon exists
			var existingDungeon mythicplus.Dungeon
			result := db.Where("id = ?", d.ID).First(&existingDungeon)

			dungeon := mythicplus.Dungeon{
				ID:              d.ID,
				ChallengeModeID: d.ChallengeModeID,
				Slug:            d.Slug,
				Name:            d.Name,
				ShortName:       d.ShortName,
				MediaURL:        d.MediaURL,
				Icon:            &d.Icon,
			}

			// Only set EncounterID if it's greater than 0
			if d.EncounterID > 0 {
				encounterID := d.EncounterID
				dungeon.EncounterID = &encounterID
			}

			if result.Error != nil {
				// Dungeon doesn't exist, create it
				if err := db.Create(&dungeon).Error; err != nil {
					return fmt.Errorf("error creating dungeon %s: %v", d.Name, err)
				}
				log.Printf("Added new dungeon: %s", d.Name)
			} else {
				// Dungeon exists, update it
				existingDungeon.ChallengeModeID = d.ChallengeModeID
				existingDungeon.Slug = d.Slug
				existingDungeon.Name = d.Name
				existingDungeon.ShortName = d.ShortName
				existingDungeon.MediaURL = d.MediaURL
				if d.Icon != "" {
					existingDungeon.Icon = &d.Icon
				}

				// Only update EncounterID if the new value is greater than 0
				// This prevents overwriting existing valid IDs with zeros
				if d.EncounterID > 0 {
					encounterID := d.EncounterID
					existingDungeon.EncounterID = &encounterID
				}

				if err := db.Save(&existingDungeon).Error; err != nil {
					return fmt.Errorf("error updating dungeon %s: %v", d.Name, err)
				}
				log.Printf("Updated existing dungeon: %s", d.Name)

				// Use the existing dungeon
				dungeon = existingDungeon
			}

			// Update keystone upgrades
			if err := db.Where("challenge_mode_id = ?", d.ChallengeModeID).Delete(&mythicplus.KeyStoneUpgrade{}).Error; err != nil {
				return fmt.Errorf("error deleting old keystone upgrades for dungeon %s: %v", d.Name, err)
			}

			// Create new KeyStoneUpgrades if they exist
			if len(d.KeystoneUpgrades) > 0 {
				for _, ku := range d.KeystoneUpgrades {
					keyStoneUpgrade := mythicplus.KeyStoneUpgrade{
						ChallengeModeID:    d.ChallengeModeID,
						QualifyingDuration: int64(ku.QualifyingDuration),
						UpgradeLevel:       ku.UpgradeLevel,
					}

					if err := db.Create(&keyStoneUpgrade).Error; err != nil {
						return fmt.Errorf("error creating keystone upgrade for dungeon %s: %v", d.Name, err)
					}
				}
			}

			// Check if dungeon is already associated with this season
			var count int64
			db.Table("season_dungeons").
				Where("season_id = ? AND dungeon_id = ?", season.ID, dungeon.ID).
				Count(&count)

			if count == 0 {
				// Associate dungeon with season if not already associated
				if err := db.Model(&season).Association("Dungeons").Append(&dungeon); err != nil {
					return fmt.Errorf("error associating dungeon %s with season %s: %v", d.Name, s.Name, err)
				}
				log.Printf("Associated dungeon %s with season %s", d.Name, s.Name)
			}
		}
	}

	return nil
}

// SeedAffixes seeds the Mythic+ affixes from a JSON file
// This is the original function which is now complemented by UpdateAffixes
func SeedAffixes(db *gorm.DB, filePath string) error {
	var affixData AffixData
	affixFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	if err := json.Unmarshal(affixFile, &affixData); err != nil {
		return fmt.Errorf("error unmarshaling file %s: %v", filePath, err)
	}

	for _, a := range affixData.AffixDetails {
		affix := mythicplus.Affix{
			ID:          a.ID,
			Name:        a.Name,
			Icon:        a.Icon,
			Description: a.Description,
			WowheadURL:  a.WowheadURL,
		}

		if err := db.FirstOrCreate(&affix, mythicplus.Affix{ID: a.ID}).Error; err != nil {
			return fmt.Errorf("error creating affix %s: %v", a.Name, err)
		}
	}

	return nil
}

// UpdateAffixes updates existing affixes and adds new ones from a JSON file
func UpdateAffixes(db *gorm.DB, filePath string) error {
	var affixData AffixData
	affixFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	if err := json.Unmarshal(affixFile, &affixData); err != nil {
		return fmt.Errorf("error unmarshaling file %s: %v", filePath, err)
	}

	for _, a := range affixData.AffixDetails {
		// Check if affix exists
		var existingAffix mythicplus.Affix
		result := db.Where("id = ?", a.ID).First(&existingAffix)

		if result.Error != nil {
			// Affix doesn't exist, create it
			affix := mythicplus.Affix{
				ID:          a.ID,
				Name:        a.Name,
				Icon:        a.Icon,
				Description: a.Description,
				WowheadURL:  a.WowheadURL,
			}

			if err := db.Create(&affix).Error; err != nil {
				return fmt.Errorf("error creating affix %s: %v", a.Name, err)
			}
			log.Printf("Added new affix: %s", a.Name)
		} else {
			// Affix exists, update it
			existingAffix.Name = a.Name
			existingAffix.Icon = a.Icon
			existingAffix.Description = a.Description
			existingAffix.WowheadURL = a.WowheadURL

			if err := db.Save(&existingAffix).Error; err != nil {
				return fmt.Errorf("error updating affix %s: %v", a.Name, err)
			}
			log.Printf("Updated existing affix: %s", a.Name)
		}
	}

	return nil
}

// parseTime parses a time string in RFC3339 format
func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		return time.Time{}
	}
	return t
}
