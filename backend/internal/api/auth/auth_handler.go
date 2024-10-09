package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"wowperf/internal/models"
	"wowperf/internal/services/auth"
	"wowperf/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// AuthHandler handles authentication routes
type AuthHandler struct {
	AuthService *auth.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// SignUp creates a new user
func (h *AuthHandler) SignUp(c *gin.Context) {
	var userCreate models.UserCreate
	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %v", err)})
		return
	}

	user := models.User{
		Username: userCreate.Username,
		Email:    userCreate.Email,
		Password: userCreate.Password,
	}

	if err := h.AuthService.SignUp(&user); err != nil {
		if err.Error() == "duplicate key value violates unique constraint" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		} else {
			log.Printf("Failed to create user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login authenticates a user and returns a JWT token pair
func (h *AuthHandler) Login(c *gin.Context) {
	var loginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	accessCookie, refreshCookie, err := h.AuthService.Login(loginInput.Username, loginInput.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	http.SetCookie(c.Writer, accessCookie)
	http.SetCookie(c.Writer, refreshCookie)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Logout invalidates the user's session token
func (h *AuthHandler) Logout(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.AuthService.Logout(accessToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	logoutCookie := h.AuthService.CreateLogoutCookie()
	http.SetCookie(c.Writer, logoutCookie)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tokenPair, err := h.AuthService.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    tokenPair.AccessToken,
		Expires:  time.Now().Add(h.AuthService.AccessExpiry),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}

	http.SetCookie(c.Writer, accessCookie)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed successfully"})
}

func (h *AuthHandler) GetCSRFToken(c *gin.Context) {
	token := csrf.Token(c.Request)
	c.JSON(http.StatusOK, gin.H{"csrf_token": token})
}

// RegisterRoutes registers the routes for the auth handler
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", middleware.CSRF(), h.SignUp)
		auth.POST("/login", middleware.CSRF(), h.Login)
		auth.POST("/logout", middleware.CSRF(), h.Logout)
		auth.POST("/refresh", middleware.CSRF(), h.RefreshToken)
		auth.GET("/csrf", h.GetCSRFToken)
	}
}
