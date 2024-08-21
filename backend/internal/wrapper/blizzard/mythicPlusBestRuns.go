package wrapper

import (
	"fmt"
	"sync"
	"wowperf/internal/models"

	"gorm.io/gorm"
)

func TransformMythicPlusBestRuns(data map[string]interface{}, db *gorm.DB) ([]models.MythicPlusRun, error) {
	bestRuns, ok := data["best_runs"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("best runs not found or not a slice")
	}

	var wg sync.WaitGroup
	runChan := make(chan models.MythicPlusRun, len(bestRuns))
	errChan := make(chan error, len(bestRuns))

	for _, run := range bestRuns {
		wg.Add(1)
		go func(run interface{}) {
			defer wg.Done()
			mythicRun, err := processMythicPlusRun(run, db)
			if err != nil {
				errChan <- err
				return
			}
			runChan <- mythicRun
		}(run)
	}

	go func() {
		wg.Wait()
		close(runChan)
		close(errChan)
	}()

	var results []models.MythicPlusRun
	for run := range runChan {
		results = append(results, run)
	}

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return results, nil
}
