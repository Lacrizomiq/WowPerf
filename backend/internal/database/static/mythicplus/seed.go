package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	mythicplus "wowperf/internal/models/mythicplus"

	"gorm.io/gorm"
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

const (
	staticMythicPlusPath = "./static/M+/"
)

// seedSeasonsFromFile seeds the Mythic+ seasons from a JSON file
func SeedSeasons(db *gorm.DB, filePath string) error {
	var seasonData SeasonData
	seasonFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	if err := json.Unmarshal(seasonFile, &seasonData); err != nil {
		return fmt.Errorf("error unmarshaling file %s: %v", filePath, err)
	}

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

		if err := db.FirstOrCreate(&season, mythicplus.Season{Slug: s.Slug}).Error; err != nil {
			return fmt.Errorf("error creating season %s: %v", s.Name, err)
		}

		for _, d := range s.Dungeons {
			dungeon := mythicplus.Dungeon{
				ID:              d.ID,
				ChallengeModeID: &d.ChallengeModeID,
				Slug:            d.Slug,
				Name:            d.Name,
				ShortName:       d.ShortName,
				MediaURL:        d.MediaURL,
				Icon:            &d.Icon,
			}

			if err := db.FirstOrCreate(&dungeon, mythicplus.Dungeon{ID: d.ID}).Error; err != nil {
				return fmt.Errorf("error creating dungeon %s: %v", d.Name, err)
			}

			if err := db.Where("challenge_mode_id = ?", d.ChallengeModeID).Delete(&mythicplus.KeyStoneUpgrade{}).Error; err != nil {
				return fmt.Errorf("error deleting old keystone upgrades for dungeon %s: %v", d.Name, err)
			}

			for _, ku := range d.KeystoneUpgrades {
				keyStoneUpgrade := mythicplus.KeyStoneUpgrade{
					ChallengeModeID:    &d.ChallengeModeID,
					QualifyingDuration: ku.QualifyingDuration,
					UpgradeLevel:       ku.UpgradeLevel,
				}

				if err := db.Create(&keyStoneUpgrade).Error; err != nil {
					return fmt.Errorf("error creating keystone upgrade for dungeon %s: %v", d.Name, err)
				}
			}

			if err := db.Model(&season).Association("Dungeons").Append(&dungeon); err != nil {
				return fmt.Errorf("error associating dungeon %s with season %s: %v", d.Name, s.Name, err)
			}
		}
	}

	return nil
}

// seedAffixes seeds the Mythic+ affixes from a JSON file
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

// parseTime parses a time string in RFC3339 format
func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		return time.Time{}
	}
	return t
}
