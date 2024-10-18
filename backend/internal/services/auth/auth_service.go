package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"wowperf/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthService struct {
	DB              *gorm.DB
	JWTSecret       []byte
	RedisClient     *redis.Client
	TokenExpiration time.Duration
	OAuthConfig     *oauth2.Config
}

type BattleNetUserInfo struct {
	Sub       string `json:"sub"`
	ID        int    `json:"id"`
	BattleTag string `json:"battletag"`
}

func NewAuthService(db *gorm.DB, jwtSecret string, redisClient *redis.Client, tokenExpiration time.Duration, clientID, clientSecret, redirectURL string) *AuthService {
	return &AuthService{
		DB:              db,
		JWTSecret:       []byte(jwtSecret),
		RedisClient:     redisClient,
		TokenExpiration: tokenExpiration,
		OAuthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "wow.profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://eu.battle.net/oauth/authorize",
				TokenURL: "https://eu.battle.net/oauth/token",
			},
		},
	}
}

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

func (s *AuthService) Login(username, password string) (string, string, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", "", errors.New("invalid credentials")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := s.generateToken(user.ID, s.TokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return accessToken, refreshToken, nil
}

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
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

func (s *AuthService) BlacklistToken(token string, expiration time.Duration) error {
	ctx := context.Background()
	return s.RedisClient.Set(ctx, "blacklist:"+token, true, expiration).Err()
}

func (s *AuthService) IsTokenBlacklisted(token string) (bool, error) {
	ctx := context.Background()
	_, err := s.RedisClient.Get(ctx, "blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *AuthService) generateToken(userID uint, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiration).Unix(),
	})

	return token.SignedString(s.JWTSecret)
}

func (s *AuthService) generateRefreshToken(userID uint) (string, error) {
	ctx := context.Background()
	refreshToken := fmt.Sprintf("%d", time.Now().UnixNano())
	err := s.RedisClient.Set(ctx, fmt.Sprintf("refresh_token:%s", refreshToken), userID, 7*24*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store refresh token: %v", err)
	}
	return refreshToken, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	ctx := context.Background()
	userIDStr, err := s.RedisClient.Get(ctx, fmt.Sprintf("refresh_token:%s", refreshToken)).Result()
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return "", fmt.Errorf("invalid user ID in refresh token")
	}

	newAccessToken, err := s.generateToken(uint(userID), s.TokenExpiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %v", err)
	}

	return newAccessToken, nil
}

// Battle.net OAuth methods

func (s *AuthService) GetBattleNetAuthURL(state string) string {
	return s.OAuthConfig.AuthCodeURL(state)
}

func (s *AuthService) ExchangeBattleNetCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.OAuthConfig.Exchange(ctx, code)
}

func (s *AuthService) GetBattleNetUserInfo(token *oauth2.Token) (*BattleNetUserInfo, error) {
	client := s.OAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://eu.battle.net/oauth/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo BattleNetUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (s *AuthService) LinkBattleNetAccount(userID uint, battleNetInfo *BattleNetUserInfo, token *oauth2.Token) error {
	return s.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"battle_net_id":         battleNetInfo.ID,
		"battle_tag":            battleNetInfo.BattleTag,
		"battle_net_token":      token.AccessToken,
		"battle_net_expires_at": token.Expiry,
	}).Error
}
