package warcraftlogsBuildsTemporalActivities_test

import (
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	activities "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/activities"
)

// TestTalentTransformation tests the transformation of talent data into TalentStatistic
func TestTalentTransformation(t *testing.T) {
	t.Log("Starting TestTalentTransformation")

	// 1. Create a PlayerBuild with the example data
	talentImport := "CAQAAAAAAAAAAAAAAAAAAAAAAADw2sMzYwyYMzMmZzsNzMzMMAAAAAAAAAAAwysMbjHwMzgZhhBjhZhtZaMxyAmZAgAMbz2GYsZD"

	// Use the example data provided
	playerBuild := &warcraftlogsBuilds.PlayerBuild{
		ID:            18933,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		PlayerName:    "Lyöko",
		Class:         "Priest",
		Spec:          "Discipline",
		ReportCode:    "wVkK6XBjt2ygD8mR",
		FightID:       11,
		ActorID:       1,
		TalentImport:  talentImport,
		TalentTree:    datatypes.JSON(`[{"id": 103678, "rank": 1, "nodeID": 82710}]`), // Pas utilisé dans le test
		ItemLevel:     636,
		EncounterID:   62286,
		KeystoneLevel: 16,
		Affixes:       pq.Int64Array{10, 152, 9, 147},
	}

	// 2. Create an instance of the activity
	t.Log("Creating the activity")
	activity := &activities.TalentStatisticActivity{}

	// 3. Call the ProcessTalentsBatch method directly
	t.Log("Calling the ProcessTalentsBatch method")
	stats, err := activity.ProcessTalentsBatch([]*warcraftlogsBuilds.PlayerBuild{playerBuild})

	// 4. Verifications
	if err != nil {
		t.Logf("ERROR: %v", err)
	} else {
		t.Logf("Transformation successful, number of statistics generated: %d", len(stats))
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, stats)
	assert.Equal(t, 1, len(stats), "One talent configuration expected")

	// 5. Verify the details of the generated talent statistic
	talentStat := stats[0]
	t.Log("Details of the generated talent statistic:")
	t.Logf("  Class: %s", talentStat.Class)
	t.Logf("  Spec: %s", talentStat.Spec)
	t.Logf("  EncounterID: %d", talentStat.EncounterID)
	t.Logf("  TalentImport: %s", talentStat.TalentImport)
	t.Logf("  UsageCount: %d", talentStat.UsageCount)
	t.Logf("  ItemLevel: avg=%.1f, min=%.1f, max=%.1f", talentStat.AvgItemLevel, talentStat.MinItemLevel, talentStat.MaxItemLevel)
	t.Logf("  KeystoneLevel: avg=%.1f, min=%d, max=%d", talentStat.AvgKeystoneLevel, talentStat.MinKeystoneLevel, talentStat.MaxKeystoneLevel)

	// 6. Assertions
	assert.Equal(t, "Priest", talentStat.Class)
	assert.Equal(t, "Discipline", talentStat.Spec)
	assert.Equal(t, uint(62286), talentStat.EncounterID)
	assert.Equal(t, talentImport, talentStat.TalentImport)
	assert.Equal(t, 1, talentStat.UsageCount)
	assert.Equal(t, 636.0, talentStat.AvgItemLevel)
	assert.Equal(t, 636.0, talentStat.MinItemLevel)
	assert.Equal(t, 636.0, talentStat.MaxItemLevel)
	assert.Equal(t, 16.0, talentStat.AvgKeystoneLevel)
	assert.Equal(t, 16, talentStat.MinKeystoneLevel)
	assert.Equal(t, 16, talentStat.MaxKeystoneLevel)

	t.Log("End of TestTalentTransformation")
}

// TestUsagePercentageCalculationForTalents tests the calculation of usage percentages for talents
func TestUsagePercentageCalculationForTalents(t *testing.T) {
	t.Log("Starting TestUsagePercentageCalculationForTalents")

	// 1. Create test data with different talent configurations
	testTalents := []*warcraftlogsBuilds.TalentStatistic{
		{
			Class:        "Priest",
			Spec:         "Discipline",
			EncounterID:  62286,
			TalentImport: "CONFIG_A",
			UsageCount:   75,
		},
		{
			Class:        "Priest",
			Spec:         "Discipline",
			EncounterID:  62286,
			TalentImport: "CONFIG_B",
			UsageCount:   25,
		},
	}

	t.Log("Talent statistics before calculation:")
	for i, stat := range testTalents {
		t.Logf("  Talent %d: Import=%s, Count=%d, %%=%.1f",
			i+1, stat.TalentImport, stat.UsageCount, stat.UsagePercentage)
	}

	// 2. Create the activity and call the CalculateUsagePercentages method
	activity := &activities.TalentStatisticActivity{}
	t.Log("Calling the CalculateUsagePercentages method")
	activity.CalculateUsagePercentages(testTalents, 100)

	t.Log("Talent statistics after calculation:")
	for i, stat := range testTalents {
		t.Logf("  Talent %d: Import=%s, Count=%d, %%=%.1f",
			i+1, stat.TalentImport, stat.UsageCount, stat.UsagePercentage)
	}

	// 3. Verify the results
	assert.Equal(t, 75.0, testTalents[0].UsagePercentage)
	assert.Equal(t, 25.0, testTalents[1].UsagePercentage)

	t.Log("End of TestUsagePercentageCalculationForTalents")
}

// TestMultipleTalentBuilds tests the processing of multiple builds with the same talent configuration
func TestMultipleTalentBuilds(t *testing.T) {
	t.Log("Starting TestMultipleTalentBuilds")

	// 1. Create builds with the same talent configuration but different levels
	talentImport := "CAQAAAAAAAAAAAAAAAAAAAAAAADw2sMzYwyYMzMmZzsNzMzMMAAAAAAAAAAAwysMbjHwMzgZhhBjhZhtZaMxyAmZAgAMbz2GYsZD"

	build1 := &warcraftlogsBuilds.PlayerBuild{
		ID:            1,
		PlayerName:    "Player1",
		Class:         "Priest",
		Spec:          "Discipline",
		TalentImport:  talentImport,
		ItemLevel:     636,
		EncounterID:   62286,
		KeystoneLevel: 16,
	}

	build2 := &warcraftlogsBuilds.PlayerBuild{
		ID:            2,
		PlayerName:    "Player2",
		Class:         "Priest",
		Spec:          "Discipline",
		TalentImport:  talentImport,
		ItemLevel:     645,
		EncounterID:   62286,
		KeystoneLevel: 18,
	}

	t.Log("Creating two builds with the same talent configuration")

	// 2. Process the builds
	activity := &activities.TalentStatisticActivity{}
	t.Log("Calling the ProcessTalentsBatch method with the two builds")
	stats, err := activity.ProcessTalentsBatch([]*warcraftlogsBuilds.PlayerBuild{build1, build2})

	// 3. Verifications
	if err != nil {
		t.Logf("ERROR: %v", err)
	} else {
		t.Logf("Transformation successful, number of statistics generated: %d", len(stats))
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, stats)
	assert.Equal(t, 1, len(stats), "One talent configuration expected for the two builds")

	// 4. Verify the aggregation of data
	talentStat := stats[0]
	t.Log("Details of the aggregated talent statistic:")
	t.Logf("  TalentImport: %s", talentStat.TalentImport)
	t.Logf("  UsageCount: %d", talentStat.UsageCount)
	t.Logf("  ItemLevel: avg=%.1f, min=%.1f, max=%.1f", talentStat.AvgItemLevel, talentStat.MinItemLevel, talentStat.MaxItemLevel)
	t.Logf("  KeystoneLevel: avg=%.1f, min=%d, max=%d", talentStat.AvgKeystoneLevel, talentStat.MinKeystoneLevel, talentStat.MaxKeystoneLevel)

	// 5. Assertions
	assert.Equal(t, 2, talentStat.UsageCount)
	assert.InDelta(t, 640.5, talentStat.AvgItemLevel, 0.01) // (636 + 645) / 2
	assert.Equal(t, 636.0, talentStat.MinItemLevel)
	assert.Equal(t, 645.0, talentStat.MaxItemLevel)
	assert.InDelta(t, 17.0, talentStat.AvgKeystoneLevel, 0.01) // (16 + 18) / 2
	assert.Equal(t, 16, talentStat.MinKeystoneLevel)
	assert.Equal(t, 18, talentStat.MaxKeystoneLevel)

	t.Log("End of TestMultipleTalentBuilds")
}
