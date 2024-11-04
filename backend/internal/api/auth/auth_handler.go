// auth_handler.go
package auth

import (
	"errors"
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

	if err := h.authService.SignUp(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var loginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.authService.Login(c, loginInput.Username, loginInput.Password); err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	// Get CSRF token after successful login
	csrfToken := c.GetString("csrf_token")

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"username": loginInput.Username,
		},
		"csrf_token": csrfToken,
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
	// Initialize middlewares
	csrfMiddleware := authMiddleware.NewCSRFMiddleware()
	jwtMiddleware := authMiddleware.JWTAuth(h.authService)

	// Public routes group
	auth := router.Group("/auth")
	{
		// Routes that don't need CSRF protection (initial auth endpoints
		auth.POST("/signup", h.SignUp)
		auth.POST("/login", h.Login)

		// Get CSRF token(needs to be public but protected by CSRF)
		auth.GET("/csrf-token", csrfMiddleware, authMiddleware.GetCSRFToken())

		// Routes that need CSRF but not JWT
		csrfProtected := auth.Group("")
		csrfProtected.Use(csrfMiddleware)
		{
			csrfProtected.POST("/refresh", h.RefreshToken)
		}

		// Protected routes that need both JWT and CSRF
		protected := auth.Group("")
		protected.Use(jwtMiddleware, csrfMiddleware)
		{
			protected.POST("/logout", h.Logout)
			protected.GET("/check", h.CheckAuth)
		}
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
