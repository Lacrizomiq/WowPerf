package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	raids "wowperf/internal/models/raids"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SeedRaids seeds the raids data from a JSON file
// This is the original function which is now complemented by UpdateRaids
func SeedRaids(db *gorm.DB, filePath string) error {
	var raidData struct {
		Raids []raids.Raid `json:"raids"`
	}

	raidFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	if err := json.Unmarshal(raidFile, &raidData); err != nil {
		return fmt.Errorf("error unmarshaling file %s: %v", filePath, err)
	}

	for _, raidInput := range raidData.Raids {
		raid := raids.Raid{
			ID:         raidInput.ID,
			Slug:       raidInput.Slug,
			Name:       raidInput.Name,
			ShortName:  raidInput.ShortName,
			Expansion:  raidInput.Expansion,
			MediaURL:   raidInput.MediaURL,
			Icon:       raidInput.Icon,
			Starts:     raidInput.Starts,
			Ends:       raidInput.Ends,
			Encounters: raidInput.Encounters,
		}

		// Use Upsert operation
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"slug", "name", "short_name", "expansion", "media_url", "icon", "starts", "ends", "encounters"}),
		}).Create(&raid).Error; err != nil {
			return fmt.Errorf("error upserting raid %s: %v", raid.ShortName, err)
		}
	}

	return nil
}

// UpdateRaids updates existing raids and adds new ones from a JSON file
// This function explicitly checks for existing records and updates them, or creates new ones if needed
func UpdateRaids(db *gorm.DB, filePath string) error {
	var raidData struct {
		Raids []raids.Raid `json:"raids"`
	}

	raidFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	if err := json.Unmarshal(raidFile, &raidData); err != nil {
		return fmt.Errorf("error unmarshaling file %s: %v", filePath, err)
	}

	for _, raidInput := range raidData.Raids {
		// Check if raid exists
		var existingRaid raids.Raid
		result := db.Where("id = ?", raidInput.ID).First(&existingRaid)

		if result.Error != nil {
			// Raid doesn't exist, create it
			if err := db.Create(&raidInput).Error; err != nil {
				return fmt.Errorf("error creating raid %s: %v", raidInput.Name, err)
			}
			log.Printf("Added new raid: %s", raidInput.Name)
		} else {
			// Raid exists, update its fields
			existingRaid.Slug = raidInput.Slug
			existingRaid.Name = raidInput.Name
			existingRaid.ShortName = raidInput.ShortName
			existingRaid.Expansion = raidInput.Expansion
			existingRaid.MediaURL = raidInput.MediaURL
			existingRaid.Icon = raidInput.Icon
			existingRaid.Starts = raidInput.Starts
			existingRaid.Ends = raidInput.Ends

			// For updating encounters, we'll use the Replace approach instead of Clear and Add
			// This requires replacing the existing record with a new one including all updated fields

			// First, prepare the raid with encounters for saving
			raidToSave := existingRaid
			raidToSave.Encounters = raidInput.Encounters

			// Save the raid with updates including encounters
			if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&raidToSave).Error; err != nil {
				// If the above approach doesn't work well, we can try an alternative approach
				// by using the original function's upsert capability
				if retryErr := db.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "id"}},
					DoUpdates: clause.AssignmentColumns([]string{"slug", "name", "short_name", "expansion", "media_url", "icon", "starts", "ends"}),
				}).Create(&raidInput).Error; retryErr != nil {
					return fmt.Errorf("error updating raid %s: %v (retry also failed: %v)", raidInput.Name, err, retryErr)
				}
			}

			log.Printf("Updated existing raid: %s", raidInput.Name)
		}
	}

	return nil
}
