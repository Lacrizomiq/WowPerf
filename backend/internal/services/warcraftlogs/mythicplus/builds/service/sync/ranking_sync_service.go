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

// SyncService handle the complete sync process
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
	return &SyncService{
		batchProcessor: NewBatchProcessor(workerPool, NewSyncMetrics(), config.Rankings.Batch.MaxAttempts, config.Rankings.Batch.RetryDelay),
		repository:     repository,
		metrics:        NewSyncMetrics(),
		config:         config,
	}
}

// StartSync starts the sync process
func (s *SyncService) StartSync(ctx context.Context) error {
	log.Printf("[INFO] Starting rankings synchronization")

	// Create batches
	batches := s.createBatches()
	s.metrics.BatchesTotal = len(batches)

	// Create a channel to receive results
	results := make(chan *BatchResult, len(batches))
	errors := make(chan error, len(batches))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, s.config.Worker.NumWorkers) // semaphore to limit the number of concurrent workers

	// Process the batches
	for _, batch := range batches {
		wg.Add(1)
		go func(b RankingBatch) {
			defer wg.Done()
			semaphore <- struct{}{}        // acquire a slot
			defer func() { <-semaphore }() // release the slot

			result, err := s.processSingleBatch(ctx, b)
			if err != nil {
				errors <- fmt.Errorf("batch processing error for %s-%s, dungeon %d: %w", b.ClassName, b.SpecName, b.DungeonID, err)
			} else {
				results <- result
			}
		}(batch)
	}

	// Go routine to wait for all the batches to be processed and close the channels
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results and errors
	return s.handleResults(ctx, results, errors)
}

// createBatches creates the batches to be processed for each combination of class, spec, and dungeon
func (s *SyncService) createBatches() []RankingBatch {
	var batches []RankingBatch

	for _, spec := range s.config.Specs {
		for _, dungeon := range s.config.Dungeons {
			batches = append(batches, RankingBatch{
				ClassName:   spec.ClassName,
				SpecName:    spec.SpecName,
				DungeonID:   dungeon.ID,
				BatchSize:   s.config.Rankings.Batch.Size,
				CurrentPage: 1,
			})
		}
	}

	return batches
}

// processSingleBatch processes a single batch of rankings and handle pagination
func (s *SyncService) processSingleBatch(ctx context.Context, batch RankingBatch) (*BatchResult, error) {
	result, err := s.batchProcessor.ProcessBatch(ctx, batch)
	if err != nil {
		return nil, err
	}

	// Check if an update is needed
	lastRanking, err := s.repository.GetLastRankingForEncounter(ctx, batch.DungeonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last ranking for dungeon %d: %w", batch.DungeonID, err)
	}

	if lastRanking != nil {
		if time.Since(lastRanking.UpdatedAt) < s.config.Rankings.UpdateInterval {
			log.Printf("[INFO] Skipping batch, last update too recent for %s-%s, dungeon: %d",
				batch.ClassName, batch.SpecName, batch.DungeonID)
			return &BatchResult{Batch: batch}, nil
		}
	}

	// If no last ranking or update is needed, save the new rankings
	if len(result.Rankings) > 0 {
		if err := s.repository.StoreRankings(ctx, batch.DungeonID, result.Rankings); err != nil {
			return nil, fmt.Errorf("failed to store rankings for dungeon %d: %w", batch.DungeonID, err)
		}
	}

	return result, nil
}

// handleResults handles the results and errors from the batches
func (s *SyncService) handleResults(ctx context.Context, results chan *BatchResult, errors chan error) error {
	var processedBatches int
	var errs []error
	var totalRankings int
	var successfulBatches int

	// Collect errors
	for err := range errors {
		errs = append(errs, err)
		s.metrics.RecordError(err)
	}

	// Collect and process results
	for result := range results {
		processedBatches++
		if result.Rankings != nil && len(result.Rankings) > 0 {
			totalRankings += len(result.Rankings)
			successfulBatches++

			// Log detailed batch information
			log.Printf("[DEBUG] Batch processed: %s-%s, dungeon: %d, rankings: %d",
				result.Batch.ClassName,
				result.Batch.SpecName,
				result.Batch.DungeonID,
				len(result.Rankings))
		}

		log.Printf("[DEBUG] Progress: %d/%d batches completed", processedBatches, s.metrics.BatchesTotal)
	}

	s.metrics.Complete()

	// Log final statistics
	log.Printf("[INFO] Sync summary: %d/%d batches successful, total rankings: %d",
		successfulBatches,
		s.metrics.BatchesTotal,
		totalRankings)

	if len(errs) > 0 {
		return fmt.Errorf("failed to process %d batches: %v", len(errs), errs)
	}

	log.Printf("[INFO] Rankings synchronization completed successfully")
	return nil
}
