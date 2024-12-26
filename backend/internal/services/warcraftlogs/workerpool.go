package warcraftlogs

import (
	"context"
	"log"
	"sync"
	"time"

	warcraftlogsBuildsMetrics "wowperf/internal/services/warcraftlogs/mythicplus/builds/service/sync/metrics"
	warcraftlogsTypes "wowperf/internal/services/warcraftlogs/types"
)

// Job represents a request to the Warcraft Logs API
type Job struct {
	Query     string
	Variables map[string]interface{}
	JobType   string
	Metadata  map[string]interface{} // Metadata pour le debugging et les métriques
	Priority  int                    // Priorité du job (plus petit = plus prioritaire)
	CreatedAt time.Time              // Pour tracker les temps de traitement
}

// JobResult represents the result of a job
type JobResult struct {
	Data      []byte
	Job       Job
	Error     error
	StartedAt time.Time
	EndedAt   time.Time
}

// WorkerPool manages a pool of workers for WarcraftLogs API requests
type WorkerPool struct {
	client      *WarcraftLogsClientService
	numWorkers  int
	jobsChan    chan Job
	resultsChan chan JobResult
	stopChan    chan struct{}
	wg          sync.WaitGroup
	metrics     *warcraftlogsBuildsMetrics.SyncMetrics

	started bool
	mu      sync.Mutex

	// Nouvelles statistiques internes
	stats struct {
		activeWorkers  int
		completedJobs  int
		failedJobs     int
		avgProcessTime time.Duration
		mu             sync.RWMutex
	}
}

// NewWorkerPool creates a new WorkerPool
func NewWorkerPool(client *WarcraftLogsClientService, numWorkers int, metrics *warcraftlogsBuildsMetrics.SyncMetrics) *WorkerPool {
	return &WorkerPool{
		client:      client,
		numWorkers:  numWorkers,
		jobsChan:    make(chan Job, numWorkers*10),
		resultsChan: make(chan JobResult, numWorkers*10),
		stopChan:    make(chan struct{}),
		metrics:     metrics,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start(ctx context.Context) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.started {
		return nil
	}

	log.Printf("[DEBUG] Starting worker pool with %d workers", wp.numWorkers)

	if wp.jobsChan == nil {
		wp.jobsChan = make(chan Job, wp.numWorkers*10)
	}
	if wp.resultsChan == nil {
		wp.resultsChan = make(chan JobResult, wp.numWorkers*10)
	}
	if wp.stopChan == nil {
		wp.stopChan = make(chan struct{})
	}

	wp.wg.Add(wp.numWorkers)
	for i := 0; i < wp.numWorkers; i++ {
		go wp.worker(ctx, i)
	}

	wp.started = true
	return nil
}

// worker processes jobs from the pool
func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()

	wp.updateActiveWorkers(1)
	defer wp.updateActiveWorkers(-1)

	log.Printf("[DEBUG] Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[DEBUG] Worker %d stopping: context cancelled", id)
			return
		case <-wp.stopChan:
			log.Printf("[DEBUG] Worker %d stopping: stop signal received", id)
			return
		case job, ok := <-wp.jobsChan:
			if !ok {
				log.Printf("[DEBUG] Worker %d stopping: jobs channel closed", id)
				return
			}

			result := JobResult{
				Job:       job,
				StartedAt: time.Now(),
			}

			log.Printf("[DEBUG] Worker %d processing %s job: %v", id, job.JobType, job.Metadata)

			data, err := wp.client.MakeRequest(ctx, job.Query, job.Variables)
			result.EndedAt = time.Now()

			if err != nil {
				result.Error = err
				wp.recordFailedJob(result)
				if wlErr, ok := err.(*warcraftlogsTypes.WarcraftLogsError); ok && wlErr.Type == warcraftlogsTypes.ErrorTypeRateLimit {
					log.Printf("[WARN] Worker %d hit rate limit, backing off", id)
					select {
					case <-ctx.Done():
						return
					case <-time.After(wlErr.RetryIn):
						continue
					}
				}
			} else {
				result.Data = data
				wp.recordCompletedJob(result)
			}

			select {
			case wp.resultsChan <- result:
			case <-ctx.Done():
				return
			case <-wp.stopChan:
				return
			}
		}
	}
}

// Submit adds a job to the worker pool
func (wp *WorkerPool) Submit(job Job) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.started {
		log.Printf("[WARN] Attempting to submit job but worker pool not started")
		return
	}

	job.CreatedAt = time.Now()

	select {
	case wp.jobsChan <- job:
		log.Printf("[DEBUG] Job submitted: %s, metadata: %v", job.JobType, job.Metadata)
	default:
		log.Printf("[WARN] Job channel is full, waiting to submit job: %s", job.JobType)
		wp.jobsChan <- job
	}
}

// Results returns the channel for receiving job results
func (wp *WorkerPool) Results() <-chan JobResult {
	return wp.resultsChan
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.started {
		return
	}

	log.Printf("[DEBUG] Stopping worker pool...")
	close(wp.stopChan)

	wp.wg.Wait()

	if wp.jobsChan != nil {
		close(wp.jobsChan)
		wp.jobsChan = nil
	}
	if wp.resultsChan != nil {
		close(wp.resultsChan)
		wp.resultsChan = nil
	}

	wp.started = false
	wp.logFinalStats()
}

// Méthodes utilitaires pour les statistiques
func (wp *WorkerPool) updateActiveWorkers(delta int) {
	wp.stats.mu.Lock()
	defer wp.stats.mu.Unlock()
	wp.stats.activeWorkers += delta
}

func (wp *WorkerPool) recordCompletedJob(result JobResult) {
	wp.stats.mu.Lock()
	defer wp.stats.mu.Unlock()

	wp.stats.completedJobs++
	processingTime := result.EndedAt.Sub(result.StartedAt)

	// Mise à jour du temps moyen de traitement
	if wp.stats.completedJobs == 1 {
		wp.stats.avgProcessTime = processingTime
	} else {
		wp.stats.avgProcessTime = (wp.stats.avgProcessTime + processingTime) / 2
	}
}

func (wp *WorkerPool) recordFailedJob(result JobResult) {
	wp.stats.mu.Lock()
	defer wp.stats.mu.Unlock()
	wp.stats.failedJobs++
}

func (wp *WorkerPool) logFinalStats() {
	wp.stats.mu.RLock()
	defer wp.stats.mu.RUnlock()

	log.Printf("[INFO] Worker pool final stats: "+
		"Completed jobs: %d, Failed jobs: %d, Avg processing time: %v",
		wp.stats.completedJobs,
		wp.stats.failedJobs,
		wp.stats.avgProcessTime)
}

// GetStats returns current worker pool statistics
func (wp *WorkerPool) GetStats() map[string]interface{} {
	wp.stats.mu.RLock()
	defer wp.stats.mu.RUnlock()

	return map[string]interface{}{
		"active_workers":      wp.stats.activeWorkers,
		"completed_jobs":      wp.stats.completedJobs,
		"failed_jobs":         wp.stats.failedJobs,
		"avg_processing_time": wp.stats.avgProcessTime.String(),
	}
}
