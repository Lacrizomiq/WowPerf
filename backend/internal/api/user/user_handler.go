package user

import (
	"net/http"
	"os"
	"wowperf/internal/services/auth"
	"wowperf/internal/services/user"
	csrfMiddleware "wowperf/pkg/middleware"
	authMiddleware "wowperf/pkg/middleware/auth"

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
		"id":                    profile.ID,
		"username":              profile.Username,
		"email":                 profile.Email,
		"battle_tag":            profile.BattleTag,
		"battle_net_id":         profile.BattleNetID,
		"favorite_character_id": profile.FavoriteCharacterID,
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
			"code":  "INVALID_INPUT",
		})
		return
	}

	if err := h.userService.UpdateUsername(userID, input.NewUsername); err != nil {
		if err.Error() == "30 days must pass before changing username again" {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": err.Error(),
				"code":  "RATE_LIMIT_EXCEEDED",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to change username",
				"code":  "UPDATE_FAILED",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Username changed successfully",
		"code":    "USERNAME_UPDATED",
	})
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
	userRoutes := router.Group("/user")

	// All user endpoints require JWT
	userRoutes.Use(authMiddleware.JWTAuth(authService))
	{
		// Read-only routes - no CSRF
		userRoutes.GET("/profile", h.GetProfile)

		// Modification routes - need CSRF
		protected := userRoutes.Group("")
		protected.Use(csrfMiddleware.InitCSRFMiddleware(csrfMiddleware.NewCSRFConfig(os.Getenv("ENVIRONMENT"))))
		{
			protected.PUT("/username", h.ChangeUsername)
			protected.PUT("/email", h.UpdateEmail)
			protected.PUT("/password", h.ChangePassword)
			protected.DELETE("/account", h.DeleteAccount)
		}
	}
}
