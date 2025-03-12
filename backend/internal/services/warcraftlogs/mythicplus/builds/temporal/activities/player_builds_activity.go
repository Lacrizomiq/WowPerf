package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"go.temporal.io/sdk/activity"
	"gorm.io/datatypes"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"
)

// PlayerDetails represents detailed player information from WarcraftLogs
type PlayerDetails struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	Specs         []string        `json:"specs"`
	MaxItemLevel  float64         `json:"maxItemLevel"`
	CombatantInfo json.RawMessage `json:"combatantInfo"`
}

// PlayerBuildsActivity manages all operations related to player builds
type PlayerBuildsActivity struct {
	repository *playerBuildsRepository.PlayerBuildsRepository
}

func NewPlayerBuildsActivity(repository *playerBuildsRepository.PlayerBuildsRepository) *PlayerBuildsActivity {
	return &PlayerBuildsActivity{
		repository: repository,
	}
}

// ProcessAllBuilds processes player builds in parallel from the provided reports
func (a *PlayerBuildsActivity) ProcessAllBuilds(ctx context.Context, reports []*warcraftlogsBuilds.Report) (*workflows.BuildsProcessingResult, error) {
	logger := activity.GetLogger(ctx)
	result := &workflows.BuildsProcessingResult{
		ProcessedAt: time.Now(),
	}

	if len(reports) == 0 {
		logger.Info("No reports to process")
		return result, nil
	}

	// Parallel processing configuration
	const (
		numWorkers      = 4  // Number of parallel workers
		reportBatchSize = 10 // Size of report batches
	)

	// Channels for communication
	reportChan := make(chan *warcraftlogsBuilds.Report)
	resultChan := make(chan struct {
		builds int
		err    error
	})

	// Start workers
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for report := range reportChan {
				// Record heartbeat for each report being processed
				activity.RecordHeartbeat(ctx, map[string]interface{}{
					"workerID":   workerID,
					"reportCode": report.Code,
					"fightID":    report.FightID,
					"status":     "processing",
				})

				logger.Info("Processing report",
					"workerID", workerID,
					"reportCode", report.Code,
					"fightID", report.FightID)

				builds, err := a.extractPlayerBuilds(report)
				if err != nil {
					logger.Error("Failed to extract builds",
						"reportCode", report.Code,
						"error", err)
					resultChan <- struct {
						builds int
						err    error
					}{0, err}
					continue
				}

				// Store extracted builds
				if len(builds) > 0 {
					if err := a.repository.StoreManyPlayerBuilds(ctx, builds); err != nil {
						logger.Error("Failed to store builds",
							"reportCode", report.Code,
							"buildsCount", len(builds),
							"error", err)
						resultChan <- struct {
							builds int
							err    error
						}{0, err}
						continue
					}

					// Record heartbeat after successful build storage
					activity.RecordHeartbeat(ctx, map[string]interface{}{
						"workerID":     workerID,
						"reportCode":   report.Code,
						"fightID":      report.FightID,
						"status":       "completed",
						"buildsStored": len(builds),
						"completedAt":  time.Now(),
					})

					logger.Info("Successfully processed report",
						"reportCode", report.Code,
						"buildsProcessed", len(builds))
				}

				resultChan <- struct {
					builds int
					err    error
				}{len(builds), nil}
			}
		}(i)
	}

	// Send reports to workers
	go func() {
		totalReports := len(reports)
		for i := 0; i < len(reports); i += reportBatchSize {
			// Record heartbeat for batch progress
			activity.RecordHeartbeat(ctx, map[string]interface{}{
				"status":         "batch_processing",
				"processedCount": i,
				"totalReports":   totalReports,
				"progress":       fmt.Sprintf("%d/%d", i, totalReports),
			})

			end := i + reportBatchSize
			if end > len(reports) {
				end = len(reports)
			}

			for _, report := range reports[i:end] {
				select {
				case <-ctx.Done():
					return
				case reportChan <- report:
				}
			}
			// Delay between batches
			time.Sleep(100 * time.Millisecond)
		}
		close(reportChan)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Process results
	for res := range resultChan {
		if res.err != nil {
			result.FailureCount++
		} else {
			result.SuccessCount++
			result.ProcessedBuilds += res.builds
		}

		// Record heartbeat for overall progress
		activity.RecordHeartbeat(ctx, map[string]interface{}{
			"status":       "progress",
			"totalBuilds":  result.ProcessedBuilds,
			"successCount": result.SuccessCount,
			"failureCount": result.FailureCount,
		})
	}

	logger.Info("Completed processing all reports",
		"totalBuilds", result.ProcessedBuilds,
		"successCount", result.SuccessCount,
		"failureCount", result.FailureCount)

	return result, nil
}

// extractPlayerBuilds extracts all player builds from a report
func (a *PlayerBuildsActivity) extractPlayerBuilds(report *warcraftlogsBuilds.Report) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	var builds []*warcraftlogsBuilds.PlayerBuild

	// Extraction from each role (DPS, Healers, Tanks)
	dpsBuilds, err := a.extractBuildsFromPlayers(report, report.PlayerDetailsDps)
	if err != nil {
		log.Printf("Error extracting DPS builds: %v", err)
	}
	builds = append(builds, dpsBuilds...)

	healerBuilds, err := a.extractBuildsFromPlayers(report, report.PlayerDetailsHealers)
	if err != nil {
		log.Printf("Error extracting healer builds: %v", err)
	}
	builds = append(builds, healerBuilds...)

	tankBuilds, err := a.extractBuildsFromPlayers(report, report.PlayerDetailsTanks)
	if err != nil {
		log.Printf("Error extracting tank builds: %v", err)
	}
	builds = append(builds, tankBuilds...)

	return builds, nil
}

// extractBuildsFromPlayers processes player data and creates builds
func (a *PlayerBuildsActivity) extractBuildsFromPlayers(report *warcraftlogsBuilds.Report, playerData []byte) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	if len(playerData) == 0 {
		return nil, nil
	}

	var players []PlayerDetails
	if err := json.Unmarshal(playerData, &players); err != nil {
		return nil, fmt.Errorf("failed to unmarshal player details: %w", err)
	}

	var builds []*warcraftlogsBuilds.PlayerBuild
	for _, player := range players {
		build, err := a.createPlayerBuild(report, player)
		if err != nil {
			log.Printf("Error creating build for player %s: %v", player.Name, err)
			continue
		}
		if build != nil {
			builds = append(builds, build)
		}
	}

	return builds, nil
}

// createPlayerBuild creates a PlayerBuild from player details
func (a *PlayerBuildsActivity) createPlayerBuild(report *warcraftlogsBuilds.Report, player PlayerDetails) (*warcraftlogsBuilds.PlayerBuild, error) {
	var talentCodes map[string]string
	if err := json.Unmarshal(report.TalentCodes, &talentCodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal talent codes: %w", err)
	}

	// Determine the active spec
	var activeSpec string
	for _, spec := range player.Specs {
		talentKey := fmt.Sprintf("%s_%s_talents", player.Type, spec)
		if _, exists := talentCodes[talentKey]; exists {
			activeSpec = spec
			break
		}
	}

	// Fallback on the first spec if no active spec found
	if activeSpec == "" && len(player.Specs) > 0 {
		activeSpec = player.Specs[0]
	}

	if activeSpec == "" {
		return nil, fmt.Errorf("no valid spec found for player %s", player.Name)
	}

	// Retrieve talent code for the active spec
	talentCode := talentCodes[fmt.Sprintf("%s_%s_talents", player.Type, activeSpec)]

	// Parse combat information
	var combatInfo struct {
		Stats      json.RawMessage `json:"stats"`
		Gear       json.RawMessage `json:"gear"`
		TalentTree json.RawMessage `json:"talentTree"`
	}

	// Handle two possible formats for combatant info
	var err error
	if err = json.Unmarshal(player.CombatantInfo, &combatInfo); err != nil {
		var combatInfoArray []json.RawMessage
		if errArray := json.Unmarshal(player.CombatantInfo, &combatInfoArray); errArray == nil && len(combatInfoArray) > 0 {
			err = json.Unmarshal(combatInfoArray[0], &combatInfo)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal combatant info: %w", err)
		}
	}

	build := &warcraftlogsBuilds.PlayerBuild{
		PlayerName:    player.Name,
		Class:         player.Type,
		Spec:          activeSpec,
		ReportCode:    report.Code,
		FightID:       report.FightID,
		ActorID:       player.ID,
		ItemLevel:     player.MaxItemLevel,
		TalentImport:  talentCode,
		TalentTree:    datatypes.JSON(combatInfo.TalentTree),
		Gear:          datatypes.JSON(combatInfo.Gear),
		Stats:         datatypes.JSON(combatInfo.Stats),
		EncounterID:   report.EncounterID,
		KeystoneLevel: report.KeystoneLevel,
		Affixes:       report.Affixes,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return build, nil
}

// CountPlayerBuilds returns the total number of player builds in the database
func (a *PlayerBuildsActivity) CountPlayerBuilds(ctx context.Context) (int64, error) {
	return a.repository.CountPlayerBuilds(ctx)
}
