package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"wowperf/internal/models"
	mythicplus "wowperf/internal/models/mythicplus"
	raiderioMythicPlus "wowperf/internal/models/raiderio/mythicrundetails"
	raids "wowperf/internal/models/raids"
	talents "wowperf/internal/models/talents"

	rankingsModels "wowperf/internal/models/warcraftlogs/mythicplus"

	"gorm.io/gorm"
)

type Migration struct {
	ID      uint   `gorm:"primaryKey"`
	Version string `gorm:"uniqueIndex"`
}

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

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
		{"006_init_team_comp", initTeamComp},
		{"007_clean_team_comp_data", cleanTeamCompData},
		{"008_add_rankings_models", addRankingsTables},
		{"009_ensure_rankings_data", ensureRankingsData},
		{"010_update_player_rankings", updatePlayerRankings},
		{"011_clean_rankings_update_state", cleanRankingsUpdateState},
		{"012_clean_and_constrain_rankings_update_state", cleanAndConstrainRankingsUpdateState},
		{"013_clean_and_fix_rankings_update_state", cleanAndFixRankingsUpdateState},
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
		&rankingsModels.RankingsUpdateState{},
		&rankingsModels.PlayerRanking{},
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

// Add dungeon stats from Raider.io API
func addDungeonStats(db *gorm.DB) error {
	return db.AutoMigrate(&raiderioMythicPlus.DungeonStats{})
}

// Add rankings models from WarcraftLogs API
func addRankingsTables(db *gorm.DB) error {
	// First create the rankings update state table
	if err := db.AutoMigrate(&rankingsModels.RankingsUpdateState{}); err != nil {
		return fmt.Errorf("failed to create rankings update state table: %v", err)
	}

	// Then create the player rankings table
	if err := db.AutoMigrate(&rankingsModels.PlayerRanking{}); err != nil {
		return fmt.Errorf("failed to create player rankings table: %v", err)
	}

	return nil
}

func initTeamComp(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Drop the existing column if it exists
		if err := tx.Exec(`ALTER TABLE dungeon_stats DROP COLUMN IF EXISTS team_comp;`).Error; err != nil {
			// Ignore the error if the column doesn't exist
			log.Printf("Note: team_comp column didn't exist or couldn't be dropped: %v", err)
		}

		// Add the new column with proper type and default value
		if err := tx.Exec(`
            ALTER TABLE dungeon_stats 
            ADD COLUMN team_comp JSONB DEFAULT '{}'::jsonb;
            
            COMMENT ON COLUMN dungeon_stats.team_comp 
            IS 'Stores the team composition statistics with their counts';
        `).Error; err != nil {
			return fmt.Errorf("failed to add team_comp column: %v", err)
		}

		return nil
	})
}

func cleanTeamCompData(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Check if the column exists before updating it
		var exists bool
		err := tx.Raw(`
            SELECT EXISTS (
                SELECT FROM information_schema.columns 
                WHERE table_name = 'dungeon_stats' 
                AND column_name = 'team_comp'
            );
        `).Scan(&exists).Error

		if err != nil {
			return fmt.Errorf("failed to check if team_comp column exists: %v", err)
		}

		if exists {
			if err := tx.Exec(`
                UPDATE dungeon_stats 
                SET team_comp = '{}'::jsonb 
                WHERE team_comp IS NULL OR team_comp = 'null'::jsonb;
            `).Error; err != nil {
				return fmt.Errorf("failed to clean team_comp data: %v", err)
			}
		}

		return nil
	})
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

func updatePlayerRankings(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Drop the existing table
		if err := tx.Exec("DROP TABLE IF EXISTS player_rankings CASCADE").Error; err != nil {
			return fmt.Errorf("failed to drop existing table: %w", err)
		}

		// Create the new table with the correct affixes definition
		createTableSQL := `
			CREATE TABLE player_rankings (
					id SERIAL PRIMARY KEY,
					created_at TIMESTAMP WITH TIME ZONE,
					updated_at TIMESTAMP WITH TIME ZONE,
					deleted_at TIMESTAMP WITH TIME ZONE,
					dungeon_id INTEGER,
					name VARCHAR(255),
					class VARCHAR(255),
					spec VARCHAR(255),
					role VARCHAR(255),
					amount DOUBLE PRECISION,
					hard_mode_level INTEGER,
					duration BIGINT,
					start_time BIGINT,
					report_code VARCHAR(255),
					report_fight_id INTEGER,
					report_start_time BIGINT,
					guild_id INTEGER,
					guild_name VARCHAR(255),
					guild_faction INTEGER,
					server_id INTEGER,
					server_name VARCHAR(255),
					server_region VARCHAR(50),
					bracket_data INTEGER,
					faction INTEGER,
					affixes INTEGER[],
					medal VARCHAR(50),
					score DOUBLE PRECISION,
					leaderboard INTEGER DEFAULT 0
			);

			CREATE INDEX IF NOT EXISTS idx_player_rankings_dungeon_id ON player_rankings(dungeon_id);
			CREATE INDEX IF NOT EXISTS idx_player_rankings_name ON player_rankings(name);
			CREATE INDEX IF NOT EXISTS idx_player_rankings_class ON player_rankings(class);
			CREATE INDEX IF NOT EXISTS idx_player_rankings_spec ON player_rankings(spec);
			CREATE INDEX IF NOT EXISTS idx_player_rankings_role ON player_rankings(role);
			CREATE INDEX IF NOT EXISTS idx_player_rankings_guild_id ON player_rankings(guild_id);
			CREATE INDEX IF NOT EXISTS idx_player_rankings_server_id ON player_rankings(server_id);
			CREATE INDEX IF NOT EXISTS idx_player_rankings_deleted_at ON player_rankings(deleted_at);
			`

		if err := tx.Exec(createTableSQL).Error; err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}

		return nil
	})
}

// CleanRankingsUpdateState cleans the rankings update state
func cleanRankingsUpdateState(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Supprime toutes les entrées existantes
		if err := tx.Exec("DELETE FROM rankings_update_states").Error; err != nil {
			return fmt.Errorf("failed to clean rankings update states: %v", err)
		}

		// Crée une nouvelle entrée avec la date d'il y a 25 heures
		state := rankingsModels.RankingsUpdateState{
			LastUpdateTime: time.Now().Add(-25 * time.Hour),
		}
		if err := tx.Create(&state).Error; err != nil {
			return fmt.Errorf("failed to create new rankings update state: %v", err)
		}

		// Ajoute une contrainte unique si elle n'existe pas déjà
		return tx.Exec(`
					DO $$
					BEGIN
							IF NOT EXISTS (
									SELECT 1 FROM information_schema.table_constraints 
									WHERE table_name = 'rankings_update_states' 
									AND constraint_type = 'UNIQUE'
							) THEN
									ALTER TABLE rankings_update_states
									ADD CONSTRAINT rankings_update_states_single_row
									UNIQUE (id);
							END IF;
					END
					$$;
			`).Error
	})
}

func cleanAndConstrainRankingsUpdateState(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Supprime la table existante
		if err := tx.Exec("DROP TABLE IF EXISTS rankings_update_states CASCADE").Error; err != nil {
			return err
		}

		// Crée la nouvelle table avec la contrainte
		if err := tx.Exec(`
					CREATE TABLE rankings_update_states (
							id INTEGER PRIMARY KEY CHECK (id = 1),
							created_at TIMESTAMP WITH TIME ZONE,
							updated_at TIMESTAMP WITH TIME ZONE,
							deleted_at TIMESTAMP WITH TIME ZONE,
							last_update_time TIMESTAMP WITH TIME ZONE
					)
			`).Error; err != nil {
			return err
		}

		// Crée l'entrée initiale
		return tx.Exec(`
					INSERT INTO rankings_update_states (id, created_at, updated_at, last_update_time)
					VALUES (1, NOW(), NOW(), NOW() - INTERVAL '25 hours')
			`).Error
	})
}

func cleanAndFixRankingsUpdateState(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Supprime la table existante
		if err := tx.Exec("DROP TABLE IF EXISTS rankings_update_states CASCADE").Error; err != nil {
			return err
		}

		// Crée la nouvelle table avec la contrainte
		if err := tx.Exec(`
					CREATE TABLE rankings_update_states (
							id INTEGER PRIMARY KEY DEFAULT 1,
							created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
							updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
							deleted_at TIMESTAMP WITH TIME ZONE,
							last_update_time TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP - INTERVAL '25 hours'),
							CONSTRAINT ensure_single_row CHECK (id = 1)
					)
			`).Error; err != nil {
			return err
		}

		// Insert initial row
		return tx.Exec(`
					INSERT INTO rankings_update_states (id)
					VALUES (1)
					ON CONFLICT (id) DO NOTHING
			`).Error
	})
}
