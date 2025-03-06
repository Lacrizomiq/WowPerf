package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"wowperf/pkg/crypto"

	"gorm.io/gorm"
)

// User is the struct for the user model
type User struct {
	gorm.Model
	ID                   uint      `gorm:"primaryKey" json:"id"`
	Username             string    `gorm:"uniqueIndex;not null" json:"username" validate:"required,min=3,max=50"`
	Email                string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password             string    `gorm:"not null" json:"-" validate:"required,min=8"`
	LastUsernameChangeAt time.Time `json:"last_username_change_at"`

	// Reset password fields
	ResetPasswordToken   *string    `gorm:"uniqueIndex" json:"-"`
	ResetPasswordExpires *time.Time `json:"-"`

	// Battle.net specific fields - now with pointers
	BattleNetID           *string   `gorm:"uniqueIndex" json:"battle_net_id"`
	BattleTag             *string   `gorm:"uniqueIndex" json:"battle_tag"`
	EncryptedAccessToken  []byte    `gorm:"type:bytea" json:"-"`
	EncryptedRefreshToken []byte    `gorm:"type:bytea" json:"-"`
	BattleNetTokenType    string    `gorm:"type:varchar(50)" json:"-"`
	BattleNetExpiresAt    time.Time `json:"-"`
	BattleNetScopes       []string  `gorm:"type:text[]" json:"-"`
	LastTokenRefresh      time.Time `json:"-"`

	// Character relationship
	FavoriteCharacterID *uint           `json:"favorite_character_id"`
	Characters          []UserCharacter `gorm:"foreignKey:UserID" json:"characters,omitempty"`
}

// UserCreate is the struct for creating a new user
type UserCreate struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Token management methods
func (u *User) SetBattleNetTokens(accessToken, refreshToken string) error {
	log.Printf("Starting SetBattleNetTokens: access_token_length=%d refresh_token_length=%d",
		len(accessToken), len(refreshToken))

	if accessToken == "" {
		return fmt.Errorf("access token is empty")
	}

	// Encrypt access token
	encryptedAccess, err := crypto.Encrypt([]byte(accessToken))
	if err != nil {
		log.Printf("Failed to encrypt access token: %v", err)
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}
	log.Printf("Access token encrypted successfully: length=%d", len(encryptedAccess))

	// Encrypt refresh token if present
	var encryptedRefresh []byte
	if refreshToken != "" {
		encryptedRefresh, err = crypto.Encrypt([]byte(refreshToken))
		if err != nil {
			log.Printf("Failed to encrypt refresh token: %v", err)
			return fmt.Errorf("failed to encrypt refresh token: %w", err)
		}
		log.Printf("Refresh token encrypted successfully: length=%d", len(encryptedRefresh))
	}

	u.EncryptedAccessToken = encryptedAccess
	u.EncryptedRefreshToken = encryptedRefresh
	u.LastTokenRefresh = time.Now()

	log.Printf("Tokens set successfully")
	return nil
}

func (u *User) GetBattleNetAccessToken() (string, error) {
	if len(u.EncryptedAccessToken) == 0 {
		return "", fmt.Errorf("no access token found")
	}
	decrypted, err := crypto.Decrypt(u.EncryptedAccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt access token: %w", err)
	}
	return string(decrypted), nil
}

// Get the refresh token
func (u *User) GetBattleNetRefreshToken() (string, error) {
	if len(u.EncryptedRefreshToken) == 0 {
		return "", fmt.Errorf("no refresh token found")
	}
	decrypted, err := crypto.Decrypt(u.EncryptedRefreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt refresh token: %w", err)
	}
	return string(decrypted), nil
}

// Battle.net account status
func (u *User) IsBattleNetLinked() bool {
	return u.BattleNetID != nil && u.BattleTag != nil
}

func (u *User) IsTokenExpired() bool {
	return time.Now().After(u.BattleNetExpiresAt)
}

// Scope validation
func (u *User) HasRequiredScopes(requiredScopes []string) bool {
	scopes := make(map[string]bool)
	for _, scope := range u.BattleNetScopes {
		scopes[scope] = true
	}

	for _, required := range requiredScopes {
		if !scopes[required] {
			return false
		}
	}
	return true
}

// GeneratePasswordResetToken generates a new reset token and sets its expiration
func (u *User) GeneratePasswordResetToken() (string, error) {
	// Generate a random 32-byte token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	// convert to hex string
	token := hex.EncodeToString(tokenBytes)

	// set expiration to 1 hour from now
	expires := time.Now().Add(1 * time.Hour)

	u.ResetPasswordToken = &token
	u.ResetPasswordExpires = &expires

	return token, nil
}

// ClearPasswordResetToken clears the reset token and expiration
func (u *User) ClearPasswordResetToken() {
	u.ResetPasswordToken = nil
	u.ResetPasswordExpires = nil
}

// IsPasswordResetTokenValid checks if the reset token is valid and not expired
func (u *User) IsPasswordResetTokenValid() bool {
	if u.ResetPasswordToken == nil || u.ResetPasswordExpires == nil {
		return false
	}
	return time.Now().Before(*u.ResetPasswordExpires)
}

// IsPasswordResetTokenExpired checks if the reset token has expired
func (u *User) IsPasswordResetTokenExpired() bool {
	if u.ResetPasswordExpires == nil {
		return true
	}

	return time.Now().After(*u.ResetPasswordExpires)
}

// ValidatePasswordResetToken checks if the provided token matches and is valid
func (u *User) ValidatePasswordResetToken(token string) bool {
	if u.ResetPasswordToken == nil || *u.ResetPasswordToken != token {
		return false
	}

	return u.IsPasswordResetTokenValid()
}
