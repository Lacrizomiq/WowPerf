package models

import (
	"time"

	"wowperf/pkg/crypto"

	"gorm.io/gorm"
)

// User is the struct for the user model
type User struct {
	gorm.Model
	ID                    uint      `gorm:"primaryKey" json:"id"`
	Username              string    `gorm:"uniqueIndex;not null" json:"username" validate:"required,min=3,max=50"`
	Email                 string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password              string    `gorm:"not null" json:"-" validate:"required,strongpassword"`
	BattleNetID           *int      `gorm:"uniqueIndex" json:"battle_net_id"`
	BattleTag             *string   `gorm:"uniqueIndex" json:"battle_tag"`
	EncryptedToken        []byte    `gorm:"type:bytea" json:"-"`
	BattleNetRefreshToken string    `gorm:"type:text" json:"-"`
	BattleNetTokenType    string    `gorm:"type:varchar(50)" json:"-"`
	BattleNetExpiresAt    time.Time `json:"battle_net_expires_at"`
	BattleNetScopes       []string  `gorm:"type:text[]" json:"-"`
	LastUsernameChangeAt  time.Time `json:"last_username_change_at"`
}

// UserCreate is the struct for creating a new user
type UserCreate struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,containsany=!@#$%^&*()_+"`
}

func (u *User) SetBattleNetToken(token string) error {
	encrypted, err := crypto.Encrypt([]byte(token))
	if err != nil {
		return err
	}
	u.EncryptedToken = encrypted
	return nil
}

func (u *User) GetBattleNetToken() (string, error) {
	if len(u.EncryptedToken) == 0 {
		return "", nil
	}
	decrypted, err := crypto.Decrypt(u.EncryptedToken)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// HasScope checks if the user has a specific scope
func (u *User) HasScope(scope string) bool {
	for _, s := range u.BattleNetScopes {
		if s == scope {
			return true
		}
	}
	return false
}

// HasRequiredScopes checks if the user has all the required scopes
func (u *User) HasRequiredScopes(requiredScopes []string) bool {
	scopeMap := make(map[string]bool)
	for _, scope := range u.BattleNetScopes {
		scopeMap[scope] = true
	}

	for _, requiredScope := range requiredScopes {
		if !scopeMap[requiredScope] {
			return false
		}
	}
	return true
}
