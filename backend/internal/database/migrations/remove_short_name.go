package migrations

import (
	"gorm.io/gorm"
)

func RemoveShortNameFromMythicDjs(db *gorm.DB) error {
	if db.Migrator().HasColumn("mythic_djs", "short_name") {
		return db.Migrator().DropColumn("mythic_djs", "short_name")
	}
	return nil
}
