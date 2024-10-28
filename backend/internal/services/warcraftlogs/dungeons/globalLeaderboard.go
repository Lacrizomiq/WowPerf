package warcraftlogs

import (
	"context"
	"fmt"
	"strings"
)

const REQUIRED_DUNGEON_COUNT = 8

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

// OrderDirection represents the sort direction
type OrderDirection string

const (
	ASC  OrderDirection = "ASC"
	DESC OrderDirection = "DESC"
)

// OrderByOption represents a column and its sort direction
type OrderByOption struct {
	Column    string
	Direction OrderDirection
}

// Création d'index optimisés pour les nouvelles requêtes
func (s *RankingsService) createOptimizedIndexes() error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_rankings_player_unique ON player_rankings(name, server_name, server_region)",
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

// Helper function to sanitize the order by columns with more options
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

	return OrderByOption{
		Column:    "total_score",
		Direction: DESC,
	}
}

// Helper to build ORDER BY clause
func (s *RankingsService) buildOrderByClause(options ...OrderByOption) string {
	if len(options) == 0 {
		return "ORDER BY total_score DESC, name ASC"
	}

	orderByClause := "ORDER BY "
	for i, opt := range options {
		if i > 0 {
			orderByClause += ", "
		}
		orderByClause += fmt.Sprintf("%s %s", opt.Column, opt.Direction)
	}
	return orderByClause
}

// Base query avec support du tri
func (s *RankingsService) getBaseLeaderboardQuery(orderBy ...OrderByOption) string {
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
			%%s  -- WHERE clause placeholder
			GROUP BY name, server_name, server_region, class, spec, role
			HAVING COUNT(DISTINCT dungeon_id) = %%d
		)
		SELECT 
			*,
			DENSE_RANK() OVER (%s) as rank
		FROM PlayerScores
		%s
		LIMIT ?
	`,
		s.buildOrderByClause(OrderByOption{Column: "total_score", Direction: DESC}),
		s.buildOrderByClause(orderBy...))
}

// GetGlobalLeaderboard avec support du tri
func (s *RankingsService) GetGlobalLeaderboard(ctx context.Context, limit int, orderBy string, direction OrderDirection) ([]LeaderboardEntry, error) {
	orderOption := s.sanitizeOrderBy(orderBy, direction)
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(orderOption),
		"",
		REQUIRED_DUNGEON_COUNT,
	)

	var entries []LeaderboardEntry
	err := s.db.WithContext(ctx).Raw(rankQuery, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get global leaderboard: %w", err)
	}

	return entries, nil
}

// GetGlobalLeaderboardByRole avec support du tri
func (s *RankingsService) GetGlobalLeaderboardByRole(ctx context.Context, role string, limit int, orderBy string, direction OrderDirection) ([]LeaderboardEntry, error) {
	orderOption := s.sanitizeOrderBy(orderBy, direction)
	whereClause := "WHERE LOWER(role) = LOWER(?)"
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(orderOption),
		whereClause,
		REQUIRED_DUNGEON_COUNT,
	)

	var entries []LeaderboardEntry
	err := s.db.WithContext(ctx).Raw(rankQuery, role, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get role leaderboard: %w", err)
	}

	return entries, nil
}

// GetGlobalLeaderboardByClass avec support du tri
func (s *RankingsService) GetGlobalLeaderboardByClass(ctx context.Context, class string, limit int, orderBy string, direction OrderDirection) ([]LeaderboardEntry, error) {
	orderOption := s.sanitizeOrderBy(orderBy, direction)
	whereClause := "WHERE LOWER(class) = LOWER(?)"
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(orderOption),
		whereClause,
		REQUIRED_DUNGEON_COUNT,
	)

	var entries []LeaderboardEntry
	err := s.db.WithContext(ctx).Raw(rankQuery, class, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get class leaderboard: %w", err)
	}

	return entries, nil
}

// GetGlobalLeaderboardBySpec avec support du tri
func (s *RankingsService) GetGlobalLeaderboardBySpec(ctx context.Context, class, spec string, limit int, orderBy string, direction OrderDirection) ([]LeaderboardEntry, error) {
	orderOption := s.sanitizeOrderBy(orderBy, direction)
	whereClause := "WHERE LOWER(class) = LOWER(?) AND LOWER(spec) = LOWER(?)"
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(orderOption),
		whereClause,
		REQUIRED_DUNGEON_COUNT,
	)

	var entries []LeaderboardEntry
	err := s.db.WithContext(ctx).Raw(rankQuery, class, spec, limit).Scan(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get spec leaderboard: %w", err)
	}

	return entries, nil
}

// Validation des entrées
func (s *RankingsService) validateInput(role, class, spec string, limit int) error {
	if role != "" {
		validRoles := map[string]bool{"tank": true, "healer": true, "dps": true}
		if !validRoles[strings.ToLower(role)] {
			return fmt.Errorf("invalid role: %s", role)
		}
	}

	if class != "" {
		validClasses := map[string]bool{
			"warrior": true, "paladin": true, "hunter": true, "rogue": true,
			"priest": true, "shaman": true, "mage": true, "warlock": true,
			"monk": true, "druid": true, "demonhunter": true, "deathknight": true,
			"evoker": true,
		}
		if !validClasses[strings.ToLower(class)] {
			return fmt.Errorf("invalid class: %s", class)
		}
	}

	if limit <= 0 || limit > 1000 {
		return fmt.Errorf("invalid limit: must be between 1 and 1000")
	}

	return nil
}
