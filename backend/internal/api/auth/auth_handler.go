package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"
	"wowperf/internal/models"
	auth "wowperf/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	AuthService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var userCreate models.UserCreate
	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	if err := models.Validate.Struct(userCreate); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": formatValidationErrors(validationErrors)})
		return
	}

	user := models.User{
		Username: userCreate.Username,
		Email:    userCreate.Email,
		Password: userCreate.Password,
	}

	if err := h.AuthService.SignUp(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	accessToken, refreshToken, err := h.AuthService.Login(loginInput.Username, loginInput.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is missing"})
		return
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization header format"})
		return
	}

	token := bearerToken[1]

	// Blacklist the token
	if err := h.AuthService.BlacklistToken(token, 24*time.Hour); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	newAccessToken, err := h.AuthService.RefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}

// Battle.net OAuth handlers

// Battle.net OAuth login
func (h *AuthHandler) HandleBattleNetLogin(c *gin.Context) {
	state := generateRandomState()
	c.SetCookie("oauth_state", state, 3600, "/", "", false, true)
	url := h.AuthService.GetBattleNetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Battle.net OAuth callback
func (h *AuthHandler) HandleBattleNetCallback(c *gin.Context) {
	state, _ := c.Cookie("oauth_state")
	if state != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	code := c.Query("code")
	token, err := h.AuthService.ExchangeBattleNetCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	userInfo, err := h.AuthService.GetBattleNetUserInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Assuming the user is already authenticated and we have their ID
	userID, _ := c.Get("user_id")
	if err := h.AuthService.LinkBattleNetAccount(userID.(uint), userInfo, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link Battle.net account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Battle.net account linked successfully"})
}

func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.SignUp)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.POST("/refresh", h.RefreshToken)
		auth.GET("/battle-net/login", h.HandleBattleNetLogin)
		auth.GET("/battle-net/callback", h.AuthService.AuthMiddleware(), h.HandleBattleNetCallback)
	}
}

func formatValidationErrors(errors validator.ValidationErrors) []string {
	var formattedErrors []string
	for _, err := range errors {
		formattedErrors = append(formattedErrors, formatValidationError(err))
	}
	return formattedErrors
}

func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email address"
	default:
		return err.Field() + " is invalid"
	}
}

func generateRandomState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
