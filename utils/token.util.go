package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// GenerateSecureToken generates a secure random string token
func GenerateSecureToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("Error generating random token:", err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}