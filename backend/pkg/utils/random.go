package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
