// activities.go
package warcraftlogsBuildsTemporalWorkflowsDefinitions

import (
	"context"

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
	FetchRankingsActivity     = "FetchAndStore"     // Fetch and store rankings
	GetStoredRankingsActivity = "GetStoredRankings" // Get stored rankings

	// Reports activities
	ProcessReportsActivity            = "ProcessReports"            // Process reports
	GetReportsBatchActivity           = "GetReportsBatch"           // Get reports batch
	CountAllReportsActivity           = "CountAllReports"           // Count all reports
	GetUniqueReportReferencesActivity = "GetUniqueReportReferences" // Get unique report references

	// Player builds activities
	ProcessBuildsActivity     = "ProcessAllBuilds"  // Process all builds
	CountPlayerBuildsActivity = "CountPlayerBuilds" // Count player builds

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
}

// ReportsActivity defines the interface for report-related activities
type ReportsActivity interface {
	ProcessReports(ctx context.Context, rankings []*warcraftlogsBuilds.ClassRanking) (*models.BatchResult, error)
	GetReportsBatch(ctx context.Context, batchSize int32, offset int32) ([]*warcraftlogsBuilds.Report, error)
	CountAllReports(ctx context.Context) (int64, error)
	GetUniqueReportReferences(ctx context.Context) ([]*warcraftlogsBuilds.ClassRanking, error)
}

// PlayerBuildsActivity defines the interface for player build-related activities
type PlayerBuildsActivity interface {
	ProcessAllBuilds(ctx context.Context, reports []*warcraftlogsBuilds.Report) (*models.BatchResult, error)
	CountPlayerBuilds(ctx context.Context) (int64, error)
}

// RateLimitActivity defines the interface for rate limiting operations
type RateLimitActivity interface {
	ReservePoints(ctx context.Context, params models.WorkflowConfig) error
	ReleasePoints(ctx context.Context, params models.WorkflowConfig) error
	CheckRemainingPoints(ctx context.Context, params models.WorkflowConfig) (float64, error)
}
