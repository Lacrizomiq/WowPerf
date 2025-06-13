// auth_handler.go
package auth

import (
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"wowperf/internal/models"
	auth "wowperf/internal/services/auth"
	csrfMiddleware "wowperf/pkg/middleware"
	authMiddleware "wowperf/pkg/middleware/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ForgotPasswordRequest represents the request structure for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the request structure for reset password
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
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

	// Message d'erreur spécifique pour le captcha, username ou email déjà existant
	if err := h.authService.SignUp(&user, userCreate.CaptchaToken); err != nil {
		log.Printf("Failed to create user: %v", err)

		// Messages d'erreur spécifiques
		if strings.Contains(err.Error(), "captcha verification failed") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Please complete the captcha verification",
				"code":  "captcha_required",
			})
			return
		}

		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(), // "user with this email already exists"
				"code":  "user_exists",
			})
			return
		}

		// Erreur générique
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
			"code":  "server_error",
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

// ForgotPassword initiates a password reset for a user
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
			"code":  "invalid_email_format",
		})
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
			"code":  "invalid_email_format",
		})
		return
	}

	if err := h.authService.InitiatePasswordReset(req.Email); err != nil {
		// Log the error but do not return it to the user to avoid email enumeration
		log.Printf("Failed to initiate password reset: %v", err)

		// Always return the same success message
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been sent",
			"code":    "reset_email_sent",
		})
		return
	}

	// Return success to prevent email enumeration
	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a password reset link has been sent",
		"code":    "reset_email_sent",
	})
}

// ValidateResetToken checks if a reset token is valid
func (h *AuthHandler) ValidateResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token is required",
			"code":  "token_required",
		})
		return
	}

	// Validate the token
	log.Printf("Validating reset token")
	if _, err := h.authService.ValidateResetToken(token); err != nil {
		log.Printf("Invalid reset token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or expired token",
			"code":  "invalid_token",
		})
		return
	}

	// Return success to prevent token enumeration
	log.Printf("Reset token validated successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "Token is valid",
		"code":    "valid_token",
	})
}

// ResetPassword handles the password reset request
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "invalid_request",
		})
		return
	}

	// Validation more complete of the password
	if len(req.NewPassword) < 8 || len(req.NewPassword) > 32 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be between 8 and 32 characters",
			"code":  "invalid_password_length",
		})
		return
	}

	// Validate the token before attempting the reset
	if _, err := h.authService.ValidateResetToken(req.Token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or expired token",
			"code":  "invalid_token",
		})
		return
	}

	// Reset the password
	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		log.Printf("Failed to reset password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to reset password",
			"code":  "reset_password_failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password has been reset successfully",
		"code":    "password_reset_success",
	})
}

// RegisterRoutes registers the authentication routes
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		// Public routes - no CSRF or JWT
		auth.POST("/signup", h.SignUp)
		auth.POST("/login", h.Login)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.GET("/validate-reset-token", h.ValidateResetToken)
		auth.POST("/reset-password", h.ResetPassword)

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
