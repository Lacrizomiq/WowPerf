package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.temporal.io/sdk/activity"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	"wowperf/internal/services/warcraftlogs"
	reportsQueries "wowperf/internal/services/warcraftlogs/mythicplus/builds/queries"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"

	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

// reportWorkItem represents a single unit of work to be processed by workers
type reportWorkItem struct {
	ranking *warcraftlogsBuilds.ClassRanking
	index   int // Preserve ordering for final results
}

// reportWorkResult represents the result of processing a single report
type reportWorkResult struct {
	report *warcraftlogsBuilds.Report
	index  int // Original position in the rankings slice
	err    error
}

// ReportsActivity handles all report-related operations
type ReportsActivity struct {
	client             *warcraftlogs.WarcraftLogsClientService
	repository         *reportsRepository.ReportRepository
	rankingsRepository *rankingsRepository.RankingsRepository
}

// NewReportsActivity creates a new instance of ReportsActivity
func NewReportsActivity(
	client *warcraftlogs.WarcraftLogsClientService,
	repository *reportsRepository.ReportRepository,
	rankingsRepository *rankingsRepository.RankingsRepository,
) *ReportsActivity {
	return &ReportsActivity{
		client:             client,
		repository:         repository,
		rankingsRepository: rankingsRepository,
	}
}

// ProcessReports processes rankings and updates corresponding reports
// It handles API fetching, storage, and synchronization
func (a *ReportsActivity) ProcessReports(
	ctx context.Context,
	rankings []*warcraftlogsBuilds.ClassRanking,
) (*models.ReportProcessingResult, error) {
	logger := activity.GetLogger(ctx)
	result := &models.ReportProcessingResult{
		ProcessedAt: time.Now(),
	}

	if len(rankings) == 0 {
		logger.Info("No rankings to process")
		return result, nil
	}

	logger.Info("Starting reports processing", "rankingsCount", len(rankings))

	// Fetch reports from API
	reports, err := a.fetchReportsFromAPI(ctx, rankings)
	if err != nil {
		logger.Error("Failed to fetch reports from API", "error", err)
		return nil, fmt.Errorf("failed to fetch reports: %w", err)
	}

	// Store fetched reports
	if len(reports) > 0 {
		if err := a.repository.StoreReports(ctx, reports); err != nil {
			logger.Error("Failed to store reports",
				"reportCount", len(reports),
				"error", err)
			return nil, fmt.Errorf("failed to store reports: %w", err)
		}
		logger.Info("Successfully stored reports", "count", len(reports))
	}

	// Synchronize with rankings
	if err := a.repository.SyncReportsWithRankings(ctx, rankings); err != nil {
		logger.Error("Failed to sync reports with rankings", "error", err)
		return nil, fmt.Errorf("failed to sync reports: %w", err)
	}

	result.ProcessedReports = reports
	result.ProcessedCount = int32(len(reports))
	result.SuccessCount = 1

	logger.Info("Completed report processing",
		"totalProcessed", len(reports),
		"duration", time.Since(result.ProcessedAt))

	// Mark rankings as processed
	var rankingIDs []uint
	for _, ranking := range rankings {
		rankingIDs = append(rankingIDs, ranking.ID)
	}

	if len(rankingIDs) > 0 {
		if err := a.rankingsRepository.MarkRankingsAsProcessedForReports(ctx, rankingIDs, "processed"); err != nil {
			logger.Error("Failed to mark rankings as processed", "error", err)
			// Continue even if marking fails
		} else {
			logger.Info("Successfully marked rankings as processed", "count", len(rankingIDs))
		}
	}

	return result, nil
}

// fetchReportsFromAPI fetches reports from the WarcraftLogs API in parallel
// It processes multiple rankings simultaneously while maintaining order and handling rate limits
func (a *ReportsActivity) fetchReportsFromAPI(
	ctx context.Context,
	rankings []*warcraftlogsBuilds.ClassRanking,
) ([]*warcraftlogsBuilds.Report, error) {
	logger := activity.GetLogger(ctx)

	// Configuration constants for parallel processing
	const (
		maxWorkers = 2  // Maximum number of concurrent workers
		batchSize  = 10 // Size of processing batches
	)

	// Pre-allocate slice to maintain order of reports
	// This allows us to preserve the relationship between rankings and reports
	reports := make([]*warcraftlogsBuilds.Report, len(rankings))

	// Channel setup for work distribution and result collection
	// Buffered channels are used to optimize throughput
	workChan := make(chan reportWorkItem, batchSize)     // Channel for distributing work
	resultChan := make(chan reportWorkResult, batchSize) // Channel for collecting results
	doneChan := make(chan struct{})                      // Channel to signal completion

	// Start worker pool
	// Each worker processes items from workChan independently
	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for work := range workChan {
				// Check for context cancellation
				select {
				case <-ctx.Done():
					return
				default:
					// Record worker activity for monitoring
					activity.RecordHeartbeat(ctx, map[string]interface{}{
						"workerID":   workerID,
						"reportCode": work.ranking.ReportCode,
						"index":      work.index,
						"status":     "processing",
					})

					// Process the report through the API
					report, err := a.getReportDetails(ctx, work.ranking)

					// Little delay to prevent 429 error
					time.Sleep(time.Millisecond * 500)

					// Send result back through result channel
					resultChan <- reportWorkResult{
						report: report,
						index:  work.index,
						err:    err,
					}
				}
			}
		}(i)
	}

	// Goroutine to close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
		close(doneChan)
	}()

	// Goroutine to distribute work to workers
	// This runs independently of result collection to maintain steady worker utilization
	go func() {
		defer close(workChan)

		for i, ranking := range rankings {
			select {
			case <-ctx.Done():
				return
			case workChan <- reportWorkItem{ranking: ranking, index: i}:
				// Work successfully queued
			}
		}
	}()

	// Collect and process results as they come in
	processedCount := 0
	failureCount := 0

	// Process results as they arrive from workers
	for result := range resultChan {
		if result.err != nil {
			logger.Error("Failed to fetch report",
				"index", result.index,
				"error", result.err)
			failureCount++
			continue
		}

		if result.report != nil {
			reports[result.index] = result.report
			processedCount++
		}

		// Record progress for monitoring
		activity.RecordHeartbeat(ctx, map[string]interface{}{
			"processedCount": processedCount,
			"failureCount":   failureCount,
			"totalCount":     len(rankings),
			"progress":       fmt.Sprintf("%d/%d", processedCount+failureCount, len(rankings)),
		})
	}

	// Wait for all processing to complete
	<-doneChan

	// Clean up the results by removing any nil entries from failed processes
	cleanReports := make([]*warcraftlogsBuilds.Report, 0, processedCount)
	for _, report := range reports {
		if report != nil {
			cleanReports = append(cleanReports, report)
		}
	}

	// Log final processing statistics
	logger.Info("Completed fetching reports",
		"processedCount", processedCount,
		"failureCount", failureCount,
		"totalRequested", len(rankings))

	return cleanReports, nil
}

// getReportDetails fetches and processes report details from WarcraftLogs API
func (a *ReportsActivity) getReportDetails(
	ctx context.Context,
	ranking *warcraftlogsBuilds.ClassRanking,
) (*warcraftlogsBuilds.Report, error) {
	response, err := a.client.MakeRequest(ctx, reportsQueries.GetReportTableQuery, map[string]interface{}{
		"code":        ranking.ReportCode,
		"fightID":     ranking.ReportFightID,
		"encounterID": ranking.EncounterID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch report %s: %w", ranking.ReportCode, err)
	}

	report, talentsQuery, err := reportsQueries.ParseReportDetailsResponse(
		response,
		ranking.ReportCode,
		ranking.ReportFightID,
		ranking.EncounterID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse report details: %w", err)
	}

	talentsResponse, err := a.client.MakeRequest(ctx, talentsQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch talents for report %s: %w", ranking.ReportCode, err)
	}

	talentCodes, err := reportsQueries.ParseReportTalentsResponse(talentsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse talents for report %s: %w", ranking.ReportCode, err)
	}

	report.TalentCodes, err = json.Marshal(talentCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal talents for report %s: %w", ranking.ReportCode, err)
	}

	return report, nil
}

// GetReportsBatch retrieves a batch of reports for builds processing
func (a *ReportsActivity) GetReportsBatch(
	ctx context.Context,
	batchSize int,
	offset int,
) ([]*warcraftlogsBuilds.Report, error) {
	return a.repository.GetReportsBatch(ctx, batchSize, offset)
}

// CountAllReports returns the total number of reports
func (a *ReportsActivity) CountAllReports(ctx context.Context) (int64, error) {
	return a.repository.CountAllReports(ctx)
}

// GetUniqueReportReferences retrieves all unique report references from rankings
func (a *ReportsActivity) GetUniqueReportReferences(ctx context.Context) ([]*warcraftlogsBuilds.ClassRanking, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Getting unique report references")

	return a.repository.GetAllUniqueReportReferences(ctx)
}

// GetRankingsNeedingReportProcessing retrieves the rankings that need report processing
// It also filters by class if a class is provided
func (a *ReportsActivity) GetRankingsNeedingReportProcessing(ctx context.Context, className string, limit int32, maxAgeDuration time.Duration) ([]*warcraftlogsBuilds.ClassRanking, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Getting rankings needing report processing",
		"className", className,
		"limit", limit,
		"maxAge", maxAgeDuration)

	if a.rankingsRepository == nil {
		logger.Error("RankingsRepository is not initialized in ReportsActivity")
		return nil, fmt.Errorf("internal error: rankingsRepository not injected")
	}

	rankings, err := a.rankingsRepository.GetRankingsNeedingReportProcessing(ctx, className, int(limit), maxAgeDuration)

	if err != nil {
		logger.Error("Failed to get rankings needing report processing from repository", "error", err)
		return nil, err // return the error to the caller
	}

	logger.Info("Retrieved rankings needing report processing", "count", len(rankings))
	return rankings, nil
}

// MarkReportsForBuildProcessing marks reports for build processing
func (a *ReportsActivity) MarkReportsForBuildProcessing(ctx context.Context, reports []*warcraftlogsBuilds.Report, batchID string) error {
	logger := activity.GetLogger(ctx)

	identifiers := make([]reportsRepository.ReportIdentifier, 0, len(reports))
	for _, report := range reports {
		identifiers = append(identifiers, reportsRepository.ReportIdentifier{
			Code:    report.Code,
			FightID: report.FightID,
		})
	}

	logger.Info("Marking reports for build processing",
		"count", len(identifiers),
		"batchID", batchID,
	)

	if a.repository == nil {
		logger.Error("ReportRepository is not initialized in ReportsActivity")
		return fmt.Errorf("internal error: reportRepository not injected")
	}

	if err := a.repository.MarkReportsForBuildProcessing(ctx, identifiers, batchID); err != nil {
		logger.Error("Failed to mark reports for build processing", "error", err)
		return err
	}

	logger.Info("Successfully marked reports for build processing", "count", len(identifiers))
	return nil
}
