// activities.go
package warcraftlogsBuildsTemporalWorkflowsDefinitions

import (
	"context"
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	models "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows/models"
)

/* This file :

- Define all activities const
- Contains all interface for each type of activity
- Define the type for each specific activity

*/

// Activity name constants - matching exactly with the activity methods
const (
	// Rankings activities
	FetchRankingsActivity         = "FetchAndStore"                   // Fetch and store rankings
	GetStoredRankingsActivity     = "GetStoredRankings"               // Get stored rankings
	MarkRankingsForReportActivity = "MarkRankingsForReportProcessing" // Mark rankings ready for reports processing

	// Reports activities
	ProcessReportsActivity                     = "ProcessReports"                     // Process reports
	GetReportsBatchActivity                    = "GetReportsBatch"                    // Get reports batch
	CountAllReportsActivity                    = "CountAllReports"                    // Count all reports
	GetUniqueReportReferencesActivity          = "GetUniqueReportReferences"          // Get unique report references
	GetRankingsNeedingReportProcessingActivity = "GetRankingsNeedingReportProcessing" // Get rankings needing report processing
	MarkReportsForBuildProcessingActivity      = "MarkReportsForBuildProcessing"      // Mark reports for build processing

	// Player builds activities
	ProcessBuildsActivity                    = "ProcessAllBuilds"                 // Process all builds
	CountPlayerBuildsActivity                = "CountPlayerBuilds"                // Count player builds
	GetReportsNeedingBuildExtractionActivity = "GetReportsNeedingBuildExtraction" // Get reports needing build extraction
	MarkReportsAsProcessedForBuildsActivity  = "MarkReportsAsProcessedForBuilds"  // Mark reports as processed for builds

	// Rate limit activities
	ReserveRateLimitPointsActivity = "ReservePoints"        // Reserve rate limit points
	ReleaseRateLimitPointsActivity = "ReleasePoints"        // Release rate limit points
	CheckRemainingPointsActivity   = "CheckRemainingPoints" // Check remaining points

	// Build statistics activities
	ProcessBuildStatisticsActivity  = "ProcessItemStatistics"   // Analyze equipment
	ProcessTalentStatisticsActivity = "ProcessTalentStatistics" // Analyze talents
	ProcessStatStatisticsActivity   = "ProcessStatStatistics"   // Analyze statistics

	// Sub-workflow names
	ProcessBuildBatchWorkflowName = "ProcessBuildBatch"     // Process build batch
	SyncWorkflowName              = "SyncWorkflow"          // Sync workflow
	AnalyzeBuildsWorkflowName     = "AnalyzeBuildsWorkflow" // Analyze builds workflow
)

/*
	Those interfaces are used to define the activities that can be used in the workflow.
	They are not used in the workflow, but they are used to define the activities that can be used in the workflow.
*/

// RankingsActivity defines the interface for rankings-related activities
type RankingsActivity interface {
	FetchAndStore(ctx context.Context, spec models.ClassSpec, dungeon models.Dungeon, batchConfig models.BatchConfig) (*models.BatchResult, error)
	GetStoredRankings(ctx context.Context, className, specName string, encounterID uint32) ([]*warcraftlogsBuilds.ClassRanking, error)
	MarkRankingsForReportProcessing(ctx context.Context, className, specName string, encounterID uint32, batchID string) error
}

// ReportsActivity defines the interface for report-related activities
type ReportsActivity interface {
	ProcessReports(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking) (*models.BatchResult, error)
	GetReportsBatch(ctx context.Context, batchSize int32, offset int32) ([]*warcraftlogsBuilds.Report, error)
	CountAllReports(ctx context.Context) (int64, error)
	GetUniqueReportReferences(ctx context.Context) ([]*warcraftlogsBuilds.ClassRanking, error)
	GetRankingsNeedingReportProcessing(ctx context.Context, limit int32, maxAgeDuration time.Duration) ([]*warcraftlogsBuilds.ClassRanking, error)
	MarkReportsForBuildProcessing(ctx context.Context, reportCodes []string, batchID string) error
}

// PlayerBuildsActivity defines the interface for player build-related activities
type PlayerBuildsActivity interface {
	ProcessAllBuilds(ctx context.Context, reports []*warcraftlogsBuilds.Report) (*models.BatchResult, error)
	CountPlayerBuilds(ctx context.Context) (int64, error)
	GetReportsNeedingBuildExtraction(ctx context.Context, limit int32, maxAgeDuration time.Duration) ([]*warcraftlogsBuilds.Report, error)
	MarkReportsAsProcessedForBuilds(ctx context.Context, reportCodes []string, batchID string) error
}

// RateLimitActivity defines the interface for rate limiting operations
type RateLimitActivity interface {
	ReservePoints(ctx context.Context, params models.WorkflowConfig) error
	ReleasePoints(ctx context.Context, params models.WorkflowConfig) error
	CheckRemainingPoints(ctx context.Context, params models.WorkflowConfig) (float64, error)
}
