package warcraftlogsBuildsTemporal

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

	// Traitement des reports par lot pour éviter la surcharge
	const batchSize = 10
	totalReports := len(rankings)
	processedReports := 0

	for i := 0; i < totalReports; i += batchSize {
		end := i + batchSize
		if end > totalReports {
			end = totalReports
		}

		batch := rankings[i:end]

		// Heartbeat pour le monitoring Temporal
		activity.RecordHeartbeat(ctx, fmt.Sprintf("Processing reports %d-%d of %d", i+1, end, totalReports))

		for _, ranking := range batch {
			if err := a.processReport(ctx, ranking); err != nil {
				logger.Error("Failed to process report",
					"reportCode", ranking.ReportCode,
					"error", err)
				result.FailureCount++
				continue
			}
			result.SuccessCount++
			processedReports++
		}

		logger.Info("Batch processed",
			"processed", processedReports,
			"total", totalReports,
			"success", result.SuccessCount,
			"failures", result.FailureCount)
	}

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
	if err := a.repository.StoreReport(ctx, report); err != nil {
		return fmt.Errorf("failed to store report: %w", err)
	}

	logger.Info("Report processed successfully",
		"reportCode", ranking.ReportCode,
		"fightID", ranking.ReportFightID)

	return nil

}
