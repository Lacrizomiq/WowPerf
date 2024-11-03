package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"wowperf/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var (
	// ErrInvalidCredentials Returned when credentials are invalid
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrTokenExpired         = errors.New("token has expired")
	ErrTokenInvalid         = errors.New("invalid token")
	ErrTokenBlacklisted     = errors.New("token has been blacklisted")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

// CookieConfig contains the parameters of the cookies configuration
type CookieConfig struct {
	Domain   string
	Path     string
	Secure   bool
	SameSite http.SameSite
}

// AuthService extend the struct
type AuthService struct {
	DB              *gorm.DB
	JWTSecret       []byte
	RedisClient     *redis.Client
	TokenExpiration time.Duration
	OAuthConfig     *oauth2.Config
	CookieConfig    CookieConfig
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
		CookieConfig: CookieConfig{
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		},
	}
}

// Method to handle the cookie
func (s *AuthService) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	c.SetCookie(
		"access_token",
		accessToken,
		int(s.TokenExpiration.Seconds()),
		s.CookieConfig.Path,
		s.CookieConfig.Domain,
		s.CookieConfig.Secure,
		true, // httpOnly always true for the tokens
	)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int((7 * 24 * time.Hour).Seconds()), // 7 days for the refresh
		s.CookieConfig.Path,
		s.CookieConfig.Domain,
		s.CookieConfig.Secure,
		true,
	)
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
		return fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return fmt.Errorf("failed to generate refresh token: %v", err)
	}

	s.setAuthCookies(c, accessToken, refreshToken)
	return nil
}

func (s *AuthService) Logout(c *gin.Context) error {
	// Retrieve the actual token to blacklist it
	if token, err := c.Cookie("access_token"); err == nil {
		if err := s.BlacklistToken(token, s.TokenExpiration); err != nil {
			return fmt.Errorf("failed to blacklist token: %v", err)
		}
	}

	// Retrieve & invalidate the refresh token
	if refreshToken, err := c.Cookie("refresh_token"); err == nil {
		ctx := context.Background()
		s.RedisClient.Del(ctx, fmt.Sprintf("refresh_token:%s", refreshToken))
	}

	// Delete the cookies
	c.SetCookie("access_token", "", -1, s.CookieConfig.Path, s.CookieConfig.Domain, false, true)
	c.SetCookie("refresh_token", "", -1, s.CookieConfig.Path, s.CookieConfig.Domain, false, true)

	return nil
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

	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {

		return false, fmt.Errorf("failed to check token blacklist: %v", err)
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

func (s *AuthService) RefreshToken(c *gin.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return errors.New("invalid refresh token")
	}

	// Verify the refresh token in Redis
	ctx := context.Background()
	userIDStr, err := s.RedisClient.Get(ctx, fmt.Sprintf("refresh_token:%s", refreshToken)).Result()
	if err != nil {
		return errors.New("invalid refresh token")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return errors.New("invalid user ID in refresh token")
	}

	// Generate new tokens
	newAccessToken, err := s.generateToken(uint(userID), s.TokenExpiration)
	if err != nil {
		return fmt.Errorf("failed to generate new access token: %v", err)
	}

	newRefreshToken, err := s.generateRefreshToken(uint(userID))
	if err != nil {
		return fmt.Errorf("failed to generate new refresh token: %v", err)
	}

	// Update the cookies
	s.setAuthCookies(c, newAccessToken, newRefreshToken)

	// Blacklist the old refresh token
	s.RedisClient.Del(ctx, fmt.Sprintf("refresh_token:%s", refreshToken))

	return nil
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
