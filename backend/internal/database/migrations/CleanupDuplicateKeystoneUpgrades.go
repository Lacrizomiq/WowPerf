package migrations

import (
	"gorm.io/gorm"
)

func CleanupDuplicateKeystoneUpgrades(db *gorm.DB) error {
	return db.Exec(`
        DELETE FROM keystone_upgrades
        WHERE id IN (
            SELECT id
            FROM (
                SELECT id,
                       ROW_NUMBER() OVER (
                           PARTITION BY mythic_dj_id, upgrade_level
                           ORDER BY id
                       ) AS row_num
                FROM keystone_upgrades
            ) t
            WHERE t.row_num > 1
        );
    `).Error
}
