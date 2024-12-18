// Package auth provides authentication services for the application
package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"wowperf/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	AccessTokenDuration  = 7 * 24 * time.Hour  // 7 days
	RefreshTokenDuration = 30 * 24 * time.Hour // 30 days
)

// Common errors returned by the auth service
var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrTokenExpired         = errors.New("token has expired")
	ErrTokenInvalid         = errors.New("invalid token")
	ErrTokenBlacklisted     = errors.New("token has been blacklisted")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

// CookieConfig contains cookie configuration parameters
type CookieConfig struct {
	Domain   string
	Path     string
	Secure   bool
	SameSite http.SameSite
}

// AuthService handles user authentication and token management
type AuthService struct {
	DB              *gorm.DB
	JWTSecret       []byte
	RedisClient     *redis.Client
	TokenExpiration time.Duration
	CookieConfig    CookieConfig
	BlizzardAuth    *BlizzardAuthService // Reference to Battle.net auth service
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(
	db *gorm.DB,
	jwtSecret string,
	redisClient *redis.Client,
	blizzardAuth *BlizzardAuthService,
) *AuthService {
	return &AuthService{
		DB:              db,
		JWTSecret:       []byte(jwtSecret),
		RedisClient:     redisClient,
		TokenExpiration: AccessTokenDuration,
		BlizzardAuth:    blizzardAuth,
		CookieConfig: CookieConfig{
			Domain:   ".localhost", // Add . to the domain to allow subdomains
			Path:     "/",
			Secure:   true,                 // Set to true in production for https
			SameSite: http.SameSiteLaxMode, // Set to SameSiteLaxMode in production
		},
	}
}

// setAuthCookies sets the authentication cookies in the response
func (s *AuthService) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	c.SetCookie(
		"access_token",
		accessToken,
		int(AccessTokenDuration.Seconds()),
		s.CookieConfig.Path,
		s.CookieConfig.Domain,
		s.CookieConfig.Secure,
		true,
	)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(RefreshTokenDuration.Seconds()),
		s.CookieConfig.Path,
		s.CookieConfig.Domain,
		s.CookieConfig.Secure,
		true,
	)
}

// SignUp registers a new user
func (s *AuthService) SignUp(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)

	if err := s.DB.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Login authenticates a user and generates tokens
func (s *AuthService) Login(c *gin.Context, username, password string) error {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}

	accessToken, err := s.generateToken(user.ID, s.TokenExpiration)
	if err != nil {
		return fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return fmt.Errorf("failed to generate refresh token: %w", err)
	}

	s.setAuthCookies(c, accessToken, refreshToken)
	return nil
}

// Logout handles user logout and token invalidation
func (s *AuthService) Logout(c *gin.Context) error {
	if token, err := c.Cookie("access_token"); err == nil {
		if err := s.BlacklistToken(token, s.TokenExpiration); err != nil {
			return fmt.Errorf("failed to blacklist token: %w", err)
		}
	}

	if refreshToken, err := c.Cookie("refresh_token"); err == nil {
		ctx := context.Background()
		s.RedisClient.Del(ctx, fmt.Sprintf("refresh_token:%s", refreshToken))
	}

	s.clearAuthCookies(c)
	return nil
}

// clearAuthCookies removes authentication cookies
func (s *AuthService) clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, s.CookieConfig.Path, s.CookieConfig.Domain, s.CookieConfig.Secure, true)
	c.SetCookie("refresh_token", "", -1, s.CookieConfig.Path, s.CookieConfig.Domain, s.CookieConfig.Secure, true)
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.JWTSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return uint(claims["user_id"].(float64)), nil
	}

	return 0, ErrTokenInvalid
}

// BlacklistToken adds a token to the blacklist
func (s *AuthService) BlacklistToken(token string, expiration time.Duration) error {
	ctx := context.Background()
	return s.RedisClient.Set(ctx, "blacklist:"+token, true, expiration).Err()
}

// IsTokenBlacklisted checks if a token is blacklisted
func (s *AuthService) IsTokenBlacklisted(token string) (bool, error) {
	ctx := context.Background()
	_, err := s.RedisClient.Get(ctx, "blacklist:"+token).Result()

	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	return true, nil
}

// generateToken creates a new JWT token
func (s *AuthService) generateToken(userID uint, expiration time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(expiration).Unix(),
		"iat":     now.Unix(),
		"type":    "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JWTSecret)
}

// generateRefreshToken creates a new refresh token
func (s *AuthService) generateRefreshToken(userID uint) (string, error) {
	ctx := context.Background()
	refreshToken := fmt.Sprintf("%d", time.Now().UnixNano())
	err := s.RedisClient.Set(ctx, fmt.Sprintf("refresh_token:%s", refreshToken), userID, RefreshTokenDuration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}
	return refreshToken, nil
}

// RefreshToken handles token refresh
func (s *AuthService) RefreshToken(c *gin.Context) error {

	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return errors.New("invalid refresh token")
	}

	// Create context for redis
	ctx := context.Background()

	// Try to acquire a lock on the refresh token
	lockKey := fmt.Sprintf("refresh_lock:%s", refreshToken)
	lock := s.RedisClient.SetNX(ctx, lockKey, "locked", 10*time.Second)
	if !lock.Val() {
		return errors.New("refresh already in progress")
	}
	defer s.RedisClient.Del(ctx, lockKey)

	// Verify refresh token
	userIDStr, err := s.RedisClient.Get(ctx, fmt.Sprintf("refresh_token:%s", refreshToken)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrRefreshTokenNotFound
		}
		return fmt.Errorf("failed to verify refresh token: %w", err)
	}

	// Parse user ID from string
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return errors.New("invalid user ID in refresh token")
	}

	// Generate new access token
	newAccessToken, err := s.generateToken(uint(userID), s.TokenExpiration)
	if err != nil {
		return fmt.Errorf("failed to generate new access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := s.generateRefreshToken(uint(userID))
	if err != nil {
		return fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	// Set new cookie
	s.setAuthCookies(c, newAccessToken, newRefreshToken)

	// Delete old refresh token
	err = s.RedisClient.Del(ctx, fmt.Sprintf("refresh_token:%s", refreshToken)).Err()
	if err != nil {
		log.Printf("Warning: Failed to delete old refresh token: %v", err)
	}

	return nil
}

// Helper method to check if token is about to expire
func (s *AuthService) isTokenNearExpiry(token string) (bool, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return s.JWTSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return false, ErrTokenInvalid
		}
	}

	if exp, ok := claims["exp"].(float64); ok {
		// Check if token expires in the next 10 minutes
		timeUntilExpiry := time.Unix(int64(exp), 0).Sub(time.Now())
		return timeUntilExpiry < 10*time.Minute, nil
	}

	return false, errors.New("invalid token claims")
}
