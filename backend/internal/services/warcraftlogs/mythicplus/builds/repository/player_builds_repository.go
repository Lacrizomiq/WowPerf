package warcraftlogsBuildsRepository

import (
	"context"
	"fmt"
	"log"
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"gorm.io/gorm"
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

// StoreManyPlayerBuilds persists a list of player builds to the database in batches
// Uses the transaction to ensure data consistency
func (r *PlayerBuildsRepository) StoreManyPlayerBuilds(ctx context.Context, playerBuilds []*warcraftlogsBuilds.PlayerBuild) error {
	if len(playerBuilds) == 0 {
		log.Printf("[DEBUG] No player builds to store")
		return nil // Return silently if no player builds to store instead of error for empty slice
	}

	reportCode := playerBuilds[0].ReportCode

	// Check if the report has already been processed
	if lastProcessed, exists := r.processedReports[reportCode]; exists {
		if time.Since(lastProcessed) < 7*24*time.Hour {
			log.Printf("[DEBUG] Skipping report %s as it was processed recently", reportCode)
			return nil
		}
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	// defer to rollback the transaction if it fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[ERROR] Recovered from panic in StoreManyPlayerBuilds: %v", r)
		}
	}()

	// Check if builds already exist for this report
	var count int64
	if err := tx.WithContext(ctx).
		Model(&warcraftlogsBuilds.PlayerBuild{}).
		Where("report_code = ?", reportCode).
		Count(&count).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to count existing builds: %w", err)
	}

	if count > 0 {
		// Update existing builds
		log.Printf("[DEBUG] Updating %d existing builds for report %s", count, reportCode)
		for _, build := range playerBuilds {
			if err := tx.WithContext(ctx).
				Model(&warcraftlogsBuilds.PlayerBuild{}).
				Where("report_code = ? AND actor_id = ?", build.ReportCode, build.ActorID).
				Updates(map[string]interface{}{
					"talent_tree": build.TalentTree,
					"gear":        build.Gear,
					"stats":       build.Stats,
				}).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update player build: %w", err)
			}
		}
	} else {
		// Create new builds
		log.Printf("[DEBUG] Creating %d new builds for report %s", len(playerBuilds), reportCode)
		if err := tx.WithContext(ctx).Create(playerBuilds).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create player builds: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit updates: %w", err)
	}

	action := "updated"
	if count == 0 {
		action = "created"
	}
	log.Printf("[DEBUG] Successfully %s %d builds for report %s", action, len(playerBuilds), reportCode)
	r.processedReports[reportCode] = time.Now()
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
