package warcraftlogs

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const REQUIRED_DUNGEON_COUNT = 8

// LeaderboardEntry represents a single entry in the leaderboard
type LeaderboardEntry struct {
	Name         string  `json:"name"`
	Class        string  `json:"class"`
	Spec         string  `json:"spec"`
	Role         string  `json:"role"`
	TotalScore   float64 `json:"total_score"`
	Rank         int     `json:"rank"`
	DungeonCount int     `json:"dungeon_count"`
	ServerName   string  `json:"server_name"`
	ServerRegion string  `json:"server_region"`
	Medal        string  `json:"medal"`
}

// OrderDirection defines the sort order for queries
type OrderDirection string

const (
	ASC  OrderDirection = "ASC"
	DESC OrderDirection = "DESC"
)

// OrderByOption represents a column and its sort direction for queries
type OrderByOption struct {
	Column    string
	Direction OrderDirection
}

// formatNameCase formats a string with proper title case handling
// using Unicode-aware casing rules
func formatNameCase(s string) string {
	if s == "" {
		return s
	}
	caser := cases.Title(language.English)
	return caser.String(strings.ToLower(strings.TrimSpace(s)))
}

// formatRole handles special case formatting for roles
// particularly the "DPS" role which should remain uppercase
func formatRole(role string) string {
	formatted := strings.ToUpper(strings.TrimSpace(role))
	if formatted == "DPS" {
		return formatted
	}
	return formatNameCase(role)
}

// createOptimizedIndexes creates database indexes for optimized query performance
func (s *RankingsService) createOptimizedIndexes() error {
	indexes := []string{
		// Index for uniquely identifying players
		"CREATE INDEX IF NOT EXISTS idx_rankings_player_unique ON player_rankings(name, server_name, server_region)",
		// Composite indexes for different query types
		"CREATE INDEX IF NOT EXISTS idx_rankings_role_score ON player_rankings(role, score DESC, name, server_name, server_region)",
		"CREATE INDEX IF NOT EXISTS idx_rankings_class_score ON player_rankings(class, score DESC, name, server_name, server_region)",
		"CREATE INDEX IF NOT EXISTS idx_rankings_class_spec_score ON player_rankings(class, spec, score DESC, name, server_name, server_region)",
		"CREATE INDEX IF NOT EXISTS idx_rankings_score_global ON player_rankings(score DESC, name, server_name, server_region)",
	}

	for _, idx := range indexes {
		if err := s.db.Exec(idx).Error; err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}
	return nil
}

// sanitizeOrderBy ensures safe column names and directions for ORDER BY clauses
func (s *RankingsService) sanitizeOrderBy(column string, direction OrderDirection) OrderByOption {
	validColumns := map[string]string{
		"score":         "total_score",
		"name":          "name",
		"rank":          "rank",
		"server":        "server_name",
		"region":        "server_region",
		"class":         "class",
		"spec":          "spec",
		"role":          "role",
		"medal":         "best_medal",
		"dungeon_count": "dungeon_count",
	}

	if direction != ASC && direction != DESC {
		direction = DESC
	}

	if sanitized, ok := validColumns[column]; ok {
		return OrderByOption{
			Column:    sanitized,
			Direction: direction,
		}
	}

	// Default sorting by total_score DESC
	return OrderByOption{
		Column:    "total_score",
		Direction: DESC,
	}
}

// getBaseLeaderboardQuery returns the base CTE query for all leaderboard types
func (s *RankingsService) getBaseLeaderboardQuery(orderBy OrderByOption) string {
	return fmt.Sprintf(`
		WITH PlayerScores AS (
			SELECT 
				name,
				server_name,
				server_region,
				class,
				spec,
				role,
				MAX(medal) as best_medal,
				ROUND(CAST(SUM(score) AS numeric), 2) as total_score,
				COUNT(DISTINCT dungeon_id) as dungeon_count
			FROM player_rankings
			WHERE deleted_at IS NULL
			%%s  -- Additional WHERE conditions placeholder
			GROUP BY name, server_name, server_region, class, spec, role
			HAVING COUNT(DISTINCT dungeon_id) = %%d
		)
		SELECT 
			*,
			DENSE_RANK() OVER (ORDER BY %s %s, name ASC, server_name ASC) as rank
		FROM PlayerScores
		ORDER BY %s %s, name ASC, server_name ASC
		LIMIT ?`,
		orderBy.Column, orderBy.Direction,
		orderBy.Column, orderBy.Direction,
	)
}

// GetGlobalLeaderboard retrieves the global leaderboard with sorting options
func (s *RankingsService) GetGlobalLeaderboard(ctx context.Context, limit int, orderBy string, direction OrderDirection) ([]LeaderboardEntry, error) {
	sanitizedOrder := s.sanitizeOrderBy(orderBy, direction)
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(sanitizedOrder),
		"", // No additional WHERE conditions
		REQUIRED_DUNGEON_COUNT,
	)

	var entries []LeaderboardEntry
	err := s.db.WithContext(ctx).Raw(rankQuery, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get global leaderboard: %w", err)
	}

	return entries, nil
}

// GetGlobalLeaderboardByRole retrieves the leaderboard filtered by role
func (s *RankingsService) GetGlobalLeaderboardByRole(ctx context.Context, role string, limit int, orderBy string, direction OrderDirection) ([]LeaderboardEntry, error) {
	formattedRole := formatRole(role)
	sanitizedOrder := s.sanitizeOrderBy(orderBy, direction)

	whereClause := "AND role = ?"
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(sanitizedOrder),
		whereClause,
		REQUIRED_DUNGEON_COUNT,
	)

	var entries []LeaderboardEntry
	err := s.db.WithContext(ctx).Raw(rankQuery, formattedRole, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get role leaderboard: %w", err)
	}

	return entries, nil
}

// GetGlobalLeaderboardByClass retrieves the leaderboard filtered by class
func (s *RankingsService) GetGlobalLeaderboardByClass(ctx context.Context, class string, limit int, orderBy string, direction OrderDirection) ([]LeaderboardEntry, error) {
	formattedClass := formatNameCase(class)
	sanitizedOrder := s.sanitizeOrderBy(orderBy, direction)

	whereClause := "AND class = ?"
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(sanitizedOrder),
		whereClause,
		REQUIRED_DUNGEON_COUNT,
	)

	var entries []LeaderboardEntry
	err := s.db.WithContext(ctx).Raw(rankQuery, formattedClass, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get class leaderboard: %w", err)
	}

	return entries, nil
}

// GetGlobalLeaderboardBySpec retrieves the leaderboard filtered by class and spec
func (s *RankingsService) GetGlobalLeaderboardBySpec(ctx context.Context, class, spec string, limit int, orderBy string, direction OrderDirection) ([]LeaderboardEntry, error) {
	formattedClass := formatNameCase(class)
	formattedSpec := formatNameCase(spec)
	sanitizedOrder := s.sanitizeOrderBy(orderBy, direction)

	whereClause := "AND class = ? AND spec = ?"
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(sanitizedOrder),
		whereClause,
		REQUIRED_DUNGEON_COUNT,
	)

	var entries []LeaderboardEntry
	err := s.db.WithContext(ctx).Raw(rankQuery, formattedClass, formattedSpec, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get spec leaderboard: %w", err)
	}

	return entries, nil
}

// validateInput validates input parameters for the leaderboard queries
func (s *RankingsService) validateInput(role, class, spec string, limit int) error {
	if role != "" {
		validRoles := map[string]bool{
			"Tank":   true,
			"Healer": true,
			"DPS":    true,
		}
		formattedRole := formatRole(role)
		if !validRoles[formattedRole] {
			return fmt.Errorf("invalid role: %s", role)
		}
	}

	if class != "" {
		validClasses := map[string]bool{
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
		}
		formattedClass := formatNameCase(class)
		if !validClasses[formattedClass] {
			return fmt.Errorf("invalid class: %s", class)
		}
	}

	if limit <= 0 || limit > 1000 {
		return fmt.Errorf("invalid limit: must be between 1 and 1000")
	}

	return nil
}
