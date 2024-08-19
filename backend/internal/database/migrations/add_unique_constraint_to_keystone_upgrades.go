package migrations

import (
	"gorm.io/gorm"
)

func AddUniqueConstraintToKeystoneUpgrades(db *gorm.DB) error {
	return db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_mythic_dj_upgrade ON keystone_upgrades (mythic_dj_id, upgrade_level)").Error
}
