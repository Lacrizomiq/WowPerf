package warcraftlogs

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// Job represents a request to the Warcraft Logs API
type Job struct {
	Query     string
	Variables map[string]interface{}
	JobType   string
	Metadata  map[string]interface{}
}

// JobResult represents the result of a job
type JobResult struct {
	Data  []byte
	Job   Job
	Error error
}

// WorkerPool manages a pool of workers to handle Warcraft Logs API requests
type WorkerPool struct {
	client      *WarcraftLogsClientService
	numWorkers  int
	jobsChan    chan Job
	resultsChan chan JobResult
	stopChan    chan struct{} // Channel to signal the stop
	wg          sync.WaitGroup
	started     bool       // Flag to track if the pool is started
	mu          sync.Mutex // Mutex to protect the started flag
}

// NewWorkerPool creates a new WorkerPool
func NewWorkerPool(client *WarcraftLogsClientService, numWorkers int) *WorkerPool {
	return &WorkerPool{
		client:      client,
		numWorkers:  numWorkers,
		jobsChan:    make(chan Job, numWorkers*10),
		resultsChan: make(chan JobResult, numWorkers*10),
		stopChan:    make(chan struct{}),
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start(ctx context.Context) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.started {
		return nil // Already started
	}

	log.Printf("Starting worker pool with %d workers", wp.numWorkers)

	// Reset the channels if necessary
	if wp.jobsChan == nil {
		wp.jobsChan = make(chan Job, wp.numWorkers*10)
	}
	if wp.resultsChan == nil {
		wp.resultsChan = make(chan JobResult, wp.numWorkers*10)
	}
	if wp.stopChan == nil {
		wp.stopChan = make(chan struct{})
	}

	// Start the workers
	wp.wg.Add(wp.numWorkers)
	for i := 0; i < wp.numWorkers; i++ {
		go wp.worker(ctx, i)
	}

	wp.started = true
	return nil
}

// worker is a single worker in the pool that handles Warcraft Logs API requests
func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()
	log.Printf("Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopping: context cancelled", id)
			return
		case <-wp.stopChan:
			log.Printf("Worker %d stopping: stop signal received", id)
			return
		case job, ok := <-wp.jobsChan:
			if !ok {
				log.Printf("Worker %d stopping: jobs channel closed", id)
				return
			}

			log.Printf("Worker %d processing job type: %s", id, job.JobType)

			result := JobResult{
				Job: job,
			}

			// Process the job with context handling
			select {
			case <-ctx.Done():
				result.Error = ctx.Err()
			default:
				data, err := wp.client.MakeRequest(ctx, job.Query, job.Variables)
				result.Data = data
				if err != nil {
					result.Error = fmt.Errorf("worker %d: request error: %w", id, err)
				}
			}

			// Send the result with context handling
			select {
			case wp.resultsChan <- result:
				log.Printf("Worker %d completed job type: %s", id, job.JobType)
			case <-ctx.Done():
				log.Printf("Worker %d stopping: context cancelled while sending result", id)
				return
			case <-wp.stopChan:
				log.Printf("Worker %d stopping: stop signal received while sending result", id)
				return
			}
		}
	}
}

// Submit adds a job to the worker pool
func (wp *WorkerPool) Submit(job Job) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.jobsChan == nil || !wp.started {
		log.Printf("[WARNING] Attempting to submit job but worker pool not ready")
		return
	}

	select {
	case wp.jobsChan <- job:
		log.Printf("[DEBUG] Job submitted: %s", job.JobType)
	default:
		log.Printf("[WARNING] Job channel is full, waiting to submit job")
		wp.jobsChan <- job
	}
}

// Results returns a channel that receives the results of the jobs
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

	log.Printf("[DEBUG] Stopping worker pool")
	close(wp.stopChan)

	// Wait for all workers to finish
	wp.wg.Wait()

	// Close the channels
	if wp.jobsChan != nil {
		close(wp.jobsChan)
		wp.jobsChan = nil
	}
	if wp.resultsChan != nil {
		close(wp.resultsChan)
		wp.resultsChan = nil
	}

	wp.started = false
	log.Printf("[DEBUG] Worker pool stopped")
}

// IsStarted returns true if the worker pool is started
func (wp *WorkerPool) IsStarted() bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	return wp.started
}
