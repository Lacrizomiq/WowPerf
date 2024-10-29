package raiderio

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"wowperf/internal/services/raiderio"
	raiderioMythicPlus "wowperf/internal/services/raiderio/mythicplus"
	raiderioMythicPlusWrapper "wowperf/internal/wrapper/raiderio"

	"github.com/gin-gonic/gin"
)

type MythicPlusRunDetailsHandler struct {
	Service *raiderio.RaiderIOService
}

func NewMythicPlusRunDetailsHandler(service *raiderio.RaiderIOService) *MythicPlusRunDetailsHandler {
	return &MythicPlusRunDetailsHandler{
		Service: service,
	}
}

func (h *MythicPlusRunDetailsHandler) GetMythicPlusRunDetails(c *gin.Context) {
	start := time.Now()
	defer func() {
		log.Printf("Total request time: %v", time.Since(start))
	}()

	season := c.Query("season")
	idStr := c.Query("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}

	if season == "" || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	// Implement retry logic
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		apiStart := time.Now()
		rawRunDetails, err := raiderioMythicPlus.GetMythicPlusRunsDetails(h.Service, season, id)
		log.Printf("API call took: %v", time.Since(apiStart))

		if err == nil {
			transformStart := time.Now()
			transformedRunDetails, err := raiderioMythicPlusWrapper.TransformMythicPlusRun(rawRunDetails)
			if err != nil {
				log.Printf("Error transforming data: %v", err)
			}
			log.Printf("Data transformation took: %v", time.Since(transformStart))
			c.JSON(http.StatusOK, transformedRunDetails)
			return
		}

		log.Printf("Attempt %d failed: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * time.Second) // Simple backoff
	}

	// If all retries fail, return an error
	log.Printf("All retries failed")
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data after multiple attempts"})
}
