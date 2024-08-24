package database

import (
	mythicplus "wowperf/internal/models/mythicplus"
	talents "wowperf/internal/models/talents"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Mythic+ migrations
	if err := db.AutoMigrate(&mythicplus.Season{}, &mythicplus.Dungeon{}, &mythicplus.Affix{}, &mythicplus.KeyStoneUpgrade{}); err != nil {
		return err
	}

	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.table_constraints WHERE table_name = 'dungeons' AND constraint_name = 'uni_dungeons_challenge_mode_id'").Scan(&count)

	if count == 0 {
		if err := db.Exec("ALTER TABLE dungeons ADD CONSTRAINT uni_dungeons_challenge_mode_id UNIQUE (challenge_mode_id)").Error; err != nil {
			return err
		}
	}

	// Talents migrations
	if err := db.AutoMigrate(
		&talents.ClassTalent{},
		&talents.SpecTalent{},
		&talents.HeroTalent{},
		&talents.SubTreeTalent{},
		&talents.TalentNode{},
		&talents.TalentEntry{},
		&talents.FullNodeOrder{},
	); err != nil {
		return err
	}

	return nil
}
