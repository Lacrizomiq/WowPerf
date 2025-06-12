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

// GetPopularItemsBySlot retrieves the most popular items for each slot for a specific class, spec and encounter
func (s *BuildAnalysisService) GetPopularItemsBySlot(ctx context.Context, class, spec string, encounterID *int) ([]ItemPopularity, error) {
	var items []ItemPopularity
	var err error

	if encounterID != nil {
		err = s.db.WithContext(ctx).Raw("SELECT * FROM get_popular_items_by_slot(?, ?, ?)", class, spec, encounterID).Scan(&items).Error
	} else {
		err = s.db.WithContext(ctx).Raw("SELECT * FROM get_popular_items_by_slot(?, ?, NULL)", class, spec).Scan(&items).Error
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get popular items: %w", err)
	}
	return items, nil
}

// GlobalItemPopularity represents global statistics about item popularity per slot
type GlobalItemPopularity struct {
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

// GetGlobalPopularItemsBySlot retrieves the most popular items for each slot across all encounters
func (s *BuildAnalysisService) GetGlobalPopularItemsBySlot(ctx context.Context, class, spec string) ([]GlobalItemPopularity, error) {
	var items []GlobalItemPopularity
	err := s.db.WithContext(ctx).Raw("SELECT * FROM get_global_popular_items_by_slot(?, ?)", class, spec).Scan(&items).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get global popular items: %w", err)
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
// Using raw SQL query to avoid deserialization issues with GORM
// Deserialization issues are due to the fact that GORM is not able to deserialize the arrays correctly
func (s *BuildAnalysisService) GetGemUsage(ctx context.Context, class, spec string) ([]GemUsage, error) {
	// Query to explicitly convert arrays to JSON
	query := `
	SELECT 
					bs.item_slot,
					bs.gems_count,
					array_to_json(COALESCE(bs.gem_ids, ARRAY[]::INTEGER[])) as gem_ids_array,
					array_to_json(COALESCE(bs.gem_icons, ARRAY[]::TEXT[])) as gem_icons_array,
					array_to_json(COALESCE(bs.gem_levels, ARRAY[]::NUMERIC[])) as gem_levels_array,
					COUNT(*)::BIGINT as usage_count,
					AVG(bs.avg_keystone_level) as avg_keystone_level,
					AVG(bs.avg_item_level) as avg_item_level,
					ROW_NUMBER() OVER (PARTITION BY bs.item_slot ORDER BY COUNT(*) DESC)::BIGINT as rank
	FROM build_statistics bs
	WHERE bs.class = $1 AND bs.spec = $2
	AND bs.has_gems = true
	GROUP BY bs.item_slot, bs.gems_count, bs.gem_ids, bs.gem_icons, bs.gem_levels
	ORDER BY bs.item_slot, ROW_NUMBER() OVER (PARTITION BY bs.item_slot ORDER BY COUNT(*) DESC)
	`

	// Get the underlying DB connection
	sqlDB, err := s.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB connection: %w", err)
	}

	// Use the standard SQL database
	rows, err := sqlDB.QueryContext(ctx, query, class, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to get gem usage: %w", err)
	}
	defer rows.Close()

	var results []GemUsage
	for rows.Next() {
		var g GemUsage
		// Temporary variables to store JSON arrays
		var gemIDsJSON, gemIconsJSON, gemLevelsJSON string

		err := rows.Scan(
			&g.ItemSlot,
			&g.GemsCount,
			&gemIDsJSON,
			&gemIconsJSON,
			&gemLevelsJSON,
			&g.UsageCount,
			&g.AvgKeystoneLevel,
			&g.AvgItemLevel,
			&g.Rank,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan gem usage: %w", err)
		}

		// Deserialize JSON to Go arrays
		if gemIDsJSON != "[]" && gemIDsJSON != "" {
			if err := json.Unmarshal([]byte(gemIDsJSON), &g.GemIDsArray); err != nil {
				return nil, fmt.Errorf("failed to unmarshal gem_ids: %w", err)
			}
		} else {
			g.GemIDsArray = []int{}
		}

		if gemIconsJSON != "[]" && gemIconsJSON != "" {
			if err := json.Unmarshal([]byte(gemIconsJSON), &g.GemIconsArray); err != nil {
				return nil, fmt.Errorf("failed to unmarshal gem_icons: %w", err)
			}
		} else {
			g.GemIconsArray = []string{}
		}

		if gemLevelsJSON != "[]" && gemLevelsJSON != "" {
			if err := json.Unmarshal([]byte(gemLevelsJSON), &g.GemLevelsArray); err != nil {
				return nil, fmt.Errorf("failed to unmarshal gem_levels: %w", err)
			}
		} else {
			g.GemLevelsArray = []float64{}
		}

		results = append(results, g)
	}

	// Check errors after iterating through results
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
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
	// 1. Get the most popular talent import
	var topTalent struct {
		TalentImport string
	}
	talentQuery := `
	SELECT talent_import
	FROM talent_statistics
	WHERE class = ? AND spec = ?
	GROUP BY talent_import
	ORDER BY SUM(usage_count) DESC
	LIMIT 1
	`
	err := s.db.WithContext(ctx).Raw(talentQuery, class, spec).Scan(&topTalent).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get top talent import: %w", err)
	}

	// 2. Get the stat priority
	var statPriority struct {
		StatPriority string
	}
	statQuery := `
	SELECT STRING_AGG(stat_name, ' > ' ORDER BY avg_value DESC) as stat_priority
	FROM (
			SELECT stat_name, AVG(avg_value) as avg_value
			FROM stat_statistics
			WHERE class = ? AND spec = ? AND stat_category = 'secondary'
			GROUP BY stat_name
	) s
	`
	err = s.db.WithContext(ctx).Raw(statQuery, class, spec).Scan(&statPriority).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get stat priority: %w", err)
	}

	// 3. Get the most popular items
	type ItemResult struct {
		ItemSlot    int    `json:"item_slot"`
		ItemName    string `json:"item_name"`
		ItemIcon    string `json:"item_icon"`
		ItemQuality int    `json:"item_quality"`
		UsageCount  int    `json:"usage_count"`
	}

	var items []ItemResult
	itemQuery := `
	SELECT DISTINCT ON (item_slot)
			item_slot, item_name, item_icon, item_quality, usage_count
	FROM build_statistics
	WHERE class = ? AND spec = ?
	ORDER BY item_slot, usage_count DESC
	`
	err = s.db.WithContext(ctx).Raw(itemQuery, class, spec).Scan(&items).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get top items: %w", err)
	}

	// Build the OptimalBuild object
	build := &OptimalBuild{
		TopTalentImport: topTalent.TalentImport,
		StatPriority:    statPriority.StatPriority,
		TopItems:        make(map[string]ItemDetails),
	}

	// Fill TopItems
	for _, item := range items {
		build.TopItems[fmt.Sprintf("%d", item.ItemSlot)] = ItemDetails{
			Name:       item.ItemName,
			Icon:       item.ItemIcon,
			Quality:    item.ItemQuality,
			UsageCount: item.UsageCount,
		}
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
