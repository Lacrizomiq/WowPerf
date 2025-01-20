package main

import (
	"context"
	"log"
	"os"
	"time"
	warcraftlogsBuildsConfig "wowperf/internal/services/warcraftlogs/mythicplus/builds/config"
	workflows "wowperf/internal/services/warcraftlogs/mythicplus/builds/temporal/workflows"

	"go.temporal.io/sdk/client"
)

const (
	defaultNamespace = "default"
	defaultTaskQueue = "warcraft-logs-sync"
)

func main() {
	log.Printf("[INFO] Starting sync workflow")

	// Initialize Temporal client
	temporalAddress := os.Getenv("TEMPORAL_ADDRESS")
	if temporalAddress == "" {
		temporalAddress = "localhost:7233"
	}

	temporalClient, err := client.Dial(client.Options{
		HostPort:  temporalAddress,
		Namespace: defaultNamespace,
	})
	if err != nil {
		log.Fatalf("[FATAL] Failed to create Temporal client: %v", err)
	}
	defer temporalClient.Close()

	// Load configuration
	cfg, err := warcraftlogsBuildsConfig.Load("configs/config_s1_tww.priest.yaml")
	if err != nil {
		log.Fatalf("[FATAL] Failed to load config: %v", err)
	}

	// Prepare workflow parameters
	workflowParams := workflows.WorkflowParams{
		Specs:    make([]workflows.ClassSpec, len(cfg.Specs)),
		Dungeons: make([]workflows.Dungeon, len(cfg.Dungeons)),
		BatchConfig: workflows.BatchConfig{
			Size:        cfg.Rankings.Batch.Size,
			MaxPages:    cfg.Rankings.Batch.MaxPages,
			RetryDelay:  cfg.Rankings.Batch.RetryDelay,
			MaxAttempts: cfg.Rankings.Batch.MaxAttempts,
		},
		Rankings: struct {
			MaxRankingsPerSpec int           `json:"max_rankings_per_spec"`
			UpdateInterval     time.Duration `json:"update_interval"`
		}{
			MaxRankingsPerSpec: cfg.Rankings.MaxRankingsPerSpec,
			UpdateInterval:     cfg.Rankings.UpdateInterval,
		},
	}

	// Convert specs and dungeons
	for i, spec := range cfg.Specs {
		workflowParams.Specs[i] = workflows.ClassSpec{
			ClassName: spec.ClassName,
			SpecName:  spec.SpecName,
		}
	}
	for i, dungeon := range cfg.Dungeons {
		workflowParams.Dungeons[i] = workflows.Dungeon{
			ID:          dungeon.ID,
			EncounterID: dungeon.EncounterID,
			Name:        dungeon.Name,
			Slug:        dungeon.Slug,
		}
	}

	// Start workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:                       "warcraft-logs-sync-" + time.Now().Format("20060102-150405"),
		TaskQueue:                defaultTaskQueue,
		WorkflowExecutionTimeout: 24 * time.Hour,
	}

	execution, err := temporalClient.ExecuteWorkflow(
		context.Background(),
		workflowOptions,
		workflows.SyncWorkflow,
		workflowParams,
	)
	if err != nil {
		log.Fatalf("[FATAL] Failed to start workflow: %v", err)
	}

	log.Printf("[INFO] Started workflow: ID=%s, RunID=%s", execution.GetID(), execution.GetRunID())

	// Wait for workflow completion
	var result workflows.WorkflowResult
	if err := execution.Get(context.Background(), &result); err != nil {
		log.Fatalf("[FATAL] Workflow failed: %v", err)
	}

	log.Printf("[INFO] Workflow completed successfully: Rankings=%d, Reports=%d, Builds=%d",
		result.RankingsProcessed,
		result.ReportsProcessed,
		result.BuildsProcessed)
}
