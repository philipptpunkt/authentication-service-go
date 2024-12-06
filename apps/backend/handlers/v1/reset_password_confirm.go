package v1

import (
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type ResetPasswordConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func ResetPasswordConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ResetPasswordConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(req.Token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(utils.GetJWTSecret()), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	userID := claims["userID"].(string)
	jti := claims["jti"].(string)

	// Check if token (jti) exists in Redis
	redisClient := utils.GetRedisClient()
	exists, err := redisClient.Exists(utils.GetRedisContext(), jti).Result()
	if err != nil {
		log.Printf("Redis error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		http.Error(w, "Token is invalid or has already been used", http.StatusUnauthorized)
		return
	}

	db := utils.GetDB()
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("Failed to hash new password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE users SET password = $1 WHERE id = $2", hashedPassword, userID)
	if err != nil {
		log.Printf("Failed to update password in the database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Mark token (jti) as used by deleting it from Redis
	err = redisClient.Del(utils.GetRedisContext(), jti).Err()
	if err != nil {
		log.Printf("Failed to delete token from Redis: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Password reset successful"}`))
}
