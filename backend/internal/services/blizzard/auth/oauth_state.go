// oauth_state.go
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// OAuthState represents the state of an OAuth request
type OAuthState struct {
	State      string    `json:"state"`       // Random state string
	UserID     uint      `json:"user_id"`     // User ID
	AutoRelink bool      `json:"auto_relink"` // 🔥 NOUVEAU: Flag pour auto-relink
	ExpiresAt  time.Time `json:"expires_at"`  // Expiration time
}

// GenerateRandomState generates a random state string
func generateRandomState(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// 🔥 MODIFIÉ: NewOAuthState avec support AutoRelink
func NewOAuthState(userID uint, autoRelink bool) (*OAuthState, error) {
	// Generate random state string
	state, err := generateRandomState(32)
	if err != nil {
		return nil, err
	}

	return &OAuthState{
		State:      state,
		UserID:     userID,
		AutoRelink: autoRelink,                       // 🔥 NOUVEAU
		ExpiresAt:  time.Now().Add(15 * time.Minute), // State expires in 15 minutes
	}, nil
}

// Marshal marshals the OAuthState to a JSON string
func (s *OAuthState) Marshal() (string, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshal OAuth state: %w", err)
	}
	return string(bytes), nil
}

// Unmarshal unmarshals the OAuthState from a JSON string
func UnmarshalOAuthState(data string) (*OAuthState, error) {
	var state OAuthState
	if err := json.Unmarshal([]byte(data), &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OAuth state: %w", err)
	}

	// Verify expiration
	if time.Now().After(state.ExpiresAt) {
		return nil, fmt.Errorf("OAuth state has expired")
	}

	return &state, nil
}
