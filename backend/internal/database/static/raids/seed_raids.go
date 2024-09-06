package database

import (
	"encoding/json"
	"fmt"
	"os"

	raids "wowperf/internal/models/raids"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
