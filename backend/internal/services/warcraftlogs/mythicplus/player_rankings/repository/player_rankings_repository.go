package player_rankings_repository

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"

	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"
)

// Batch size constant for insertions
const batchSize = 100

type PlayerRankingsRepository struct {
	db *gorm.DB
}

func NewPlayerRankingsRepository(db *gorm.DB) *PlayerRankingsRepository {
	return &PlayerRankingsRepository{
		db: db,
	}
}

// DeleteExistingRankings removes all existing rankings
func (r *PlayerRankingsRepository) DeleteExistingRankings(ctx context.Context) error {
	log.Println("Deleting existing rankings")
	return r.db.WithContext(ctx).Exec("DELETE FROM player_rankings").Error
}

// StoreRankingsByBatches inserts rankings in batches
func (r *PlayerRankingsRepository) StoreRankingsByBatches(ctx context.Context, rankings []playerRankingModels.PlayerRanking) error {
	if len(rankings) == 0 {
		log.Println("No rankings to store")
		return nil
	}

	log.Printf("Storing %d rankings in batches", len(rankings))

	baseSQL := `
    INSERT INTO player_rankings (
        created_at, updated_at, dungeon_id, name, class, spec, role, 
        amount, hard_mode_level, duration, start_time, report_code, 
        report_fight_id, report_start_time, guild_id, guild_name, 
        guild_faction, server_id, server_name, server_region, 
        bracket_data, faction, affixes, medal, score, leaderboard
    ) VALUES `

	totalRankings := len(rankings)
	batches := int(math.Ceil(float64(totalRankings) / float64(batchSize)))

	for i := 0; i < batches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > totalRankings {
			end = totalRankings
		}

		batch := rankings[start:end]
		valueStrings := make([]string, len(batch))
		valueArgs := make([]interface{}, 0, len(batch)*26)

		for j, ranking := range batch {
			valueStrings[j] = fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				j*26+1, j*26+2, j*26+3, j*26+4, j*26+5, j*26+6, j*26+7, j*26+8, j*26+9, j*26+10,
				j*26+11, j*26+12, j*26+13, j*26+14, j*26+15, j*26+16, j*26+17, j*26+18, j*26+19, j*26+20,
				j*26+21, j*26+22, j*26+23, j*26+24, j*26+25, j*26+26,
			)

			now := time.Now()
			valueArgs = append(valueArgs,
				now, // created_at
				now, // updated_at
				ranking.DungeonID,
				ranking.Name,
				ranking.Class,
				ranking.Spec,
				ranking.Role,
				ranking.Amount,
				ranking.HardModeLevel,
				ranking.Duration,
				ranking.StartTime,
				ranking.ReportCode,
				ranking.ReportFightID,
				ranking.ReportStartTime,
				ranking.GuildID,
				ranking.GuildName,
				ranking.GuildFaction,
				ranking.ServerID,
				ranking.ServerName,
				ranking.ServerRegion,
				ranking.BracketData,
				ranking.Faction,
				pq.Array(ranking.Affixes),
				ranking.Medal,
				ranking.Score,
				ranking.Leaderboard,
			)
		}

		query := baseSQL + strings.Join(valueStrings, ",")
		if err := r.db.WithContext(ctx).Exec(query, valueArgs...).Error; err != nil {
			return fmt.Errorf("failed to insert batch %d: %w", i+1, err)
		}

		log.Printf("Batch %d/%d inserted (%d rankings)", i+1, batches, len(batch))
	}

	return nil
}

// CalculateDailySpecMetrics calculates daily metrics for all specializations
func (r *PlayerRankingsRepository) CalculateDailySpecMetrics(ctx context.Context) error {
	// Unique date for the entire processing
	processingDate := time.Now().Truncate(24 * time.Hour)
	log.Printf("Calculating spec metrics for date: %s", processingDate.Format("2006-01-02"))

	// Fixed number of dungeons (always 8)
	const totalDungeonCount = 8
	// Number of top players to consider for averages
	const topPlayersCount = 10

	// Use a transaction to ensure integrity
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Delete existing metrics for this date
		if err := tx.Exec("DELETE FROM daily_spec_metrics_mythic_plus WHERE capture_date = ?", processingDate).Error; err != nil {
			return fmt.Errorf("error deleting existing metrics: %w", err)
		}

		// 2. Calculate metrics per dungeon in a single query
		type DungeonMetric struct {
			Spec        string
			Class       string
			Role        string
			DungeonID   int     `gorm:"column:dungeon_id"`
			AvgScore    float64 `gorm:"column:avg_score"`
			MaxScore    float64 `gorm:"column:max_score"`
			MinScore    float64 `gorm:"column:min_score"`
			AvgKeyLevel float64 `gorm:"column:avg_key_level"`
			MaxKeyLevel int     `gorm:"column:max_key_level"`
			MinKeyLevel int     `gorm:"column:min_key_level"`
			Count       int
		}

		var dungeonMetrics []DungeonMetric

		// Optimized query with GROUP BY
		if err := tx.Raw(`
			SELECT 
				spec, 
				class, 
				role, 
				dungeon_id,
				COALESCE(AVG(score), 0) AS avg_score,
				COALESCE(MAX(score), 0) AS max_score,
				COALESCE(MIN(score), 0) AS min_score,
				COALESCE(AVG(hard_mode_level), 0) AS avg_key_level,
				COALESCE(MAX(hard_mode_level), 0) AS max_key_level,
				COALESCE(MIN(hard_mode_level), 0) AS min_key_level,
				COUNT(*) AS count
			FROM player_rankings
			WHERE server_region != 'CN'
			GROUP BY spec, class, role, dungeon_id
		`).Scan(&dungeonMetrics).Error; err != nil {
			return fmt.Errorf("error calculating dungeon metrics: %w", err)
		}

		log.Printf("Calculated metrics for %d dungeon+spec combinations", len(dungeonMetrics))

		// Convert to DailySpecMetricMythicPlus
		var metrics []playerRankingModels.DailySpecMetricMythicPlus
		for _, dm := range dungeonMetrics {
			metrics = append(metrics, playerRankingModels.DailySpecMetricMythicPlus{
				CaptureDate: processingDate,
				Spec:        dm.Spec,
				Class:       dm.Class,
				Role:        dm.Role,
				EncounterID: dm.DungeonID,
				IsGlobal:    false,
				AvgScore:    dm.AvgScore,
				MaxScore:    dm.MaxScore,
				MinScore:    dm.MinScore,
				AvgKeyLevel: dm.AvgKeyLevel,
				MaxKeyLevel: dm.MaxKeyLevel,
				MinKeyLevel: dm.MinKeyLevel,
				RoleRank:    0, // Will be calculated later per dungeon
				OverallRank: 0, // Will be calculated later per dungeon
			})
		}

		// 3. Calculate rankings for dungeon metrics (grouped by dungeon)
		calculateDungeonMetricsRankings(&metrics)

		// 4. Calculate global metrics based on TOP player performances
		var globalMetrics []playerRankingModels.DailySpecMetricMythicPlus

		// Structure to store global scores by player
		type PlayerGlobalScore struct {
			Spec         string
			Class        string
			Role         string
			Name         string
			ServerName   string
			TotalScore   float64 `gorm:"column:total_score"`
			AvgKeyLevel  float64 `gorm:"column:avg_key_level"`
			MaxKeyLevel  int     `gorm:"column:max_key_level"`
			MinKeyLevel  int     `gorm:"column:min_key_level"`
			DungeonCount int     `gorm:"column:dungeon_count"`
		}

		var playerScores []PlayerGlobalScore

		// SQL query to calculate global scores by player
		// Filter to keep only players who have completed all 8 dungeons
		if err := tx.Raw(`
			SELECT 
    spec, 
    class, 
    role, 
    name,
    server_name,
    SUM(best_score) AS total_score,
    AVG(avg_key_level) AS avg_key_level,
    MAX(max_key_level) AS max_key_level,
    MIN(min_key_level) AS min_key_level,
    COUNT(DISTINCT dungeon_id) AS dungeon_count
FROM (
    SELECT 
        spec, 
        class, 
        role, 
        name,
        server_name,
        dungeon_id,
        MAX(score) as best_score,
        AVG(hard_mode_level) as avg_key_level,
        MAX(hard_mode_level) as max_key_level,
        MIN(hard_mode_level) as min_key_level
    FROM player_rankings
    WHERE server_region != 'CN'
    GROUP BY spec, class, role, name, server_name, dungeon_id
) AS best_scores
GROUP BY spec, class, role, name, server_name
HAVING COUNT(DISTINCT dungeon_id) = 8
		`).Scan(&playerScores).Error; err != nil {
			return fmt.Errorf("error calculating player global scores: %w", err)
		}

		log.Printf("Calculated global scores for %d players who completed all 8 dungeons", len(playerScores))

		// Group players by spec/class/role
		specPlayerMap := make(map[string][]PlayerGlobalScore)
		for _, ps := range playerScores {
			key := fmt.Sprintf("%s-%s-%s", ps.Spec, ps.Class, ps.Role)
			specPlayerMap[key] = append(specPlayerMap[key], ps)
		}

		// Calculate global metrics based on TOP 10 players of each spec
		for key, players := range specPlayerMap {
			parts := strings.Split(key, "-")
			spec, class, role := parts[0], parts[1], parts[2]

			// Sort players by total score (from highest to lowest)
			sort.Slice(players, func(i, j int) bool {
				return players[i].TotalScore > players[j].TotalScore
			})

			// Take only the top 10 (or less if not enough players)
			topCount := topPlayersCount
			if len(players) < topCount {
				topCount = len(players)
			}
			topPlayers := players[:topCount]

			// Calculate metrics based on top players
			var totalPlayerScores float64
			var maxPlayerScore float64 = 0
			var minPlayerScore float64 = math.MaxFloat64
			var totalAvgKeyLevel float64
			var maxKeyLevel int = 0
			var minKeyLevel int = math.MaxInt32

			// Iterate through top players of this spec
			for _, p := range topPlayers {
				totalPlayerScores += p.TotalScore
				totalAvgKeyLevel += p.AvgKeyLevel

				// Find the player with the highest score
				if p.TotalScore > maxPlayerScore {
					maxPlayerScore = p.TotalScore
				}

				// Find the player with the lowest score
				if p.TotalScore < minPlayerScore {
					minPlayerScore = p.TotalScore
				}

				// Find the max key level
				if p.MaxKeyLevel > maxKeyLevel {
					maxKeyLevel = p.MaxKeyLevel
				}

				// Find the min key level
				if p.MinKeyLevel < minKeyLevel {
					minKeyLevel = p.MinKeyLevel
				}
			}

			if topCount > 0 {
				globalMetrics = append(globalMetrics, playerRankingModels.DailySpecMetricMythicPlus{
					CaptureDate: processingDate,
					Spec:        spec,
					Class:       class,
					Role:        role,
					EncounterID: 0, // 0 for global metric
					IsGlobal:    true,
					AvgScore:    totalPlayerScores / float64(topCount),
					MaxScore:    maxPlayerScore,
					MinScore:    minPlayerScore,
					AvgKeyLevel: totalAvgKeyLevel / float64(topCount),
					MaxKeyLevel: maxKeyLevel,
					MinKeyLevel: minKeyLevel,
					RoleRank:    0,
					OverallRank: 0,
				})
			}
		}

		log.Printf("Calculated global metrics for %d specializations", len(globalMetrics))

		// 5. Calculate rankings for the global metrics (only for global metrics)
		calculateGlobalMetricsRankings(&globalMetrics)

		// 6. Persist all metrics
		allMetrics := append(metrics, globalMetrics...)
		if err := tx.CreateInBatches(allMetrics, 100).Error; err != nil {
			return fmt.Errorf("error persisting metrics: %w", err)
		}

		log.Printf("Successfully stored %d metrics records", len(allMetrics))
		return nil
	})
}

// calculateDungeonMetricsRankings calculates rankings separately for each dungeon
func calculateDungeonMetricsRankings(metrics *[]playerRankingModels.DailySpecMetricMythicPlus) {
	// Group metrics by dungeon
	dungeonGroups := make(map[int][]int) // map[dungeonID][]metricIndex

	for i, metric := range *metrics {
		dungeonGroups[metric.EncounterID] = append(dungeonGroups[metric.EncounterID], i)
	}

	// For each dungeon, independently calculate the rankings
	for _, indices := range dungeonGroups {
		// Global ranking based on avgScore (only for this dungeon)
		sort.Slice(indices, func(i, j int) bool {
			return (*metrics)[indices[i]].AvgScore > (*metrics)[indices[j]].AvgScore
		})

		// Assign global rankings for this dungeon
		for rank, idx := range indices {
			(*metrics)[idx].OverallRank = rank + 1
		}

		// Role ranking for this dungeon
		roleGroups := make(map[string][]int) // map[role][]metricIndex

		for _, idx := range indices {
			role := (*metrics)[idx].Role
			roleGroups[role] = append(roleGroups[role], idx)
		}

		// For each role in this dungeon
		for _, roleIndices := range roleGroups {
			// Sort by avgScore
			sort.Slice(roleIndices, func(i, j int) bool {
				return (*metrics)[roleIndices[i]].AvgScore > (*metrics)[roleIndices[j]].AvgScore
			})

			// Assign rankings for this role and dungeon
			for rank, idx := range roleIndices {
				(*metrics)[idx].RoleRank = rank + 1
			}
		}
	}
}

// calculateGlobalMetricsRankings calculates rankings only for global metrics
func calculateGlobalMetricsRankings(metrics *[]playerRankingModels.DailySpecMetricMythicPlus) {
	// Global ranking based on avgScore
	sort.Slice(*metrics, func(i, j int) bool {
		return (*metrics)[i].AvgScore > (*metrics)[j].AvgScore
	})

	for i := range *metrics {
		(*metrics)[i].OverallRank = i + 1
	}

	// Role ranking
	roleGroups := make(map[string][]int)

	for i, metric := range *metrics {
		roleGroups[metric.Role] = append(roleGroups[metric.Role], i)
	}

	for _, indices := range roleGroups {
		// Sort indices by avgScore
		sort.Slice(indices, func(i, j int) bool {
			return (*metrics)[indices[i]].AvgScore > (*metrics)[indices[j]].AvgScore
		})

		// Assign rankings
		for rank, idx := range indices {
			(*metrics)[idx].RoleRank = rank + 1
		}
	}
}

// GetGlobalRankings retrieves global rankings by role
func (r *PlayerRankingsRepository) GetGlobalRankings(ctx context.Context) (*playerRankingModels.GlobalRankings, error) {
	log.Println("Retrieving global rankings")

	// Structure to store results
	result := &playerRankingModels.GlobalRankings{
		Tanks:   playerRankingModels.RoleRankings{Players: make([]playerRankingModels.PlayerScore, 0)},
		Healers: playerRankingModels.RoleRankings{Players: make([]playerRankingModels.PlayerScore, 0)},
		DPS:     playerRankingModels.RoleRankings{Players: make([]playerRankingModels.PlayerScore, 0)},
	}

	// Structure to store player scores
	type PlayerData struct {
		Name       string
		Class      string
		Spec       string
		Role       string
		TotalScore float64
		GuildID    int
		GuildName  string
		Faction    int
		ServerID   int
		ServerName string
		Region     string
	}

	// Retrieve all players with their total scores
	var playersData []PlayerData
	err := r.db.WithContext(ctx).Raw(`
		SELECT
    name,
    class,
    spec,
    role,
    SUM(best_score) as total_score,
    guild_id,
    guild_name,
    guild_faction as faction,
    server_id,
    server_name,
    server_region as region
FROM (
    SELECT
        name,
        class,
        spec,
        role,
        dungeon_id,
        MAX(score) as best_score,
        guild_id,
        guild_name,
        guild_faction,
        server_id,
        server_name,
        server_region
    FROM player_rankings
    GROUP BY name, class, spec, role, dungeon_id, guild_id, guild_name, guild_faction, server_id, server_name, server_region
) as best_runs
GROUP BY name, class, spec, role, guild_id, guild_name, guild_faction, server_id, server_name, server_region
ORDER BY total_score DESC
	`).Scan(&playersData).Error

	if err != nil {
		return nil, fmt.Errorf("error retrieving player data: %w", err)
	}

	// Structure to store unique runs per player
	type RunData struct {
		Name        string
		DungeonID   int
		Score       float64
		Duration    int64
		StartTime   int64
		HardMode    int
		BracketData int
		Medal       string
		ReportCode  string
		FightID     int
		ReportTime  int64
		Affixes     pq.Int64Array
	}

	// Retrieve all runs
	var runsData []RunData
	err = r.db.WithContext(ctx).Raw(`
		SELECT
			name,
			dungeon_id,
			score,
			duration,
			start_time,
			hard_mode_level as hard_mode,
			bracket_data,
			medal,
			report_code,
			report_fight_id as fight_id,
			report_start_time as report_time,
			affixes
		FROM player_rankings
	`).Scan(&runsData).Error

	if err != nil {
		return nil, fmt.Errorf("error retrieving run data: %w", err)
	}

	// Group runs by player
	runsMap := make(map[string][]playerRankingModels.Run)
	for _, run := range runsData {
		r := playerRankingModels.Run{
			DungeonID:     run.DungeonID,
			Score:         run.Score,
			Duration:      run.Duration,
			StartTime:     run.StartTime,
			HardModeLevel: run.HardMode,
			BracketData:   run.BracketData,
			Medal:         run.Medal,
			Affixes:       []int{}, // Convert affixes
			Report: playerRankingModels.Report{
				Code:      run.ReportCode,
				FightID:   run.FightID,
				StartTime: run.ReportTime,
			},
		}

		// Convert affixes from pq.Int64Array to []int
		if run.Affixes != nil {
			for _, affix := range run.Affixes {
				r.Affixes = append(r.Affixes, int(affix))
			}
		}

		runsMap[run.Name] = append(runsMap[run.Name], r)
	}

	// Build PlayerScore objects and add them to corresponding roles
	for _, player := range playersData {
		// Create PlayerScore object
		playerScore := playerRankingModels.PlayerScore{
			Name:  player.Name,
			Class: player.Class,
			Spec:  player.Spec,
			Role:  player.Role,
			Guild: playerRankingModels.Guild{
				ID:      player.GuildID,
				Name:    player.GuildName,
				Faction: player.Faction,
			},
			Server: playerRankingModels.Server{
				ID:     player.ServerID,
				Name:   player.ServerName,
				Region: player.Region,
			},
			Faction:    player.Faction,
			TotalScore: player.TotalScore,
			Runs:       runsMap[player.Name],
		}

		// Add player to corresponding role
		switch player.Role {
		case "Tank":
			result.Tanks.Players = append(result.Tanks.Players, playerScore)
		case "Healer":
			result.Healers.Players = append(result.Healers.Players, playerScore)
		case "DPS":
			result.DPS.Players = append(result.DPS.Players, playerScore)
		}
	}

	// Sort players by TotalScore in each role
	sort.Slice(result.Tanks.Players, func(i, j int) bool {
		return result.Tanks.Players[i].TotalScore > result.Tanks.Players[j].TotalScore
	})
	sort.Slice(result.Healers.Players, func(i, j int) bool {
		return result.Healers.Players[i].TotalScore > result.Healers.Players[j].TotalScore
	})
	sort.Slice(result.DPS.Players, func(i, j int) bool {
		return result.DPS.Players[i].TotalScore > result.DPS.Players[j].TotalScore
	})

	// Update counters
	result.Tanks.Count = len(result.Tanks.Players)
	result.Healers.Count = len(result.Healers.Players)
	result.DPS.Count = len(result.DPS.Players)

	log.Printf("Rankings retrieval completed: Tanks: %d, Healers: %d, DPS: %d",
		result.Tanks.Count,
		result.Healers.Count,
		result.DPS.Count)

	return result, nil
}
