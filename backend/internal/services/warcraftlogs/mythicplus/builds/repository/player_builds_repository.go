package warcraftlogsBuildsRepository

import (
	"context"
	"fmt"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"

	"gorm.io/gorm"
)

// PlayerBuildsRepository handles database operations for player builds
type PlayerBuildsRepository struct {
	db *gorm.DB
}

// NewPlayerBuildsRepository creates a new instance of PlayerBuildsRepository
func NewPlayerBuildsRepository(db *gorm.DB) *PlayerBuildsRepository {
	return &PlayerBuildsRepository{
		db: db,
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
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(playerBuilds).Error; err != nil {
			return fmt.Errorf("failed to store player builds: %w", err)
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
