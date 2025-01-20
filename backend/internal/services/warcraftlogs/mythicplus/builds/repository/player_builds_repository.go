package warcraftlogsBuildsRepository

import (
	"context"
	"fmt"
	"log"
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PlayerBuildsRepository handles database operations for player builds
type PlayerBuildsRepository struct {
	db               *gorm.DB
	processedReports map[string]time.Time
}

// NewPlayerBuildsRepository creates a new instance of PlayerBuildsRepository
func NewPlayerBuildsRepository(db *gorm.DB) *PlayerBuildsRepository {
	return &PlayerBuildsRepository{
		db:               db,
		processedReports: make(map[string]time.Time),
	}
}

// StorePlayerBuilds stores a list of player builds to the database
// Returns an error if the operation fails
func (r *PlayerBuildsRepository) StorePlayerBuilds(ctx context.Context, playerBuilds []*warcraftlogsBuilds.PlayerBuild) error {
	if result := r.db.WithContext(ctx).Create(playerBuilds); result.Error != nil {
		return fmt.Errorf("failed to store player builds: %w", result.Error)
	}

	return nil
}

// StoreManyPlayerBuilds persists multiple player builds to the database
// It handles batching to avoid memory issues and uses UPSERT for conflict resolution
// If a build already exists (same report_code, fight_id, actor_id), it will be updated
func (r *PlayerBuildsRepository) StoreManyPlayerBuilds(ctx context.Context, newBuilds []*warcraftlogsBuilds.PlayerBuild) error {
	if len(newBuilds) == 0 {
		log.Printf("[DEBUG] No player builds to store")
		return nil
	}

	// Get the encounter ID from the first build
	encounterID := newBuilds[0].EncounterID

	// Log start of transaction
	log.Printf("[DEBUG] Starting transaction for encounter %d with %d builds",
		encounterID, len(newBuilds))

	// Process builds in batches to avoid memory issues
	// Smaller batch size to prevent Temporal timeouts and memory issues
	const batchSize = 10
	processedBuilds := 0

	for i := 0; i < len(newBuilds); i += batchSize {
		end := i + batchSize
		if end > len(newBuilds) {
			end = len(newBuilds)
		}

		batch := newBuilds[i:end]

		// Process each batch in its own transaction
		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			for _, build := range batch {
				// Validate essential fields
				if build.ReportCode == "" {
					log.Printf("[WARN] Skipping build with empty report_code for player %s (%s-%s)",
						build.PlayerName, build.Class, build.Spec)
					continue
				}

				// Set timestamps
				now := time.Now()
				build.CreatedAt = now
				build.UpdatedAt = now

				// Attempt to create/update the build using UPSERT
				result := tx.Clauses(clause.OnConflict{
					Columns: []clause.Column{
						{Name: "report_code"},
						{Name: "fight_id"},
						{Name: "actor_id"},
					},
					DoUpdates: clause.AssignmentColumns([]string{
						"player_name",
						"class",
						"spec",
						"talent_import",
						"talent_tree",
						"item_level",
						"gear",
						"stats",
						"encounter_id",
						"keystone_level",
						"affixes",
						"updated_at",
					}),
				}).Create(build)

				if result.Error != nil {
					return fmt.Errorf("failed to store build for player %s in report %s: %w",
						build.PlayerName, build.ReportCode, result.Error)
				}

				log.Printf("[TRACE] Stored build for player %s in report %s",
					build.PlayerName, build.ReportCode)
			}
			return nil
		})

		if err != nil {
			log.Printf("[ERROR] Failed to store builds batch %d-%d: %v",
				i, end, err)
			return fmt.Errorf("failed to store builds batch: %w", err)
		}

		processedBuilds += len(batch)
		log.Printf("[DEBUG] Processed batch %d/%d, total builds processed: %d",
			(i/batchSize)+1, (len(newBuilds)+batchSize-1)/batchSize, processedBuilds)

		// Add a small delay between batches to prevent database overload
		time.Sleep(time.Millisecond * 100)
	}

	log.Printf("[INFO] Successfully processed all builds for encounter %d: stored %d builds",
		encounterID, processedBuilds)

	return nil
}

// GetByReportsAndActor retrieves a player build by its report code and actor ID
// Return the player build if found, nil otherwise
func (r *PlayerBuildsRepository) GetByReportsAndActor(ctx context.Context, reportCode string, actorID int) (*warcraftlogsBuilds.PlayerBuild, error) {
	var playerBuild warcraftlogsBuilds.PlayerBuild

	if err := r.db.WithContext(ctx).
		Where("report_code = ? AND actor_id = ?", reportCode, actorID).
		First(&playerBuild).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get player build: %w", err)
	}

	return &playerBuild, nil
}

// GetExistingBuilds retrieves all existing builds for a given report
// Used to avoid duplicates when processing reports
func (r *PlayerBuildsRepository) GetExistingBuilds(ctx context.Context, reportCode string) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	var playerBuilds []*warcraftlogsBuilds.PlayerBuild

	if err := r.db.WithContext(ctx).
		Where("report_code = ?", reportCode).
		Find(&playerBuilds).Error; err != nil {
		return nil, fmt.Errorf("failed to get existing builds: %w", err)
	}

	return playerBuilds, nil
}

// CountPlayerBuilds returns the total count of player builds
func (r *PlayerBuildsRepository) CountPlayerBuilds(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.Model(&warcraftlogsBuilds.PlayerBuild{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count player builds: %w", err)
	}
	return count, nil
}
