package database

import (
	"fmt"
	"log"
	"strings"

	"wowperf/internal/models"
	mythicplus "wowperf/internal/models/mythicplus"
	raiderioMythicPlus "wowperf/internal/models/raiderio/mythicrundetails"
	raids "wowperf/internal/models/raids"
	talents "wowperf/internal/models/talents"

	"gorm.io/gorm"
)

type Migration struct {
	ID      uint   `gorm:"primaryKey"`
	Version string `gorm:"uniqueIndex"`
}

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Create migrations table if it doesn't exist
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	migrations := []struct {
		version string
		up      func(*gorm.DB) error
	}{
		{"001_initial_schema", initialSchema},
		{"002_add_user_columns", addUserColumns},
		{"003_add_dungeon_stats", addDungeonStats},
		{"004_add_constraints", addConstraints},
		{"005_update_foreign_keys", updateForeignKeys},
		{"006_add_team_comp_to_dungeon_stats", addTeamCompToDungeonStats},
	}

	for _, m := range migrations {
		var migration Migration
		if err := db.Where("version = ?", m.version).First(&migration).Error; err == gorm.ErrRecordNotFound {
			log.Printf("Running migration: %s", m.version)
			if err := m.up(db); err != nil {
				return fmt.Errorf("failed to run migration %s: %v", m.version, err)
			}
			db.Create(&Migration{Version: m.version})
		} else if err != nil {
			return fmt.Errorf("error checking migration %s: %v", m.version, err)
		}
	}

	log.Println("Database migrations completed successfully.")
	return nil
}

func initialSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&mythicplus.Dungeon{},
		&mythicplus.Season{},
		&mythicplus.Affix{},
		&raiderioMythicPlus.DungeonStats{},
		&models.User{},
		&mythicplus.KeyStoneUpgrade{},
		&raids.Raid{},
		&talents.TalentTree{},
		&talents.TalentNode{},
		&talents.TalentEntry{},
		&talents.HeroNode{},
		&talents.HeroEntry{},
		&talents.SubTreeNode{},
		&talents.SubTreeEntry{},
		&raiderioMythicPlus.UpdateState{},
	)
}

func addUserColumns(db *gorm.DB) error {
	return db.Exec(`
	ALTER TABLE users 
	ADD COLUMN IF NOT EXISTS battle_net_id INTEGER UNIQUE,
	ADD COLUMN IF NOT EXISTS battle_tag VARCHAR(255) UNIQUE,
	ADD COLUMN IF NOT EXISTS encrypted_token BYTEA,
	ADD COLUMN IF NOT EXISTS battle_net_expires_at TIMESTAMP,
	ADD COLUMN IF NOT EXISTS last_username_change_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

	ALTER TABLE users
	ALTER COLUMN battle_net_id DROP NOT NULL,
	ALTER COLUMN battle_tag DROP NOT NULL,
	ALTER COLUMN battle_net_expires_at DROP NOT NULL;
	`).Error
}

func addDungeonStats(db *gorm.DB) error {
	return db.AutoMigrate(&raiderioMythicPlus.DungeonStats{})
}

func addTeamCompToDungeonStats(db *gorm.DB) error {
	return db.Exec(`
	ALTER TABLE dungeon_stats
	ADD COLUMN IF NOT EXISTS team_comp JSONB;
	`).Error
}

func addConstraints(db *gorm.DB) error {
	constraints := []struct {
		tableName            string
		constraintName       string
		constraintDefinition string
	}{
		{"dungeons", "uni_dungeons_challenge_mode_id", "UNIQUE (challenge_mode_id)"},
		{"talent_nodes", "uni_talent_nodes_node_tree_spec", "UNIQUE (node_id, talent_tree_id, spec_id)"},
		{"talent_entries", "uni_talent_entries_entry_node_tree_spec", "UNIQUE (entry_id, node_id, talent_tree_id, spec_id)"},
		{"talent_trees", "uni_talent_trees_trait_tree_spec_id", "UNIQUE (trait_tree_id, spec_id)"},
		{"sub_tree_nodes", "uni_sub_tree_nodes_id_tree_spec", "UNIQUE (sub_tree_node_id, talent_tree_id, spec_id)"},
		{"hero_nodes", "uni_hero_nodes_node_tree_spec", "UNIQUE (node_id, talent_tree_id, spec_id)"},
		{"hero_entries", "uni_hero_entries_entry_node_tree_spec", "UNIQUE (entry_id, node_id, talent_tree_id, spec_id)"},
		{"dungeon_stats", "uni_dungeon_stats_season_region_dungeon", "UNIQUE (season, region, dungeon_slug)"},
	}

	for _, c := range constraints {
		if err := addOrUpdateConstraint(db, c.tableName, c.constraintName, c.constraintDefinition); err != nil {
			return err
		}
	}

	return nil
}

func updateForeignKeys(db *gorm.DB) error {
	foreignKeys := []struct {
		table      string
		constraint string
		query      string
	}{
		{"talent_nodes", "fk_talent_trees", "ALTER TABLE talent_nodes ADD CONSTRAINT fk_talent_trees FOREIGN KEY (talent_tree_id, spec_id) REFERENCES talent_trees(trait_tree_id, spec_id)"},
		{"talent_entries", "fk_talent_trees_entries", "ALTER TABLE talent_entries ADD CONSTRAINT fk_talent_trees_entries FOREIGN KEY (talent_tree_id, spec_id) REFERENCES talent_trees(trait_tree_id, spec_id)"},
		{"sub_tree_nodes", "fk_talent_trees_sub", "ALTER TABLE sub_tree_nodes ADD CONSTRAINT fk_talent_trees_sub FOREIGN KEY (talent_tree_id, spec_id) REFERENCES talent_trees(trait_tree_id, spec_id)"},
		{"sub_tree_entries", "fk_sub_tree_nodes_entries", "ALTER TABLE sub_tree_entries ADD CONSTRAINT fk_sub_tree_nodes_entries FOREIGN KEY (sub_tree_node_id) REFERENCES sub_tree_nodes(sub_tree_node_id)"},
		{"hero_nodes", "fk_talent_trees_hero", "ALTER TABLE hero_nodes ADD CONSTRAINT fk_talent_trees_hero FOREIGN KEY (talent_tree_id, spec_id) REFERENCES talent_trees(trait_tree_id, spec_id)"},
		{"hero_entries", "fk_hero_nodes_entries", "ALTER TABLE hero_entries ADD CONSTRAINT fk_hero_nodes_entries FOREIGN KEY (node_id, talent_tree_id, spec_id) REFERENCES hero_nodes(node_id, talent_tree_id, spec_id)"},
	}

	for _, fk := range foreignKeys {
		if err := db.Exec(fk.query).Error; err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}
	}

	return nil
}

func addOrUpdateConstraint(db *gorm.DB, tableName, constraintName, constraintDefinition string) error {
	var constraintExists int64
	db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM information_schema.table_constraints WHERE table_name = '%s' AND constraint_name = '%s'", tableName, constraintName)).Scan(&constraintExists)

	if constraintExists == 0 {
		if err := db.Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s", tableName, constraintName, constraintDefinition)).Error; err != nil {
			return err
		}
	} else {
		if err := db.Exec(fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s", tableName, constraintName)).Error; err != nil {
			return err
		}
		if err := db.Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s", tableName, constraintName, constraintDefinition)).Error; err != nil {
			return err
		}
	}
	return nil
}
