package raiderioMythicPlusRunsActivities

// Activities is a struct that contains all the activities for the Temporal worker
type Activities struct {
	MythicPlusRuns *MythicPlusRunsActivity
}

// NewActivities creates a new instance of Activities
func NewActivities(
	mythicPlusRunsActivity *MythicPlusRunsActivity,
) *Activities {
	return &Activities{
		MythicPlusRuns: mythicPlusRunsActivity,
	}
}
