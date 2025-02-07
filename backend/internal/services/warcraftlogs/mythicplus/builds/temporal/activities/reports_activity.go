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

// ProcessReports processes rankings and updates corresponding reports
// It handles API fetching, storage, and synchronization
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
	result.ProcessedCount = len(reports)
	result.SuccessCount = 1

	logger.Info("Completed report processing",
		"totalProcessed", len(reports),
		"duration", time.Since(result.ProcessedAt))

	return result, nil
}

// fetchReportsFromAPI fetches reports data from WarcraftLogs API
func (a *ReportsActivity) fetchReportsFromAPI(
	ctx context.Context,
	rankings []*warcraftlogsBuilds.ClassRanking,
) ([]*warcraftlogsBuilds.Report, error) {
	var reports []*warcraftlogsBuilds.Report
	logger := activity.GetLogger(ctx)

	for i, ranking := range rankings {
		activity.RecordHeartbeat(ctx, fmt.Sprintf("Fetching report %d/%d", i+1, len(rankings)))

		report, err := a.getReportDetails(ctx, ranking)
		if err != nil {
			logger.Error("Failed to fetch report details",
				"reportCode", ranking.ReportCode,
				"error", err)
			continue
		}

		if report != nil {
			reports = append(reports, report)
		}
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
		return nil, fmt.Errorf("failed to fetch talents for report %s: %w", err)
	}

	talentCodes, err := reportsQueries.ParseReportTalentsResponse(talentsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse talents for report %s: %w", err)
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
