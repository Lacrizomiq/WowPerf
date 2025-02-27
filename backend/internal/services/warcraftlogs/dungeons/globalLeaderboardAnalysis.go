package warcraftlogs

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// GlobalLeaderboardAnalysisService handles analysis-specific queries for high-key M+ data
type GlobalLeaderboardAnalysisService struct {
	db *gorm.DB
}

// NewGlobalLeaderboardAnalysisService creates a new instance of GlobalLeaderboardAnalysisService
func NewGlobalLeaderboardAnalysisService(db *gorm.DB) *GlobalLeaderboardAnalysisService {
	return &GlobalLeaderboardAnalysisService{db: db}
}

// SpecGlobalScore represents the average global score for a spec from the view
type SpecGlobalScore struct {
	Class          string  `json:"class"`
	Spec           string  `json:"spec"`
	AvgGlobalScore float64 `json:"avg_global_score"`
	PlayerCount    int     `json:"player_count"`
	Role           string  `json:"role"`
	OverallRank    int     `json:"overall_rank"`
	RoleRank       int     `json:"role_rank"`
}

// GetSpecGlobalScores retrieves the average global score per spec from the spec_global_score_averages view
func (s *GlobalLeaderboardAnalysisService) GetSpecGlobalScores(ctx context.Context) ([]SpecGlobalScore, error) {
	var results []SpecGlobalScore
	err := s.db.WithContext(ctx).Raw("SELECT * FROM spec_global_score_averages").Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get spec global scores: %w", err)
	}
	return results, nil
}

// ClassGlobalScore represents the average global score for a class from the view
type ClassGlobalScore struct {
	Class          string  `json:"class"`
	AvgGlobalScore float64 `json:"avg_global_score"`
	PlayerCount    int     `json:"player_count"`
}

// GetClassGlobalScores retrieves the average global score per class from the class_global_score_averages view
func (s *GlobalLeaderboardAnalysisService) GetClassGlobalScores(ctx context.Context) ([]ClassGlobalScore, error) {
	var results []ClassGlobalScore
	err := s.db.WithContext(ctx).Raw("SELECT * FROM class_global_score_averages").Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get class global scores: %w", err)
	}
	return results, nil
}

// SpecDungeonMaxKeyLevel represents the max key level for a spec per dungeon from the view
type SpecDungeonMaxKeyLevel struct {
	Class       string `json:"class"`
	Spec        string `json:"spec"`
	DungeonName string `json:"dungeon_name"`
	DungeonSlug string `json:"dungeon_slug"`
	MaxKeyLevel int    `json:"max_key_level"`
}

// GetSpecDungeonMaxKeyLevels retrieves the max key levels per spec and dungeon from the spec_dungeon_max_key_levels view
func (s *GlobalLeaderboardAnalysisService) GetSpecDungeonMaxKeyLevels(ctx context.Context) ([]SpecDungeonMaxKeyLevel, error) {
	var results []SpecDungeonMaxKeyLevel
	err := s.db.WithContext(ctx).Raw("SELECT * FROM spec_dungeon_max_key_levels").Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get spec dungeon max key levels: %w", err)
	}
	return results, nil
}

// DungeonMedia
type DungeonMedia struct {
	Slug        string `json:"dungeon_slug"`
	MediaURL    string `json:"media_url"`
	Icon        string `json:"icon"`
	EncounterID int    `json:"encounter_id"`
}

// GetDungeonMedia retrieves the dungeons media
func (s *GlobalLeaderboardAnalysisService) GetDungeonMedia(ctx context.Context) ([]DungeonMedia, error) {
	var results []DungeonMedia
	err := s.db.WithContext(ctx).Raw("SELECT slug, media_url, icon, encounter_id FROM dungeons WHERE encounter_id IS NOT NULL").Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get dungeon media: %w", err)
	}
	return results, nil
}

// DungeonAvgKeyLevel represents the average key level for a dungeon from the view
type DungeonAvgKeyLevel struct {
	DungeonName string  `json:"dungeon_name"`
	DungeonSlug string  `json:"dungeon_slug"`
	AvgKeyLevel float64 `json:"avg_key_level"`
	RunCount    int     `json:"run_count"`
}

// GetDungeonAvgKeyLevels retrieves the average key levels per dungeon from the dungeon_avg_key_levels view
func (s *GlobalLeaderboardAnalysisService) GetDungeonAvgKeyLevels(ctx context.Context) ([]DungeonAvgKeyLevel, error) {
	var results []DungeonAvgKeyLevel
	err := s.db.WithContext(ctx).Raw("SELECT * FROM dungeon_avg_key_levels").Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get dungeon avg key levels: %w", err)
	}
	return results, nil
}

// TopPlayerPerSpec represents a top player per spec from the view
type TopPlayerPerSpec struct {
	Class        string  `json:"class"`
	Spec         string  `json:"spec"`
	Name         string  `json:"name"`
	ServerName   string  `json:"server_name"`
	ServerRegion string  `json:"server_region"`
	TotalScore   float64 `json:"total_score"`
	Rank         int     `json:"rank"`
}

// GetTop10PlayersPerSpec retrieves the top 10 players per spec from the top_10_players_per_spec view
func (s *GlobalLeaderboardAnalysisService) GetTop10PlayersPerSpec(ctx context.Context) ([]TopPlayerPerSpec, error) {
	var results []TopPlayerPerSpec
	err := s.db.WithContext(ctx).Raw("SELECT * FROM top_10_players_per_spec").Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get top 10 players per spec: %w", err)
	}
	return results, nil
}

// TopPlayerPerRole represents a top player per role from the view
type TopPlayerPerRole struct {
	Name         string  `json:"name"`
	ServerName   string  `json:"server_name"`
	ServerRegion string  `json:"server_region"`
	Class        string  `json:"class"`
	Spec         string  `json:"spec"`
	Role         string  `json:"role"`
	TotalScore   float64 `json:"total_score"`
	Rank         int     `json:"rank"`
}

// GetTop5PlayersPerRole retrieves the top 5 players per role from the top_5_players_per_role view
func (s *GlobalLeaderboardAnalysisService) GetTop5PlayersPerRole(ctx context.Context) ([]TopPlayerPerRole, error) {
	var results []TopPlayerPerRole
	err := s.db.WithContext(ctx).Raw("SELECT * FROM top_5_players_per_role").Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get top 5 players per role: %w", err)
	}
	return results, nil
}
