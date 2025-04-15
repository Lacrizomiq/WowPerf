// analysis_result.go
package warcraftlogsBuildsTemporalWorkflowsModels

import "time"

// EquipmentAnalysisWorkflowResult represents the complete results of the equipment analysis
// It contains statistics on analyzed player equipment.
type EquipmentAnalysisWorkflowResult struct {
	ItemsAnalyzed     int32     `json:"items_analyzed"`     // Equipment statistics
	SpecsProcessed    int32     `json:"specs_processed"`    // Specializations processed
	DungeonsProcessed int32     `json:"dungeons_processed"` // Dungeons processed
	TotalBuilds       int32     `json:"total_builds"`       // Total number of builds analyzed
	StartedAt         time.Time `json:"started_at"`
	CompletedAt       time.Time `json:"completed_at"`
}

// TalentAnalysisWorkflowResult represents the complete results of the talent analysis
// It contains statistics on analyzed player talent builds.
type TalentAnalysisWorkflowResult struct {
	TalentsAnalyzed   int32     `json:"talents_analyzed"`   // Talent configurations
	SpecsProcessed    int32     `json:"specs_processed"`    // Specializations processed
	DungeonsProcessed int32     `json:"dungeons_processed"` // Dungeons processed
	TotalBuilds       int32     `json:"total_builds"`       // Total number of builds analyzed
	StartedAt         time.Time `json:"started_at"`
	CompletedAt       time.Time `json:"completed_at"`
}

// StatAnalysisWorkflowResult represents the complete results of the stats analysis
// It contains statistics on analyzed player stats distribution.
type StatAnalysisWorkflowResult struct {
	StatsAnalyzed     int32     `json:"stats_analyzed"`     // Character statistics
	SpecsProcessed    int32     `json:"specs_processed"`    // Specializations processed
	DungeonsProcessed int32     `json:"dungeons_processed"` // Dungeons processed
	TotalBuilds       int32     `json:"total_builds"`       // Total number of builds analyzed
	StartedAt         time.Time `json:"started_at"`
	CompletedAt       time.Time `json:"completed_at"`
}

// AnalysisWorkflowResult represents the complete results of the analysis
// Note: This will not be used anymore.
// It is kept here for reference and in case we need to use it again.
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
