package warcraftlogsBuildsSync

import (
	"context"
	"fmt"
	"log"
	"time"

	warcraftlogsBuilds "wowperf/internal/models/warcraftlogs/mythicplus/builds"
	"wowperf/internal/services/warcraftlogs"
	warcraftlogsBuildsConfig "wowperf/internal/services/warcraftlogs/mythicplus/builds/config"
	rankingsQueries "wowperf/internal/services/warcraftlogs/mythicplus/builds/queries"
)

// RankingBatch represents a batch of rankings to process
type RankingBatch struct {
	ClassName   string // ex: "Priest"
	SpecName    string // ex: "Discipline"
	DungeonID   uint   // ex: 123
	BatchSize   uint   // ex: 100
	CurrentPage uint   // ex: 1
}

// BatchProcessor manages the processing of ranking batches
type BatchProcessor struct {
	workerPool *warcraftlogs.WorkerPool // worker pool to fetch rankings
	metrics    *SyncMetrics             // metrics to update
	config     *warcraftlogsBuildsConfig.Config
}

// BatchResult represents the result of processing a batch
type BatchResult struct {
	Batch    RankingBatch                       // batch to process
	Rankings []*warcraftlogsBuilds.ClassRanking // fetched rankings
	HasMore  bool                               // true if there are more pages to fetch
	Error    error                              // error that occurred during processing
}

// NewBatchProcessor creates a new BatchProcessor instance
func NewBatchProcessor(
	workerPool *warcraftlogs.WorkerPool,
	metrics *SyncMetrics,
	config *warcraftlogsBuildsConfig.Config,
) *BatchProcessor {
	return &BatchProcessor{
		workerPool: workerPool,
		metrics:    metrics,
		config:     config,
	}
}

// ProcessBatch processes a batch and returns the fetched rankings
func (p *BatchProcessor) ProcessBatch(ctx context.Context, batch RankingBatch) (*BatchResult, error) {
	result := &BatchResult{
		Batch: batch, // batch to process
	}

	log.Printf("[DEBUG] Processing batch for %s-%s, dungeon: %d, page: %d",
		batch.ClassName, batch.SpecName, batch.DungeonID, batch.CurrentPage)

	var err error
	// Retry attempts with retry
	for attempt := 1; attempt <= p.config.Rankings.Batch.MaxAttempts; attempt++ {
		result.Rankings, result.HasMore, err = p.fetchRankings(ctx, batch)
		if err == nil {
			break
		}

		if attempt == p.config.Rankings.Batch.MaxAttempts {
			return result, fmt.Errorf("failed to process batch after %d attempts: %w", p.config.Rankings.Batch.MaxAttempts, err)
		}

		log.Printf("[WARN] Attempt %d failed, retrying in %v...", attempt, p.config.Rankings.Batch.RetryDelay)
		time.Sleep(p.config.Rankings.Batch.RetryDelay)
	}

	// Update metrics
	p.updateMetrics(batch, len(result.Rankings))

	return result, nil
}

// fetchRankings fetches rankings from the API and returns the rankings, if there are more pages, and an error if it occurs
func (p *BatchProcessor) fetchRankings(ctx context.Context, batch RankingBatch) ([]*warcraftlogsBuilds.ClassRanking, bool, error) {
	job := warcraftlogs.Job{
		Query: rankingsQueries.ClassRankingsQuery,
		Variables: map[string]interface{}{
			"encounterId": batch.DungeonID,
			"className":   batch.ClassName,
			"specName":    batch.SpecName,
			"page":        batch.CurrentPage,
		},
		JobType: "rankings", // type of job to fetch rankings
	}

	p.workerPool.Submit(job)

	select {
	case result := <-p.workerPool.Results():
		if result.Error != nil {
			return nil, false, fmt.Errorf("worker pool error: %w", result.Error)
		}
		return rankingsQueries.ParseRankingsResponse(result.Data, batch.DungeonID, batch.DungeonID)
	case <-ctx.Done():
		return nil, false, ctx.Err()
	}
}

// updateMetrics updates the metrics for a specific batch
func (p *BatchProcessor) updateMetrics(batch RankingBatch, rankingsCount int) {
	p.metrics.RecordBatchProcessed(batch.SpecName, batch.DungeonID)
	// Update metrics for this specific batch
	p.metrics.RecordRankingChanges(
		rankingsCount, // total
		rankingsCount, // new
		0,             // updated
		0,             // deleted
		0,             // unchanged
	)
}
