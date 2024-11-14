// auth_handler.go
package auth

import (
	"errors"
	"log"
	"net/http"
	"wowperf/internal/models"
	auth "wowperf/internal/services/auth"
	authMiddleware "wowperf/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Handlers struct to hold all auth handlers
type Handlers struct {
	AuthHandler         *AuthHandler
	BlizzardAuthHandler *BlizzardAuthHandler
}

// NewHandlers creates all auth handlers
func NewHandlers(authService *auth.AuthService, blizzardAuthService *auth.BlizzardAuthService) *Handlers {
	return &Handlers{
		AuthHandler:         NewAuthHandler(authService),
		BlizzardAuthHandler: NewBlizzardAuthHandler(blizzardAuthService, authService),
	}
}

// AuthHandler handles user authentication endpoints
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SignUp handles user registration
func (h *AuthHandler) SignUp(c *gin.Context) {
	var userCreate models.UserCreate
	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "code": "invalid_input"})
		return
	}

	// Validate the user create struct
	if len(userCreate.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long", "code": "password_too_short"})
		return
	}

	if err := models.Validate.Struct(userCreate); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "code": "validation_error"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": formatValidationErrors(validationErrors)})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.authService.DB.Where("username = ? OR email = ?",
		userCreate.Username, userCreate.Email).First(&existingUser).Error; err == nil {
		if existingUser.Username == userCreate.Username {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Username already exists",
				"code":  "username_exists",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Email already exists",
				"code":  "email_exists",
			})
		}
		return
	}

	user := models.User{
		Username: userCreate.Username,
		Email:    userCreate.Email,
		Password: userCreate.Password,
	}

	if err := h.authService.SignUp(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "code": "server_error"})
		return
	}

	// Auto-login the user after signup
	if err := h.authService.Login(c, user.Username, userCreate.Password); err != nil {
		log.Println("Failed to auto-login user after signup:", err)
		c.JSON(http.StatusOK, gin.H{
			"message": "User created successfully",
			"code":    "signup_success_login_required",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "code": "signup_success"})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var loginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
			"code":  "invalid_input",
		})
		return
	}

	if err := h.authService.Login(c, loginInput.Username, loginInput.Password); err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
				"code":  "invalid_credentials",
			})
			return
		}
		log.Printf("Login error for user %s: %v", loginInput.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to login",
			"code":  "login_error",
		})
		return
	}

	// Get the CSRF token from the request
	authMiddleware.GetTokenFromRequest(c, func(token string) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"user": gin.H{
				"username": loginInput.Username,
			},
			"csrf_token": token,
		})
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.authService.Logout(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	if err := h.authService.RefreshToken(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Refresh token successful"})
}

// RegisterRoutes registers the authentication routes
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	// Group for the auth
	auth := router.Group("/auth")

	// public route
	auth.GET("/csrf-token", authMiddleware.GetCSRFToken())

	// Base auth routes
	auth.POST("/signup", h.SignUp)
	auth.POST("/login", h.Login)

	// Protected routes with JWT
	protected := auth.Group("")
	protected.Use(authMiddleware.JWTAuth(h.authService))
	{
		protected.POST("/refresh", h.RefreshToken)
		protected.POST("/logout", h.Logout)
		protected.GET("/check", h.CheckAuth)
	}
}

// CheckAuth checks if the user is authenticated
func (h *AuthHandler) CheckAuth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"authenticated": true})
}

// Helper functions
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
