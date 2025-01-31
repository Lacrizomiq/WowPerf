package warcraftlogsBuildsTemporalActivities

import (
	"context"
	"encoding/json"
	"fmt"
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

// GetReportsForEncounterBatch retrieves a batch of reports for an encounter
// Used by the rebuild process to fetch reports in manageable chunks
func (a *ReportsActivity) GetReportsForEncounterBatch(
	ctx context.Context,
	encounterID uint,
	limit int,
	offset int,
) ([]warcraftlogsBuilds.Report, error) {
	logger := activity.GetLogger(ctx)

	// Record heartbeat to prevent timeout
	activity.RecordHeartbeat(ctx, fmt.Sprintf("Fetching reports %d-%d for encounter %d",
		offset, offset+limit, encounterID))

	reports, err := a.repository.GetReportsForEncounter(ctx, encounterID, limit, offset)
	if err != nil {
		logger.Error("Failed to fetch reports batch",
			"encounterID", encounterID,
			"offset", offset,
			"limit", limit,
			"error", err)
		return nil, fmt.Errorf("failed to fetch reports batch: %w", err)
	}

	logger.Info("Successfully fetched reports batch",
		"encounterID", encounterID,
		"batchSize", len(reports),
		"offset", offset)

	return reports, nil
}

// CountReportsForEncounter counts the total number of reports for an encounter
// Used to determine the total number of batches needed for processing
func (a *ReportsActivity) CountReportsForEncounter(
	ctx context.Context,
	encounterID uint,
) (int64, error) {
	logger := activity.GetLogger(ctx)

	count, err := a.repository.CountReportsForEncounter(ctx, encounterID)
	if err != nil {
		logger.Error("Failed to count reports for encounter",
			"encounterID", encounterID,
			"error", err)
		return 0, fmt.Errorf("failed to count reports: %w", err)
	}

	logger.Info("Counted reports for encounter",
		"encounterID", encounterID,
		"count", count)

	return count, nil
}

// ProcessReports processes a list of rankings and creates or updates the corresponding reports
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

	var allProcessedReports []*warcraftlogsBuilds.Report

	// Group rankings by encounterID for efficient batch processing
	rankingsByEncounter := make(map[uint][]*warcraftlogsBuilds.ClassRanking)
	for _, ranking := range rankings {
		rankingsByEncounter[ranking.EncounterID] = append(
			rankingsByEncounter[ranking.EncounterID],
			ranking,
		)
	}

	// Process rankings in batches by encounter
	const batchSize = 10
	for encounterID, encounterRankings := range rankingsByEncounter {
		logger.Info("Processing encounter",
			"encounterID", encounterID,
			"rankingsCount", len(encounterRankings))

		for i := 0; i < len(encounterRankings); i += batchSize {
			end := min(i+batchSize, len(encounterRankings))
			batch := encounterRankings[i:end]

			activity.RecordHeartbeat(ctx, fmt.Sprintf(
				"Processing reports %d-%d for encounter %d",
				i+1, end, encounterID))

			reports, err := a.processRankingsBatch(ctx, batch)
			if err != nil {
				logger.Error("Failed to process rankings batch",
					"encounterID", encounterID,
					"startIndex", i,
					"endIndex", end,
					"error", err)
				result.FailureCount++
				continue
			}

			if err := a.repository.StoreReports(ctx, reports); err != nil {
				logger.Error("Failed to store reports batch",
					"encounterID", encounterID,
					"batchSize", len(reports),
					"error", err)
				result.FailureCount++
				continue
			}

			allProcessedReports = append(allProcessedReports, reports...)
			result.SuccessCount++
			result.ProcessedCount += len(reports)

			logger.Info("Successfully processed reports batch",
				"encounterID", encounterID,
				"batchProcessed", len(reports),
				"totalProcessed", result.ProcessedCount)

			time.Sleep(time.Millisecond * 100)
		}
	}

	// Assign the processed reports to the result
	result.ProcessedReports = allProcessedReports

	return result, nil
}

// processRankingsBatch processes a batch of rankings and fetches their report details
func (a *ReportsActivity) processRankingsBatch(
	ctx context.Context,
	rankings []*warcraftlogsBuilds.ClassRanking,
) ([]*warcraftlogsBuilds.Report, error) {
	var reports []*warcraftlogsBuilds.Report

	for _, ranking := range rankings {
		report, err := a.getReportDetails(ctx, ranking)
		if err != nil {
			return nil, fmt.Errorf("failed to process report %s: %w",
				ranking.ReportCode, err)
		}
		reports = append(reports, report)
	}

	return reports, nil
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
		return nil, fmt.Errorf("failed to fetch report details: %w", err)
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
		return nil, fmt.Errorf("failed to fetch talents: %w", err)
	}

	talentCodes, err := reportsQueries.ParseReportTalentsResponse(talentsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse talents: %w", err)
	}

	// Convert map to JSON
	talentCodesJSON, err := json.Marshal(talentCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal talent codes: %w", err)
	}
	report.TalentCodes = talentCodesJSON

	return report, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetProcessedReports retrieves reports based on a list of rankings
func (a *ReportsActivity) GetProcessedReports(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking) ([]*warcraftlogsBuilds.Report, error) {
	logger := activity.GetLogger(ctx)

	if len(rankings) == 0 {
		logger.Info("No rankings provided to fetch reports")
		return nil, nil
	}

	// Extract unique report codes
	reportCodes := make(map[string]bool)
	for _, ranking := range rankings {
		reportCodes[ranking.ReportCode] = true
	}

	var allReports []*warcraftlogsBuilds.Report
	processedReports := 0

	// Process reports in batches
	const batchSize = 50
	uniqueCodes := make([]string, 0, len(reportCodes))
	for code := range reportCodes {
		uniqueCodes = append(uniqueCodes, code)
	}

	for i := 0; i < len(uniqueCodes); i += batchSize {
		end := i + batchSize
		if end > len(uniqueCodes) {
			end = len(uniqueCodes)
		}

		activity.RecordHeartbeat(ctx, fmt.Sprintf("Fetching reports %d-%d of %d", i+1, end, len(uniqueCodes)))

		batchCodes := uniqueCodes[i:end]
		reports, err := a.repository.GetReportsByCode(ctx, batchCodes)
		if err != nil {
			logger.Error("Failed to fetch reports batch",
				"startIndex", i,
				"endIndex", end,
				"error", err)
			continue
		}

		allReports = append(allReports, reports...)
		processedReports += len(reports)

		// Add small delay between batches
		time.Sleep(time.Millisecond * 100)
	}

	logger.Info("Finished fetching reports",
		"totalReports", len(allReports),
		"uniqueCodes", len(reportCodes))

	return allReports, nil
}

// GetReportsForEncounter retrieves all reports for an encounter
func (a *ReportsActivity) GetReportsForEncounter(ctx context.Context, encounterID uint) ([]warcraftlogsBuilds.Report, error) {
	logger := activity.GetLogger(ctx)

	// Get total count first
	count, err := a.CountReportsForEncounter(ctx, encounterID)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		logger.Info("No reports found for encounter", "encounterID", encounterID)
		return nil, nil
	}

	const batchSize = 5
	var allReports []warcraftlogsBuilds.Report

	for offset := 0; offset < int(count); offset += batchSize {
		activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing reports %d-%d of %d for encounter %d",
			offset, min(offset+batchSize, int(count)), count, encounterID))

		reports, err := a.GetReportsForEncounterBatch(ctx, encounterID, batchSize, offset)
		if err != nil {
			logger.Error("Failed to get reports batch",
				"encounterID", encounterID,
				"offset", offset,
				"error", err)
			continue
		}

		allReports = append(allReports, reports...)
		logger.Info("Processed reports batch",
			"encounterID", encounterID,
			"batchSize", len(reports),
			"totalProcessed", len(allReports),
			"totalExpected", count)

		// Add small delay between batches
		time.Sleep(time.Millisecond * 100)
	}

	return allReports, nil
}
