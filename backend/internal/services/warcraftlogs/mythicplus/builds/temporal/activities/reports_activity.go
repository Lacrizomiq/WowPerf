package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.temporal.io/sdk/activity"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	"wowperf/internal/services/warcraftlogs"
	reportsQueries "wowperf/internal/services/warcraftlogs/mythicplus/builds/queries"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"
)

// ReportsActivity handles all report-related operations
type ReportsActivity struct {
	client     *warcraftlogs.WarcraftLogsClientService
	repository *reportsRepository.ReportRepository
}

// NewReportsActivity creates a new instance of ReportsActivity
func NewReportsActivity(
	client *warcraftlogs.WarcraftLogsClientService,
	repository *reportsRepository.ReportRepository,
) *ReportsActivity {
	return &ReportsActivity{
		client:     client,
		repository: repository,
	}
}

// ProcessReports processes a list of rankings and creates or updates the corresponding reports
// It handles batching and rate limiting internally
func (a *ReportsActivity) ProcessReports(
	ctx context.Context,
	rankings []*warcraftlogsBuilds.ClassRanking,
) (*workflows.ReportProcessingResult, error) {
	logger := activity.GetLogger(ctx)
	result := &workflows.ReportProcessingResult{
		ProcessedAt: time.Now(),
	}

	if len(rankings) == 0 {
		logger.Info("No rankings to process")
		return result, nil
	}

	// Create a map to track unique reports and avoid duplicates
	processedReports := make(map[string]*warcraftlogsBuilds.Report)
	uniqueRankings := make([]*warcraftlogsBuilds.ClassRanking, 0)

	// First, deduplicate rankings based on report code and fight ID
	for _, ranking := range rankings {
		key := fmt.Sprintf("%s-%d", ranking.ReportCode, ranking.ReportFightID)
		if _, exists := processedReports[key]; !exists {
			uniqueRankings = append(uniqueRankings, ranking)
		}
	}

	logger.Info("Filtered unique rankings",
		"totalRankings", len(rankings),
		"uniqueRankings", len(uniqueRankings))

	// Process unique rankings in batches
	const batchSize = 10
	for i := 0; i < len(uniqueRankings); i += batchSize {
		end := i + batchSize
		if end > len(uniqueRankings) {
			end = len(uniqueRankings)
		}

		batch := uniqueRankings[i:end]
		activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing rankings batch %d-%d of %d",
			i+1, end, len(uniqueRankings)))

		// Process batch of rankings
		reports, err := a.processRankingsBatch(ctx, batch)
		if err != nil {
			logger.Error("Failed to process rankings batch",
				"startIndex", i,
				"endIndex", end,
				"error", err)
			result.FailureCount++
			continue
		}

		// Track processed reports
		for _, report := range reports {
			key := fmt.Sprintf("%s-%d", report.Code, report.FightID)
			if _, exists := processedReports[key]; !exists {
				processedReports[key] = report
				result.ProcessedReports = append(result.ProcessedReports, report)
			}
		}

		// Store only new reports
		if len(reports) > 0 {
			if err := a.repository.StoreReports(ctx, reports); err != nil {
				logger.Error("Failed to store reports batch",
					"batchSize", len(reports),
					"error", err)
				result.FailureCount++
				continue
			}
		}

		result.SuccessCount++
		result.ProcessedCount += len(reports)

		logger.Info("Successfully processed reports batch",
			"batchProcessed", len(reports),
			"totalProcessed", result.ProcessedCount,
			"progress", fmt.Sprintf("%d/%d", end, len(uniqueRankings)))

		// Add delay between batches
		time.Sleep(time.Millisecond * 100)
	}

	logger.Info("Completed reports processing",
		"totalReports", len(processedReports),
		"successCount", result.SuccessCount,
		"failureCount", result.FailureCount,
		"duration", time.Since(result.ProcessedAt))

	return result, nil
}

// GetProcessedReports retrieves reports based on a list of rankings from the database
func (a *ReportsActivity) GetProcessedReports(
	ctx context.Context,
	rankings []*warcraftlogsBuilds.ClassRanking,
) ([]*warcraftlogsBuilds.Report, error) {
	logger := activity.GetLogger(ctx)

	// Get all reports in a single query
	reports, err := a.repository.GetReportsByRankings(ctx, rankings)
	if err != nil {
		logger.Error("Failed to get reports", "error", err)
		return nil, err
	}

	// Create a map for quick lookup
	existingReports := make(map[string]*warcraftlogsBuilds.Report)
	for _, report := range reports {
		key := fmt.Sprintf("%s-%d", report.Code, report.FightID)
		existingReports[key] = report
	}

	return reports, nil
}

func (a *ReportsActivity) processRankingsBatch(
	ctx context.Context,
	rankings []*warcraftlogsBuilds.ClassRanking,
) ([]*warcraftlogsBuilds.Report, error) {
	logger := activity.GetLogger(ctx)

	// Get existing reports first
	existingReports, err := a.repository.GetReportsByRankings(ctx, rankings)
	if err != nil {
		logger.Error("Failed to fetch existing reports", "error", err)
	}

	// Create map of existing reports
	existingMap := make(map[string]*warcraftlogsBuilds.Report)
	for _, report := range existingReports {
		key := fmt.Sprintf("%s-%d", report.Code, report.FightID)
		existingMap[key] = report
	}

	var reports []*warcraftlogsBuilds.Report
	var processingErrors []string

	// Process each ranking
	for _, ranking := range rankings {
		key := fmt.Sprintf("%s-%d", ranking.ReportCode, ranking.ReportFightID)

		// Use existing report if available
		if report, exists := existingMap[key]; exists {
			reports = append(reports, report)
			continue
		}

		// Process new report
		report, err := a.getReportDetails(ctx, ranking)
		if err != nil {
			logger.Error("Failed to process report",
				"reportCode", ranking.ReportCode,
				"error", err)
			processingErrors = append(processingErrors,
				fmt.Sprintf("report %s: %v", ranking.ReportCode, err))
			continue
		}

		reports = append(reports, report)
	}

	// Log processing results
	if len(processingErrors) > 0 {
		logger.Error("Some reports failed to process",
			"successCount", len(reports),
			"failureCount", len(processingErrors),
			"errors", strings.Join(processingErrors, "; "))
	}

	return reports, nil
}

// getReportDetails fetches and processes report details from WarcraftLogs API
func (a *ReportsActivity) getReportDetails(
	ctx context.Context,
	ranking *warcraftlogsBuilds.ClassRanking,
) (*warcraftlogsBuilds.Report, error) {

	// Fetch report details from API
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

	// Fetch and process talents with separate query
	talentsResponse, err := a.client.MakeRequest(ctx, talentsQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch talents for report %s: %w", err)
	}

	talentCodes, err := reportsQueries.ParseReportTalentsResponse(talentsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse talents for report %s: %w", err)
	}

	// Convert talent codes to JSON
	report.TalentCodes, err = json.Marshal(talentCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal talents for report %s: %w", ranking.ReportCode, err)
	}

	return report, nil
}

// GetReportsBatch retrieves a batch of reports from the database without filtering
func (a *ReportsActivity) GetReportsBatch(
	ctx context.Context,
	batchSize int,
	offset int,
) ([]*warcraftlogsBuilds.Report, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Fetching reports batch for builds processing",
		"batchSize", batchSize,
		"offset", offset)

	// Retrieve reports via the repository
	reports, err := a.repository.GetReportsBatch(ctx, batchSize, offset)
	if err != nil {
		logger.Error("Failed to fetch reports batch",
			"error", err,
			"batchSize", batchSize,
			"offset", offset)
		return nil, fmt.Errorf("failed to fetch reports batch: %w", err)
	}

	logger.Info("Successfully fetched reports batch",
		"reportsCount", len(reports),
		"batchSize", batchSize,
		"offset", offset)

	return reports, nil
}

// CountAllReports returns the total number of reports
func (a *ReportsActivity) CountAllReports(ctx context.Context) (int64, error) {
	return a.repository.CountAllReports(ctx)
}
