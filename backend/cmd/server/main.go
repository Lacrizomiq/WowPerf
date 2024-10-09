package main

import (
	"log"
	"os"
	"time"
	apiBlizzard "wowperf/internal/api/blizzard"
	"wowperf/internal/api/raiderio"
	"wowperf/internal/database"
	"wowperf/pkg/cache"
	"wowperf/pkg/middleware"

	apiAuth "wowperf/internal/api/auth"
	authService "wowperf/internal/services/auth"

	mythicPlusRaiderioCache "wowperf/internal/api/raiderio/mythicplus"
	raidsRaiderioCache "wowperf/internal/api/raiderio/raids"
	raiderioMythicPlus "wowperf/internal/services/raiderio/mythicplus"

	models "wowperf/internal/models/raiderio/mythicrundetails"

	serviceBlizzard "wowperf/internal/services/blizzard"
	serviceRaiderio "wowperf/internal/services/raiderio"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func checkUpdateState(db *gorm.DB) {
	var state models.UpdateState
	result := db.First(&state)
	if result.Error != nil {
		log.Printf("Error fetching update state: %v", result.Error)
		return
	}
	log.Printf("Current update state: Last update was %v ago", time.Since(state.LastUpdateTime))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return
	}

	// Init Cache
	cache.InitCache()

	// Wait for Redis to be ready
	waitForRedis()

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Use the new migration function
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	if err := database.SeedDatabase(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Services

	// Auth Service
	accessExpiry := 15 * time.Minute
	refreshExpiry := 7 * 24 * time.Hour
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")
	if accessSecret == "" || refreshSecret == "" {
		log.Fatal("ACCESS_SECRET and REFRESH_SECRET must be set in the environment")
	}
	authService := authService.NewAuthService(db, accessSecret, refreshSecret, accessExpiry, refreshExpiry)

	// Blizzard Service
	blizzardService, err := serviceBlizzard.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize blizzard service: %v", err)
	}

	// Raider.io Service
	rioService, err := serviceRaiderio.NewRaiderIOService()
	if err != nil {
		log.Fatalf("Failed to initialize raiderio service: %v", err)
	}

	// Cache Updater
	startCacheUpdater(blizzardService, rioService)

	// Handlers
	// Auth Handler
	authHandler := apiAuth.NewAuthHandler(authService)

	// Raider.io Handler
	rioHandler := raiderio.NewHandler(rioService, db)

	// Blizzard Handler
	blizzardHandler := apiBlizzard.NewHandler(blizzardService, db)

	r := gin.Default()

	// CSRF Protection
	r.Use(middleware.CSRF())

	// CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Todo : Replace with the actual frontend URL like https://wowperf.com when deployed
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))

	// API

	// Auth API
	authHandler.RegisterRoutes(r)

	// Raider.io API
	rioHandler.RegisterRoutes(r)

	// Blizzard API
	blizzardHandler.RegisterRoutes(r)

	checkUpdateState(db)

	go func() {
		log.Println("Checking if dungeon stats update is needed...")
		if err := raiderioMythicPlus.UpdateDungeonStats(db, rioService); err != nil {
			log.Printf("Error during dungeon stats update check: %v", err)
		} else {
			log.Println("Dungeon stats update check completed")
		}
		checkUpdateState(db)

		log.Println("Setting up weekly dungeon stats update...")
		raiderioMythicPlus.StartWeeklyDungeonStatsUpdate(db, rioService)
	}()

	// Protected Routes
	protected := r.Group("/user")
	protected.Use(middleware.JWTAuth(authService))
	{
		// Add your protected routes here
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{"message": "Profile accessed", "user_id": userID})
		})
	}

	log.Println("Server is starting on :8080")
	log.Fatal(r.Run(":8080"))
}

func startCacheUpdater(blizzardService *serviceBlizzard.Service, rioService *serviceRaiderio.RaiderIOService) {
	raidsRaiderioCache.StartRaidLeaderboardCacheUpdater(rioService)
	mythicPlusRaiderioCache.StartMythicPlusBestRunsCacheUpdater(rioService)
}

func waitForRedis() {
	for i := 0; i < 30; i++ {
		err := cache.Ping()
		if err == nil {
			log.Println("Redis is ready")
			return
		}
		log.Println("Waiting for Redis to be ready")
		time.Sleep(time.Second)
	}
	log.Println("Redis is not ready")
}
