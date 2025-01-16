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

// StoreManyPlayerBuilds persists a list of player builds to the database
func (r *PlayerBuildsRepository) StoreManyPlayerBuilds(ctx context.Context, newBuilds []*warcraftlogsBuilds.PlayerBuild) error {
	if len(newBuilds) == 0 {
		log.Printf("[DEBUG] No player builds to store")
		return nil
	}

	encounterID := newBuilds[0].EncounterID
	reportCode := newBuilds[0].ReportCode

	if reportCode == "" {
		return fmt.Errorf("invalid report code: empty")
	}

	// Check if the report has already been processed recently
	if lastProcessed, exists := r.processedReports[reportCode]; exists {
		if time.Since(lastProcessed) < 7*24*time.Hour {
			log.Printf("[DEBUG] Skipping report %s as it was processed recently", reportCode)
			return nil
		}
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		log.Printf("[DEBUG] Starting transaction for encounter %d with %d builds", encounterID, len(newBuilds))

		// Récupérer les builds existants
		var existingBuilds []*warcraftlogsBuilds.PlayerBuild
		if err := tx.WithContext(ctx).
			Where("encounter_id = ?", encounterID).
			Find(&existingBuilds).Error; err != nil {
			return fmt.Errorf("failed to fetch existing builds: %w", err)
		}

		// On ne procède à la mise à jour que si on a de nouveaux builds
		if len(newBuilds) > 0 {
			// Créer une map des builds existants
			existingMap := make(map[string]*warcraftlogsBuilds.PlayerBuild)
			for _, build := range existingBuilds {
				key := fmt.Sprintf("%s_%d_%d", build.ReportCode, build.FightID, build.ActorID)
				existingMap[key] = build
			}

			// Traiter les nouveaux builds
			for _, newBuild := range newBuilds {
				if newBuild.ReportCode == "" {
					log.Printf("[WARN] Skipping build with empty report_code for player %s (%s-%s)",
						newBuild.PlayerName, newBuild.Class, newBuild.Spec)
					continue
				}

				key := fmt.Sprintf("%s_%d_%d", newBuild.ReportCode, newBuild.FightID, newBuild.ActorID)
				delete(existingMap, key)

				// Update or create the build
				result := tx.Clauses(clause.OnConflict{
					Columns: []clause.Column{
						{Name: "report_code"},
						{Name: "fight_id"},
						{Name: "actor_id"},
					},
					DoUpdates: clause.AssignmentColumns([]string{
						"talent_tree",
						"gear",
						"stats",
						"updated_at",
					}),
				}).Create(newBuild)

				if result.Error != nil {
					return fmt.Errorf("failed to store player build: %w", result.Error)
				}

				log.Printf("[DEBUG] Stored build for player %s in report %s", newBuild.PlayerName, newBuild.ReportCode)
			}

			// Supprimer uniquement les builds qui ne sont plus dans la nouvelle liste
			// et qui correspondent au même encounter_id
			for _, oldBuild := range existingMap {
				if err := tx.Delete(oldBuild).Error; err != nil {
					return fmt.Errorf("failed to delete old build: %w", err)
				}
				log.Printf("[DEBUG] Deleted obsolete build for player %s in report %s",
					oldBuild.PlayerName, oldBuild.ReportCode)
			}

			// Mark the report as processed
			r.processedReports[reportCode] = time.Now()

			log.Printf("[INFO] Successfully processed builds for encounter %d: stored %d, deleted %d",
				encounterID, len(newBuilds), len(existingMap))
		} else {
			log.Printf("[DEBUG] No new builds to process, keeping existing builds for encounter %d", encounterID)
		}

		return nil
	})
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
