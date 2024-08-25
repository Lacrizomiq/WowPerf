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
		&talents.TalentTree{},
		&talents.TalentNode{},
		&talents.TalentEntry{},
		&talents.SubTreeNode{},
		&talents.SubTreeEntry{},
	); err != nil {
		return err
	}

	// Add unique constraint on trait_tree_id
	if err := db.Exec("ALTER TABLE talent_trees ADD CONSTRAINT uni_talent_trees_trait_tree_id UNIQUE (trait_tree_id)").Error; err != nil {
		return err
	}

	// Add unique constraint on node_id in talent_nodes
	if err := db.Exec("ALTER TABLE talent_nodes ADD CONSTRAINT uni_talent_nodes_node_id UNIQUE (node_id)").Error; err != nil {
		return err
	}

	// Add foreign key constraints
	constraints := []struct {
		table      string
		constraint string
		query      string
	}{
		{"talent_nodes", "fk_talent_trees", "ALTER TABLE talent_nodes ADD CONSTRAINT fk_talent_trees FOREIGN KEY (talent_tree_id) REFERENCES talent_trees(trait_tree_id)"},
		{"talent_entries", "fk_talent_nodes", "ALTER TABLE talent_entries ADD CONSTRAINT fk_talent_nodes FOREIGN KEY (node_id) REFERENCES talent_nodes(node_id)"},
		{"sub_tree_nodes", "fk_talent_trees_sub", "ALTER TABLE sub_tree_nodes ADD CONSTRAINT fk_talent_trees_sub FOREIGN KEY (talent_tree_id) REFERENCES talent_trees(trait_tree_id)"},
		{"sub_tree_entries", "fk_sub_tree_nodes_entries", "ALTER TABLE sub_tree_entries ADD CONSTRAINT fk_sub_tree_nodes_entries FOREIGN KEY (sub_tree_node_id) REFERENCES sub_tree_nodes(sub_tree_node_id)"},
	}

	for _, c := range constraints {
		var count int64
		db.Raw("SELECT COUNT(*) FROM information_schema.table_constraints WHERE table_name = ? AND constraint_name = ?", c.table, c.constraint).Scan(&count)
		if count == 0 {
			if err := db.Exec(c.query).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
