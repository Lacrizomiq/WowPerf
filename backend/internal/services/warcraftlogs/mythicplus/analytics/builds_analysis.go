package WarcraftLogsMythicPlusBuildAnalysis

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

// BuildAnalysisService handles analysis for build optimization and statistics
type BuildAnalysisService struct {
	db *gorm.DB
}

// NewBuildAnalysisService creates a new BuildAnalysisService
func NewBuildAnalysisService(db *gorm.DB) *BuildAnalysisService {
	return &BuildAnalysisService{db: db}
}

// ItemPopularity represents statistics about item popularity per slot
type ItemPopularity struct {
	EncounterID      int     `json:"encounter_id"`
	ItemSlot         int     `json:"item_slot"`
	ItemID           int     `json:"item_id"`
	ItemName         string  `json:"item_name"`
	ItemIcon         string  `json:"item_icon"`
	ItemQuality      int     `json:"item_quality"`
	ItemLevel        float64 `json:"item_level"`
	UsageCount       int     `json:"usage_count"`
	UsagePercentage  float64 `json:"usage_percentage"`
	AvgKeystoneLevel float64 `json:"avg_keystone_level"`
	Rank             int64   `json:"rank"`
}

// GetPopularItemsBySlot retrieves the most popular items for each slot for a specific class and spec
func (s *BuildAnalysisService) GetPopularItemsBySlot(ctx context.Context, class, spec string) ([]ItemPopularity, error) {
	var items []ItemPopularity
	err := s.db.WithContext(ctx).Raw("SELECT * FROM get_popular_items_by_slot(?, ?)", class, spec).Scan(&items).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get popular items: %w", err)
	}
	return items, nil
}

// EnchantUsage represents statistics about enchant usage per slot
type EnchantUsage struct {
	ItemSlot             int     `json:"item_slot"`
	PermanentEnchantID   int     `json:"permanent_enchant_id"`
	PermanentEnchantName string  `json:"permanent_enchant_name"`
	UsageCount           int64   `json:"usage_count"`
	AvgKeystoneLevel     float64 `json:"avg_keystone_level"`
	AvgItemLevel         float64 `json:"avg_item_level"`
	MaxKeystoneLevel     int64   `json:"max_keystone_level"`
	Rank                 int64   `json:"rank"`
}

// GetEnchantUsage retrieves enchant usage statistics for a specific class and spec
func (s *BuildAnalysisService) GetEnchantUsage(ctx context.Context, class, spec string) ([]EnchantUsage, error) {
	var enchants []EnchantUsage
	err := s.db.WithContext(ctx).Raw("SELECT * FROM get_enchant_usage(?, ?)", class, spec).Scan(&enchants).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get enchant usage: %w", err)
	}
	return enchants, nil
}

// GemUsage represents statistics about gem usage per slot
type GemUsage struct {
	ItemSlot         int       `json:"item_slot"`
	GemsCount        int       `json:"gems_count"`
	GemIDsArray      []int     `json:"gem_ids_array"`
	GemIconsArray    []string  `json:"gem_icons_array"`
	GemLevelsArray   []float64 `json:"gem_levels_array"`
	UsageCount       int64     `json:"usage_count"`
	AvgKeystoneLevel float64   `json:"avg_keystone_level"`
	AvgItemLevel     float64   `json:"avg_item_level"`
	Rank             int64     `json:"rank"`
}

// GetGemUsage retrieves gem usage statistics for a specific class and spec
func (s *BuildAnalysisService) GetGemUsage(ctx context.Context, class, spec string) ([]GemUsage, error) {
	var gems []GemUsage
	err := s.db.WithContext(ctx).Raw("SELECT * FROM get_gem_usage(?, ?)", class, spec).Scan(&gems).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get gem usage: %w", err)
	}
	return gems, nil
}

// TalentBuild represents statistics about a talent build
type TalentBuild struct {
	TalentImport       string  `json:"talent_import"`
	TotalUsage         int64   `json:"total_usage"`
	AvgUsagePercentage float64 `json:"avg_usage_percentage"`
	AvgKeystoneLevel   float64 `json:"avg_keystone_level"`
}

// GetTopTalentBuilds retrieves the top talent builds for a specific class and spec
func (s *BuildAnalysisService) GetTopTalentBuilds(ctx context.Context, class, spec string) ([]TalentBuild, error) {
	var builds []TalentBuild
	err := s.db.WithContext(ctx).Raw("SELECT * FROM get_top_talent_builds(?, ?)", class, spec).Scan(&builds).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get top talent builds: %w", err)
	}
	return builds, nil
}

// DungeonTalentBuild represents statistics about talent builds per dungeon
type DungeonTalentBuild struct {
	Class              string  `json:"class"`
	Spec               string  `json:"spec"`
	EncounterID        int     `json:"encounter_id"`
	DungeonName        string  `json:"dungeon_name"`
	TalentImport       string  `json:"talent_import"`
	TotalUsage         int64   `json:"total_usage"`
	AvgUsagePercentage float64 `json:"avg_usage_percentage"`
	AvgKeystoneLevel   float64 `json:"avg_keystone_level"`
}

// GetTalentBuildsByDungeon retrieves talent build statistics per dungeon for a specific class and spec
func (s *BuildAnalysisService) GetTalentBuildsByDungeon(ctx context.Context, class, spec string) ([]DungeonTalentBuild, error) {
	var builds []DungeonTalentBuild
	err := s.db.WithContext(ctx).Raw("SELECT * FROM get_talent_builds_by_dungeon(?, ?)", class, spec).Scan(&builds).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get talent builds by dungeon: %w", err)
	}
	return builds, nil
}

// StatPriority represents statistics about stat priorities
type StatPriority struct {
	StatName         string  `json:"stat_name"`
	StatCategory     string  `json:"stat_category"`
	AvgValue         float64 `json:"avg_value"`
	MinValue         float64 `json:"min_value"`
	MaxValue         float64 `json:"max_value"`
	TotalSamples     int64   `json:"total_samples"`
	AvgKeystoneLevel float64 `json:"avg_keystone_level"`
	PriorityRank     int64   `json:"priority_rank"`
}

// GetStatPriorities retrieves stat priority statistics for a specific class and spec
func (s *BuildAnalysisService) GetStatPriorities(ctx context.Context, class, spec string) ([]StatPriority, error) {
	var stats []StatPriority
	err := s.db.WithContext(ctx).Raw("SELECT * FROM get_stat_priorities(?, ?)", class, spec).Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get stat priorities: %w", err)
	}
	return stats, nil
}

// ItemDetails represents details about a specific item in the optimal build
type ItemDetails struct {
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Quality    int    `json:"quality"`
	UsageCount int    `json:"usage_count"`
}

// OptimalBuild represents the optimal build for a specific class and spec
type OptimalBuild struct {
	TopTalentImport string                 `json:"top_talent_import"`
	StatPriority    string                 `json:"stat_priority"`
	TopItems        map[string]ItemDetails `json:"top_items"`
}

// GetOptimalBuild retrieves the optimal build for a specific class and spec
func (s *BuildAnalysisService) GetOptimalBuild(ctx context.Context, class, spec string) (*OptimalBuild, error) {
	// This is a special case because it returns JSON in one of the columns
	var result struct {
		TopTalentImport string      `json:"top_talent_import"`
		StatPriority    string      `json:"stat_priority"`
		TopItems        interface{} `json:"top_items"` // Start with interface{} for the JSON column
	}

	err := s.db.WithContext(ctx).Raw("SELECT * FROM get_optimal_build(?, ?)", class, spec).Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get optimal build: %w", err)
	}

	// Create the return object
	build := &OptimalBuild{
		TopTalentImport: result.TopTalentImport,
		StatPriority:    result.StatPriority,
		TopItems:        make(map[string]ItemDetails),
	}

	// Handle JSON to map conversion
	jsonData, err := json.Marshal(result.TopItems)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := json.Unmarshal(jsonData, &build.TopItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return build, nil
}

// SpecComparison represents comparison statistics between specs
type SpecComparison struct {
	Spec             string  `json:"spec"`
	AvgKeystoneLevel float64 `json:"avg_keystone_level"`
	MaxKeystoneLevel int64   `json:"max_keystone_level"`
	AvgItemLevel     float64 `json:"avg_item_level"`
	DungeonsCount    int64   `json:"dungeons_count"`
	StatPriority     string  `json:"stat_priority"`
}

// GetSpecComparison retrieves comparison statistics for all specs of a class
func (s *BuildAnalysisService) GetSpecComparison(ctx context.Context, class string) ([]SpecComparison, error) {
	var specs []SpecComparison
	err := s.db.WithContext(ctx).Raw("SELECT * FROM get_spec_comparison(?)", class).Scan(&specs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get spec comparison: %w", err)
	}
	return specs, nil
}

// ClassSpecSummary represents a summary of class and spec performance
type ClassSpecSummary struct {
	AvgKeystoneLevel float64 `json:"avg_keystone_level"`
	MaxKeystoneLevel int64   `json:"max_keystone_level"`
	AvgItemLevel     float64 `json:"avg_item_level"`
	TopTalentImport  string  `json:"top_talent_import"`
	StatPriority     string  `json:"stat_priority"`
	DungeonsCount    int64   `json:"dungeons_count"`
}

// GetClassSpecSummary retrieves summary statistics for a specific class and spec
func (s *BuildAnalysisService) GetClassSpecSummary(ctx context.Context, class, spec string) (*ClassSpecSummary, error) {
	var summary ClassSpecSummary
	err := s.db.WithContext(ctx).Raw("SELECT * FROM get_class_spec_summary(?, ?)", class, spec).Scan(&summary).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get class spec summary: %w", err)
	}
	return &summary, nil
}
