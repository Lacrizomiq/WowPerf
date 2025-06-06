package warcraftlogs

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"
	service "wowperf/internal/services/warcraftlogs"
	middleware "wowperf/middleware/cache"
	cache "wowperf/pkg/cache"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

const (
	MinimumUpdateInterval = 20 * time.Hour
	DefaultUpdateInterval = 24 * time.Hour
	updateLockKey         = "warcraftlogs:rankings:update:lock"
	batchSize             = 100
)

// RankingsUpdater is responsible for updating rankings in the database and calculating daily spec metric
type RankingsUpdater struct {
	db           *gorm.DB
	service      *service.WarcraftLogsClientService
	cache        cache.CacheService
	cacheManager *middleware.CacheManager
}

func NewRankingsUpdater(db *gorm.DB, service *service.WarcraftLogsClientService, cache cache.CacheService, cacheManager *middleware.CacheManager) *RankingsUpdater {
	return &RankingsUpdater{
		db:           db,
		service:      service,
		cache:        cache,
		cacheManager: cacheManager,
	}
}

// StartPeriodicUpdate starts the periodic updates
func (r *RankingsUpdater) StartPeriodicUpdate(ctx context.Context) {
	log.Println("Starting WarcraftLogs rankings periodic update...")

	// Initial immediate check
	log.Println("Performing initial check...")
	if err := r.checkAndUpdate(ctx); err != nil {
		log.Printf("Initial check error: %v", err)
	}

	ticker := time.NewTicker(MinimumUpdateInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				log.Println("Ticker triggered, checking for updates...")
				if err := r.checkAndUpdate(ctx); err != nil {
					log.Printf("Periodic check error: %v", err)
				}
			}
		}
	}()
}

// initializeUpdateState initializes the update state if necessary
func (r *RankingsUpdater) initializeUpdateState() error {
	var state playerRankingModels.RankingsUpdateState
	result := r.db.First(&state)

	if result.Error == gorm.ErrRecordNotFound {
		state = playerRankingModels.RankingsUpdateState{
			LastUpdateTime: time.Now(),
		}
		if err := r.db.Create(&state).Error; err != nil {
			return fmt.Errorf("failed to create initial state: %w", err)
		}
		log.Println("Created initial rankings update state")
	} else if result.Error != nil {
		return fmt.Errorf("failed to check update state: %w", result.Error)
	}

	return nil
}

func (r *RankingsUpdater) updateLastUpdateTime(ctx context.Context) error {
	return r.db.Exec(`
			INSERT INTO rankings_update_states (id, last_update_time, updated_at)
			VALUES (1, NOW(), NOW())
			ON CONFLICT (id) DO UPDATE 
			SET last_update_time = EXCLUDED.last_update_time,
					updated_at = EXCLUDED.updated_at
	`).Error
}

// checkAndUpdate checks and performs the update if necessary
func (r *RankingsUpdater) checkAndUpdate(ctx context.Context) error {
	// Distributed lock check
	locked, err := r.cache.SetNX(ctx, updateLockKey, time.Now().String(), 1*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to check update lock: %w", err)
	}
	if !locked {
		return fmt.Errorf("update already in progress")
	}
	defer r.cache.Delete(ctx, updateLockKey)

	// Unique state retrieval
	state, err := playerRankingModels.GetOrCreateRankingsUpdateState(r.db)
	if err != nil {
		return fmt.Errorf("failed to get rankings state: %w", err)
	}

	timeSinceLastUpdate := time.Since(state.LastUpdateTime)
	log.Printf("Time since last rankings update: %v", timeSinceLastUpdate)

	if timeSinceLastUpdate >= MinimumUpdateInterval {
		log.Printf("Rankings update needed, last update was %v ago", timeSinceLastUpdate)

		if err := r.UpdateRankings(ctx); err != nil {
			return fmt.Errorf("failed to update rankings: %w", err)
		}

		if err := r.updateLastUpdateTime(ctx); err != nil {
			return fmt.Errorf("failed to update timestamp: %w", err)
		}

		log.Println("Rankings update completed successfully")
	} else {
		log.Printf("Skipping update: last update was %v ago", timeSinceLastUpdate)
	}

	return nil
}

// UpdateRankings updates the rankings in the database
func (r *RankingsUpdater) UpdateRankings(ctx context.Context) error {
	dungeonIDs := []int{
		DungeonCinderbrew,
		DungeonDarkflame,
		DungeonFloodgate,
		DungeonMechagon,
		DungeonPriory,
		DungeonMotherlode,
		DungeonRookery,
		DungeonTheaterPain,
	}
	pagesPerDungeon := 2

	rankings, err := GetGlobalRankings(r.service, ctx, dungeonIDs, pagesPerDungeon)
	if err != nil {
		return fmt.Errorf("failed to get global rankings: %w", err)
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing data
		if err := tx.Exec("DELETE FROM player_rankings").Error; err != nil {
			return fmt.Errorf("failed to delete existing rankings: %w", err)
		}

		// Prepare new data
		var newRankings []playerRankingModels.PlayerRanking
		processRoleRankings := func(players []PlayerScore) {
			for _, player := range players {
				for _, run := range player.Runs {
					newRankings = append(newRankings, playerRankingModels.PlayerRanking{
						DungeonID:       run.DungeonID,
						Name:            player.Name,
						Class:           player.Class,
						Spec:            player.Spec,
						Role:            player.Role,
						Amount:          player.Amount,
						HardModeLevel:   run.HardModeLevel,
						Duration:        run.Duration,
						StartTime:       run.StartTime,
						ReportCode:      run.Report.Code,
						ReportFightID:   run.Report.FightID,
						ReportStartTime: run.Report.StartTime,
						GuildID:         player.Guild.ID,
						GuildName:       player.Guild.Name,
						GuildFaction:    player.Guild.Faction,
						ServerID:        player.Server.ID,
						ServerName:      player.Server.Name,
						ServerRegion:    player.Server.Region,
						BracketData:     run.BracketData,
						Faction:         player.Faction,
						Affixes:         run.Affixes,
						Medal:           run.Medal,
						Score:           run.Score,
						Leaderboard:     0,
					})
				}
			}
		}

		processRoleRankings(rankings.Tanks.Players)
		processRoleRankings(rankings.Healers.Players)
		processRoleRankings(rankings.DPS.Players)

		// Insert new data by batches
		return r.insertRankingsInBatches(tx, newRankings)
	})

	if err != nil {
		return err
	}

	// After updating the rankings, calculate the daily spec metrics
	log.Println("Rankings updated successfully, now calculating daily spec metrics...")
	if err := r.calculateDailySpecMetrics(ctx); err != nil {
		return fmt.Errorf("failed to calculate daily spec metrics: %w", err)
	}

	// Invalidate the caches
	if err := r.invalidateCache(ctx); err != nil {
		log.Printf("Warning: Failed to invalidate caches: %v", err)
	}

	return nil
}

// insertRankingsInBatches inserts the rankings by batches
func (r *RankingsUpdater) insertRankingsInBatches(tx *gorm.DB, rankings []playerRankingModels.PlayerRanking) error {
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
		if err := tx.Exec(query, valueArgs...).Error; err != nil {
			return fmt.Errorf("failed to insert batch %d: %w", i+1, err)
		}

		log.Printf("Inserted batch %d/%d (%d rankings)", i+1, batches, len(batch))
	}

	return nil
}

// invalidateCache invalidates the caches related to rankings
func (r *RankingsUpdater) invalidateCache(ctx context.Context) error {
	tags := []string{
		"rankings",
		"leaderboard",
		"global-rankings",
		"class-rankings",
		"spec-rankings",
		"daily-spec-metrics", // Added to invalidate also the metrics caches
	}
	return r.cacheManager.InvalidateByTags(ctx, tags)
}

// calculateDailySpecMetrics calculates the daily metrics for all specializations
func (r *RankingsUpdater) calculateDailySpecMetrics(ctx context.Context) error {
	// Unique date for the entire processing
	processingDate := time.Now().Truncate(24 * time.Hour)
	log.Printf("Calculating spec metrics for date: %s", processingDate.Format("2006-01-02"))

	// Fixed number of dungeons (always 8)
	const totalDungeonCount = 8
	// Number of top players to consider for averages
	const topPlayersCount = 10

	// Use a transaction to ensure integrity
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Delete existing metrics for this date (Hard delete to avoid problems)
		if err := tx.Exec("DELETE FROM daily_spec_metrics_mythic_plus WHERE capture_date = ?", processingDate).Error; err != nil {
			return fmt.Errorf("error deleting existing metrics: %w", err)
		}

		// 2. Calculate the metrics per dungeon in a single query
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
    WITH ranked_players AS (
        SELECT 
            spec,
            class,
            role,
            dungeon_id,
            score,
            hard_mode_level,
            ROW_NUMBER() OVER (
                PARTITION BY spec, class, role, dungeon_id 
                ORDER BY score DESC
            ) as player_rank
        FROM player_rankings
        WHERE server_region != 'CN'
    )
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
    FROM ranked_players
    WHERE player_rank <= 10
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

		// 3. Calculate the rankings for the dungeon metrics (groupé par donjon)
		r.calculateDungeonMetricsRankings(&metrics)

		// 4. Calculate the global metrics based on TOP player performances
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
				SUM(score) AS total_score,
				AVG(hard_mode_level) AS avg_key_level,
				MAX(hard_mode_level) AS max_key_level,
				MIN(hard_mode_level) AS min_key_level,
				COUNT(DISTINCT dungeon_id) AS dungeon_count
			FROM player_rankings
			WHERE server_region != 'CN'
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

		// 5. Calculate the rankings for the global metrics (only for global metrics)
		r.calculateGlobalMetricsRankings(&globalMetrics)

		// 6. Persist all metrics
		allMetrics := append(metrics, globalMetrics...)
		if err := tx.CreateInBatches(allMetrics, 100).Error; err != nil {
			return fmt.Errorf("error persisting metrics: %w", err)
		}

		log.Printf("Successfully stored %d metrics records", len(allMetrics))
		return nil
	})
}

// calculateDungeonMetricsRankings calculates the rankings separately for each dungeon
func (r *RankingsUpdater) calculateDungeonMetricsRankings(metrics *[]playerRankingModels.DailySpecMetricMythicPlus) {
	// Group metrics by dungeon
	dungeonGroups := make(map[int][]int) // map[dungeonID][]metricIndex

	for i, metric := range *metrics {
		dungeonGroups[metric.EncounterID] = append(dungeonGroups[metric.EncounterID], i)
	}

	// For each dungeon, calculate independently the rankings
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

// calculateGlobalMetricsRankings calculates the rankings only for the global metrics
func (r *RankingsUpdater) calculateGlobalMetricsRankings(metrics *[]playerRankingModels.DailySpecMetricMythicPlus) {
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
