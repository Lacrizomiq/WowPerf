package warcraftlogsBuildsTemporalActivities

// Activities is a struct that contains all the activities for the Temporal worker
type Activities struct {
	Rankings         *RankingsActivity
	Reports          *ReportsActivity
	PlayerBuilds     *PlayerBuildsActivity
	RateLimit        *RateLimitActivity
	BuildStatistics  *BuildsStatisticsActivity
	TalentStatistics *TalentStatisticActivity
	StatStatistics   *StatStatisticsActivity
}

// NewActivities creates a new instance of Activities
func NewActivities(
	rankingsActivity *RankingsActivity,
	reportsActivity *ReportsActivity,
	playerBuildsActivity *PlayerBuildsActivity,
	rateLimitActivity *RateLimitActivity,
	buildStatisticsActivity *BuildsStatisticsActivity,
	statStatisticsActivity *StatStatisticsActivity,
	talentStatisticsActivity *TalentStatisticActivity,
) *Activities {
	return &Activities{
		Rankings:         rankingsActivity,
		Reports:          reportsActivity,
		PlayerBuilds:     playerBuildsActivity,
		RateLimit:        rateLimitActivity,
		BuildStatistics:  buildStatisticsActivity,
		StatStatistics:   statStatisticsActivity,
		TalentStatistics: talentStatisticsActivity,
	}
}
