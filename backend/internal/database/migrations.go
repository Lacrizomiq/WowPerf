package database

import (
	"fmt"
	mythicplus "wowperf/internal/models/mythicplus"
	raids "wowperf/internal/models/raids"
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
		{"talent_nodes", "uni_talent_nodes_node_tree_spec", "UNIQUE (node_id, talent_tree_id, spec_id)"},
		{"talent_entries", "uni_talent_entries_entry_node_tree_spec", "UNIQUE (entry_id, node_id, talent_tree_id, spec_id)"},
		{"talent_trees", "uni_talent_trees_trait_tree_spec_id", "UNIQUE (trait_tree_id, spec_id)"},
		{"sub_tree_nodes", "uni_sub_tree_nodes_id_tree_spec", "UNIQUE (sub_tree_node_id, talent_tree_id, spec_id)"},
		{"hero_nodes", "uni_hero_nodes_node_tree_spec", "UNIQUE (node_id, talent_tree_id, spec_id)"},                     // Added constraint for HeroNodes
		{"hero_entries", "uni_hero_entries_entry_node_tree_spec", "UNIQUE (entry_id, node_id, talent_tree_id, spec_id)"}, // Added constraint for HeroEntries
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
		{"hero_nodes", "fk_talent_trees_hero", "ALTER TABLE hero_nodes ADD CONSTRAINT fk_talent_trees_hero FOREIGN KEY (talent_tree_id, spec_id) REFERENCES talent_trees(trait_tree_id, spec_id)"},                        // Added foreign key for HeroNodes
		{"hero_entries", "fk_hero_nodes_entries", "ALTER TABLE hero_entries ADD CONSTRAINT fk_hero_nodes_entries FOREIGN KEY (node_id, talent_tree_id, spec_id) REFERENCES hero_nodes(node_id, talent_tree_id, spec_id)"}, // Added foreign key for HeroEntries
	}

	for _, fk := range foreignKeys {
		db.Exec(fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS %s", fk.table, fk.constraint))
		if err := db.Exec(fk.query).Error; err != nil {
			return err
		}
	}

	return nil
}
