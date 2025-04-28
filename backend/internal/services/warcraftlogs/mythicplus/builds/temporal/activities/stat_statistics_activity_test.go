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

// TestStatTransformation tests the transformation of stats data to StatStatistic
func TestStatTransformation(t *testing.T) {
	t.Log("Starting TestStatTransformation")

	// 1. Create a PlayerBuild with example data
	statsJSON := `{"Crit": {"max": 3762, "min": 3762}, "Haste": {"max": 24047, "min": 24047}, "Leech": {"max": 2995, "min": 2995}, "Speed": {"max": 0, "min": 0}, "Mastery": {"max": 3152, "min": 3152}, "Stamina": {"max": 331917, "min": 331917}, "Avoidance": {"max": 0, "min": 0}, "Intellect": {"max": 65128, "min": 65128}, "Item Level": {"max": 636, "min": 636}, "Versatility": {"max": 11708, "min": 11708}}`

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
		ItemLevel:     636,
		Stats:         datatypes.JSON(statsJSON),
		EncounterID:   62286,
		KeystoneLevel: 16,
		Affixes:       pq.Int64Array{10, 152, 9, 147},
	}

	// 2. Create an instance of the activity
	t.Log("Creating the activity")
	activity := &activities.StatStatisticsActivity{}

	// 3. Create the map to store the aggregated data
	statData := make(map[string]*activities.StatAggregation)

	// 4. Call the ProcessStatsBatch method
	t.Log("Calling the ProcessStatsBatch method")
	err := activity.ProcessStatsBatch([]*warcraftlogsBuilds.PlayerBuild{playerBuild}, statData)

	// 5. Verifications
	assert.NoError(t, err)
	assert.NotEmpty(t, statData)

	// 6. Convert the aggregated data to statistics
	t.Log("Converting the aggregated data to statistics")
	stats := activity.ConvertToStatStatistics(statData, "Priest", "Discipline", 62286)

	// Verify the number of statistics extracted (secondary + minor)
	// We expect: Crit, Haste, Mastery, Versatility (secondary) + Leech, Avoidance, Speed (minor)
	expectedStatsCount := 7
	assert.Equal(t, expectedStatsCount, len(stats), "Incorrect number of statistics extracted")

	// 7. Verify the details of each statistic
	t.Log("Details of the extracted statistics:")

	// Map to verify the presence of expected statistics
	statMap := make(map[string]*warcraftlogsBuilds.StatStatistic)
	for _, stat := range stats {
		statMap[stat.StatName] = stat
		t.Logf("  Stat: %s (category: %s)", stat.StatName, stat.StatCategory)
		t.Logf("    Values: avg=%.1f, min=%.1f, max=%.1f", stat.AvgValue, stat.MinValue, stat.MaxValue)
		t.Logf("    Sample size: %d", stat.SampleSize)
		t.Logf("    Item Level: avg=%.1f, min=%.1f, max=%.1f", stat.AvgItemLevel, stat.MinItemLevel, stat.MaxItemLevel)
		t.Logf("    Keystone: avg=%.1f, min=%d, max=%d", stat.AvgKeystoneLevel, stat.MinKeystoneLevel, stat.MaxKeystoneLevel)
	}

	// 8. Assertions for secondary stats
	secondaryStats := []string{"Crit", "Haste", "Mastery", "Versatility"}
	for _, statName := range secondaryStats {
		stat, exists := statMap[statName]
		assert.True(t, exists, "The stat %s should be extracted", statName)
		if exists {
			assert.Equal(t, "secondary", stat.StatCategory)
			assert.Equal(t, 1, stat.SampleSize)
			assert.Equal(t, "Priest", stat.Class)
			assert.Equal(t, "Discipline", stat.Spec)
			assert.Equal(t, uint(62286), stat.EncounterID)
			assert.Equal(t, 636.0, stat.AvgItemLevel)
			assert.Equal(t, 16.0, stat.AvgKeystoneLevel)
		}
	}

	// Verify some specific values
	critStat := statMap["Crit"]
	if critStat != nil {
		assert.InDelta(t, 3762.0, critStat.AvgValue, 1.0)
		assert.InDelta(t, 3762.0, critStat.MinValue, 1.0)
		assert.InDelta(t, 3762.0, critStat.MaxValue, 1.0)
	}

	hasteStat := statMap["Haste"]
	if hasteStat != nil {
		assert.InDelta(t, 24047.0, hasteStat.AvgValue, 1.0)
		assert.InDelta(t, 24047.0, hasteStat.MinValue, 1.0)
		assert.InDelta(t, 24047.0, hasteStat.MaxValue, 1.0)
	}

	// 9. Assertions for minor stats
	minorStats := []string{"Leech", "Avoidance", "Speed"}
	for _, statName := range minorStats {
		stat, exists := statMap[statName]
		assert.True(t, exists, "The stat %s should be extracted", statName)
		if exists {
			assert.Equal(t, "minor", stat.StatCategory)
		}
	}

	leechStat := statMap["Leech"]
	if leechStat != nil {
		assert.InDelta(t, 2995.0, leechStat.AvgValue, 1.0)
	}

	t.Log("End of TestStatTransformation")
}

// TestMultipleStatBuilds tests the processing of multiple builds for the aggregation of stats
func TestMultipleStatBuilds(t *testing.T) {
	t.Log("Starting TestMultipleStatBuilds")

	// 1. Create two builds with different stats values
	statsJSON1 := `{"Crit": {"max": 3000, "min": 3000}, "Haste": {"max": 20000, "min": 20000}}`
	statsJSON2 := `{"Crit": {"max": 4000, "min": 4000}, "Haste": {"max": 30000, "min": 30000}}`

	build1 := &warcraftlogsBuilds.PlayerBuild{
		ID:            1,
		PlayerName:    "Player1",
		Class:         "Priest",
		Spec:          "Discipline",
		ItemLevel:     630,
		Stats:         datatypes.JSON(statsJSON1),
		EncounterID:   62286,
		KeystoneLevel: 15,
	}

	build2 := &warcraftlogsBuilds.PlayerBuild{
		ID:            2,
		PlayerName:    "Player2",
		Class:         "Priest",
		Spec:          "Discipline",
		ItemLevel:     650,
		Stats:         datatypes.JSON(statsJSON2),
		EncounterID:   62286,
		KeystoneLevel: 20,
	}

	t.Log("Creating two builds with different stats values")

	// 2. Create the map to store the aggregated data
	statData := make(map[string]*activities.StatAggregation)

	// 3. Process the builds
	activity := &activities.StatStatisticsActivity{}
	t.Log("Calling the ProcessStatsBatch method with the two builds")
	err := activity.ProcessStatsBatch([]*warcraftlogsBuilds.PlayerBuild{build1, build2}, statData)
	assert.NoError(t, err)

	// 4. Convert the aggregated data to statistics
	stats := activity.ConvertToStatStatistics(statData, "Priest", "Discipline", 62286)
	assert.Equal(t, 2, len(stats), "Two types of stats expected: Crit and Haste")

	// 5. Map of stats by name
	statMap := make(map[string]*warcraftlogsBuilds.StatStatistic)
	for _, stat := range stats {
		statMap[stat.StatName] = stat
		t.Logf("  Stat: %s", stat.StatName)
		t.Logf("    Values: avg=%.1f, min=%.1f, max=%.1f", stat.AvgValue, stat.MinValue, stat.MaxValue)
		t.Logf("    Sample size: %d", stat.SampleSize)
		t.Logf("    Item Level: avg=%.1f, min=%.1f, max=%.1f", stat.AvgItemLevel, stat.MinItemLevel, stat.MaxItemLevel)
		t.Logf("    Keystone: avg=%.1f, min=%d, max=%d", stat.AvgKeystoneLevel, stat.MinKeystoneLevel, stat.MaxKeystoneLevel)
	}

	// 6. Verify the aggregations for Crit
	critStat := statMap["Crit"]
	if critStat != nil {
		assert.Equal(t, 2, critStat.SampleSize)
		assert.InDelta(t, 3500.0, critStat.AvgValue, 1.0) // (3000 + 4000) / 2
		assert.Equal(t, 3000.0, critStat.MinValue)
		assert.Equal(t, 4000.0, critStat.MaxValue)
		assert.InDelta(t, 640.0, critStat.AvgItemLevel, 1.0) // (630 + 650) / 2
		assert.Equal(t, 630.0, critStat.MinItemLevel)
		assert.Equal(t, 650.0, critStat.MaxItemLevel)
		assert.InDelta(t, 17.5, critStat.AvgKeystoneLevel, 0.1) // (15 + 20) / 2
		assert.Equal(t, 15, critStat.MinKeystoneLevel)
		assert.Equal(t, 20, critStat.MaxKeystoneLevel)
	}

	// 7. Verify the aggregations for Haste
	hasteStat := statMap["Haste"]
	if hasteStat != nil {
		assert.Equal(t, 2, hasteStat.SampleSize)
		assert.InDelta(t, 25000.0, hasteStat.AvgValue, 1.0) // (20000 + 30000) / 2
		assert.Equal(t, 20000.0, hasteStat.MinValue)
		assert.Equal(t, 30000.0, hasteStat.MaxValue)
	}

	t.Log("End of TestMultipleStatBuilds")
}

// TestSkippedStats tests that stats that are neither secondary nor minor are ignored
func TestSkippedStats(t *testing.T) {
	t.Log("Starting TestSkippedStats")

	// 1. Create a build with primary and secondary stats
	statsJSON := `{
		"Crit": {"max": 3762, "min": 3762},
		"Intellect": {"max": 65128, "min": 65128},
		"Stamina": {"max": 331917, "min": 331917}
	}`

	playerBuild := &warcraftlogsBuilds.PlayerBuild{
		ID:            1,
		PlayerName:    "Player1",
		Class:         "Priest",
		Spec:          "Discipline",
		ItemLevel:     636,
		Stats:         datatypes.JSON(statsJSON),
		EncounterID:   62286,
		KeystoneLevel: 16,
	}

	// 2. Create the map to store the aggregated data
	statData := make(map[string]*activities.StatAggregation)

	// 3. Process the build
	activity := &activities.StatStatisticsActivity{}
	t.Log("Calling the ProcessStatsBatch method")
	err := activity.ProcessStatsBatch([]*warcraftlogsBuilds.PlayerBuild{playerBuild}, statData)
	assert.NoError(t, err)

	// 4. Convert the aggregated data to statistics
	stats := activity.ConvertToStatStatistics(statData, "Priest", "Discipline", 62286)

	// We expect only Crit (secondary) and not Intellect or Stamina (primary)
	assert.Equal(t, 1, len(stats), "Only one stat (Crit) should be extracted")

	if len(stats) > 0 {
		assert.Equal(t, "Crit", stats[0].StatName)
	}

	t.Log("End of TestSkippedStats")
}
