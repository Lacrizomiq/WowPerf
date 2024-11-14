package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	migrations "wowperf/internal/database/migrations"
	models "wowperf/internal/models/raiderio/mythicrundetails"

	// Internal Packages - Utils
	cacheMiddleware "wowperf/middleware/cache"
	"wowperf/pkg/cache"
	csrf "wowperf/pkg/middleware"
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

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		if os.Getenv("ENV") == "production" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}

func setupHealthCheck(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})
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
	// Health check endpoint
	setupHealthCheck(r)

	// Initialize the CSRF middleware
	csrf.InitCSRFMiddleware()

	// Security headers middleware with HTTPS support
	r.Use(securityHeaders())

	// Get environment variables
	environment := os.Getenv("ENVIRONMENT")
	frontendURL := os.Getenv("FRONTEND_URL")

	// Default allowed origins
	allowedOrigins := []string{
		"https://localhost",
		"https://*.localhost",
		"https://api.localhost",
		frontendURL,
	}

	// Add production domain if in production
	if environment == "production" {
		// In production, only allow the specific frontend URL
		allowedOrigins = []string{frontendURL}
	}

	// CORS Configuration
	corsConfig := cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-CSRF-Token",
			"X-Requested-With",
			"Cookie",
			"Set-Cookie",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Set-Cookie",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Credentials",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		// Important: Add custom handling for OPTIONS requests
		AllowWildcard: true,
	}
	r.Use(cors.New(corsConfig))

	// Add OPTIONS handler for preflight requests
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(http.StatusOK)
	})
	// Logger middleware
	r.Use(gin.Logger())

	r.Use(func(c *gin.Context) {
		log.Printf("📨 Request: %s %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("📝 Headers: %v", c.Request.Header)

		c.Next()

		log.Printf("📤 Response Status: %d", c.Writer.Status())
		log.Printf("📤 Response Headers: %v", c.Writer.Header())
	})

	// Initialize CSRF middleware
	r.Use(csrf.NewCSRFHandler())

	// CSRF Token endpoint
	r.GET("/api/csrf-token", csrf.GetCSRFToken())

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
	userHandler.RegisterRoutes(r, authService)
}

func main() {

	// Verify required environment variables
	requiredEnvVars := []string{
		"JWT_SECRET",
		"CSRF_SECRET",
		"REDIS_URL",
	}
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
		return
	}

	// Verify required environment variables
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("%s is not set", envVar)
		}
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

	if err := migrations.RunMigrations(db); err != nil {
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
