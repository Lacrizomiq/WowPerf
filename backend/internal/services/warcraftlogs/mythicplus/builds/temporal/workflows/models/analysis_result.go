// analysis_result.go
package warcraftlogsBuildsTemporalWorkflowsModels

import "time"

// AnalysisWorkflowResult represents the complete results of the analysis
type AnalysisWorkflowResult struct {
	ItemsAnalyzed     int32     `json:"items_analyzed"`     // Equipment statistics
	TalentsAnalyzed   int32     `json:"talents_analyzed"`   // Talent configurations
	StatsAnalyzed     int32     `json:"stats_analyzed"`     // Character statistics
	SpecsProcessed    int32     `json:"specs_processed"`    // Specializations processed
	DungeonsProcessed int32     `json:"dungeons_processed"` // Dungeons processed
	TotalBuilds       int32     `json:"total_builds"`       // Total number of builds analyzed
	StartedAt         time.Time `json:"started_at"`
	CompletedAt       time.Time `json:"completed_at"`
}
