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
	wg          sync.WaitGroup
}

// NewWorkerPool creates a new WorkerPool
func NewWorkerPool(client *WarcraftLogsClientService, numWorkers int) *WorkerPool {
	return &WorkerPool{
		client:      client,
		numWorkers:  numWorkers,
		jobsChan:    make(chan Job, numWorkers*2),
		resultsChan: make(chan JobResult, numWorkers*2),
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start(ctx context.Context) {
	log.Printf("Starting worker pool with %d workers", wp.numWorkers)
	wp.wg.Add(wp.numWorkers)

	for i := 0; i < wp.numWorkers; i++ {
		go wp.worker(ctx, i)
	}
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
		case job, ok := <-wp.jobsChan:
			if !ok {
				log.Printf("Worker %d stopping: jobs channel closed", id)
				return
			}

			result := JobResult{}
			data, err := wp.client.MakeRequest(ctx, job.Query, job.Variables)
			if err != nil {
				result.Error = fmt.Errorf("worker %d: request error: %w", id, err)
			} else {
				result.Data = data
			}

			select {
			case wp.resultsChan <- result:
			case <-ctx.Done():
				log.Printf("Worker %d stopping: context cancelled", id)
				return
			}
		}
	}
}

// Submit adds a job to the worker pool
func (wp *WorkerPool) Submit(job Job) {
	wp.jobsChan <- job
}

// Results returns a channel that receives the results of the jobs
func (wp *WorkerPool) Results() <-chan JobResult {
	return wp.resultsChan
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.jobsChan)
	wp.wg.Wait()
	close(wp.resultsChan)
	log.Printf("Worker pool stopped")
}
