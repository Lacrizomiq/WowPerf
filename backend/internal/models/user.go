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

	// Battle.net specific fields - NOTE: Battle.net doesn't provide refresh tokens
	BattleNetID          *string   `gorm:"uniqueIndex" json:"battle_net_id"`
	BattleTag            *string   `gorm:"uniqueIndex" json:"battle_tag"`
	EncryptedAccessToken []byte    `gorm:"type:bytea" json:"-"`
	BattleNetTokenType   string    `gorm:"type:varchar(50)" json:"-"`
	BattleNetExpiresAt   time.Time `json:"-"`
	BattleNetScopes      []string  `gorm:"type:text[]" json:"-"`
	LastTokenRefresh     time.Time `json:"-"`

	// Google OAuth fields
	GoogleID    *string `gorm:"uniqueIndex" json:"google_id"`
	GoogleEmail *string `json:"google_email"`

	// Character relationship
	FavoriteCharacterID *uint           `json:"favorite_character_id"`
	Characters          []UserCharacter `gorm:"foreignKey:UserID" json:"characters,omitempty"`
}

// UserCreate is the struct for creating a new user
type UserCreate struct {
	Username     string `json:"username" validate:"required,min=3,max=14"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=8"`
	CaptchaToken string `json:"captcha_token"`
}

// Token management methods
func (u *User) SetBattleNetTokens(accessToken, refreshToken string) error {
	log.Printf("Starting SetBattleNetTokens: access_token_length=%d", len(accessToken))

	if accessToken == "" {
		return fmt.Errorf("access token is empty")
	}

	// Note: Battle.net doesn't provide refresh tokens, so we ignore the refreshToken parameter
	if refreshToken != "" {
		log.Printf("Warning: Battle.net doesn't provide refresh tokens, ignoring refresh token parameter")
	}

	// Encrypt access token
	encryptedAccess, err := crypto.Encrypt([]byte(accessToken))
	if err != nil {
		log.Printf("Failed to encrypt access token: %v", err)
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}
	log.Printf("Access token encrypted successfully: length=%d", len(encryptedAccess))

	u.EncryptedAccessToken = encryptedAccess
	u.LastTokenRefresh = time.Now()

	log.Printf("Battle.net access token set successfully")
	return nil
}

// SetBattleNetAccessToken sets only the access token (simplified version for future use)
func (u *User) SetBattleNetAccessToken(accessToken string) error {
	log.Printf("Setting Battle.net access token: length=%d", len(accessToken))

	if accessToken == "" {
		return fmt.Errorf("access token is empty")
	}

	// Encrypt access token
	encryptedAccess, err := crypto.Encrypt([]byte(accessToken))
	if err != nil {
		log.Printf("Failed to encrypt access token: %v", err)
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}

	u.EncryptedAccessToken = encryptedAccess
	u.LastTokenRefresh = time.Now()

	log.Printf("Battle.net access token encrypted and set successfully")
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

// DEPRECATED: Battle.net doesn't provide refresh tokens
// This method is kept for backward compatibility but will always return an error
func (u *User) GetBattleNetRefreshToken() (string, error) {
	return "", fmt.Errorf("battle.net doesn't provide refresh tokens")
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

// === GOOGLE OAUTH METHODS ===

// IsGoogleLinked vérifie si le compte Google est lié à un utilisateur
func (u *User) IsGoogleLinked() bool {
	return u.GoogleID != nil
}

// CanLoginWithGoogle vérifie si l'utilisateur peut se connecter avec Google
func (u *User) CanLoginWithGoogle() bool {
	return u.IsGoogleLinked() && u.GoogleEmail != nil
}

// LinkGoogleAccount lie un compte Google à l'utilisateur
func (u *User) LinkGoogleAccount(googleID, googleEmail string) {
	u.GoogleID = &googleID
	u.GoogleEmail = &googleEmail
}

// UnlinkGoogleAccount supprime le lien Google de l'utilisateur
func (u *User) UnlinkGoogleAccount() {
	u.GoogleID = nil
	u.GoogleEmail = nil
}

// HasMultipleAuthMethods vérifie si l'utilisateur a plusieurs méthodes d'auth
func (u *User) HasMultipleAuthMethods() bool {
	methods := 0
	if u.Password != "" {
		methods++
	}
	if u.IsGoogleLinked() {
		methods++
	}
	if u.IsBattleNetLinked() {
		methods++
	}
	return methods > 1
}
