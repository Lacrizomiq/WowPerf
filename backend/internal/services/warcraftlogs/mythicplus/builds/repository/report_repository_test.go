// report_repository_test.go
package warcraftlogsBuildsRepository

import (
	"context"
	"os"
	"testing"
	"time"
	"wowperf/internal/database"
	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupTestDatabase initializes database connection for tests
func setupTestDatabase(t *testing.T) *gorm.DB {
	// Skip in CI/automated tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Try to load .env from multiple possible locations
	err := godotenv.Load()
	if err != nil {
		err = godotenv.Load("../../../../../../.env")
		require.NoError(t, err, "Failed to load .env file")
	}

	// Override DB_HOST for local testing
	os.Setenv("DB_HOST", "localhost")

	// Initialize database connection
	db, err := database.InitDB()
	require.NoError(t, err, "Failed to initialize database")

	return db
}

func TestReportRepository_GetReportsBatch(t *testing.T) {
	// Setup test database and repository
	db := setupTestDatabase(t)
	repo := NewReportRepository(db)
	ctx := context.Background()

	// Test 1: Check total number of reports in database
	t.Run("count total reports", func(t *testing.T) {
		var count int64
		err := db.Model(&warcraftlogsBuilds.Report{}).
			Where("deleted_at IS NULL").
			Count(&count).Error

		assert.NoError(t, err)
		t.Logf("Total reports in database: %d", count)
	})

	t.Run("investigate data discrepancy", func(t *testing.T) {
		// Count without deleted_at filter
		var totalCount int64
		err := db.Model(&warcraftlogsBuilds.Report{}).
			Count(&totalCount).Error
		assert.NoError(t, err)
		t.Logf("Total reports (including deleted): %d", totalCount)

		// Count with deleted_at filter
		var activeCount int64
		err = db.Model(&warcraftlogsBuilds.Report{}).
			Where("deleted_at IS NULL").
			Count(&activeCount).Error
		assert.NoError(t, err)
		t.Logf("Active reports (not deleted): %d", activeCount)

		// Sample some deleted records if they exist
		if totalCount > activeCount {
			var deletedSample []warcraftlogsBuilds.Report
			err = db.Model(&warcraftlogsBuilds.Report{}).
				Where("deleted_at IS NOT NULL").
				Limit(5).
				Find(&deletedSample).Error
			assert.NoError(t, err)
			for _, report := range deletedSample {
				t.Logf("Deleted report sample: code=%s, deleted_at=%v",
					report.Code, report.DeletedAt)
			}
		}
	})

	// Test 2: Verify first batch retrieval
	t.Run("first batch should return records", func(t *testing.T) {
		reports, err := repo.GetReportsBatch(ctx, 50, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, reports)

		t.Logf("First batch contains %d reports", len(reports))

		// Verify report structure
		for _, report := range reports {
			assert.NotEmpty(t, report.Code)
			assert.NotZero(t, report.FightID)
		}
	})

	t.Run("check duplicate handling", func(t *testing.T) {
		// 1. Récupérer un exemple de report actif
		var sampleReport warcraftlogsBuilds.Report
		err := db.Where("deleted_at IS NULL").First(&sampleReport).Error
		require.NoError(t, err)

		// 2. Compter combien de fois ce report apparaît dans la base
		var duplicateCount int64
		err = db.Model(&warcraftlogsBuilds.Report{}).
			Where("code = ? AND fight_id = ?", sampleReport.Code, sampleReport.FightID).
			Count(&duplicateCount).Error
		require.NoError(t, err)

		t.Logf("Report %s (fight_id: %d) appears %d times in database",
			sampleReport.Code, sampleReport.FightID, duplicateCount)

		// 3. Vérifier l'historique des versions
		var versions []struct {
			DeletedAt time.Time
			UpdatedAt time.Time
		}
		err = db.Model(&warcraftlogsBuilds.Report{}).
			Where("code = ? AND fight_id = ?", sampleReport.Code, sampleReport.FightID).
			Order("updated_at DESC").
			Select("deleted_at, updated_at").
			Find(&versions).Error
		require.NoError(t, err)

		for i, v := range versions {
			t.Logf("Version %d - Updated: %v, Deleted: %v",
				i+1, v.UpdatedAt, v.DeletedAt)
		}
	})

	// Test 3: Check behavior at offset 2900 where issues were reported
	t.Run("check offset 2900 behavior", func(t *testing.T) {
		reports, err := repo.GetReportsBatch(ctx, 50, 2900)
		assert.NoError(t, err)

		t.Logf("Batch at offset 2900 contains %d reports", len(reports))
		if len(reports) == 0 {
			// Get total count to understand the empty result
			var count int64
			db.Model(&warcraftlogsBuilds.Report{}).
				Where("deleted_at IS NULL").
				Count(&count)
			t.Logf("Total reports: %d", count)
		}
	})

	// Test 4: Scan through data to find where batches become empty
	t.Run("scan for data boundaries", func(t *testing.T) {
		batchSize := 50
		lastNonEmptyOffset := 0

		for offset := 0; offset <= 3000; offset += batchSize {
			reports, err := repo.GetReportsBatch(ctx, batchSize, offset)
			assert.NoError(t, err)

			if len(reports) == 0 {
				t.Logf("Found empty batch at offset %d", offset)
				t.Logf("Last non-empty batch was at offset %d", lastNonEmptyOffset)
				break
			}

			lastNonEmptyOffset = offset
			t.Logf("Offset %d: found %d reports", offset, len(reports))
		}
	})

	// Test 5: Verify data consistency in reports
	t.Run("verify data consistency", func(t *testing.T) {
		reports, err := repo.GetReportsBatch(ctx, 50, 0)
		assert.NoError(t, err)

		if len(reports) > 0 {
			// Check first report details
			firstReport := reports[0]
			t.Logf("First report details: code=%s, fight_id=%d",
				firstReport.Code,
				firstReport.FightID)

			// Verify required fields
			assert.NotEmpty(t, firstReport.Code)
			assert.NotZero(t, firstReport.FightID)
			// Add more field verifications as needed
		}
	})
}
