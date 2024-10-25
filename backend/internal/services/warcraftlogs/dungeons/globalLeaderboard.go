package warcraftlogs

import (
	"context"
)

type LeaderboardEntry struct {
	PlayerID   int     `json:"player_id"`
	Name       string  `json:"name"`
	Class      string  `json:"class"`
	Spec       string  `json:"spec"`
	Role       string  `json:"role"`
	TotalScore float64 `json:"total_score"`
	Rank       int     `json:"rank"`
}

// Get the global leaderboard by role
func (s *RankingsService) GetGlobalLeaderboardByRole(ctx context.Context, role string, limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry

	rankQuery := `
			SELECT 
					player_id,
					name,
					class,
					spec,
					role,
					SUM(score) as total_score,
					DENSE_RANK() OVER (ORDER BY SUM(score) DESC, name ASC) as rank
			FROM player_rankings
			WHERE role = ?
			GROUP BY player_id, name, class, spec, role
			ORDER BY total_score DESC, name ASC
			LIMIT ?
	`

	err := s.db.Raw(rankQuery, role, limit).Scan(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// Get the global leaderboard by class
func (s *RankingsService) GetGlobalLeaderboardByClass(ctx context.Context, class string, limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry

	rankQuery := `
			SELECT 
					player_id,
					name,
					class,
					spec,
					role,
					SUM(score) as total_score,
					DENSE_RANK() OVER (ORDER BY SUM(score) DESC, name ASC) as rank
			FROM player_rankings
			WHERE class = ?
			GROUP BY player_id, name, class, spec, role
			ORDER BY total_score DESC, name ASC
			LIMIT ?
	`

	err := s.db.Raw(rankQuery, class, limit).Scan(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// Get the global leaderboard by spec
func (s *RankingsService) GetGlobalLeaderboardBySpec(ctx context.Context, spec string, limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry

	rankQuery := `
			SELECT 
					player_id,
					name,
					class,
					spec,
					role,
					SUM(score) as total_score,
					DENSE_RANK() OVER (ORDER BY SUM(score) DESC, name ASC) as rank
			FROM player_rankings
			WHERE spec = ?
			GROUP BY player_id, name, class, spec, role
			ORDER BY total_score DESC, name ASC
			LIMIT ?
	`

	err := s.db.Raw(rankQuery, spec, limit).Scan(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// Get the global leaderboard in every role
func (s *RankingsService) GetGlobalLeaderboard(ctx context.Context, limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry

	rankQuery := `
			SELECT 
					player_id,
					name,
					class,
					spec,
					role,
					SUM(score) as total_score,
					DENSE_RANK() OVER (ORDER BY SUM(score) DESC, name ASC) as rank
			FROM player_rankings
			GROUP BY player_id, name, class, spec, role
			ORDER BY total_score DESC, name ASC
			LIMIT ?
	`

	err := s.db.Raw(rankQuery, limit).Scan(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}
