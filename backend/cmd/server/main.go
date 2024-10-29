package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	// Internal Packages - API Handlers
	authHandler "wowperf/internal/api/auth"
	apiBlizzard "wowperf/internal/api/blizzard"
	"wowperf/internal/api/raiderio"
	mythicPlusRaiderioCache "wowperf/internal/api/raiderio/mythicplus"
	raidsRaiderioCache "wowperf/internal/api/raiderio/raids"
	userHandler "wowperf/internal/api/user"
	apiWarcraftlogs "wowperf/internal/api/warcraftlogs"

	// Internal Packages - Services
	auth "wowperf/internal/services/auth"
	serviceBlizzard "wowperf/internal/services/blizzard"
	serviceRaiderio "wowperf/internal/services/raiderio"
	mythicplusUpdate "wowperf/internal/services/raiderio/mythicplus"
	userService "wowperf/internal/services/user"
	warcraftlogs "wowperf/internal/services/warcraftlogs"
	warcraftLogsRankings "wowperf/internal/services/warcraftlogs/dungeons"

	// Internal Packages - Models & Database
	"wowperf/internal/database"
	models "wowperf/internal/models/raiderio/mythicrundetails"
	playerRankingModels "wowperf/internal/models/warcraftlogs/mythicplus"

	// Internal Packages - Utils
	middlewares "wowperf/middlewares"
	"wowperf/pkg/cache"
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

func startCacheUpdater(blizzardService *serviceBlizzard.Service, rioService *serviceRaiderio.RaiderIOService) {
	raidsRaiderioCache.StartRaidLeaderboardCacheUpdater(rioService)
	mythicPlusRaiderioCache.StartMythicPlusBestRunsCacheUpdater(rioService)
}

// Redis initialization
func initializeCacheService() (cache.CacheService, error) {
	cacheService, err := cache.NewRedisCache(&cache.Config{
		URL:      os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"), // as an optionnal parameter
		DB:       0,                           // Use the default DB
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache service: %w", err)
	}

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test the connection
	if err := cacheService.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to cache: %w", err)
	}

	return cacheService, nil
}

func initializeServices(db *gorm.DB, cacheService cache.CacheService) (
	*auth.AuthService,
	*userService.UserService,
	*serviceBlizzard.Service,
	*serviceRaiderio.RaiderIOService,
	*warcraftLogsRankings.RankingsService,
	*warcraftLogsRankings.RankingsUpdater,
	*warcraftLogsRankings.DungeonService,
	error,
) {
	// Auth Service Setup
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("JWT_SECRET must be set in the environment")
	}

	jwtExpirationStr := os.Getenv("JWT_EXPIRATION")
	jwtExpiration := 24 * time.Hour
	if jwtExpirationStr != "" {
		if duration, err := time.ParseDuration(jwtExpirationStr); err == nil {
			jwtExpiration = duration
		}
	}

	authService := auth.NewAuthService(
		db,
		jwtSecret,
		jwtExpiration,
		os.Getenv("BATTLE_NET_CLIENT_ID"),
		os.Getenv("BATTLE_NET_CLIENT_SECRET"),
		os.Getenv("BATTLE_NET_REDIRECT_URL"),
		cacheService,
	)

	// User Service
	userSvc := userService.NewUserService(db)

	// Blizzard Service
	blizzardService, err := serviceBlizzard.NewService()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to initialize blizzard service: %v", err)
	}

	// Raider.io Service
	rioService, err := serviceRaiderio.NewRaiderIOService()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to initialize raiderio service: %v", err)
	}

	// WarcraftLogs Service Setup
	warcraftLogsClient, err := warcraftlogs.NewClient()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to initialize warcraftlogs client: %v", err)
	}

	dungeonService := warcraftLogsRankings.NewDungeonService(warcraftLogsClient)
	rankingsService := warcraftLogsRankings.NewRankingsService(dungeonService, db)
	rankingsUpdater := warcraftLogsRankings.NewRankingsUpdater(db, rankingsService, cacheService)

	return authService, userSvc, blizzardService, rioService, rankingsService, rankingsUpdater, dungeonService, nil
}

func setupRoutes(
	r *gin.Engine,
	authService *auth.AuthService,
	authHandler *authHandler.AuthHandler,
	userHandler *userHandler.UserHandler,
	rioHandler *raiderio.Handler,
	blizzardHandler *apiBlizzard.Handler,
	warcraftlogsHandler *apiWarcraftlogs.Handler,
	cacheManager *middlewares.CacheManager,
) {
	// CORS Configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))
	r.Use(gin.Logger())

	// Auth Routes
	// r.GET("/csrf-token", middlewares.CSRFToken())
	authHandler.RegisterRoutes(r)

	// API Routes
	rioHandler.RegisterRoutes(r)
	blizzardHandler.RegisterRoutes(r)
	warcraftlogsHandler.RegisterRoutes(r)

	// Protected Routes
	userHandler.RegisterRoutes(r, authService.AuthMiddleware())
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
		return
	}

	// Initialize Cache and Redis
	cacheService, err := initializeCacheService()
	if err != nil {
		log.Fatalf("Failed to initialize cache service: %v", err)
	}
	defer cacheService.Close()

	// Initialize Database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	if err := database.InitializeDatabase(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Initialize Services
	authService, userSvc, blizzardService, rioService, rankingsService, rankingsUpdater, dungeonService, err := initializeServices(db, cacheService)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// Initialize Handlers
	authHandler := authHandler.NewAuthHandler(authService)
	userHandler := userHandler.NewUserHandler(userSvc)
	rioHandler := raiderio.NewHandler(rioService, db)
	blizzardHandler := apiBlizzard.NewHandler(blizzardService, db, cacheService)
	warcraftlogsHandler := apiWarcraftlogs.NewHandler(rankingsService, dungeonService, db, cacheService)

	// Initialize Cache Manager
	cacheManager := middlewares.NewCacheManager(middlewares.CacheConfig{
		Cache:      cacheService,
		Expiration: 8 * time.Hour,
		KeyPrefix:  "warcraftlogs", // TODO: change this for all routes
	})

	// Start Cache Updater
	startCacheUpdater(blizzardService, rioService)

	// Start Periodic Updates (Mythic+ Dungeon Stats)
	go func() {
		log.Println("Setting up dungeon stats update...")
		if mythicplusUpdate.CheckAndSetUpdateLock(db) {
			log.Println("Performing initial dungeon stats update...")
			if err := mythicplusUpdate.UpdateDungeonStats(db, rioService); err != nil {
				log.Printf("Error during initial dungeon stats update: %v", err)
			} else {
				log.Println("Initial dungeon stats update completed")
			}
		} else {
			log.Println("Dungeon stats are up to date")
		}
		checkUpdateState(db)

		log.Println("Setting up weekly dungeon stats update...")
		mythicplusUpdate.StartWeeklyDungeonStatsUpdate(db, rioService)
	}()

	// Start Periodic Updates (WarcraftLogs Dungeon Rankings)
	go func() {
		log.Println("Setting up WarcraftLogs rankings update scheduler...")

		// Attendre que la base de données soit bien initialisée
		time.Sleep(10 * time.Second)

		// Vérifier l'état initial sans forcer de mise à jour
		var state playerRankingModels.RankingsUpdateState
		if err := db.First(&state).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Println("Creating initial rankings update state...")
				state = playerRankingModels.RankingsUpdateState{
					LastUpdateTime: time.Now(),
				}
				if err := db.Create(&state).Error; err != nil {
					log.Printf("Error creating initial state: %v", err)
				}
			}
		}

		// Start the update scheduler
		log.Printf("Starting rankings update scheduler (minimum interval: %v)", warcraftLogsRankings.MinimumUpdateInterval)
		rankingsUpdater.StartPeriodicUpdate()
	}()

	// Setup and Start Server
	r := gin.Default()
	setupRoutes(r, authService, authHandler, userHandler, rioHandler, blizzardHandler, warcraftlogsHandler, cacheManager)

	log.Println("Server is starting on :8080")
	log.Fatal(r.Run(":8080"))
}
