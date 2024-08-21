package database

import (
	mythicplus "wowperf/internal/database/static/mythicplus"

	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	return mythicplus.SeedDatabase(db)
}
