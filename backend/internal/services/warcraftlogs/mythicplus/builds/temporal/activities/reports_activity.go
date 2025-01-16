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

type ReportsActivity struct {
	client     *warcraftlogs.WarcraftLogsClientService
	repository *reportsRepository.ReportRepository
}

func NewReportsActivity(client *warcraftlogs.WarcraftLogsClientService, repository *reportsRepository.ReportRepository) *ReportsActivity {
	return &ReportsActivity{
		client:     client,
		repository: repository,
	}
}

func (a *ReportsActivity) ProcessReports(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking) (*workflows.ReportProcessingResult, error) {
	logger := activity.GetLogger(ctx)
	result := &workflows.ReportProcessingResult{
		ProcessedAt: time.Now(),
	}

	if len(rankings) == 0 {
		logger.Info("No rankings to process")
		return result, nil
	}

	logger.Info("Starting report processing", "rankingsCount", len(rankings))

	const batchSize = 10
	processedReports := 0

	// Group rankings by encounterID for batch processing
	rankingsByEncounter := make(map[uint][]*warcraftlogsBuilds.ClassRanking)
	for _, ranking := range rankings {
		rankingsByEncounter[ranking.EncounterID] = append(rankingsByEncounter[ranking.EncounterID], ranking)
	}

	for encounterID, encounterRankings := range rankingsByEncounter {
		logger.Info("Processing encounter",
			"encounterID", encounterID,
			"rankingsCount", len(encounterRankings))

		// Process each encounter in batches
		for i := 0; i < len(encounterRankings); i += batchSize {
			end := i + batchSize
			if end > len(encounterRankings) {
				end = len(encounterRankings)
			}
			batch := encounterRankings[i:end]

			activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing reports %d-%d for encounter %d", i+1, end, encounterID))

			// Fetch all reports from the batch at once
			var reports []*warcraftlogsBuilds.Report
			for _, ranking := range batch {
				// Check if the report already exists
				existingReport, err := a.repository.GetReportByCodeAndFightID(ctx, ranking.ReportCode, ranking.ReportFightID)
				if err != nil {
					logger.Error("Failed to check existing report",
						"error", err,
						"reportCode", ranking.ReportCode,
						"fightID", ranking.ReportFightID)
					result.FailureCount++
					continue
				}

				if existingReport != nil {
					logger.Debug("Report already exists, will be updated",
						"reportCode", ranking.ReportCode,
						"fightID", ranking.ReportFightID)
				}

				// Fetch the report details
				response, err := a.client.MakeRequest(ctx, reportsQueries.GetReportTableQuery, map[string]interface{}{
					"code":        ranking.ReportCode,
					"fightID":     ranking.ReportFightID,
					"encounterID": ranking.EncounterID,
				})
				if err != nil {
					logger.Error("Failed to get report details",
						"error", err,
						"reportCode", ranking.ReportCode,
						"fightID", ranking.ReportFightID)
					result.FailureCount++
					continue
				}

				report, talentsQuery, err := reportsQueries.ParseReportDetailsResponse(
					response,
					ranking.ReportCode,
					ranking.ReportFightID,
					ranking.EncounterID,
				)
				if err != nil {
					logger.Error("Failed to parse report details",
						"error", err,
						"reportCode", ranking.ReportCode)
					result.FailureCount++
					continue
				}

				// Fetch the talents
				talentsResponse, err := a.client.MakeRequest(ctx, talentsQuery, nil)
				if err != nil {
					logger.Error("Failed to fetch talents",
						"error", err,
						"reportCode", ranking.ReportCode)
					result.FailureCount++
					continue
				}

				talentsCode, err := reportsQueries.ParseReportTalentsResponse(talentsResponse)
				if err != nil {
					logger.Error("Failed to parse talents",
						"error", err,
						"reportCode", ranking.ReportCode)
					result.FailureCount++
					continue
				}

				talentCodeJSON, err := json.Marshal(talentsCode)
				if err != nil {
					logger.Error("Failed to marshal talents",
						"error", err,
						"reportCode", ranking.ReportCode)
					result.FailureCount++
					continue
				}
				report.TalentCodes = talentCodeJSON

				reports = append(reports, report)
				result.SuccessCount++
			}

			// Store the batch of reports
			if len(reports) > 0 {
				logger.Info("Storing batch of reports",
					"batchSize", len(reports),
					"encounterID", encounterID)

				if err := a.repository.StoreReports(ctx, reports); err != nil {
					logger.Error("Failed to store reports batch",
						"error", err,
						"encounterID", encounterID,
						"batchSize", len(reports))
					return nil, fmt.Errorf("failed to store reports batch: %w", err)
				}

				processedReports += len(reports)
				logger.Info("Successfully stored batch of reports",
					"batchSize", len(reports),
					"totalProcessed", processedReports,
					"encounterID", encounterID)
			}
		}
	}

	logger.Info("Completed report processing",
		"totalProcessed", processedReports,
		"successCount", result.SuccessCount,
		"failureCount", result.FailureCount)

	result.ProcessedReports = processedReports
	return result, nil
}

func (a *ReportsActivity) processReport(ctx context.Context, ranking *warcraftlogsBuilds.ClassRanking) error {
	logger := activity.GetLogger(ctx)

	// Check if report already exists
	existingReport, err := a.repository.GetReportByCodeAndFightID(ctx, ranking.ReportCode, ranking.ReportFightID)
	if err != nil {
		return fmt.Errorf("failed to check if report exists: %w", err)
	}

	if existingReport != nil {
		logger.Info("Report already exists, skipping",
			"reportCode", ranking.ReportCode,
			"fightID", ranking.ReportFightID)
		return nil
	}

	// Get the reports details
	response, err := a.client.MakeRequest(ctx, reportsQueries.GetReportTableQuery, map[string]interface{}{
		"code":        ranking.ReportCode,
		"fightID":     ranking.ReportFightID,
		"encounterID": ranking.EncounterID,
	})
	if err != nil {
		return fmt.Errorf("failed to get report details: %w", err)
	}

	// Parse the response
	report, talentsQuery, err := reportsQueries.ParseReportDetailsResponse(
		response,
		ranking.ReportCode,
		ranking.ReportFightID,
		ranking.EncounterID,
	)
	if err != nil {
		return fmt.Errorf("failed to parse report details: %w", err)
	}

	// Retrieve the talents
	talentsResponse, err := a.client.MakeRequest(ctx, talentsQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch talents: %w", err)
	}

	talentsCode, err := reportsQueries.ParseReportTalentsResponse(talentsResponse)
	if err != nil {
		return fmt.Errorf("failed to parse talents: %w", err)
	}

	// Convert and store talents
	talentCodeJSON, err := json.Marshal(talentsCode)
	if err != nil {
		return fmt.Errorf("failed to marshal talents: %w", err)
	}
	report.TalentCodes = talentCodeJSON

	// store the report
	if err := a.repository.StoreReports(ctx, []*warcraftlogsBuilds.Report{report}); err != nil {
		return fmt.Errorf("failed to store report: %w", err)
	}

	logger.Info("Report processed successfully",
		"reportCode", ranking.ReportCode,
		"fightID", ranking.ReportFightID)

	return nil

}

func (a *ReportsActivity) GetProcessedReports(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking) ([]*warcraftlogsBuilds.Report, error) {
	logger := activity.GetLogger(ctx)

	if len(rankings) == 0 {
		logger.Info("No rankings provided to fetch reports")
		return nil, nil
	}

	// Get all reports for the given rankings
	reportCodes := make(map[string]bool)
	for _, ranking := range rankings {
		reportCodes[ranking.ReportCode] = true
	}

	var allReports []*warcraftlogsBuilds.Report

	// Get all reports by batch
	const batchSize = 50
	processedReports := 0

	// convert the map to a slice to process in batches
	uniqueCodes := make([]string, 0, len(reportCodes))
	for code := range reportCodes {
		uniqueCodes = append(uniqueCodes, code)
	}

	for i := 0; i < len(uniqueCodes); i += batchSize {
		end := i + batchSize
		if end > len(uniqueCodes) {
			end = len(uniqueCodes)
		}

		// record heartbeat for monitoring
		activity.RecordHeartbeat(ctx, fmt.Sprintf("Fetching reports %d-%d of %d", i+1, end, len(uniqueCodes)))

		// Get the batch of reports
		batchCodes := uniqueCodes[i:end]
		reports, err := a.repository.GetReportsByCode(ctx, batchCodes)
		if err != nil {
			logger.Error("Failed to fetch reports",
				"startIndex", i,
				"endIndex", end,
				"error", err)
			continue
		}

		allReports = append(allReports, reports...)
		processedReports += len(reports)
	}

	logger.Info("Finished fetching reports",
		"totalReports", len(allReports),
		"uniqueCodes", len(reportCodes))

	return allReports, nil
}

func (a *ReportsActivity) GetReportsForEncounter(ctx context.Context, encounterID uint) ([]warcraftlogsBuilds.Report, error) {
	logger := activity.GetLogger(ctx)

	const batchSize = 10
	offset := 0
	var allReports []warcraftlogsBuilds.Report

	for {
		reports, err := a.repository.GetReportsForEncounter(ctx, encounterID, batchSize, offset)
		if err != nil {
			logger.Error("Failed to get reports batch",
				"encounterID", encounterID,
				"offset", offset,
				"error", err)
			return nil, err
		}

		if len(reports) == 0 {
			break
		}

		allReports = append(allReports, reports...)
		offset += batchSize

		// Record heartbeat
		activity.RecordHeartbeat(ctx, fmt.Sprintf("Retrieved %d reports for encounter %d", len(allReports), encounterID))
	}

	logger.Info("Retrieved reports for encounter",
		"encounterID", encounterID,
		"count", len(allReports))

	return allReports, nil
}
