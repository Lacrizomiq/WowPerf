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

// PlayerBuildsRepository handles all database operations for player builds.
type PlayerBuildsRepository struct {
	db *gorm.DB
}

// NewPlayerBuildsRepository creates a new instance of PlayerBuildsRepository.
func NewPlayerBuildsRepository(db *gorm.DB) *PlayerBuildsRepository {
	return &PlayerBuildsRepository{
		db: db,
	}
}

// StoreManyPlayerBuilds persists multiple player builds to the database.
// It handles batching to avoid memory issues and uses UPSERT for conflict resolution.
// If a build already exists (same report_code, fight_id, actor_id), it will be updated.
func (r *PlayerBuildsRepository) StoreManyPlayerBuilds(ctx context.Context, builds []*warcraftlogsBuilds.PlayerBuild) error {
	if len(builds) == 0 {
		log.Printf("[DEBUG] No player builds to store")
		return nil
	}

	// Process by batch of 5 builds at a time
	const batchSize = 5
	for i := 0; i < len(builds); i += batchSize {
		end := i + batchSize
		if end > len(builds) {
			end = len(builds)
		}

		batch := builds[i:end]
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

				// Create/update build with UPSERT
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
			return fmt.Errorf("failed to store builds batch %d-%d: %w", i, end, err)
		}

		log.Printf("[DEBUG] Processed batch %d/%d, total builds processed: %d",
			(i/batchSize)+1, (len(builds)+batchSize-1)/batchSize, end)

		// Small delay between batches to avoid overload
		time.Sleep(time.Millisecond * 100)
	}

	log.Printf("[INFO] Successfully processed all builds: stored %d builds", len(builds))
	return nil
}

// CountPlayerBuilds returns the total count of player builds in the database.
func (r *PlayerBuildsRepository) CountPlayerBuilds(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.Model(&warcraftlogsBuilds.PlayerBuild{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count player builds: %w", err)
	}
	return count, nil
}

// CountPlayerBuildsByFilter count the builds for a specific class, spec and encounter_id
func (r *PlayerBuildsRepository) CountPlayerBuildsByFilter(
	ctx context.Context,
	class, spec string,
	encounterID uint,
) (int64, error) {
	var count int64
	query := r.db.Model(&warcraftlogsBuilds.PlayerBuild{})

	if class != "" {
		query = query.Where("class = ?", class)
	}

	if spec != "" {
		query = query.Where("spec = ?", spec)
	}

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count player builds: %w", err)
	}

	return count, nil
}

// GetPlayerBuildsByFilter get the builds for a specific class, spec and encounter_id with pagination
func (r *PlayerBuildsRepository) GetPlayerBuildsByFilter(
	ctx context.Context,
	class, spec string,
	encounterID uint,
	limit, offset int,
) ([]*warcraftlogsBuilds.PlayerBuild, error) {
	var builds []*warcraftlogsBuilds.PlayerBuild
	query := r.db.WithContext(ctx)

	if class != "" {
		query = query.Where("class = ?", class)
	}

	if spec != "" {
		query = query.Where("spec = ?", spec)
	}

	if encounterID > 0 {
		query = query.Where("encounter_id = ?", encounterID)
	}

	if err := query.Limit(limit).Offset(offset).Find(&builds).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch player builds: %w", err)
	}

	return builds, nil
}
