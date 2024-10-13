package user

import (
	"errors"
	"time"
	"wowperf/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (s *UserService) GetUserProfile(userID uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateEmail updates the user's email, can only be done once every 30 days
func (s *UserService) UpdateEmail(userID uint, newEmail string) error {
	return s.DB.Model(&models.User{}).Where("id = ?", userID).Update("email", newEmail).Error
}

// UpdateUsername updates the user's username
func (s *UserService) UpdateUsername(userID uint, newUsername string) error {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return err
	}

	// Check if 30 days have passed since the last username change
	if time.Since(user.LastUsernameChangeAt) < 30*24*time.Hour {
		return errors.New("30 days must pass before changing username again")
	}

	// Update the user's username and last username change time
	return s.DB.Model(&user).Where("id = ?", userID).Updates(map[string]interface{}{
		"username":                newUsername,
		"last_username_change_at": time.Now(),
	}).Error
}

// ChangePassword changes the user's password
func (s *UserService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.DB.Model(&user).Update("password", string(hashedPassword)).Error
}

func (s *UserService) DeleteAccount(userID uint) error {
	return s.DB.Delete(&models.User{}, userID).Error
}
