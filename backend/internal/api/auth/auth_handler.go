// auth_handler.go
package auth

import (
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"wowperf/internal/models"
	auth "wowperf/internal/services/auth"
	csrfMiddleware "wowperf/pkg/middleware"
	authMiddleware "wowperf/pkg/middleware/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

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
		log.Printf("SignUp binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input format",
			"code":    "invalid_input",
			"details": err.Error(),
		})
		return
	}

	// Validating password length separately
	if len(userCreate.Password) < 8 {
		log.Printf("Password length validation failed: got %d chars, need at least 8", len(userCreate.Password))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 characters long",
			"code":  "invalid_password",
		})
		return
	}

	// Check if user already exists BEFORE other validations
	var existingUser models.User
	if err := h.authService.DB.Where("username = ? OR email = ?",
		userCreate.Username, userCreate.Email).First(&existingUser).Error; err == nil {
		if existingUser.Username == userCreate.Username {
			log.Printf("Username already exists: %s", userCreate.Username)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Username already exists",
				"code":  "username_exists",
			})
		} else {
			log.Printf("Email already exists: %s", userCreate.Email)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Email already exists",
				"code":  "email_exists",
			})
		}
		return
	}

	// Add all field validations
	// Validate username
	if len(userCreate.Username) < 3 || len(userCreate.Username) > 12 {
		log.Printf("Username length validation failed: got %d chars, need between 3 and 12", len(userCreate.Username))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username must be between 3 and 12 characters",
			"code":  "invalid_username",
		})
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(userCreate.Email) {
		log.Printf("Email format validation failed: %s", userCreate.Email)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
			"code":  "invalid_email",
		})
		return
	}

	// Create the user
	user := models.User{
		Username: userCreate.Username,
		Email:    userCreate.Email,
		Password: userCreate.Password,
	}

	if err := h.authService.SignUp(&user); err != nil {
		log.Printf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"code":    "server_error",
			"details": err.Error(),
		})
		return
	}

	// Auto-login the user after signup
	if err := h.authService.Login(c, user.Email, userCreate.Password); err != nil {
		log.Printf("Failed to auto-login user after signup: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"message": "User created successfully, but login failed",
			"code":    "signup_success_login_failed",
		})
		return
	}

	log.Printf("User created successfully: %s, %s", user.Username, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"code":    "signup_success",
		"user":    gin.H{"username": user.Username, "email": user.Email},
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var loginInput struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
			"code":  "invalid_input",
		})
		return
	}

	if err := h.authService.Login(c, loginInput.Email, loginInput.Password); err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
				"code":  "invalid_credentials",
			})
			return
		}
		log.Printf("Login error for user %s: %v", loginInput.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to login",
			"code":  "login_error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"code":    "LOGIN_SUCCESS",
		"user":    gin.H{"email": loginInput.Email},
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.authService.Logout(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout", "code": "logout_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful", "code": "logout_success"})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	if err := h.authService.RefreshToken(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token", "code": "refresh_token_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Refresh token successful", "code": "refresh_token_success"})
}

// RegisterRoutes registers the authentication routes
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		// Public routes - no CSRF or JWT
		auth.POST("/signup", h.SignUp)
		auth.POST("/login", h.Login)

		// Routes protected by JWT only
		jwtProtected := auth.Group("")
		jwtProtected.Use(authMiddleware.JWTAuth(h.authService))
		{
			jwtProtected.GET("/check", h.CheckAuth)
		}

		// Routes protected by JWT and CSRF
		csrfProtected := auth.Group("")
		csrfProtected.Use(authMiddleware.JWTAuth(h.authService))
		csrfProtected.Use(csrfMiddleware.InitCSRFMiddleware(csrfMiddleware.NewCSRFConfig(os.Getenv("ENVIRONMENT"))))
		{
			csrfProtected.POST("/logout", h.Logout)
			csrfProtected.POST("/refresh", h.RefreshToken)
		}
	}
}

// CheckAuth checks if the user is authenticated
func (h *AuthHandler) CheckAuth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"authenticated": true, "code": "check_auth_success"})
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
