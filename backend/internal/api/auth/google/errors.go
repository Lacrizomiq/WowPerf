package googleauth

import (
	"log"

	"github.com/gin-gonic/gin"
)

// GoogleOAuthErrorType représente les types d'erreurs OAuth Google
type GoogleOAuthErrorType string

// Types d'erreurs OAuth Google
const (
	AccessDenied     GoogleOAuthErrorType = "access_denied"
	InvalidRequest   GoogleOAuthErrorType = "invalid_request"
	ServerError      GoogleOAuthErrorType = "server_error"
	InvalidGrant     GoogleOAuthErrorType = "invalid_grant"
	UnsupportedGrant GoogleOAuthErrorType = "unsupported_grant_type"
)

// GoogleOAuthErrorInfo contient les infos d'une erreur OAuth
type GoogleOAuthErrorInfo struct {
	ErrorCode   string
	UserMessage string
}

// oauthErrorMap mappe les erreurs Google vers des messages utilisateur
var oauthErrorMap = map[GoogleOAuthErrorType]GoogleOAuthErrorInfo{
	AccessDenied: {
		ErrorCode:   "auth_cancelled",
		UserMessage: "Google authentication was cancelled",
	},
	InvalidRequest: {
		ErrorCode:   "invalid_request",
		UserMessage: "Invalid authentication request",
	},
	ServerError: {
		ErrorCode:   "server_error",
		UserMessage: "Google authentication server error",
	},
	InvalidGrant: {
		ErrorCode:   "invalid_grant",
		UserMessage: "Invalid authorization grant",
	},
	UnsupportedGrant: {
		ErrorCode:   "unsupported_grant",
		UserMessage: "Unsupported grant type",
	},
}

// handleOAuthError gère les erreurs OAuth de manière centralisée
func (h *GoogleAuthHandler) handleOAuthError(c *gin.Context, errorType, description string) {
	log.Printf("OAuth error: type=%s, description=%s", errorType, description)

	// Chercher l'erreur dans la map
	if errorInfo, exists := oauthErrorMap[GoogleOAuthErrorType(errorType)]; exists {
		h.redirectToFrontendWithError(c, errorInfo.ErrorCode, errorInfo.UserMessage)
	} else {
		// Erreur inconnue - fallback
		log.Printf("Unknown OAuth error type: %s", errorType)
		h.redirectToFrontendWithError(c, "auth_failed", "Google authentication failed")
	}
}

// handleGenericError gère les erreurs génériques avec logging
func (h *GoogleAuthHandler) handleGenericError(c *gin.Context, err error, context string) {
	log.Printf("Error in %s: %v", context, err)
	h.redirectToFrontendWithError(c, "auth_processing_failed", "Authentication processing failed")
}
