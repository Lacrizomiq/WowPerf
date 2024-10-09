package models

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()

	// Add custom validation for password strength
	Validate.RegisterValidation("strongpassword", validateStrongPassword)
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// At least 8 characters
	if len(password) < 8 {
		return false
	}

	// At least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}

	// At least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}

	// At least one number
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	}

	// At least one special character
	if !regexp.MustCompile(`[!@#$%^&*()_+]`).MatchString(password) {
		return false
	}

	return true
}
