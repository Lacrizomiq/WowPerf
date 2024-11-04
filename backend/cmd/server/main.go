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
	userHandler "wowperf/internal/api/user"
	apiWarcraftlogs "wowperf/internal/api/warcraftlogs"

	// Internal Packages - Services
	auth "wowperf/internal/services/auth"
	serviceBlizzard "wowperf/internal/services/blizzard"
	serviceRaiderio "wowperf/internal/services/raiderio"
	mythicplusUpdate "wowperf/internal/services/raiderio/mythicplus"
	userService "wowperf/internal/services/user"
	warcraftlogs "wowperf/internal/services/warcraftlogs"
	warcraftLogsLeaderboard "wowperf/internal/services/warcraftlogs/dungeons"

	// Internal Packages - Models & Database
	"wowperf/internal/database"
	models "wowperf/internal/models/raiderio/mythicrundetails"

	// Internal Packages - Utils
	cacheMiddleware "wowperf/middleware/cache"
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

func initializeServices(db *gorm.DB, cacheService cache.CacheService, cacheManagers struct {
	raiderio     *cacheMiddleware.CacheManager
	blizzard     *cacheMiddleware.CacheManager
	warcraftlogs *cacheMiddleware.CacheManager
}) (
	*auth.AuthService,
	*auth.BlizzardAuthService,
	*userService.UserService,
	*serviceBlizzard.Service,
	*serviceRaiderio.RaiderIOService,
	*warcraftlogs.WarcraftLogsClientService,
	*warcraftLogsLeaderboard.GlobalLeaderboardService,
	*warcraftLogsLeaderboard.RankingsUpdater,
	error,
) {
	// Auth Service Setup
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("JWT_SECRET must be set in the environment")
	}

	jwtExpirationStr := os.Getenv("JWT_EXPIRATION")
	jwtExpiration := 24 * time.Hour
	if jwtExpirationStr != "" {
		if duration, err := time.ParseDuration(jwtExpirationStr); err == nil {
			jwtExpiration = duration
		}
	}

	redisClient := cacheService.GetRedisClient()

	// Blizzard Auth Service for OAuth2
	blizzardAuthService := auth.NewBlizzardAuthService(db, auth.BlizzardAuthConfig{
		ClientID:     os.Getenv("BLIZZARD_CLIENT_ID"),
		ClientSecret: os.Getenv("BLIZZARD_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("BLIZZARD_REDIRECT_URL"),
		Region:       "eu",
	})

	// Auth Service
	authService := auth.NewAuthService(
		db,
		jwtSecret,
		redisClient,
		jwtExpiration,
		blizzardAuthService,
	)

	// User Service
	userSvc := userService.NewUserService(db)

	// Blizzard Service
	blizzardService, err := serviceBlizzard.NewService()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to initialize blizzard service: %v", err)
	}

	// Raider.io Service
	rioService, err := serviceRaiderio.NewRaiderIOService()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to initialize raiderio service: %v", err)
	}

	// WarcraftLogs Service Setup
	warcraftLogsService, err := warcraftlogs.NewWarcraftLogsClientService()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to initialize warcraftlogs service: %v", err)
	}
	globalLeaderboardService := warcraftLogsLeaderboard.NewGlobalLeaderboardService(db)
	rankingsUpdater := warcraftLogsLeaderboard.NewRankingsUpdater(db, warcraftLogsService, cacheService, cacheManagers.warcraftlogs)
	return authService, blizzardAuthService, userSvc, blizzardService, rioService, warcraftLogsService, globalLeaderboardService, rankingsUpdater, nil
}

func setupRoutes(
	r *gin.Engine,
	authService *auth.AuthService,
	blizzardAuthService *auth.BlizzardAuthService,
	authHandler *authHandler.AuthHandler,
	userHandler *userHandler.UserHandler,
	blizzardAuthHandler *authHandler.BlizzardAuthHandler,
	rioHandler *raiderio.Handler,
	blizzardHandler *apiBlizzard.Handler,
	warcraftlogsHandler *apiWarcraftlogs.Handler,
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
	authHandler.RegisterRoutes(r)
	if blizzardAuthHandler != nil {
		blizzardAuthHandler.RegisterRoutes(r)
	}

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

	// Initialize Cache Manager with the cache service and the routes prefix
	cacheManagers := struct {
		raiderio     *cacheMiddleware.CacheManager
		blizzard     *cacheMiddleware.CacheManager
		warcraftlogs *cacheMiddleware.CacheManager
	}{
		raiderio: cacheMiddleware.NewCacheManager(cacheMiddleware.CacheConfig{
			Cache:      cacheService,
			Expiration: 24 * time.Hour,
			KeyPrefix:  "raiderio",
			Metrics:    true,
		}),
		blizzard: cacheMiddleware.NewCacheManager(cacheMiddleware.CacheConfig{
			Cache:      cacheService,
			Expiration: 24 * time.Hour,
			KeyPrefix:  "blizzard",
			Metrics:    true,
		}),
		warcraftlogs: cacheMiddleware.NewCacheManager(cacheMiddleware.CacheConfig{
			Cache:      cacheService,
			Expiration: 2 * time.Hour,
			KeyPrefix:  "warcraftlogs",
			Tags:       []string{"rankings", "leaderboard"},
			Metrics:    true,
		}),
	}

	// Initialize Services
	authService, blizzardAuthService, userSvc, blizzardService, rioService, warcraftLogsService, globalLeaderboardService, rankingsUpdater, err := initializeServices(db, cacheService, cacheManagers)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// Initialize Handlers
	authHandlers := authHandler.NewHandlers(authService, blizzardAuthService)
	userHandler := userHandler.NewUserHandler(userSvc)
	rioHandler := raiderio.NewHandler(rioService, db, cacheService, cacheManagers.raiderio)
	blizzardHandler := apiBlizzard.NewHandler(blizzardService, db, cacheService, cacheManagers.blizzard)
	warcraftlogsHandler := apiWarcraftlogs.NewHandler(globalLeaderboardService, warcraftLogsService, db, cacheService, cacheManagers.warcraftlogs)

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

		// Wait for the database to be ready
		time.Sleep(10 * time.Second)

		// Start the periodic updates
		rankingsUpdater.StartPeriodicUpdate(context.Background())
	}()

	// Setup and Start Server
	r := gin.Default()
	setupRoutes(r, authService, blizzardAuthService, authHandlers.AuthHandler, userHandler, authHandlers.BlizzardAuthHandler, rioHandler, blizzardHandler, warcraftlogsHandler)
	log.Println("Server is starting on :8080")
	log.Fatal(r.Run(":8080"))
}
