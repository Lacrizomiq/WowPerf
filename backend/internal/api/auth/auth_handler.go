package auth

import (
	"net/http"
	"wowperf/internal/models"
	auth "wowperf/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
)

type AuthHandler struct {
	AuthService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// SignUp creates a new user
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

	user, token, err := h.AuthService.Login(loginInput.Username, loginInput.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := models.Validate.Struct(loginInput); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": formatValidationErrors(validationErrors)})
		return
	}

	if err := h.AuthService.CreateSession(c.Writer, c.Request, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.SetCookie("jwt", token, int(h.AuthService.AccessExpiry.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Logout logs out a user
func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token found"})
		return
	}

	if err := h.AuthService.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.SetCookie("jwt", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// RefreshToken refreshes the access token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	token, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token found"})
		return
	}

	newToken, err := h.AuthService.RefreshToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.SetCookie("jwt", newToken, int(h.AuthService.AccessExpiry.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed successfully"})
}

// RegisterRoutes registers the routes for the auth handler
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.SignUp)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.POST("/refresh", h.RefreshToken)
	}
}

func (h *AuthHandler) CSRFToken(c *gin.Context) {
	c.Header("X-CSRF-Token", csrf.Token(c.Request))
	c.JSON(http.StatusOK, gin.H{"csrf_token": csrf.Token(c.Request)})
}

// formatValidationErrors formats validation errors
func formatValidationErrors(errors validator.ValidationErrors) map[string]string {
	errorMap := make(map[string]string)
	for _, err := range errors {
		switch err.Tag() {
		case "required":
			errorMap[err.Field()] = "This field is required"
		case "email":
			errorMap[err.Field()] = "Invalid email format"
		case "min":
			errorMap[err.Field()] = "This field must be at least " + err.Param() + " characters long"
		case "max":
			errorMap[err.Field()] = "This field must not exceed " + err.Param() + " characters"
		case "strongpassword":
			errorMap[err.Field()] = "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character"
		default:
			errorMap[err.Field()] = "Invalid value"
		}
	}
	return errorMap
}
