package warcraftlogs

import (
	"context"
	"fmt"
	"strings"
)

const REQUIRED_DUNGEON_COUNT = 8

type LeaderboardEntry struct {
	PlayerID     int     `json:"player_id"`
	Name         string  `json:"name"`
	Class        string  `json:"class"`
	Spec         string  `json:"spec"`
	Role         string  `json:"role"`
	TotalScore   float64 `json:"total_score"`
	Rank         int     `json:"rank"`
	DungeonCount int     `json:"dungeon_count"`
}

// Base query helper
func (s *RankingsService) getBaseLeaderboardQuery() string {
	return `
        SELECT 
            player_id,
            name,
            class,
            spec,
            role,
            ROUND(CAST(SUM(score) AS numeric), 2) as total_score,
            COUNT(DISTINCT dungeon_id) as dungeon_count,
            DENSE_RANK() OVER (ORDER BY SUM(score) DESC, name ASC) as rank
        FROM player_rankings
        %s  -- WHERE clause placeholder
        GROUP BY player_id, name, class, spec, role
        HAVING COUNT(DISTINCT dungeon_id) = %d
        ORDER BY total_score DESC, name ASC
        LIMIT ?
    `
}

// Get the global leaderboard in every role
func (s *RankingsService) GetGlobalLeaderboard(ctx context.Context, limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(),
		"", // No WHERE clause needed
		REQUIRED_DUNGEON_COUNT,
	)

	err := s.db.Raw(rankQuery, limit).Scan(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// Get the global leaderboard by role
func (s *RankingsService) GetGlobalLeaderboardByRole(ctx context.Context, role string, limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	whereClause := "WHERE LOWER(role) = LOWER(?)"
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(),
		whereClause,
		REQUIRED_DUNGEON_COUNT,
	)

	err := s.db.Raw(rankQuery, role, limit).Scan(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// Get the global leaderboard by class
func (s *RankingsService) GetGlobalLeaderboardByClass(ctx context.Context, class string, limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	whereClause := "WHERE LOWER(class) = LOWER(?)"
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(),
		whereClause,
		REQUIRED_DUNGEON_COUNT,
	)

	err := s.db.Raw(rankQuery, class, limit).Scan(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// Get the global leaderboard by spec
func (s *RankingsService) GetGlobalLeaderboardBySpec(ctx context.Context, class, spec string, limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry
	whereClause := "WHERE LOWER(class) = LOWER(?) AND LOWER(spec) = LOWER(?)"
	rankQuery := fmt.Sprintf(
		s.getBaseLeaderboardQuery(),
		whereClause,
		REQUIRED_DUNGEON_COUNT,
	)

	err := s.db.Raw(rankQuery, class, spec, limit).Scan(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// function to validate the input
func (s *RankingsService) validateInput(role, class, spec string, limit int) error {
	// Validate the role if it is provided
	roleLower := strings.ToLower(role)

	if role != "" {
		validRoles := map[string]bool{"tank": true, "healer": true, "dps": true}
		if !validRoles[roleLower] {
			return fmt.Errorf("invalid role: %s", role)
		}
	}

	// Validate the class if it is provided
	if class != "" {
		validClasses := map[string]bool{
			"warrior": true, "paladin": true, "hunter": true, "rogue": true,
			"priest": true, "shaman": true, "mage": true, "warlock": true,
			"monk": true, "druid": true, "demonhunter": true, "deathknight": true,
			"evoker": true,
		}
		if !validClasses[class] {
			return fmt.Errorf("invalid class: %s", class)
		}
	}

	// Validate the limit
	if limit <= 0 || limit > 1000 {
		return fmt.Errorf("invalid limit: must be between 1 and 1000")
	}

	return nil
}

// Helper function to sanitize the order by columns
func (s *RankingsService) sanitizeOrderBy(column string) string {
	validColumns := map[string]string{
		"score": "score",
		"name":  "name",
		"rank":  "rank",
	}

	if sanitized, ok := validColumns[column]; ok {
		return sanitized
	}
	return "score" // default
}
