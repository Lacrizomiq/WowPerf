// internal/services/warcraftlogs/mythicplus/builds/temporal/init.go
package temporal

import (
	"wowperf/internal/services/warcraftlogs"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"gorm.io/gorm"

	// Repositories
	buildsStatisticsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	playerBuildsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	reportsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	statStatisticsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	talentStatisticsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
	workflowStatesRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"

	// Activities
	activities "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/activities"

	// Definitions
	definitions "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/definitions"

	// Workflows
	buildsWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/builds"
	equipmentAnalysisWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/builds_statistics/equipment_statistics"
	statAnalysisWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/builds_statistics/stats_statistics"
	talentAnalysisWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/builds_statistics/talent_statistics"
	rankingsWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/rankings"
	reportsWorkflow "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/reports"
)

// InitBuilds initialise les repositories et activités pour la feature builds
func InitBuilds(db *gorm.DB, warcraftLogsClient *warcraftlogs.WarcraftLogsClientService, maxRankingsPerSpec int) (
	*reportsRepository.ReportRepository,
	*rankingsRepository.RankingsRepository,
	*playerBuildsRepository.PlayerBuildsRepository,
	*buildsStatisticsRepository.BuildsStatisticsRepository,
	*talentStatisticsRepository.TalentStatisticsRepository,
	*statStatisticsRepository.StatStatisticsRepository,
	*workflowStatesRepository.WorkflowStateRepository,
	*activities.Activities,
) {
	// Initialiser les repositories
	reportsRepo := reportsRepository.NewReportRepository(db)
	rankingsRepo := rankingsRepository.NewRankingsRepository(db, maxRankingsPerSpec)
	playerBuildsRepo := playerBuildsRepository.NewPlayerBuildsRepository(db)
	buildsStatsRepo := buildsStatisticsRepository.NewBuildsStatisticsRepository(db)
	talentStatsRepo := talentStatisticsRepository.NewTalentStatisticsRepository(db)
	statStatsRepo := statStatisticsRepository.NewStatStatisticsRepository(db)
	workflowStatesRepo := workflowStatesRepository.NewWorkflowStateRepository(db)

	// Initialiser les activités
	rankingsActivity := activities.NewRankingsActivity(warcraftLogsClient, rankingsRepo)
	reportsActivity := activities.NewReportsActivity(warcraftLogsClient, reportsRepo, rankingsRepo)
	playerBuildsActivity := activities.NewPlayerBuildsActivity(playerBuildsRepo, reportsRepo)
	rateLimitActivity := activities.NewRateLimitActivity(warcraftLogsClient)
	workflowStatesActivity := activities.NewWorkflowStateActivity(workflowStatesRepo)

	// Activities pour les analyses de builds
	buildsStatisticsActivity := activities.NewBuildsStatisticsActivity(
		playerBuildsRepo,
		buildsStatsRepo,
	)
	talentStatisticActivity := activities.NewTalentStatisticActivity(
		playerBuildsRepo,
		talentStatsRepo,
	)
	statStatisticsActivity := activities.NewStatStatisticsActivity(
		playerBuildsRepo,
		statStatsRepo,
	)

	// Créer le service d'activités
	activitiesService := &activities.Activities{
		Rankings:         rankingsActivity,
		Reports:          reportsActivity,
		PlayerBuilds:     playerBuildsActivity,
		RateLimit:        rateLimitActivity,
		BuildStatistics:  buildsStatisticsActivity,
		TalentStatistics: talentStatisticActivity,
		StatStatistics:   statStatisticsActivity,
		WorkflowState:    workflowStatesActivity,
	}

	return reportsRepo, rankingsRepo, playerBuildsRepo, buildsStatsRepo, talentStatsRepo, statStatsRepo, workflowStatesRepo, activitiesService
}

// RegisterBuilds enregistre les workflows et activités avec le worker
func RegisterBuilds(w worker.Worker, activitiesService *activities.Activities) {
	// Enregistrer les workflows
	rankingsWorkflowImpl := rankingsWorkflow.NewRankingsWorkflow()
	reportsWorkflowImpl := reportsWorkflow.NewReportsWorkflow()
	buildsBatchWorkflowImpl := buildsWorkflow.NewBuildsBatchWorkflow()
	buildsWorkflowImpl := buildsWorkflow.NewBuildsWorkflow()
	equipmentAnalysisWorkflowImpl := equipmentAnalysisWorkflow.NewEquipmentAnalysisWorkflow()
	talentAnalysisWorkflowImpl := talentAnalysisWorkflow.NewTalentAnalysisWorkflow()
	statAnalysisWorkflowImpl := statAnalysisWorkflow.NewStatAnalysisWorkflow()

	// Enregistrer les workflows
	w.RegisterWorkflowWithOptions(rankingsWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: definitions.RankingsWorkflowName,
	})
	w.RegisterWorkflowWithOptions(reportsWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: definitions.ReportsWorkflowName,
	})
	w.RegisterWorkflowWithOptions(buildsBatchWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: definitions.ProcessBuildsBatchWorkflow,
	})
	w.RegisterWorkflowWithOptions(buildsWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: definitions.BuildsWorkflowName,
	})
	w.RegisterWorkflowWithOptions(equipmentAnalysisWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: definitions.AnalyzeBuildsWorkflowName,
	})
	w.RegisterWorkflowWithOptions(talentAnalysisWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: definitions.AnalyzeTalentsWorkflowName,
	})
	w.RegisterWorkflowWithOptions(statAnalysisWorkflowImpl.Execute, workflow.RegisterOptions{
		Name: definitions.AnalyzeStatStatisticsWorkflowName,
	})

	// Enregistrer les activities
	// Rankings activities
	w.RegisterActivity(activitiesService.Rankings.FetchAndStore)
	w.RegisterActivity(activitiesService.Rankings.GetStoredRankings)
	w.RegisterActivity(activitiesService.Rankings.MarkRankingsForReportProcessing)

	// Reports activities
	w.RegisterActivity(activitiesService.Reports.ProcessReports)
	w.RegisterActivity(activitiesService.Reports.GetReportsBatch)
	w.RegisterActivity(activitiesService.Reports.CountAllReports)
	w.RegisterActivity(activitiesService.Reports.GetUniqueReportReferences)
	w.RegisterActivity(activitiesService.Reports.GetRankingsNeedingReportProcessing)
	w.RegisterActivity(activitiesService.Reports.MarkReportsForBuildProcessing)

	// Player builds activities
	w.RegisterActivity(activitiesService.PlayerBuilds.ProcessAllBuilds)
	w.RegisterActivity(activitiesService.PlayerBuilds.CountPlayerBuilds)
	w.RegisterActivity(activitiesService.PlayerBuilds.GetReportsNeedingBuildExtraction)
	w.RegisterActivity(activitiesService.PlayerBuilds.MarkReportsAsProcessedForBuilds)
	w.RegisterActivity(activitiesService.PlayerBuilds.CountReportsNeedingBuildExtraction)

	// Rate limit activities
	w.RegisterActivity(activitiesService.RateLimit.ReservePoints)
	w.RegisterActivity(activitiesService.RateLimit.ReleasePoints)
	w.RegisterActivity(activitiesService.RateLimit.CheckRemainingPoints)

	// Build statistics activities
	w.RegisterActivity(activitiesService.BuildStatistics.ProcessItemStatistics)
	w.RegisterActivity(activitiesService.TalentStatistics.ProcessTalentStatistics)
	w.RegisterActivity(activitiesService.StatStatistics.ProcessStatStatistics)

	// Workflow state activities
	w.RegisterActivity(activitiesService.WorkflowState.CreateWorkflowState)
	w.RegisterActivity(activitiesService.WorkflowState.UpdateWorkflowState)
	w.RegisterActivity(activitiesService.WorkflowState.GetLastWorkflowRun)
	w.RegisterActivity(activitiesService.WorkflowState.GetWorkflowStatistics)
	w.RegisterActivity(activitiesService.WorkflowState.GetWorkflowStateByID)
	w.RegisterActivity(activitiesService.WorkflowState.DeleteOldWorkflowStates)
}
