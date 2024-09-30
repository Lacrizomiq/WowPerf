package cache

import (
	"log"
	"time"
)

type UpdateFunc func() error

func StartPeriodicUpdate(key string, updateFunc UpdateFunc, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			err := updateFunc()
			if err != nil {
				log.Printf("Error updating cache for key %s: %v", key, err)
			}
		}
	}()
}
