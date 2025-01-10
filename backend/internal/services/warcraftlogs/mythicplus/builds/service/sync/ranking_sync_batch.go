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
	warcraftlogsBuildsMetrics "wowperf/internal/services/warcraftlogs/mythicplus/builds/service/sync/metrics"
	warcraftlogsTypes "wowperf/internal/services/warcraftlogs/types"
)

// RankingBatch represents a batch of rankings to process
type RankingBatch struct {
	ClassName   string // ex: "Priest"
	SpecName    string // ex: "Discipline"
	EncounterID uint   // ex: 12660
	BatchSize   uint   // ex: 100
	CurrentPage uint   // ex: 1
}

// BatchProcessor manages the processing of ranking batches
type BatchProcessor struct {
	workerPool *warcraftlogs.WorkerPool               // worker pool to fetch rankings
	metrics    *warcraftlogsBuildsMetrics.SyncMetrics // metrics to update
	config     *warcraftlogsBuildsConfig.Config
}

// BatchResult represents the result of processing a batch
type BatchResult struct {
	Batch    RankingBatch                       // batch to process
	Rankings []*warcraftlogsBuilds.ClassRanking // fetched rankings
	HasMore  bool                               // true if there are more pages to fetch
	Error    error                              // error that occurred during processing
	RetryAt  time.Time                          // time at which the batch should be retried
	Skipped  bool                               // true if the batch has been skipped
}

// NewBatchProcessor creates a new BatchProcessor instance
func NewBatchProcessor(
	workerPool *warcraftlogs.WorkerPool,
	metrics *warcraftlogsBuildsMetrics.SyncMetrics,
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
		Batch: batch,
	}

	log.Printf("[DEBUG] Processing batch for %s-%s, encounter: %d, page: %d",
		batch.ClassName, batch.SpecName, batch.EncounterID, batch.CurrentPage)

	var lastErr error
	for attempt := 1; attempt <= p.config.Rankings.Batch.MaxAttempts; attempt++ {
		result.Rankings, result.HasMore, lastErr = p.fetchRankings(ctx, batch)
		if lastErr == nil {
			break
		}

		// Gestion spécifique selon le type d'erreur
		if wlErr, ok := lastErr.(*warcraftlogsTypes.WarcraftLogsError); ok {
			switch wlErr.Type {
			case warcraftlogsTypes.ErrorTypeRateLimit, warcraftlogsTypes.ErrorTypeQuotaExceeded:
				// En cas de rate limit, on stocke quand réessayer
				if info := warcraftlogsTypes.GetRateLimitInfo(wlErr); info != nil {
					result.RetryAt = time.Now().Add(info.ResetIn)
					log.Printf("[INFO] Rate limit reached for %s-%s, retry after: %v",
						batch.ClassName, batch.SpecName, result.RetryAt)
					p.metrics.RecordRateLimit(info)
					return result, wlErr
				}
			case warcraftlogsTypes.ErrorTypeAPI:
				if !wlErr.Retryable {
					log.Printf("[ERROR] Non-retryable API error for %s-%s: %v",
						batch.ClassName, batch.SpecName, wlErr)
					return result, wlErr
				}
			}
		}

		if attempt == p.config.Rankings.Batch.MaxAttempts {
			log.Printf("[ERROR] Failed to process batch for %s-%s after %d attempts, last error: %v",
				batch.ClassName, batch.SpecName, attempt, lastErr)
			return result, fmt.Errorf("max retry attempts reached: %w", lastErr)
		}

		retryDelay := warcraftlogsTypes.GetRetryDelay(lastErr, attempt)
		log.Printf("[WARN] Attempt %d failed for %s-%s, retrying in %v...",
			attempt, batch.ClassName, batch.SpecName, retryDelay)

		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(retryDelay):
			continue
		}
	}

	p.updateMetrics(batch, len(result.Rankings))
	return result, nil
}

// fetchRankings fetches rankings from the API and returns the rankings, if there are more pages, and an error if it occurs
func (p *BatchProcessor) fetchRankings(ctx context.Context, batch RankingBatch) ([]*warcraftlogsBuilds.ClassRanking, bool, error) {
	job := warcraftlogs.Job{
		Query: rankingsQueries.ClassRankingsQuery,
		Variables: map[string]interface{}{
			"encounterId": batch.EncounterID,
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
			if wlErr, ok := result.Error.(*warcraftlogsTypes.WarcraftLogsError); ok {
				wlErr.Message = fmt.Sprintf("%s (class: %s, spec: %s, encounter: %d, page: %d)",
					wlErr.Message, batch.ClassName, batch.SpecName, batch.EncounterID, batch.CurrentPage)
			}
			return nil, false, result.Error
		}
		return rankingsQueries.ParseRankingsResponse(result.Data, batch.EncounterID)
	case <-ctx.Done():
		return nil, false, ctx.Err()
	}
}

// updateMetrics updates the metrics for a specific batch
func (p *BatchProcessor) updateMetrics(batch RankingBatch, rankingsCount int) {
	p.metrics.RecordBatchProcessed(batch.SpecName, batch.EncounterID)
	// Update metrics for this specific batch
	p.metrics.RecordRankingChanges(
		rankingsCount, // total
		rankingsCount, // new
		0,             // updated
		0,             // deleted
		0,             // unchanged
	)
}
