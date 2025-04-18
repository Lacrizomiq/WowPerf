package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
	"gorm.io/datatypes"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
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
	repository        *playerBuildsRepository.PlayerBuildsRepository
	reportsRepository *reportsRepository.ReportRepository
}

func NewPlayerBuildsActivity(repository *playerBuildsRepository.PlayerBuildsRepository, reportsRepository *reportsRepository.ReportRepository) *PlayerBuildsActivity {
	return &PlayerBuildsActivity{
		repository:        repository,
		reportsRepository: reportsRepository,
	}
}

// ProcessAllBuilds processes player builds AND marks corresponding reports as processed using ReportIdentifier.
// This activity is designed to be called by the BuildsWorkflow.
func (a *PlayerBuildsActivity) ProcessAllBuilds(ctx context.Context, reports []*warcraftlogsBuilds.Report) (*models.BuildsActivityResult, error) {
	logger := activity.GetLogger(ctx)
	activityStart := time.Now() // Keep track of the start time for duration
	result := &models.BuildsActivityResult{
		ProcessedAt:       activityStart, // Initial timestamp
		BuildsByClassSpec: make(map[string]int32),
	}

	if len(reports) == 0 {
		logger.Info("No reports received in ProcessAllBuilds")
		result.ProcessedAt = time.Now()
		return result, nil
	}

	// Generate a unique batchID for this specific activity execution
	// This batchID will be used to mark the reports.
	batchID := fmt.Sprintf("builds-activity-%s", uuid.New().String())

	logger.Info("Starting ProcessAllBuilds activity", "reportCount", len(reports), "activityBatchID", batchID)

	// Parallel processing logic (workers, channels)
	const (
		numWorkers = 4
		// Note: reportBatchSize here controls the distribution to workers,
		// it is NOT the size of the 'reports' parameter received.
		reportBatchSizeInternal = 10
	)
	reportChan := make(chan *warcraftlogsBuilds.Report, len(reports))
	// The result channel now contains the entire report for identification
	resultChan := make(chan struct {
		report *warcraftlogsBuilds.Report // <- Keep the entire report
		builds []*warcraftlogsBuilds.PlayerBuild
		err    error
	}, len(reports))

	var wg sync.WaitGroup
	processingCtx, cancel := context.WithCancel(ctx) // Context to cancel workers
	defer cancel()

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for report := range reportChan {
				select {
				case <-processingCtx.Done():
					logger.Warn("Worker cancelled", "workerID", workerID)
					return
				default:
					activity.RecordHeartbeat(processingCtx, map[string]interface{}{"workerID": workerID, "reportCode": report.Code, "status": "processing"})
					logger.Info("Processing report", "workerID", workerID, "reportCode", report.Code, "fightID", report.FightID)

					builds, errExtract := a.extractPlayerBuilds(report)
					if errExtract != nil {
						logger.Error("Failed to extract builds", "reportCode", report.Code, "error", errExtract)
						resultChan <- struct {
							report *warcraftlogsBuilds.Report
							builds []*warcraftlogsBuilds.PlayerBuild
							err    error
						}{report, nil, errExtract}
						continue
					}

					if len(builds) == 0 {
						logger.Info("No builds extracted (but no error)", "reportCode", report.Code)
						resultChan <- struct {
							report *warcraftlogsBuilds.Report
							builds []*warcraftlogsBuilds.PlayerBuild
							err    error
						}{report, builds, nil}
						continue
					}

					errStore := a.repository.StoreManyPlayerBuilds(processingCtx, builds)
					if errStore != nil {
						logger.Error("Failed to store builds", "reportCode", report.Code, "buildsCount", len(builds), "error", errStore)
						resultChan <- struct {
							report *warcraftlogsBuilds.Report
							builds []*warcraftlogsBuilds.PlayerBuild
							err    error
						}{report, nil, errStore}
						continue
					}

					activity.RecordHeartbeat(processingCtx, map[string]interface{}{"workerID": workerID, "reportCode": report.Code, "status": "completed", "buildsStored": len(builds)})
					logger.Info("Successfully stored builds for report", "reportCode", report.Code, "buildsProcessed", len(builds))

					resultChan <- struct {
						report *warcraftlogsBuilds.Report
						builds []*warcraftlogsBuilds.PlayerBuild
						err    error
					}{report, builds, nil} // Succès
				}
			}
		}(i)
	}

	// Distribute the work
	go func() {
		defer close(reportChan)
		totalReports := len(reports)
		for i := 0; i < totalReports; i += reportBatchSizeInternal {
			// No need for heartbeat here, already done in the collection
			end := i + reportBatchSizeInternal
			if end > totalReports {
				end = totalReports
			}
			for _, report := range reports[i:end] {
				select {
				case <-processingCtx.Done():
					return
				case reportChan <- report:
				}
			}
			// Optional: time.Sleep(50 * time.Millisecond) if distribution too fast
		}
	}()

	// Wait for workers to finish and close the result channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Processing results and aggregation
	buildsByClassSpec := make(map[string]int32)
	successfulReportIdentifiers := make([]reportsRepository.ReportIdentifier, 0, len(reports)) // <- Slice for identifiers
	var firstProcessingError error                                                             // To store the first extraction/storage error

	processedReportCount := 0
	totalReportCount := len(reports)

	for res := range resultChan {
		processedReportCount++
		if res.err != nil {
			result.FailureCount++
			if firstProcessingError == nil {
				firstProcessingError = res.err // Register the first internal error
			}
			logger.Warn("Processing failed for report", "reportCode", res.report.Code, "fightID", res.report.FightID, "error", res.err)
		} else {
			// Successful extraction/storage for this report
			result.SuccessCount++
			result.ProcessedBuildsCount += int32(len(res.builds))

			// ADDITION: Add the identifier (Code+FightID) for marking
			successfulReportIdentifiers = append(successfulReportIdentifiers, reportsRepository.ReportIdentifier{
				Code:    res.report.Code,
				FightID: res.report.FightID,
			})

			// AGGREGATION by class/spec
			for _, build := range res.builds {
				classSpecKey := fmt.Sprintf("%s-%s", build.Class, build.Spec)
				buildsByClassSpec[classSpecKey]++
			}
		}

		// Global progress heartbeat
		activity.RecordHeartbeat(ctx, map[string]interface{}{
			"status":           "progress_collection",
			"reportsCollected": processedReportCount,
			"totalReports":     totalReportCount,
			"currentSuccess":   result.SuccessCount,
			"currentFailure":   result.FailureCount,
			"totalBuildsSoFar": result.ProcessedBuildsCount,
		})
	}
	result.BuildsByClassSpec = buildsByClassSpec

	logger.Info("Completed collecting results from workers",
		"totalBuilds", result.ProcessedBuildsCount,
		"successReports", result.SuccessCount,
		"failedReports", result.FailureCount)

	//  Marking successfully processed reports
	var markingErr error // To store the marking error
	if len(successfulReportIdentifiers) > 0 {
		logger.Info("Marking successfully processed reports as completed for builds",
			"count", len(successfulReportIdentifiers),
			"activityBatchID", batchID)

		if a.reportsRepository == nil {
			errMsg := "internal error: reportsRepository not injected in PlayerBuildsActivity"
			logger.Error(errMsg)
			markingErr = fmt.Errorf(errMsg) // Critical error
		} else {
			// Use the activity's main context (ctx)
			// CORRECT CALL: Pass the slice of ReportIdentifier and the activity's batchID
			err := a.reportsRepository.MarkReportsAsProcessedForBuilds(ctx, successfulReportIdentifiers, batchID)
			if err != nil {
				logger.Error("Failed to mark reports as processed for builds", "error", err, "activityBatchID", batchID)
				markingErr = err // Store the marking error
				// Optional: Adjust counters if marking fails
				result.FailureCount += result.SuccessCount
				result.SuccessCount = 0
			} else {
				logger.Info("Successfully marked reports as processed for builds", "count", len(successfulReportIdentifiers), "activityBatchID", batchID)
			}
		}
	} else {
		logger.Info("No reports were successfully processed to be marked.")
	}

	// Finalization and Return
	result.ProcessedAt = time.Now() // Final timestamp of the activity
	logger.Info("Finished activity ProcessAllBuilds",
		"duration", result.ProcessedAt.Sub(activityStart),
		"buildsStored", result.ProcessedBuildsCount,
		"successReports", result.SuccessCount,
		"failedReports", result.FailureCount)

	if markingErr != nil {
		return result, markingErr
	}

	return result, nil // The activity has finished its cycle, returns nil as main error
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

// GetReportsNeedingBuildExtraction retrieves reports that need build extraction
func (a *PlayerBuildsActivity) GetReportsNeedingBuildExtraction(ctx context.Context, limit int32, maxAgeDuration time.Duration) ([]*warcraftlogsBuilds.Report, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Getting reports needing build extraction", "limit", limit)

	// If maxAgeDuration is 0, set a default (e.g., 10 days)
	if maxAgeDuration == 0 {
		maxAgeDuration = 10 * 24 * time.Hour // 10 days
		logger.Info("Using default maxAge", "maxAge", maxAgeDuration)
	}

	if a.reportsRepository == nil {
		logger.Error("ReportRepository is not initialized in PlayerBuildsActivity")
		return nil, fmt.Errorf("internal error: reportRepository not injected")
	}

	reports, err := a.reportsRepository.GetReportsNeedingBuildExtraction(ctx, int(limit), maxAgeDuration)
	if err != nil {
		logger.Error("Failed to get reports needing build extraction from repository", "error", err)
		return nil, err
	}

	logger.Info("Retrieved reports needing build extraction", "count", len(reports))
	return reports, nil
}

// MarkReportsAsProcessedForBuilds updates the status of reports after build extraction
func (a *PlayerBuildsActivity) MarkReportsAsProcessedForBuilds(ctx context.Context, reports []*warcraftlogsBuilds.Report, batchID string) error {
	logger := activity.GetLogger(ctx)

	identifiers := make([]reportsRepository.ReportIdentifier, 0, len(reports))
	for _, report := range reports {
		identifiers = append(identifiers, reportsRepository.ReportIdentifier{
			Code:    report.Code,
			FightID: report.FightID,
		})
	}

	logger.Info("Marking reports as processed for builds",
		"count", len(identifiers),
		"batchID", batchID)

	if a.reportsRepository == nil {
		logger.Error("ReportRepository is not initialized in PlayerBuildsActivity")
		return fmt.Errorf("internal error: reportRepository not injected")
	}

	if err := a.reportsRepository.MarkReportsAsProcessedForBuilds(ctx, identifiers, batchID); err != nil {
		logger.Error("Failed to mark reports as processed for builds", "error", err)
		return err
	}

	logger.Info("Successfully marked reports as processed for builds", "count", len(identifiers))
	return nil
}
