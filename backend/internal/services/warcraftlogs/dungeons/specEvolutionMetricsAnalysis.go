// backend/internal/services/warcraftlogs/dungeons/specEvolutionMetricsAnalysis.go
package warcraftlogs

import (
	"context"
	"fmt"
	"time"

	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"

	"gorm.io/gorm"
)

// SpecEvolutionMetricsAnalysisService handles analysis-specific queries for spec evolution metrics
type SpecEvolutionMetricsAnalysisService struct {
	db *gorm.DB
}

// NewSpecEvolutionMetricsAnalysisService creates a new instance of SpecEvolutionMetricsAnalysisService
func NewSpecEvolutionMetricsAnalysisService(db *gorm.DB) *SpecEvolutionMetricsAnalysisService {
	return &SpecEvolutionMetricsAnalysisService{db: db}
}

// buildEvolutionQuery builds a SQL query for the evolution of metrics
func (s *SpecEvolutionMetricsAnalysisService) buildEvolutionQuery(daysInterval int, isGlobal bool, includeClassFilter bool) string {
	query := `
	SELECT
		current.capture_date AS end_date,
		current.capture_date - INTERVAL '%d days' AS start_date,
		current.spec,
		current.class,
		current.role,
		current.encounter_id,
		current.is_global,
		current.avg_score AS current_avg_score,
		prev.avg_score AS prev_avg_score,
		(current.avg_score - prev.avg_score) AS score_change,
		CASE 
			WHEN prev.avg_score IS NULL OR prev.avg_score = 0 THEN NULL 
			ELSE (current.avg_score - prev.avg_score) / prev.avg_score * 100 
		END AS score_change_percent,
		current.max_score AS current_max_score,
		prev.max_score AS prev_max_score,
		(current.max_score - prev.max_score) AS max_score_change,
		current.avg_key_level AS current_avg_key,
		prev.avg_key_level AS prev_avg_key,
		(current.avg_key_level - prev.avg_key_level) AS key_level_change,
		current.overall_rank AS current_overall_rank,
		prev.overall_rank AS prev_overall_rank,
		(prev.overall_rank - current.overall_rank) AS overall_rank_change,
		current.role_rank AS current_role_rank,
		prev.role_rank AS prev_role_rank,
		(prev.role_rank - current.role_rank) AS role_rank_change
	FROM
		daily_spec_metrics_mythic_plus current
	LEFT JOIN
		daily_spec_metrics_mythic_plus prev
		ON current.spec = prev.spec
		AND current.class = prev.class
		AND current.role = prev.role
		AND current.encounter_id = prev.encounter_id
		AND current.is_global = prev.is_global
		AND prev.capture_date = current.capture_date - INTERVAL '%d days'
	WHERE
		current.spec = ?
		AND current.capture_date = ?
		AND current.is_global = ?
	`

	if includeClassFilter {
		query += " AND current.class = ?"
	}

	if !isGlobal {
		query += " AND current.encounter_id = ?"
	}

	return fmt.Sprintf(query, daysInterval, daysInterval)
}

// GetSpecEvolution retrieves the evolution of metrics for a specialization
func (s *SpecEvolutionMetricsAnalysisService) GetSpecEvolution(ctx context.Context, spec string, class *string, daysInterval int, encounterID *int, targetDate time.Time, isGlobal bool) ([]playerRankingModels.SpecEvolutionMythicPlus, error) {
	includeClassFilter := class != nil
	query := s.buildEvolutionQuery(daysInterval, isGlobal, includeClassFilter)
	var results []playerRankingModels.SpecEvolutionMythicPlus
	var err error

	if isGlobal {
		if includeClassFilter {
			err = s.db.WithContext(ctx).Raw(query, spec, targetDate, isGlobal, *class).Scan(&results).Error
		} else {
			err = s.db.WithContext(ctx).Raw(query, spec, targetDate, isGlobal).Scan(&results).Error
		}
	} else {
		if includeClassFilter {
			err = s.db.WithContext(ctx).Raw(query, spec, targetDate, isGlobal, *class, *encounterID).Scan(&results).Error
		} else {
			err = s.db.WithContext(ctx).Raw(query, spec, targetDate, isGlobal, *encounterID).Scan(&results).Error
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get spec evolution: %w", err)
	}

	return results, nil
}

// GetCurrentRanking retrieves the current ranking of all specializations
func (s *SpecEvolutionMetricsAnalysisService) GetCurrentRanking(ctx context.Context, role *string, targetDate time.Time, isGlobal bool) ([]playerRankingModels.DailySpecMetricMythicPlus, error) {
	query := `
	SELECT * FROM daily_spec_metrics_mythic_plus 
	WHERE capture_date = ? AND is_global = ?
	`

	if role != nil {
		query += " AND role = ?"
	}

	query += " ORDER BY overall_rank"

	var results []playerRankingModels.DailySpecMetricMythicPlus
	var err error

	if role != nil {
		err = s.db.WithContext(ctx).Raw(query, targetDate, isGlobal, *role).Scan(&results).Error
	} else {
		err = s.db.WithContext(ctx).Raw(query, targetDate, isGlobal).Scan(&results).Error
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get current ranking: %w", err)
	}

	return results, nil
}

// GetLatestMetricsDate retrieves the latest date for which metrics are available
func (s *SpecEvolutionMetricsAnalysisService) GetLatestMetricsDate(ctx context.Context) (time.Time, error) {
	var result struct {
		LatestDate time.Time
	}

	err := s.db.WithContext(ctx).Raw(`
		SELECT MAX(capture_date) as latest_date 
		FROM daily_spec_metrics_mythic_plus
	`).Scan(&result).Error

	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get latest metrics date: %w", err)
	}

	// If no metrics are available yet, return current date
	if result.LatestDate.IsZero() {
		return time.Now().Truncate(24 * time.Hour), nil
	}

	return result.LatestDate, nil
}

// GetClassSpecs retrieves all available specs for a class
func (s *SpecEvolutionMetricsAnalysisService) GetClassSpecs(ctx context.Context, className string) ([]string, error) {
	var specs []string

	err := s.db.WithContext(ctx).Raw(`
		SELECT DISTINCT spec
		FROM daily_spec_metrics_mythic_plus
		WHERE class = ?
		ORDER BY spec
	`, className).Scan(&specs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get class specs: %w", err)
	}

	return specs, nil
}

// GetAvailableClasses retrieves all available classes
func (s *SpecEvolutionMetricsAnalysisService) GetAvailableClasses(ctx context.Context) ([]string, error) {
	var classes []string

	err := s.db.WithContext(ctx).Raw(`
		SELECT DISTINCT class
		FROM daily_spec_metrics_mythic_plus
		ORDER BY class
	`).Scan(&classes).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get available classes: %w", err)
	}

	return classes, nil
}

// GetRoleSpecs retrieves all specs for a given role
func (s *SpecEvolutionMetricsAnalysisService) GetRoleSpecs(ctx context.Context, role string) ([]playerRankingModels.DailySpecMetricMythicPlus, error) {
	query := `
	SELECT DISTINCT spec, class, role
	FROM daily_spec_metrics_mythic_plus
	WHERE role = ? AND is_global = true
	ORDER BY class, spec
	`

	var results []playerRankingModels.DailySpecMetricMythicPlus
	err := s.db.WithContext(ctx).Raw(query, role).Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get role specs: %w", err)
	}

	return results, nil
}

// GetSpecHistoricalData retrieves historical data for a spec over multiple dates
func (s *SpecEvolutionMetricsAnalysisService) GetSpecHistoricalData(ctx context.Context, spec string, class *string, days int, isGlobal bool, encounterID *int) ([]playerRankingModels.DailySpecMetricMythicPlus, error) {
	query := `
	SELECT *
	FROM daily_spec_metrics_mythic_plus
	WHERE spec = ?
	AND is_global = ?
	AND capture_date >= CURRENT_DATE - INTERVAL '%d days'
	`

	if class != nil {
		query += " AND class = ?"
	}

	if encounterID != nil {
		query += " AND encounter_id = ?"
	} else {
		query += " AND encounter_id = 0"
	}

	query += " ORDER BY capture_date"

	query = fmt.Sprintf(query, days)

	var results []playerRankingModels.DailySpecMetricMythicPlus
	var err error

	params := []interface{}{spec, isGlobal}
	if class != nil {
		params = append(params, *class)
	}
	if encounterID != nil {
		params = append(params, *encounterID)
	}

	err = s.db.WithContext(ctx).Raw(query, params...).Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get spec historical data: %w", err)
	}

	return results, nil
}

// GetAvailableDungeons retrieves all available dungeons with metrics
func (s *SpecEvolutionMetricsAnalysisService) GetAvailableDungeons(ctx context.Context) ([]struct {
	EncounterID int    `json:"encounter_id"`
	Name        string `json:"name"`
}, error) {
	query := `
	SELECT DISTINCT d.encounter_id, d.name
	FROM daily_spec_metrics_mythic_plus m
	JOIN dungeons d ON m.encounter_id = d.encounter_id
	WHERE m.encounter_id > 0
	ORDER BY d.name
	`

	var results []struct {
		EncounterID int    `json:"encounter_id"`
		Name        string `json:"name"`
	}

	err := s.db.WithContext(ctx).Raw(query).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get available dungeons: %w", err)
	}

	return results, nil
}

// GetTopSpecsForDungeon retrieves the top performing specs for a specific dungeon
func (s *SpecEvolutionMetricsAnalysisService) GetTopSpecsForDungeon(ctx context.Context, dungeonID int, limit int) ([]playerRankingModels.DailySpecMetricMythicPlus, error) {
	query := `
	SELECT *
	FROM daily_spec_metrics_mythic_plus
	WHERE encounter_id = ?
	AND is_global = false
	ORDER BY avg_score DESC
	LIMIT ?
	`

	var results []playerRankingModels.DailySpecMetricMythicPlus
	err := s.db.WithContext(ctx).Raw(query, dungeonID, limit).Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top specs for dungeon: %w", err)
	}

	return results, nil
}
