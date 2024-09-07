package raids

type ExpansionRaids struct {
	Expansions []ExpansionWithRaids `json:"expansions"`
}

type ExpansionWithRaids struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Raids []Raids `json:"raids"`
}

type Raids struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Modes []Mode `json:"modes"`
}

type Mode struct {
	Difficulty string   `json:"difficulty"`
	Progress   Progress `json:"progress"`
	Status     string   `json:"status"`
}

type Progress struct {
	CompletedCount int                 `json:"completed_count"`
	TotalCount     int                 `json:"total_count"`
	Encounters     []EncounterProgress `json:"encounters"`
}

type EncounterProgress struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	CompletedCount    int    `json:"completed_count"`
	LastKillTimestamp int64  `json:"last_kill_timestamp"`
}
