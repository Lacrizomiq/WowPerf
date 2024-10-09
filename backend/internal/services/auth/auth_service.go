package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"wowperf/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB            *gorm.DB
	RedisClient   *redis.Client
	AccessSecret  []byte
	RefreshSecret []byte
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthService(db *gorm.DB, accessSecret, refreshSecret string, accessExpiry, refreshExpiry time.Duration) *AuthService {
	return &AuthService{
		DB:            db,
		RedisClient:   redis.NewClient(&redis.Options{Addr: "localhost:6379"}),
		AccessSecret:  []byte(accessSecret),
		RefreshSecret: []byte(refreshSecret),
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
	}
}

// SignUp creates a new user in the database
func (s *AuthService) SignUp(user *models.User) error {
	if user.Password == "" {
		return errors.New("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = string(hashedPassword)

	if err := s.DB.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return errors.New("duplicate key value violates unique constraint")
		}
		return fmt.Errorf("failed to create user in database: %v", err)
	}

	return nil
}

// Login authenticates a user and returns a JWT token pair
func (s *AuthService) Login(username, password string) (*http.Cookie, *http.Cookie, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	tokenPair, err := s.createTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    tokenPair.AccessToken,
		Expires:  time.Now().Add(s.AccessExpiry),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Expires:  time.Now().Add(s.RefreshExpiry),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/auth/refresh",
	}

	return accessCookie, refreshCookie, nil
}

// Logout invalidates the user's session token
func (s *AuthService) Logout(accessToken string) error {
	claims, err := s.validateToken(accessToken, s.AccessSecret)
	if err != nil {
		return err
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return errors.New("invalid jti claim")
	}

	return s.RedisClient.Set(context.Background(), "blacklist:"+jti, "true", s.AccessExpiry).Err()
}

// CreateLogoutCookie creates a logout cookie
func (s *AuthService) CreateLogoutCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	}
}

// RefreshToken refreshes the user's access token
func (s *AuthService) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := s.validateToken(refreshToken, s.RefreshSecret)
	if err != nil {
		return nil, err
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID claim")
	}

	return s.createTokenPair(uint(userID))
}

func (s *AuthService) createTokenPair(userID uint) (*TokenPair, error) {
	accessJTI := generateJTI()
	refreshJTI := generateJTI()

	accessToken, err := s.createToken(userID, accessJTI, s.AccessSecret, s.AccessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.createToken(userID, refreshJTI, s.RefreshSecret, s.RefreshExpiry)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// createToken creates a JWT token
func (s *AuthService) createToken(userID uint, jti string, secret []byte, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"jti":     jti,
		"exp":     time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// validateToken validates a JWT token
func (s *AuthService) validateToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// IsTokenBlacklisted checks if a JWT token is blacklisted
func (s *AuthService) IsTokenBlacklisted(jti string) (bool, error) {
	result, err := s.RedisClient.Exists(context.Background(), "blacklist:"+jti).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}

func generateJTI() string {
	jti := make([]byte, 16)
	_, err := rand.Read(jti)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(jti)
}
