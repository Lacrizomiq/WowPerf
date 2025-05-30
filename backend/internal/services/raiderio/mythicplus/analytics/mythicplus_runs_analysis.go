package mythicPlusRunsAnalysis

import (
	"fmt"

	"gorm.io/gorm"
)

// MythicPlusRunsAnalysisService gère toutes les analyses statistiques des runs Mythic+
//
// PARAMÈTRES IMPORTANTS:
// - topN = 0 : Retourne TOUS les résultats (recommandé pour analyses complètes)
// - topN > 0 : Limite aux top N résultats (pour vues simplifiées)
// - minUsage > 0 : Filtre les données avec moins de X utilisations (évite le bruit)
type MythicPlusRunsAnalysisService struct {
	db *gorm.DB
}

// NewAnalyticsService crée un nouveau service d'analyse
func NewMythicPlusRunsAnalysisService(db *gorm.DB) *MythicPlusRunsAnalysisService {
	return &MythicPlusRunsAnalysisService{db: db}
}

// ========================================
// STRUCTURES DE DONNÉES
// ========================================

// SpecializationStats représente les stats d'une spécialisation
type SpecializationStats struct {
	Class      string   `json:"class"`
	Spec       string   `json:"spec"`
	Display    string   `json:"display"`
	UsageCount int      `json:"usage_count"`
	Percentage float64  `json:"percentage"`
	Rank       int      `json:"rank"`
	AvgScore   *float64 `json:"avg_score,omitempty"`
}

// CompositionStats représente les stats d'une composition complète
type CompositionStats struct {
	Tank       string  `json:"tank"`
	Healer     string  `json:"healer"`
	DPS1       string  `json:"dps1"`
	DPS2       string  `json:"dps2"`
	DPS3       string  `json:"dps3"`
	UsageCount int     `json:"usage_count"`
	Percentage float64 `json:"percentage"`
	Rank       int     `json:"rank"`
	AvgScore   float64 `json:"avg_score"`
}

// DungeonSpecStats représente les stats d'une spéc par donjon
type DungeonSpecStats struct {
	DungeonSlug   string  `json:"dungeon_slug"`
	DungeonName   string  `json:"dungeon_name"`
	Class         string  `json:"class"`
	Spec          string  `json:"spec"`
	Display       string  `json:"display"`
	UsageCount    int     `json:"usage_count"`
	Percentage    float64 `json:"percentage"`
	RankInDungeon int     `json:"rank_in_dungeon"`
}

// DungeonCompositionStats représente les stats d'une composition par donjon
type DungeonCompositionStats struct {
	DungeonSlug   string  `json:"dungeon_slug"`
	DungeonName   string  `json:"dungeon_name"`
	Tank          string  `json:"tank"`
	Healer        string  `json:"healer"`
	DPS1          string  `json:"dps1"`
	DPS2          string  `json:"dps2"`
	DPS3          string  `json:"dps3"`
	UsageCount    int     `json:"usage_count"`
	Percentage    float64 `json:"percentage"`
	AvgScore      float64 `json:"avg_score"`
	RankInDungeon int     `json:"rank_in_dungeon"`
}

// KeyLevelStats représente les stats par niveau de clé
type KeyLevelStats struct {
	Role            string  `json:"role"`
	KeyLevelBracket string  `json:"key_level_bracket"`
	Class           string  `json:"class"`
	Spec            string  `json:"spec"`
	Display         string  `json:"display"`
	UsageCount      int     `json:"usage_count"`
	Percentage      float64 `json:"percentage"`
	Rank            int     `json:"rank"`
	AvgScore        float64 `json:"avg_score"`
}

// RegionStats représente les stats par région
type RegionStats struct {
	Role               string  `json:"role"`
	Region             string  `json:"region"`
	Class              string  `json:"class"`
	Spec               string  `json:"spec"`
	Display            string  `json:"display"`
	UsageCount         int     `json:"usage_count"`
	PercentageInRegion float64 `json:"percentage_in_region"`
	RankInRegion       int     `json:"rank_in_region"`
}

// OverallStats représente les stats générales
type OverallStats struct {
	TotalRuns          int     `json:"total_runs"`
	RunsWithScore      int     `json:"runs_with_score"`
	UniqueCompositions int     `json:"unique_compositions"`
	UniqueDungeons     int     `json:"unique_dungeons"`
	UniqueRegions      int     `json:"unique_regions"`
	OldestRun          string  `json:"oldest_run"`
	NewestRun          string  `json:"newest_run"`
	AvgScore           float64 `json:"avg_score"`
	AvgKeyLevel        float64 `json:"avg_key_level"`
}

// KeyLevelDistribution représente la distribution des niveaux de clés
type KeyLevelDistribution struct {
	MythicLevel int     `json:"mythic_level"`
	Count       int     `json:"count"`
	Percentage  float64 `json:"percentage"`
	Rank        int     `json:"rank"`
	AvgScore    float64 `json:"avg_score"`
	MinScore    float64 `json:"min_score"`
	MaxScore    float64 `json:"max_score"`
}

// DungeonDistribution représente la distribution des runs par donjon
type DungeonDistribution struct {
	DungeonSlug string  `json:"dungeon_slug"`
	DungeonName string  `json:"dungeon_name"`
	RunCount    int     `json:"run_count"`
	Percentage  float64 `json:"percentage"`
	Rank        int     `json:"rank"`
	AvgScore    float64 `json:"avg_score"`
	AvgKeyLevel float64 `json:"avg_key_level"`
	MinScore    float64 `json:"min_score"`
	MaxScore    float64 `json:"max_score"`
}

// RegionDistribution représente la distribution des runs par région
type RegionDistribution struct {
	Region      string  `json:"region"`
	RunCount    int     `json:"run_count"`
	Percentage  float64 `json:"percentage"`
	Rank        int     `json:"rank"`
	AvgScore    float64 `json:"avg_score"`
	AvgKeyLevel float64 `json:"avg_key_level"`
}

// ========================================
// MÉTHODES D'ANALYSE GLOBALE
// ========================================

// GetTankSpecializations retourne les spécialisations Tank les plus utilisées
func (s *MythicPlusRunsAnalysisService) GetTankSpecializations() ([]SpecializationStats, error) {
	var results []SpecializationStats

	query := `
		SELECT 
			tc.tank_class as class,
			tc.tank_spec as spec,
			CONCAT(tc.tank_class, ' - ', tc.tank_spec) as display,
			COUNT(*) as usage_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage,
			ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank
		FROM mythicplus_runs r
		JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
		GROUP BY tc.tank_class, tc.tank_spec
		ORDER BY rank`

	if err := s.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get tank specializations: %w", err)
	}

	return results, nil
}

// GetHealerSpecializations retourne les spécialisations Healer les plus utilisées
func (s *MythicPlusRunsAnalysisService) GetHealerSpecializations() ([]SpecializationStats, error) {
	var results []SpecializationStats

	query := `
		SELECT 
			tc.healer_class as class,
			tc.healer_spec as spec,
			CONCAT(tc.healer_class, ' - ', tc.healer_spec) as display,
			COUNT(*) as usage_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage,
			ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank
		FROM mythicplus_runs r
		JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
		GROUP BY tc.healer_class, tc.healer_spec
		ORDER BY rank`

	if err := s.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get healer specializations: %w", err)
	}

	return results, nil
}

// GetDPSSpecializations retourne les spécialisations DPS les plus utilisées
func (s *MythicPlusRunsAnalysisService) GetDPSSpecializations() ([]SpecializationStats, error) {
	var results []SpecializationStats

	query := `
		SELECT 
			dps_class as class,
			dps_spec as spec,
			CONCAT(dps_class, ' - ', dps_spec) as display,
			COUNT(*) as usage_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage,
			ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank
		FROM (
			SELECT tc.dps1_class as dps_class, tc.dps1_spec as dps_spec FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
			
			UNION ALL
			
			SELECT tc.dps2_class, tc.dps2_spec FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
			
			UNION ALL
			
			SELECT tc.dps3_class, tc.dps3_spec FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
		) dps_data
		GROUP BY dps_class, dps_spec
		ORDER BY rank`

	if err := s.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get DPS specializations: %w", err)
	}

	return results, nil
}

// GetTopCompositions retourne les compositions les plus utilisées
func (s *MythicPlusRunsAnalysisService) GetTopCompositions(limit int, minUsage int) ([]CompositionStats, error) {
	var results []CompositionStats

	if limit <= 0 {
		limit = 20
	}
	if minUsage <= 0 {
		minUsage = 5
	}

	query := `
		SELECT 
			CONCAT(tc.tank_class, ' - ', tc.tank_spec) as tank,
			CONCAT(tc.healer_class, ' - ', tc.healer_spec) as healer,
			CONCAT(tc.dps1_class, ' - ', tc.dps1_spec) as dps1,
			CONCAT(tc.dps2_class, ' - ', tc.dps2_spec) as dps2,
			CONCAT(tc.dps3_class, ' - ', tc.dps3_spec) as dps3,
			COUNT(*) as usage_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage,
			ROUND(AVG(r.score), 1) as avg_score,
			ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank
		FROM mythicplus_runs r
		JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
		GROUP BY tc.id, tc.tank_class, tc.tank_spec, tc.healer_class, tc.healer_spec,
				 tc.dps1_class, tc.dps1_spec, tc.dps2_class, tc.dps2_spec, tc.dps3_class, tc.dps3_spec
		HAVING COUNT(*) >= ?
		ORDER BY rank
		LIMIT ?`

	if err := s.db.Raw(query, minUsage, limit).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get top compositions: %w", err)
	}

	return results, nil
}

// ========================================
// MÉTHODES D'ANALYSE PAR DONJON
// ========================================

// GetTankSpecsByDungeon retourne les spécialisations Tank par donjon
// Si topN = 0, retourne toutes les spécs. Si topN > 0, limite aux top N
func (s *MythicPlusRunsAnalysisService) GetTankSpecsByDungeon(topN int) ([]DungeonSpecStats, error) {
	var results []DungeonSpecStats

	var query string
	if topN <= 0 {
		// Retourne TOUTES les spécialisations
		query = `
			SELECT 
				r.dungeon_slug,
				r.dungeon_name,
				tc.tank_class as class,
				tc.tank_spec as spec,
				CONCAT(tc.tank_class, ' - ', tc.tank_spec) as display,
				COUNT(*) as usage_count,
				ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY r.dungeon_slug), 2) as percentage,
				ROW_NUMBER() OVER (PARTITION BY r.dungeon_slug ORDER BY COUNT(*) DESC) as rank_in_dungeon
			FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
			GROUP BY r.dungeon_slug, r.dungeon_name, tc.tank_class, tc.tank_spec
			ORDER BY dungeon_slug, rank_in_dungeon`
	} else {
		// Limite aux top N
		query = `
			WITH tank_specs_by_dungeon AS (
				SELECT 
					r.dungeon_slug,
					r.dungeon_name,
					tc.tank_class as class,
					tc.tank_spec as spec,
					CONCAT(tc.tank_class, ' - ', tc.tank_spec) as display,
					COUNT(*) as usage_count,
					ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY r.dungeon_slug), 2) as percentage,
					ROW_NUMBER() OVER (PARTITION BY r.dungeon_slug ORDER BY COUNT(*) DESC) as rank_in_dungeon
				FROM mythicplus_runs r
				JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
				GROUP BY r.dungeon_slug, r.dungeon_name, tc.tank_class, tc.tank_spec
			)
			SELECT *
			FROM tank_specs_by_dungeon
			WHERE rank_in_dungeon <= ?
			ORDER BY dungeon_slug, rank_in_dungeon`
	}

	if topN <= 0 {
		if err := s.db.Raw(query).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("failed to get tank specs by dungeon: %w", err)
		}
	} else {
		if err := s.db.Raw(query, topN).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("failed to get tank specs by dungeon: %w", err)
		}
	}

	return results, nil
}

// GetHealerSpecsByDungeon retourne les spécialisations Healer par donjon
// Si topN = 0, retourne toutes les spécs. Si topN > 0, limite aux top N
func (s *MythicPlusRunsAnalysisService) GetHealerSpecsByDungeon(topN int) ([]DungeonSpecStats, error) {
	var results []DungeonSpecStats

	var query string
	if topN <= 0 {
		// Retourne TOUTES les spécialisations
		query = `
			SELECT 
				r.dungeon_slug,
				r.dungeon_name,
				tc.healer_class as class,
				tc.healer_spec as spec,
				CONCAT(tc.healer_class, ' - ', tc.healer_spec) as display,
				COUNT(*) as usage_count,
				ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY r.dungeon_slug), 2) as percentage,
				ROW_NUMBER() OVER (PARTITION BY r.dungeon_slug ORDER BY COUNT(*) DESC) as rank_in_dungeon
			FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
			GROUP BY r.dungeon_slug, r.dungeon_name, tc.healer_class, tc.healer_spec
			ORDER BY dungeon_slug, rank_in_dungeon`
	} else {
		// Limite aux top N
		query = `
			WITH healer_specs_by_dungeon AS (
				SELECT 
					r.dungeon_slug,
					r.dungeon_name,
					tc.healer_class as class,
					tc.healer_spec as spec,
					CONCAT(tc.healer_class, ' - ', tc.healer_spec) as display,
					COUNT(*) as usage_count,
					ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY r.dungeon_slug), 2) as percentage,
					ROW_NUMBER() OVER (PARTITION BY r.dungeon_slug ORDER BY COUNT(*) DESC) as rank_in_dungeon
				FROM mythicplus_runs r
				JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
				GROUP BY r.dungeon_slug, r.dungeon_name, tc.healer_class, tc.healer_spec
			)
			SELECT *
			FROM healer_specs_by_dungeon
			WHERE rank_in_dungeon <= ?
			ORDER BY dungeon_slug, rank_in_dungeon`
	}

	if topN <= 0 {
		if err := s.db.Raw(query).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("failed to get healer specs by dungeon: %w", err)
		}
	} else {
		if err := s.db.Raw(query, topN).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("failed to get healer specs by dungeon: %w", err)
		}
	}

	return results, nil
}

// GetDPSSpecsByDungeon retourne les spécialisations DPS par donjon
// Si topN = 0, retourne toutes les spécs. Si topN > 0, limite aux top N
func (s *MythicPlusRunsAnalysisService) GetDPSSpecsByDungeon(topN int) ([]DungeonSpecStats, error) {
	var results []DungeonSpecStats

	var query string
	if topN <= 0 {
		// Retourne TOUTES les spécialisations DPS
		query = `
			SELECT 
				dungeon_slug,
				dungeon_name,
				dps_class as class,
				dps_spec as spec,
				CONCAT(dps_class, ' - ', dps_spec) as display,
				COUNT(*) as usage_count,
				ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY dungeon_slug), 2) as percentage,
				ROW_NUMBER() OVER (PARTITION BY dungeon_slug ORDER BY COUNT(*) DESC) as rank_in_dungeon
			FROM (
				SELECT r.dungeon_slug, r.dungeon_name, tc.dps1_class as dps_class, tc.dps1_spec as dps_spec 
				FROM mythicplus_runs r
				JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
				
				UNION ALL
				
				SELECT r.dungeon_slug, r.dungeon_name, tc.dps2_class, tc.dps2_spec 
				FROM mythicplus_runs r
				JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
				
				UNION ALL
				
				SELECT r.dungeon_slug, r.dungeon_name, tc.dps3_class, tc.dps3_spec 
				FROM mythicplus_runs r
				JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
			) dps_data
			GROUP BY dungeon_slug, dungeon_name, dps_class, dps_spec
			ORDER BY dungeon_slug, rank_in_dungeon`
	} else {
		// Limite aux top N
		query = `
			WITH dps_specs_by_dungeon AS (
				SELECT 
					dungeon_slug,
					dungeon_name,
					dps_class as class,
					dps_spec as spec,
					CONCAT(dps_class, ' - ', dps_spec) as display,
					COUNT(*) as usage_count,
					ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY dungeon_slug), 2) as percentage,
					ROW_NUMBER() OVER (PARTITION BY dungeon_slug ORDER BY COUNT(*) DESC) as rank_in_dungeon
				FROM (
					SELECT r.dungeon_slug, r.dungeon_name, tc.dps1_class as dps_class, tc.dps1_spec as dps_spec 
					FROM mythicplus_runs r
					JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
					
					UNION ALL
					
					SELECT r.dungeon_slug, r.dungeon_name, tc.dps2_class, tc.dps2_spec 
					FROM mythicplus_runs r
					JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
					
					UNION ALL
					
					SELECT r.dungeon_slug, r.dungeon_name, tc.dps3_class, tc.dps3_spec 
					FROM mythicplus_runs r
					JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
				) dps_data
				GROUP BY dungeon_slug, dungeon_name, dps_class, dps_spec
			)
			SELECT *
			FROM dps_specs_by_dungeon
			WHERE rank_in_dungeon <= ?
			ORDER BY dungeon_slug, rank_in_dungeon`
	}

	if topN <= 0 {
		if err := s.db.Raw(query).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("failed to get DPS specs by dungeon: %w", err)
		}
	} else {
		if err := s.db.Raw(query, topN).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("failed to get DPS specs by dungeon: %w", err)
		}
	}

	return results, nil
}

// GetTopCompositionsByDungeon retourne les compositions par donjon
// Si topN = 0, retourne toutes les compositions. Si topN > 0, limite aux top N
func (s *MythicPlusRunsAnalysisService) GetTopCompositionsByDungeon(topN int, minUsage int) ([]DungeonCompositionStats, error) {
	var results []DungeonCompositionStats

	if minUsage <= 0 {
		minUsage = 3
	}

	var query string
	if topN <= 0 {
		// Retourne TOUTES les compositions (avec filtre minUsage)
		query = `
			SELECT 
				r.dungeon_slug,
				r.dungeon_name,
				CONCAT(tc.tank_class, ' - ', tc.tank_spec) as tank,
				CONCAT(tc.healer_class, ' - ', tc.healer_spec) as healer,
				CONCAT(tc.dps1_class, ' - ', tc.dps1_spec) as dps1,
				CONCAT(tc.dps2_class, ' - ', tc.dps2_spec) as dps2,
				CONCAT(tc.dps3_class, ' - ', tc.dps3_spec) as dps3,
				COUNT(*) as usage_count,
				ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY r.dungeon_slug), 2) as percentage,
				ROUND(AVG(r.score), 1) as avg_score,
				ROW_NUMBER() OVER (PARTITION BY r.dungeon_slug ORDER BY COUNT(*) DESC) as rank_in_dungeon
			FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
			GROUP BY r.dungeon_slug, r.dungeon_name, tc.tank_class, tc.tank_spec, tc.healer_class, tc.healer_spec, 
					 tc.dps1_class, tc.dps1_spec, tc.dps2_class, tc.dps2_spec, tc.dps3_class, tc.dps3_spec
			HAVING COUNT(*) >= ?
			ORDER BY dungeon_slug, rank_in_dungeon`
	} else {
		// Limite aux top N
		query = `
			WITH compositions_by_dungeon AS (
				SELECT 
					r.dungeon_slug,
					r.dungeon_name,
					CONCAT(tc.tank_class, ' - ', tc.tank_spec) as tank,
					CONCAT(tc.healer_class, ' - ', tc.healer_spec) as healer,
					CONCAT(tc.dps1_class, ' - ', tc.dps1_spec) as dps1,
					CONCAT(tc.dps2_class, ' - ', tc.dps2_spec) as dps2,
					CONCAT(tc.dps3_class, ' - ', tc.dps3_spec) as dps3,
					COUNT(*) as usage_count,
					ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY r.dungeon_slug), 2) as percentage,
					ROUND(AVG(r.score), 1) as avg_score,
					ROW_NUMBER() OVER (PARTITION BY r.dungeon_slug ORDER BY COUNT(*) DESC) as rank_in_dungeon
				FROM mythicplus_runs r
				JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
				GROUP BY r.dungeon_slug, r.dungeon_name, tc.tank_class, tc.tank_spec, tc.healer_class, tc.healer_spec, 
						 tc.dps1_class, tc.dps1_spec, tc.dps2_class, tc.dps2_spec, tc.dps3_class, tc.dps3_spec
				HAVING COUNT(*) >= ?
			)
			SELECT *
			FROM compositions_by_dungeon
			WHERE rank_in_dungeon <= ?
			ORDER BY dungeon_slug, rank_in_dungeon`
	}

	if topN <= 0 {
		if err := s.db.Raw(query, minUsage).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("failed to get compositions by dungeon: %w", err)
		}
	} else {
		if err := s.db.Raw(query, minUsage, topN).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("failed to get compositions by dungeon: %w", err)
		}
	}

	return results, nil
}

// ========================================
// MÉTHODES D'ANALYSE AVANCÉE
// ========================================

// GetSpecsByKeyLevel retourne les spécialisations par niveau de clé (tous rôles)
func (s *MythicPlusRunsAnalysisService) GetSpecsByKeyLevel(minUsage int) ([]KeyLevelStats, error) {
	var results []KeyLevelStats

	if minUsage <= 0 {
		minUsage = 5
	}

	query := `
		-- Tank par niveau de clé
		SELECT 
			'Tank' as role,
			CASE 
				WHEN r.mythic_level >= 20 THEN 'Very High Keys (20+)'
				WHEN r.mythic_level >= 18 THEN 'High Keys (18-19)'
				WHEN r.mythic_level >= 16 THEN 'Mid Keys (16-17)'
				ELSE 'Other Keys (<16)'
			END as key_level_bracket,
			tc.tank_class as class,
			tc.tank_spec as spec,
			CONCAT(tc.tank_class, ' - ', tc.tank_spec) as display,
			COUNT(*) as usage_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY 
				CASE 
					WHEN r.mythic_level >= 20 THEN 'Very High Keys (20+)'
					WHEN r.mythic_level >= 18 THEN 'High Keys (18-19)'
					WHEN r.mythic_level >= 16 THEN 'Mid Keys (16-17)'
					ELSE 'Other Keys (<16)'
				END
			), 2) as percentage,
			ROW_NUMBER() OVER (PARTITION BY 
				CASE 
					WHEN r.mythic_level >= 20 THEN 'Very High Keys (20+)'
					WHEN r.mythic_level >= 18 THEN 'High Keys (18-19)'
					WHEN r.mythic_level >= 16 THEN 'Mid Keys (16-17)'
					ELSE 'Other Keys (<16)'
				END
				ORDER BY COUNT(*) DESC
			) as rank,
			ROUND(AVG(r.score), 1) as avg_score
		FROM mythicplus_runs r
		JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
		GROUP BY 
			CASE 
				WHEN r.mythic_level >= 20 THEN 'Very High Keys (20+)'
				WHEN r.mythic_level >= 18 THEN 'High Keys (18-19)'
				WHEN r.mythic_level >= 16 THEN 'Mid Keys (16-17)'
				ELSE 'Other Keys (<16)'
			END,
			tc.tank_class, tc.tank_spec
		HAVING COUNT(*) >= ?

		UNION ALL

		-- Healer par niveau de clé
		SELECT 
			'Healer' as role,
			CASE 
				WHEN r.mythic_level >= 20 THEN 'Very High Keys (20+)'
				WHEN r.mythic_level >= 18 THEN 'High Keys (18-19)'
				WHEN r.mythic_level >= 16 THEN 'Mid Keys (16-17)'
				ELSE 'Other Keys (<16)'
			END as key_level_bracket,
			tc.healer_class as class,
			tc.healer_spec as spec,
			CONCAT(tc.healer_class, ' - ', tc.healer_spec) as display,
			COUNT(*) as usage_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY 
				CASE 
					WHEN r.mythic_level >= 20 THEN 'Very High Keys (20+)'
					WHEN r.mythic_level >= 18 THEN 'High Keys (18-19)'
					WHEN r.mythic_level >= 16 THEN 'Mid Keys (16-17)'
					ELSE 'Other Keys (<16)'
				END
			), 2) as percentage,
			ROW_NUMBER() OVER (PARTITION BY 
				CASE 
					WHEN r.mythic_level >= 20 THEN 'Very High Keys (20+)'
					WHEN r.mythic_level >= 18 THEN 'High Keys (18-19)'
					WHEN r.mythic_level >= 16 THEN 'Mid Keys (16-17)'
					ELSE 'Other Keys (<16)'
				END
				ORDER BY COUNT(*) DESC
			) as rank,
			ROUND(AVG(r.score), 1) as avg_score
		FROM mythicplus_runs r
		JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
		GROUP BY 
			CASE 
				WHEN r.mythic_level >= 20 THEN 'Very High Keys (20+)'
				WHEN r.mythic_level >= 18 THEN 'High Keys (18-19)'
				WHEN r.mythic_level >= 16 THEN 'Mid Keys (16-17)'
				ELSE 'Other Keys (<16)'
			END,
			tc.healer_class, tc.healer_spec
		HAVING COUNT(*) >= ?

		UNION ALL

		-- DPS par niveau de clé
		SELECT 
			'DPS' as role,
			key_level_bracket,
			dps_class as class,
			dps_spec as spec,
			CONCAT(dps_class, ' - ', dps_spec) as display,
			COUNT(*) as usage_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY key_level_bracket), 2) as percentage,
			ROW_NUMBER() OVER (PARTITION BY key_level_bracket ORDER BY COUNT(*) DESC) as rank,
			ROUND(AVG(score), 1) as avg_score
		FROM (
			SELECT 
				CASE 
					WHEN r.mythic_level >= 20 THEN 'Very High Keys (20+)'
					WHEN r.mythic_level >= 18 THEN 'High Keys (18-19)'
					WHEN r.mythic_level >= 16 THEN 'Mid Keys (16-17)'
					ELSE 'Other Keys (<16)'
				END as key_level_bracket,
				tc.dps1_class as dps_class, 
				tc.dps1_spec as dps_spec,
				r.score
			FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
			
			UNION ALL
			
			SELECT 
				CASE 
					WHEN r.mythic_level >= 20 THEN 'Very High Keys (20+)'
					WHEN r.mythic_level >= 18 THEN 'High Keys (18-19)'
					WHEN r.mythic_level >= 16 THEN 'Mid Keys (16-17)'
					ELSE 'Other Keys (<16)'
				END as key_level_bracket,
				tc.dps2_class, 
				tc.dps2_spec,
				r.score
			FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
			
			UNION ALL
			
			SELECT 
				CASE 
					WHEN r.mythic_level >= 20 THEN 'Very High Keys (20+)'
					WHEN r.mythic_level >= 18 THEN 'High Keys (18-19)'
					WHEN r.mythic_level >= 16 THEN 'Mid Keys (16-17)'
					ELSE 'Other Keys (<16)'
				END as key_level_bracket,
				tc.dps3_class, 
				tc.dps3_spec,
				r.score
			FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
		) dps_data
		GROUP BY key_level_bracket, dps_class, dps_spec
		HAVING COUNT(*) >= ?

		ORDER BY role, key_level_bracket, rank`

	if err := s.db.Raw(query, minUsage, minUsage, minUsage).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get specs by key level: %w", err)
	}

	return results, nil
}

// GetSpecsByRegion retourne les spécialisations par région (tous rôles)
func (s *MythicPlusRunsAnalysisService) GetSpecsByRegion() ([]RegionStats, error) {
	var results []RegionStats

	query := `
		-- Tank par région
		SELECT 
			'Tank' as role,
			r.region,
			tc.tank_class as class,
			tc.tank_spec as spec,
			CONCAT(tc.tank_class, ' - ', tc.tank_spec) as display,
			COUNT(*) as usage_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY r.region), 2) as percentage_in_region,
			ROW_NUMBER() OVER (PARTITION BY r.region ORDER BY COUNT(*) DESC) as rank_in_region
		FROM mythicplus_runs r
		JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
		GROUP BY r.region, tc.tank_class, tc.tank_spec

		UNION ALL

		-- Healer par région
		SELECT 
			'Healer' as role,
			r.region,
			tc.healer_class as class,
			tc.healer_spec as spec,
			CONCAT(tc.healer_class, ' - ', tc.healer_spec) as display,
			COUNT(*) as usage_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY r.region), 2) as percentage_in_region,
			ROW_NUMBER() OVER (PARTITION BY r.region ORDER BY COUNT(*) DESC) as rank_in_region
		FROM mythicplus_runs r
		JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
		GROUP BY r.region, tc.healer_class, tc.healer_spec

		UNION ALL

		-- DPS par région
		SELECT 
			'DPS' as role,
			region,
			dps_class as class,
			dps_spec as spec,
			CONCAT(dps_class, ' - ', dps_spec) as display,
			COUNT(*) as usage_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY region), 2) as percentage_in_region,
			ROW_NUMBER() OVER (PARTITION BY region ORDER BY COUNT(*) DESC) as rank_in_region
		FROM (
			SELECT r.region, tc.dps1_class as dps_class, tc.dps1_spec as dps_spec 
			FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
			
			UNION ALL
			
			SELECT r.region, tc.dps2_class, tc.dps2_spec 
			FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
			
			UNION ALL
			
			SELECT r.region, tc.dps3_class, tc.dps3_spec 
			FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
		) dps_data
		GROUP BY region, dps_class, dps_spec

		ORDER BY role, region, rank_in_region`

	if err := s.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get specs by region: %w", err)
	}

	return results, nil
}

// ========================================
// MÉTHODES UTILITAIRES
// ========================================

// GetOverallStats retourne les statistiques générales
func (s *MythicPlusRunsAnalysisService) GetOverallStats() (*OverallStats, error) {
	var result OverallStats

	query := `
		SELECT 
			COUNT(*) as total_runs,
			COUNT(CASE WHEN r.score > 0 THEN 1 END) as runs_with_score,
			COUNT(DISTINCT tc.id) as unique_compositions,
			COUNT(DISTINCT r.dungeon_slug) as unique_dungeons,
			COUNT(DISTINCT r.region) as unique_regions,
			MIN(r.completed_at) as oldest_run,
			MAX(r.completed_at) as newest_run,
			ROUND(AVG(r.score), 1) as avg_score,
			ROUND(AVG(r.mythic_level), 1) as avg_key_level
		FROM mythicplus_runs r
		JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id`

	if err := s.db.Raw(query).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get overall stats: %w", err)
	}

	return &result, nil
}

// GetKeyLevelDistribution retourne la distribution des niveaux de clés
func (s *MythicPlusRunsAnalysisService) GetKeyLevelDistribution() ([]KeyLevelDistribution, error) {
	var results []KeyLevelDistribution

	query := `
		SELECT 
			r.mythic_level,
			COUNT(*) as count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage,
			ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank,
			ROUND(AVG(r.score), 1) as avg_score,
			ROUND(MIN(r.score), 1) as min_score,
			ROUND(MAX(r.score), 1) as max_score
		FROM mythicplus_runs r
		GROUP BY r.mythic_level
		ORDER BY r.mythic_level`

	if err := s.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get key level distribution: %w", err)
	}

	return results, nil
}

// GetDungeonDistribution retourne la distribution des runs par donjon
func (s *MythicPlusRunsAnalysisService) GetDungeonDistribution() ([]DungeonDistribution, error) {
	var results []DungeonDistribution

	query := `
		SELECT 
			r.dungeon_slug,
			r.dungeon_name,
			COUNT(*) as run_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage,
			ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank,
			ROUND(AVG(r.score), 1) as avg_score,
			ROUND(AVG(r.mythic_level), 1) as avg_key_level,
			ROUND(MIN(r.score), 1) as min_score,
			ROUND(MAX(r.score), 1) as max_score
		FROM mythicplus_runs r
		GROUP BY r.dungeon_slug, r.dungeon_name
		ORDER BY rank`

	if err := s.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get dungeon distribution: %w", err)
	}

	return results, nil
}

// GetRegionDistribution retourne la distribution des runs par région
func (s *MythicPlusRunsAnalysisService) GetRegionDistribution() ([]RegionDistribution, error) {
	var results []RegionDistribution

	query := `
		SELECT 
			r.region,
			COUNT(*) as run_count,
			ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage,
			ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank,
			ROUND(AVG(r.score), 1) as avg_score,
			ROUND(AVG(r.mythic_level), 1) as avg_key_level
		FROM mythicplus_runs r
		GROUP BY r.region
		ORDER BY rank`

	if err := s.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get region distribution: %w", err)
	}

	return results, nil
}

// ========================================
// MÉTHODES DE CONVENANCE
// ========================================

// GetAllTankSpecsByDungeon retourne TOUTES les spécialisations Tank par donjon
func (s *MythicPlusRunsAnalysisService) GetAllTankSpecsByDungeon() ([]DungeonSpecStats, error) {
	return s.GetTankSpecsByDungeon(0) // topN = 0 = tout retourner
}

// GetAllHealerSpecsByDungeon retourne TOUTES les spécialisations Healer par donjon
func (s *MythicPlusRunsAnalysisService) GetAllHealerSpecsByDungeon() ([]DungeonSpecStats, error) {
	return s.GetHealerSpecsByDungeon(0) // topN = 0 = tout retourner
}

// GetAllDPSSpecsByDungeon retourne TOUTES les spécialisations DPS par donjon
func (s *MythicPlusRunsAnalysisService) GetAllDPSSpecsByDungeon() ([]DungeonSpecStats, error) {
	return s.GetDPSSpecsByDungeon(0) // topN = 0 = tout retourner
}

// GetAllCompositionsByDungeon retourne TOUTES les compositions par donjon (avec filtre minUsage)
func (s *MythicPlusRunsAnalysisService) GetAllCompositionsByDungeon(minUsage int) ([]DungeonCompositionStats, error) {
	return s.GetTopCompositionsByDungeon(0, minUsage) // topN = 0 = tout retourner
}

// GetTopCompositionsClean retourne les top compositions avec des paramètres recommandés
func (s *MythicPlusRunsAnalysisService) GetTopCompositionsClean() ([]CompositionStats, error) {
	return s.GetTopCompositions(20, 10) // Top 20, min 10 utilisations
}

// GetTopSpecsByDungeonClean retourne les top 3 spécs par donjon pour un rôle
func (s *MythicPlusRunsAnalysisService) GetTopSpecsByDungeonClean(role string) ([]DungeonSpecStats, error) {
	switch role {
	case "tank":
		return s.GetTankSpecsByDungeon(3)
	case "healer":
		return s.GetHealerSpecsByDungeon(3)
	case "dps":
		return s.GetDPSSpecsByDungeon(5) // Plus de DPS différents
	default:
		return nil, fmt.Errorf("invalid role: %s", role)
	}
}

// ========================================
// MÉTHODES DE FILTRAGE SPÉCIALISÉES
// ========================================

// GetSpecializationsByRole retourne toutes les spécialisations d'un rôle donné
func (s *MythicPlusRunsAnalysisService) GetSpecializationsByRole(role string) ([]SpecializationStats, error) {
	switch role {
	case "tank":
		return s.GetTankSpecializations()
	case "healer":
		return s.GetHealerSpecializations()
	case "dps":
		return s.GetDPSSpecializations()
	default:
		return nil, fmt.Errorf("invalid role: %s (must be 'tank', 'healer', or 'dps')", role)
	}
}

// GetSpecsByDungeonAndRole retourne les spécialisations d'un rôle pour un donjon spécifique
// Si topN = 0, retourne toutes les spécs. Si topN > 0, limite aux top N
func (s *MythicPlusRunsAnalysisService) GetSpecsByDungeonAndRole(dungeonSlug string, role string, topN int) ([]DungeonSpecStats, error) {
	var results []DungeonSpecStats

	var roleColumn, roleField string
	switch role {
	case "tank":
		roleColumn = "tc.tank_class"
		roleField = "tc.tank_spec"
	case "healer":
		roleColumn = "tc.healer_class"
		roleField = "tc.healer_spec"
	case "dps":
		// Pour DPS, on utilise la requête UNION comme dans GetDPSSpecsByDungeon
		// mais filtrée par donjon
		var query string
		if topN <= 0 {
			// Retourne TOUTES les spécialisations DPS pour ce donjon
			query = `
				SELECT 
					dungeon_slug,
					dungeon_name,
					dps_class as class,
					dps_spec as spec,
					CONCAT(dps_class, ' - ', dps_spec) as display,
					COUNT(*) as usage_count,
					ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage,
					ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank_in_dungeon
				FROM (
					SELECT r.dungeon_slug, r.dungeon_name, tc.dps1_class as dps_class, tc.dps1_spec as dps_spec 
					FROM mythicplus_runs r
					JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
					WHERE r.dungeon_slug = ?
					
					UNION ALL
					
					SELECT r.dungeon_slug, r.dungeon_name, tc.dps2_class, tc.dps2_spec 
					FROM mythicplus_runs r
					JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
					WHERE r.dungeon_slug = ?
					
					UNION ALL
					
					SELECT r.dungeon_slug, r.dungeon_name, tc.dps3_class, tc.dps3_spec 
					FROM mythicplus_runs r
					JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
					WHERE r.dungeon_slug = ?
				) dps_data
				GROUP BY dungeon_slug, dungeon_name, dps_class, dps_spec
				ORDER BY rank_in_dungeon`

			if err := s.db.Raw(query, dungeonSlug, dungeonSlug, dungeonSlug).Scan(&results).Error; err != nil {
				return nil, fmt.Errorf("failed to get DPS specs for dungeon %s: %w", dungeonSlug, err)
			}
		} else {
			// Limite aux top N
			query = `
				WITH dps_specs_by_specific_dungeon AS (
					SELECT 
						dungeon_slug,
						dungeon_name,
						dps_class as class,
						dps_spec as spec,
						CONCAT(dps_class, ' - ', dps_spec) as display,
						COUNT(*) as usage_count,
						ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage,
						ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank_in_dungeon
					FROM (
						SELECT r.dungeon_slug, r.dungeon_name, tc.dps1_class as dps_class, tc.dps1_spec as dps_spec 
						FROM mythicplus_runs r
						JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
						WHERE r.dungeon_slug = ?
						
						UNION ALL
						
						SELECT r.dungeon_slug, r.dungeon_name, tc.dps2_class, tc.dps2_spec 
						FROM mythicplus_runs r
						JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
						WHERE r.dungeon_slug = ?
						
						UNION ALL
						
						SELECT r.dungeon_slug, r.dungeon_name, tc.dps3_class, tc.dps3_spec 
						FROM mythicplus_runs r
						JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
						WHERE r.dungeon_slug = ?
					) dps_data
					GROUP BY dungeon_slug, dungeon_name, dps_class, dps_spec
				)
				SELECT *
				FROM dps_specs_by_specific_dungeon
				WHERE rank_in_dungeon <= ?
				ORDER BY rank_in_dungeon`

			if err := s.db.Raw(query, dungeonSlug, dungeonSlug, dungeonSlug, topN).Scan(&results).Error; err != nil {
				return nil, fmt.Errorf("failed to get DPS specs for dungeon %s: %w", dungeonSlug, err)
			}
		}
		return results, nil
	default:
		return nil, fmt.Errorf("invalid role: %s (must be 'tank', 'healer', or 'dps')", role)
	}

	// Pour Tank et Healer
	var query string
	if topN <= 0 {
		// Retourne TOUTES les spécialisations
		query = fmt.Sprintf(`
			SELECT 
				r.dungeon_slug,
				r.dungeon_name,
				%s as class,
				%s as spec,
				CONCAT(%s, ' - ', %s) as display,
				COUNT(*) as usage_count,
				ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage,
				ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank_in_dungeon
			FROM mythicplus_runs r
			JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
			WHERE r.dungeon_slug = ?
			GROUP BY r.dungeon_slug, r.dungeon_name, %s, %s
			ORDER BY rank_in_dungeon`, roleColumn, roleField, roleColumn, roleField, roleColumn, roleField)

		if err := s.db.Raw(query, dungeonSlug).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("failed to get %s specs for dungeon %s: %w", role, dungeonSlug, err)
		}
	} else {
		// Limite aux top N
		query = fmt.Sprintf(`
			WITH specs_by_specific_dungeon AS (
				SELECT 
					r.dungeon_slug,
					r.dungeon_name,
					%s as class,
					%s as spec,
					CONCAT(%s, ' - ', %s) as display,
					COUNT(*) as usage_count,
					ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage,
					ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank_in_dungeon
				FROM mythicplus_runs r
				JOIN mythicplus_team_compositions tc ON r.team_composition_id = tc.id
				WHERE r.dungeon_slug = ?
				GROUP BY r.dungeon_slug, r.dungeon_name, %s, %s
			)
			SELECT *
			FROM specs_by_specific_dungeon
			WHERE rank_in_dungeon <= ?
			ORDER BY rank_in_dungeon`, roleColumn, roleField, roleColumn, roleField, roleColumn, roleField)

		if err := s.db.Raw(query, dungeonSlug, topN).Scan(&results).Error; err != nil {
			return nil, fmt.Errorf("failed to get %s specs for dungeon %s: %w", role, dungeonSlug, err)
		}
	}

	return results, nil
}
