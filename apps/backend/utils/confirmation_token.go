package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

func GenerateConfirmationToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func StoreConfirmationToken(userID int, token string) error {
	db := GetDB()
	expiresAt := time.Now().Add(24 * time.Hour) // Token valid 24 hours
	_, err := db.Exec(`
		INSERT INTO email_confirmation_tokens (token, user_id, expires_at) 
		VALUES ($1, $2, $3)`, token, userID, expiresAt)
	return err
}
