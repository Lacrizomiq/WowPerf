package auth

import (
	"context"
	"errors"
	"time"
	"wowperf/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AuthService struct {
	DB          *gorm.DB
	JWTKey      []byte
	RedisClient *redis.Client
	TokenExpiry time.Duration
}

func NewAuthService(db *gorm.DB, jwtKey string, redis *redis.Client, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		DB:          db,
		JWTKey:      []byte(jwtKey),
		RedisClient: redis,
		TokenExpiry: tokenExpiry,
	}
}

// SignUp creates a new user in the database
func (s *AuthService) SignUp(user *models.User) error {
	return s.DB.Create(user).Error
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(username, password string) (string, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", err
	}

	if err := user.ComparePassword(password); err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(s.TokenExpiry)

	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.JWTKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Logout invalidates the user's session token
func (s *AuthService) Logout(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.JWTKey, nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		duration := time.Until(expirationTime)

		if duration > 0 {
			err := s.RedisClient.Set(context.Background(), "blacklist:"+tokenString, "true", duration).Err()
			if err != nil {
				return err
			}
		}
		return nil
	}

	return errors.New("invalid token")
}

// IsTokenBlacklisted checks if a token is blacklisted
func (s *AuthService) IsTokenBlacklisted(tokenString string) (bool, error) {
	result, err := s.RedisClient.Exists(context.Background(), "blacklist:"+tokenString).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
