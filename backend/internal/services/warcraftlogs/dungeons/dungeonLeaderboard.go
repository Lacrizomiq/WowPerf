package warcraftlogs

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Constants for security
const (
	MaxLimit     = 1000 // Maximum results per query
	DefaultLimit = 100  // Default results if not specified
)

// validDungeonIDs contains all valid dungeon IDs
var validDungeonIDs = map[int]bool{
	DungeonAraKara:       true,
	DungeonCityOfThreads: true,
	DungeonGrimBatol:     true,
	DungeonMists:         true,
	DungeonSiege:         true,
	DungeonDawnbreaker:   true,
	DungeonNecroticWake:  true,
	DungeonStonevault:    true,
}

// DungeonLeaderboardEntry represents a single entry in the dungeon leaderboard
type DungeonLeaderboardEntry struct {
	Name         string        `json:"name"`
	Class        string        `json:"class"`
	Spec         string        `json:"spec"`
	Role         string        `json:"role"`
	Score        float64       `json:"score"`
	Rank         int           `json:"rank"`
	Duration     int           `json:"duration"`
	Level        int           `json:"hard_mode_level"`
	Medal        string        `json:"medal"`
	ServerName   string        `json:"server_name"`
	ServerRegion string        `json:"server_region"`
	StartTime    int64         `json:"start_time"`
	GuildName    string        `json:"guild_name"`
	GuildFaction int           `json:"guild_faction"`
	Faction      int           `json:"faction"`
	Affixes      pq.Int64Array `json:"affixes"`
	BracketData  int           `json:"bracket_data"`
}

// DungeonLeaderboardService is a service for querying the dungeon leaderboard
type DungeonLeaderboardService struct {
	db           *gorm.DB
	validRoles   map[string]bool
	validClasses map[string]bool
}

// NewDungeonLeaderboardService initializes service with security validations
func NewDungeonLeaderboardService(db *gorm.DB) *DungeonLeaderboardService {
	return &DungeonLeaderboardService{
		db: db,
		validRoles: map[string]bool{
			"Tank":   true,
			"Healer": true,
			"DPS":    true,
		},
		validClasses: map[string]bool{
			"Warrior":     true,
			"Paladin":     true,
			"Hunter":      true,
			"Rogue":       true,
			"Priest":      true,
			"Shaman":      true,
			"Mage":        true,
			"Warlock":     true,
			"Monk":        true,
			"Druid":       true,
			"Demonhunter": true,
			"Deathknight": true,
			"Evoker":      true,
		},
	}
}

// createDungeonIndexes creates database indexes for optimized query performance
func (s *DungeonLeaderboardService) createDungeonIndexes() error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_rankings_dungeon_score ON player_rankings(dungeon_id, score DESC, hard_mode_level DESC)",
		"CREATE INDEX IF NOT EXISTS idx_rankings_dungeon_role ON player_rankings(dungeon_id, role, score DESC)",
		"CREATE INDEX IF NOT EXISTS idx_rankings_dungeon_class ON player_rankings(dungeon_id, class, score DESC)",
		"CREATE INDEX IF NOT EXISTS idx_rankings_dungeon_spec ON player_rankings(dungeon_id, class, spec, score DESC)",
	}

	for _, idx := range indexes {
		if err := s.db.Exec(idx).Error; err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// getBaseDungeonLeaderboardQuery returns the base query for the dungeon leaderboard
func (s *DungeonLeaderboardService) getBaseDungeonLeaderboardQuery() string {
	return `
		WITH DungeonScores AS (
			SELECT 
				name,
				server_name,
				server_region,
				class,
				spec,
				role,
				score,
				medal,
				hard_mode_level,
				duration,
				start_time,
				guild_name,
				guild_faction,
				faction,
				affixes,
				bracket_data
			FROM player_rankings
			WHERE deleted_at IS NULL
				AND dungeon_id = ?
				AND score > 0  -- Protect against invalid scores
				%s  -- Additional WHERE conditions placeholder
		)
		SELECT 
			*,
			DENSE_RANK() OVER (
				ORDER BY score DESC, 
				hard_mode_level DESC, 
				duration ASC
			) as rank
		FROM DungeonScores
		ORDER BY rank ASC, duration ASC
		LIMIT ?`
}

// GetDungeonLeaderboardByPlayer returns the global dungeon leaderboard
func (s *DungeonLeaderboardService) GetDungeonLeaderboardByPlayer(
	ctx context.Context,
	dungeonID int,
	limit int,
) ([]DungeonLeaderboardEntry, error) {
	if err := s.validateDungeonInput(dungeonID, "", "", "", limit); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	query := fmt.Sprintf(s.getBaseDungeonLeaderboardQuery(), "")
	var entries []DungeonLeaderboardEntry
	err := s.db.WithContext(ctx).Raw(query, dungeonID, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dungeon leaderboard: %w", err)
	}

	return entries, nil
}

// GetDungeonLeaderboardByRole returns the dungeon leaderboard filtered by role
func (s *DungeonLeaderboardService) GetDungeonLeaderboardByRole(
	ctx context.Context,
	dungeonID int,
	role string,
	limit int,
) ([]DungeonLeaderboardEntry, error) {
	if err := s.validateDungeonInput(dungeonID, role, "", "", limit); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	query := fmt.Sprintf(s.getBaseDungeonLeaderboardQuery(), "AND role = ?")
	var entries []DungeonLeaderboardEntry
	err := s.db.WithContext(ctx).Raw(query, dungeonID, role, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch role leaderboard: %w", err)
	}

	return entries, nil
}

// GetDungeonLeaderboardByClass returns the dungeon leaderboard filtered by class
func (s *DungeonLeaderboardService) GetDungeonLeaderboardByClass(
	ctx context.Context,
	dungeonID int,
	class string,
	limit int,
) ([]DungeonLeaderboardEntry, error) {
	if err := s.validateDungeonInput(dungeonID, "", class, "", limit); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	query := fmt.Sprintf(s.getBaseDungeonLeaderboardQuery(), "AND class = ?")
	var entries []DungeonLeaderboardEntry
	err := s.db.WithContext(ctx).Raw(query, dungeonID, class, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch class leaderboard: %w", err)
	}

	return entries, nil
}

// GetDungeonLeaderboardBySpec returns the dungeon leaderboard filtered by class and spec
func (s *DungeonLeaderboardService) GetDungeonLeaderboardBySpec(
	ctx context.Context,
	dungeonID int,
	class string,
	spec string,
	limit int,
) ([]DungeonLeaderboardEntry, error) {
	if err := s.validateDungeonInput(dungeonID, "", class, spec, limit); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	query := fmt.Sprintf(s.getBaseDungeonLeaderboardQuery(), "AND class = ? AND spec = ?")
	var entries []DungeonLeaderboardEntry
	err := s.db.WithContext(ctx).Raw(query, dungeonID, class, spec, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spec leaderboard: %w", err)
	}

	return entries, nil
}

// validateDungeonInput validates the input for the dungeon leaderboard
func (s *DungeonLeaderboardService) validateDungeonInput(dungeonID int, role, class, spec string, limit int) error {
	// Validate dungeon ID against whitelist
	if !validDungeonIDs[dungeonID] {
		return fmt.Errorf("invalid dungeon ID: %d", dungeonID)
	}

	// Apply default limit if not specified
	if limit <= 0 {
		limit = DefaultLimit
	}

	// Enforce maximum limit
	if limit > MaxLimit {
		return fmt.Errorf("invalid limit: must be between 1 and %d", MaxLimit)
	}

	// Validate role if provided
	if role != "" && !s.validRoles[formatRole(role)] {
		return fmt.Errorf("invalid role: %s", role)
	}

	// Validate class if provided
	if class != "" && !s.validClasses[formatNameCase(class)] {
		return fmt.Errorf("invalid class: %s", class)
	}

	return nil
}
