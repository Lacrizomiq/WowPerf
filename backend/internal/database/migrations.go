package database

import (
	"fmt"
	"strings"

	"wowperf/internal/models"
	mythicplus "wowperf/internal/models/mythicplus"
	raiderioMythicPlus "wowperf/internal/models/raiderio/mythicrundetails"
	raids "wowperf/internal/models/raids"
	talents "wowperf/internal/models/talents"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error; err != nil {
		return err
	}

	if err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error; err != nil {
		return err
	}

	// Mythic+ migrations
	if err := db.AutoMigrate(&mythicplus.Dungeon{}, &mythicplus.Season{}, &mythicplus.Affix{}, &raiderioMythicPlus.DungeonStats{}); err != nil {
		return err
	}

	// Migrate User model
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate User model: %v", err)
	}

	// Add new columns to User model
	if err := db.Exec(`
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
`).Error; err != nil {
		return fmt.Errorf("failed to add new columns to User model: %v", err)
	}

	// Add unique constraints
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username)").Error; err != nil {
		return fmt.Errorf("failed to create unique index on username: %v", err)
	}
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
		return fmt.Errorf("failed to create unique index on email: %v", err)
	}

	// Create KeyStoneUpgrade table without foreign key constraint
	if err := db.AutoMigrate(&mythicplus.KeyStoneUpgrade{}); err != nil {
		return err
	}

	// Raids migrations
	if err := db.AutoMigrate(&raids.Raid{}); err != nil {
		return err
	}

	// Talents migrations
	if err := db.AutoMigrate(
		&talents.TalentTree{},
		&talents.TalentNode{},
		&talents.TalentEntry{},
		&talents.HeroNode{},
		&talents.HeroEntry{},
		&talents.SubTreeNode{},
		&talents.SubTreeEntry{},
	); err != nil {
		return err
	}

	// UpdateState migration
	if err := db.AutoMigrate(&raiderioMythicPlus.UpdateState{}); err != nil {
		return err
	}

	// Add foreign key constraint for KeyStoneUpgrade if it doesn't exist
	if err := db.Exec(`
DO $$ 
BEGIN
		IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_dungeons_key_stone_upgrades'
		) THEN
				ALTER TABLE key_stone_upgrades 
				ADD CONSTRAINT fk_dungeons_key_stone_upgrades 
				FOREIGN KEY (challenge_mode_id) 
				REFERENCES dungeons(challenge_mode_id);
		END IF;
END $$;
`).Error; err != nil {
		return fmt.Errorf("failed to add foreign key constraint: %v", err)
	}

	// Helper function to check and add/update constraints
	addOrUpdateConstraint := func(tableName, constraintName, constraintDefinition string) error {
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

	// Add or update constraints
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
		{"dungeon_stats", "uni_dungeon_stats_season_region_dungeon", "UNIQUE (season, region, dungeon_slug)"}, // Nouvelle contrainte
	}

	for _, c := range constraints {
		if err := addOrUpdateConstraint(c.tableName, c.constraintName, c.constraintDefinition); err != nil {
			return err
		}
	}

	// Update foreign key constraints
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
