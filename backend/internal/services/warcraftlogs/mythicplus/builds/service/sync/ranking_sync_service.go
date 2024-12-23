package warcraftlogsBuildsSync

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"wowperf/internal/services/warcraftlogs"
	warcraftlogsBuildsConfig "wowperf/internal/services/warcraftlogs/mythicplus/builds/config"
	rankingsRepository "wowperf/internal/services/warcraftlogs/mythicplus/builds/repository"
)

// SyncService handles the complete sync process for WarcraftLogs rankings
type SyncService struct {
	batchProcessor *BatchProcessor
	repository     *rankingsRepository.RankingsRepository
	metrics        *SyncMetrics
	config         *warcraftlogsBuildsConfig.Config
}

// NewSyncService creates a new SyncService instance
func NewSyncService(
	workerPool *warcraftlogs.WorkerPool,
	repository *rankingsRepository.RankingsRepository,
	config *warcraftlogsBuildsConfig.Config,
) *SyncService {
	metrics := NewSyncMetrics()
	return &SyncService{
		batchProcessor: NewBatchProcessor(
			workerPool,
			metrics,
			config,
		),
		repository: repository,
		metrics:    metrics,
		config:     config,
	}
}

// StartSync initiates the synchronization process
func (s *SyncService) StartSync(ctx context.Context) error {
	log.Printf("[INFO] Starting rankings synchronization with %d workers", s.config.Worker.NumWorkers)
	s.metrics.StartTime = time.Now()

	batches := s.createBatches()
	s.metrics.BatchesTotal = len(batches)

	results := make(chan *BatchResult, len(batches))
	errors := make(chan error, len(batches))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, s.config.Worker.NumWorkers)

	// Process batches with worker pool
	for _, batch := range batches {
		wg.Add(1)
		go func(b RankingBatch) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			result, err := s.processSingleBatch(ctx, b)
			if err != nil {
				log.Printf("[ERROR] Batch processing error for %s-%s, dungeon %d: %v",
					b.ClassName, b.SpecName, b.DungeonID, err)
				errors <- fmt.Errorf("batch processing error for %s-%s, dungeon %d: %w",
					b.ClassName, b.SpecName, b.DungeonID, err)
				s.metrics.RecordError(err)
				return
			}

			results <- result
			s.metrics.RecordBatchProcessed(b.SpecName, b.DungeonID)

			// Respect API rate limiting
			time.Sleep(s.config.Worker.RequestDelay)
		}(batch)
	}

	// Wait for completion
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	return s.handleResults(ctx, results, errors)
}

// createBatches creates batches based on configuration
func (s *SyncService) createBatches() []RankingBatch {
	var batches []RankingBatch
	for _, spec := range s.config.Specs {
		for _, dungeon := range s.config.Dungeons {
			batches = append(batches, RankingBatch{
				ClassName:   spec.ClassName,
				SpecName:    spec.SpecName,
				DungeonID:   dungeon.EncounterID,
				BatchSize:   s.config.Rankings.Batch.Size,
				CurrentPage: 1,
			})
			log.Printf("[DEBUG] Created batch for %s-%s, dungeon: %d",
				spec.ClassName, spec.SpecName, dungeon.EncounterID)
		}
	}
	return batches
}

// processSingleBatch processes an individual batch
func (s *SyncService) processSingleBatch(ctx context.Context, batch RankingBatch) (*BatchResult, error) {
	log.Printf("[DEBUG] Processing batch for %s-%s, dungeon: %d, page: %d",
		batch.ClassName, batch.SpecName, batch.DungeonID, batch.CurrentPage)

	lastRanking, err := s.repository.GetLastRankingForEncounter(ctx, batch.DungeonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last ranking: %w", err)
	}

	if lastRanking != nil {
		timeSinceLastUpdate := time.Since(lastRanking.UpdatedAt)
		if timeSinceLastUpdate < s.config.Rankings.UpdateInterval {
			log.Printf("[INFO] Skipping update for %s-%s, dungeon: %d. Last update was %v ago",
				batch.ClassName, batch.SpecName, batch.DungeonID, timeSinceLastUpdate)
			return &BatchResult{Batch: batch}, nil
		}
	}

	result, err := s.batchProcessor.ProcessBatch(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("batch processor error: %w", err)
	}

	if len(result.Rankings) > 0 {
		if len(result.Rankings) > s.config.Rankings.MaxRankingsPerSpec {
			log.Printf("[WARN] Truncating rankings for %s-%s from %d to %d",
				batch.ClassName, batch.SpecName, len(result.Rankings), s.config.Rankings.MaxRankingsPerSpec)
			result.Rankings = result.Rankings[:s.config.Rankings.MaxRankingsPerSpec]
		}

		if err := s.repository.StoreRankings(ctx, batch.DungeonID, result.Rankings); err != nil {
			return nil, fmt.Errorf("failed to store rankings: %w", err)
		}
	}

	return result, nil
}

// handleResults processes the results of all batches
func (s *SyncService) handleResults(ctx context.Context, results chan *BatchResult, errors chan error) error {
	var processedBatches int
	var errs []error
	var totalRankings int
	var successfulBatches int

	// Collect errors
	for err := range errors {
		errs = append(errs, err)
	}

	// Process results
	for result := range results {
		processedBatches++
		if result.Rankings != nil {
			rankingsCount := len(result.Rankings)
			totalRankings += rankingsCount
			successfulBatches++

			s.metrics.RecordRankingChanges(
				rankingsCount, // total
				rankingsCount, // new
				0,             // updated
				0,             // deleted
				0,             // unchanged
			)

			log.Printf("[DEBUG] Processed batch for %s-%s, dungeon: %d, rankings: %d",
				result.Batch.ClassName,
				result.Batch.SpecName,
				result.Batch.DungeonID,
				rankingsCount)
		}

		log.Printf("[DEBUG] Progress: %d/%d batches completed", processedBatches, s.metrics.BatchesTotal)
	}

	// Complete metrics and log summary
	s.metrics.Complete()
	summary := s.metrics.GetSummary()

	log.Printf("[INFO] Sync completed: %d/%d batches successful, total rankings: %d, duration: %s",
		successfulBatches,
		s.metrics.BatchesTotal,
		totalRankings,
		summary["duration"])

	if len(errs) > 0 {
		return fmt.Errorf("failed to process %d batches: %v", len(errs), errs)
	}

	return nil
}
