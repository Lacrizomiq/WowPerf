package warcraftlogsBuildsRepository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
)

// createTestBuildStatistics creates a test build statistic for testing
func createTestBuildStatistics() *warcraftlogsBuilds.BuildStatistic {
	return &warcraftlogsBuilds.BuildStatistic{
		Class:            "Priest",
		Spec:             "Discipline",
		EncounterID:      12661, // Corresponding to Cinderbrew Meadery
		ItemSlot:         1,     // Neck slot
		ItemID:           221060,
		ItemName:         "Strapped Rescue-Keg",
		ItemIcon:         "inv_11_0_earthen_earthennecklace01_color2.jpg",
		ItemQuality:      3,
		ItemLevel:        658,
		HasGems:          true,
		GemsCount:        2,
		GemIDs:           []int64{213494, 213746},
		GemIcons:         []string{"inv_jewelcrafting_cut-standart-gem-hybrid_color2_3.jpg", "inv_misc_gem_x4_metagem_cut.jpg"},
		GemLevels:        []float64{610, 610},
		BonusIDs:         []int64{10390, 6652, 10383, 10879, 10396, 11988, 3150, 10255},
		UsageCount:       1,
		UsagePercentage:  100,
		AvgItemLevel:     658,
		MinItemLevel:     658,
		MaxItemLevel:     658,
		AvgKeystoneLevel: 14,
		MinKeystoneLevel: 14,
		MaxKeystoneLevel: 14,
	}
}

// clearTestData removes all test data from the database
func clearTestData(t *testing.T, repo *BuildsStatisticsRepository, class, spec string, encounterID uint) {
	ctx := context.Background()
	err := repo.DeleteBuildStatistics(ctx, class, spec, encounterID)
	require.NoError(t, err, "Failed to clear test data")
}

// TestDeleteBuildStatistics tests the DeleteBuildStatistics method
func TestDeleteBuildStatistics(t *testing.T) {
	db := setupTestDatabase(t)
	repo := NewBuildsStatisticsRepository(db)
	ctx := context.Background()

	// Create test data
	buildStat := createTestBuildStatistics()

	// Store the test data first
	err := repo.StoreManyBuildStatistics(ctx, []*warcraftlogsBuilds.BuildStatistic{buildStat})
	require.NoError(t, err, "Failed to store test data")

	// Verify that data was stored
	stats, err := repo.GetBuildStatistics(ctx, buildStat.Class, buildStat.Spec, buildStat.EncounterID)
	require.NoError(t, err)
	require.NotEmpty(t, stats, "Test data should be stored")

	// Test: Delete the data
	err = repo.DeleteBuildStatistics(ctx, buildStat.Class, buildStat.Spec, buildStat.EncounterID)
	assert.NoError(t, err)

	// Verify data was deleted
	stats, err = repo.GetBuildStatistics(ctx, buildStat.Class, buildStat.Spec, buildStat.EncounterID)
	assert.NoError(t, err)
	assert.Empty(t, stats, "Data should be deleted")
}

// TestStoreManyBuildStatistics tests the StoreManyBuildStatistics method
func TestStoreManyBuildStatistics(t *testing.T) {
	db := setupTestDatabase(t)
	repo := NewBuildsStatisticsRepository(db)
	ctx := context.Background()

	// Create test data
	buildStat1 := createTestBuildStatistics()
	buildStat2 := createTestBuildStatistics()
	buildStat2.ItemSlot = 2
	buildStat2.ItemID = 212081
	buildStat2.ItemName = "Living Luster's Dominion"

	// Clean up before testing
	clearTestData(t, repo, buildStat1.Class, buildStat1.Spec, buildStat1.EncounterID)

	// Test 1: Store multiple build statistics
	t.Run("store multiple build statistics", func(t *testing.T) {
		err := repo.StoreManyBuildStatistics(ctx, []*warcraftlogsBuilds.BuildStatistic{buildStat1, buildStat2})
		assert.NoError(t, err)

		// Verify data was stored
		stats, err := repo.GetBuildStatistics(ctx, buildStat1.Class, buildStat1.Spec, buildStat1.EncounterID)
		assert.NoError(t, err)
		assert.Len(t, stats, 2, "Should have stored 2 statistics")
	})

	// Test 2: Update existing statistic
	t.Run("update existing statistic", func(t *testing.T) {
		// Modify usage count and percentage
		buildStat1.UsageCount = 2
		buildStat1.UsagePercentage = 75
		buildStat2.UsageCount = 2
		buildStat2.UsagePercentage = 25

		err := repo.StoreManyBuildStatistics(ctx, []*warcraftlogsBuilds.BuildStatistic{buildStat1, buildStat2})
		assert.NoError(t, err)

		// Verify data was updated
		stats, err := repo.GetBuildStatistics(ctx, buildStat1.Class, buildStat1.Spec, buildStat1.EncounterID)
		assert.NoError(t, err)

		// Find our updated buildStat1
		var updatedStat *warcraftlogsBuilds.BuildStatistic
		for _, s := range stats {
			if s.ItemSlot == buildStat1.ItemSlot && s.ItemID == buildStat1.ItemID {
				updatedStat = s
				break
			}
		}

		require.NotNil(t, updatedStat, "Updated statistic should exist")
		assert.Equal(t, 2, updatedStat.UsageCount)
		assert.Equal(t, float64(75), updatedStat.UsagePercentage)
	})

	// Test 3: Empty batch
	t.Run("store empty batch", func(t *testing.T) {
		err := repo.StoreManyBuildStatistics(ctx, []*warcraftlogsBuilds.BuildStatistic{})
		assert.NoError(t, err)
	})

	// Clean up after tests
	clearTestData(t, repo, buildStat1.Class, buildStat1.Spec, buildStat1.EncounterID)
}

// TestGetBuildStatistics tests the GetBuildStatistics method
func TestGetBuildStatistics(t *testing.T) {
	db := setupTestDatabase(t)
	repo := NewBuildsStatisticsRepository(db)
	ctx := context.Background()

	// Create test data
	buildStat := createTestBuildStatistics()

	// Clean up before testing
	clearTestData(t, repo, buildStat.Class, buildStat.Spec, buildStat.EncounterID)

	// Store test data
	err := repo.StoreManyBuildStatistics(ctx, []*warcraftlogsBuilds.BuildStatistic{buildStat})
	require.NoError(t, err, "Failed to store test data")

	// Test: Get data with encounterID
	t.Run("get data with encounterID", func(t *testing.T) {
		stats, err := repo.GetBuildStatistics(ctx, buildStat.Class, buildStat.Spec, buildStat.EncounterID)
		assert.NoError(t, err)
		assert.NotEmpty(t, stats)
		assert.Equal(t, buildStat.ItemName, stats[0].ItemName)
		assert.Equal(t, buildStat.ItemIcon, stats[0].ItemIcon)
		assert.Equal(t, buildStat.HasGems, stats[0].HasGems)
		assert.Equal(t, buildStat.GemsCount, stats[0].GemsCount)
	})

	// Test: Get data without encounterID
	t.Run("get data without encounterID", func(t *testing.T) {
		stats, err := repo.GetBuildStatistics(ctx, buildStat.Class, buildStat.Spec, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, stats)
	})

	// Clean up after tests
	clearTestData(t, repo, buildStat.Class, buildStat.Spec, buildStat.EncounterID)
}

// TestCountBuildStatistics tests the CountBuildStatistics method
func TestCountBuildStatistics(t *testing.T) {
	db := setupTestDatabase(t)
	repo := NewBuildsStatisticsRepository(db)
	ctx := context.Background()

	// Create test data
	buildStat1 := createTestBuildStatistics()
	buildStat2 := createTestBuildStatistics()
	buildStat2.ItemSlot = 2
	buildStat2.ItemID = 212081
	buildStat2.ItemName = "Living Luster's Dominion"

	// Clean up before testing
	clearTestData(t, repo, buildStat1.Class, buildStat1.Spec, buildStat1.EncounterID)

	// Store test data
	err := repo.StoreManyBuildStatistics(ctx, []*warcraftlogsBuilds.BuildStatistic{buildStat1, buildStat2})
	require.NoError(t, err, "Failed to store test data")

	// Test: Count with encounterID
	t.Run("count with encounterID", func(t *testing.T) {
		count, err := repo.CountBuildStatistics(ctx, buildStat1.Class, buildStat1.Spec, buildStat1.EncounterID)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})

	// Test: Count without encounterID
	t.Run("count without encounterID", func(t *testing.T) {
		count, err := repo.CountBuildStatistics(ctx, buildStat1.Class, buildStat1.Spec, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})

	// Test: Count with different class
	t.Run("count with different class", func(t *testing.T) {
		count, err := repo.CountBuildStatistics(ctx, "Mage", buildStat1.Spec, buildStat1.EncounterID)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	// Clean up after tests
	clearTestData(t, repo, buildStat1.Class, buildStat1.Spec, buildStat1.EncounterID)
}
