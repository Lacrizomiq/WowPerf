package authboss

import (
	"context"
	"fmt"
	"strconv"

	"github.com/volatiletech/authboss/v3"
	"gorm.io/gorm"

	"wowperf/internal/models"
)

// ServerStorer is a storer for the authboss server
type ServerStorer struct {
	db *gorm.DB
}

// NewServerStorer creates a new ServerStorer
func NewServerStorer(db *gorm.DB) *ServerStorer {
	return &ServerStorer{db: db}
}

// Load implements autboss.ServerStorer
func (s *ServerStorer) Load(ctx context.Context, key string) (authboss.User, error) {
	var user models.User

	// Convert the PID (string) into a uuint for my model
	id, err := strconv.ParseUint(key, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID : %w", err)
	}

	if err := s.db.First(&user, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, authboss.ErrUserNotFound
		}
		return nil, err
	}

	return &AuthbossUser{&user}, nil
}

// Save implements authboss.ServerStore
func (s *ServerStorer) Save(ctx context.Context, user authboss.User) error {
	u, ok := user.(*AuthbossUser)
	if !ok {
		return fmt.Errorf("user is not the correct type")
	}

	return s.db.Save(u.User).Error
}

// AuthbossUser is a wrapper to adapter my user model to the authboss interface
type AuthbossUser struct {
	*models.User
}

// GetPID from the authboss.User interface
func (u *AuthbossUser) GetPID() string {
	return strconv.FormatUint(uint64(u.ID), 10)
}

// PutPID from the authboss.User interface
func (u *AuthbossUser) PutPID(pid string) {
	id, _ := strconv.ParseUint(pid, 10, 32)
	u.ID = uint(id)
}

// GetPassword from the autboss.User interface
func (u *AuthbossUser) GetPassword() string {
	return u.Password
}

// PutPassword from the authboss.User interface
func (u *AuthbossUser) PutPassword(password string) {
	u.Password = password
}

// GetEmail from the authboss.User interface
func (u *AuthbossUser) GetEmail() string {
	return u.Email
}

// PutEmail from the authboss.User interface
func (u *AuthbossUser) PutEmail(email string) {
	u.Email = email
}

// GetUsername from the authboss.User interface
func (u *AuthbossUser) GetUsername() string {
	return u.Username
}

// PutUsername from the authboss.User interface
func (u *AuthbossUser) PutUsername(username string) {
	u.Username = username
}
