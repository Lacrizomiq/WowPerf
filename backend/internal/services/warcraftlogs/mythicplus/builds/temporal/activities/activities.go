package warcraftlogsBuildsTemporalActivities

// Activities is a struct that contains all the activities for the Temporal worker
type Activities struct {
	Rankings     *RankingsActivity
	Reports      *ReportsActivity
	PlayerBuilds *PlayerBuildsActivity
}

// NewActivities creates a new instance of Activities
func NewActivities(
	rankingsActivity *RankingsActivity,
	reportsActivity *ReportsActivity,
	playerBuildsActivity *PlayerBuildsActivity,
) *Activities {
	return &Activities{
		Rankings:     rankingsActivity,
		Reports:      reportsActivity,
		PlayerBuilds: playerBuildsActivity,
	}
}
