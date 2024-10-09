package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"wowperf/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB            *gorm.DB
	RedisClient   *redis.Client
	SessionStore  *sessions.CookieStore
	JWTSecret     []byte
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func NewAuthService(db *gorm.DB, redisClient *redis.Client, sessionSecret, jwtSecret string, accessExpiry, refreshExpiry time.Duration) *AuthService {
	return &AuthService{
		DB:            db,
		RedisClient:   redisClient,
		SessionStore:  sessions.NewCookieStore([]byte(sessionSecret)),
		JWTSecret:     []byte(jwtSecret),
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
	}
}

// SignUp creates a new user in the database
func (s *AuthService) SignUp(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = string(hashedPassword)

	if err := s.DB.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

// Login authenticates a user and returns a JWT token pair
func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.GenerateJWT(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %v", err)
	}

	return &user, token, nil
}

// Logout blacklists the token and clears the session
func (s *AuthService) Logout(tokenString string) error {
	return s.BlacklistToken(tokenString)
}

// RefreshToken refreshes the access token using the refresh token
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return s.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return "", errors.New("invalid user ID in token")
	}

	var user models.User
	if err := s.DB.First(&user, uint(userID)).Error; err != nil {
		return "", errors.New("user not found")
	}

	newToken, err := s.GenerateJWT(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %v", err)
	}

	return newToken, nil
}

// GenerateJWT generates a JWT token for a user
func (s *AuthService) GenerateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.AccessExpiry).Unix(),
	})

	return token.SignedString(s.JWTSecret)
}

// BlacklistToken blacklists a token by adding it to Redis
func (s *AuthService) BlacklistToken(tokenString string) error {
	return s.RedisClient.Set(context.Background(), "blacklist:"+tokenString, "true", s.AccessExpiry).Err()
}

// IsTokenBlacklisted checks if a token is blacklisted
func (s *AuthService) IsTokenBlacklisted(tokenString string) (bool, error) {
	result, err := s.RedisClient.Exists(context.Background(), "blacklist:"+tokenString).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// CreateSession creates a session for a user
func (s *AuthService) CreateSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	session, _ := s.SessionStore.Get(r, "session-name")
	session.Values["user_id"] = user.ID
	return session.Save(r, w)
}

// GetUserFromSession retrieves the user from the session
func (s *AuthService) GetUserFromSession(r *http.Request) (*models.User, error) {
	session, _ := s.SessionStore.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(uint)
	if !ok {
		return nil, errors.New("no user in session")
	}

	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
