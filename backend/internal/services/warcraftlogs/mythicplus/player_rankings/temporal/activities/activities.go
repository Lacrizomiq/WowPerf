package warcraftlogsPlayerRankingsActivities

// Activities is a struct that contains all the activities for the Temporal worker
type Activities struct {
	PlayerRankings *PlayerRankingsActivity
}

// NewActivities creates a new instance of Activities
func NewActivities(
	playerRankingsActivity *PlayerRankingsActivity,
) *Activities {
	return &Activities{
		PlayerRankings: playerRankingsActivity,
	}
}
