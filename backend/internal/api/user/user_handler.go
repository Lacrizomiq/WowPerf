package user

import (
	"net/http"
	"wowperf/internal/services/auth"
	"wowperf/internal/services/user"
	"wowperf/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *user.UserService
}

func NewUserHandler(userService *user.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile returns the user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	profile, err := h.userService.GetUserProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       profile.ID,
		"username": profile.Username,
		"email":    profile.Email,
	})
}

// UpdateEmail updates the user's email
func (h *UserHandler) UpdateEmail(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input struct {
		NewEmail string `json:"new_email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.UpdateEmail(userID, input.NewEmail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email updated successfully"})
}

// ChangePassword changes the user's password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8,containsany=!@#$%^&*()_+"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.ChangePassword(userID, input.CurrentPassword, input.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ChangeUsername changes the user's username
func (h *UserHandler) ChangeUsername(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input struct {
		NewUsername string `json:"new_username" binding:"required,min=3,max=50"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.userService.UpdateUsername(userID, input.NewUsername); err != nil {
		if err.Error() == "username can only be changed once every 30 days" {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change username"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Username changed successfully"})
}

// DeleteAccount deletes the user's account
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := h.userService.DeleteAccount(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine, authService *auth.AuthService) {
	// Initialize JWT middleware
	jwtMiddleware := middleware.JWTAuth(authService)

	// Group for user routes
	user := router.Group("/user")
	user.Use(jwtMiddleware)

	// Read-only routes (no CSRF)
	user.GET("/profile", h.GetProfile)

	// Routes that modify data (with CSRF)
	csrfProtected := user.Group("")
	csrfProtected.Use(middleware.NewCSRFHandler())
	{
		csrfProtected.PUT("/email", h.UpdateEmail)
		csrfProtected.PUT("/password", h.ChangePassword)
		csrfProtected.PUT("/username", h.ChangeUsername)
		csrfProtected.DELETE("/account", h.DeleteAccount)
	}
}
