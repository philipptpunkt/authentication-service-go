package v1

import (
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decoding refresh request:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the refresh token and fetch associated user data
	db := utils.GetDB()
	var userID int
	var email string
	var expiresAt time.Time
	err := db.QueryRow(`
		SELECT rt.user_id, u.email, rt.expires_at 
		FROM refresh_tokens rt 
		JOIN users u ON rt.user_id = u.id 
		WHERE rt.token = $1`, req.RefreshToken).Scan(&userID, &email, &expiresAt)
	if err != nil {
		log.Printf("Error validating refresh token: %v\n", err)
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// Check if the refresh token has expired
	if time.Now().After(expiresAt) {
		http.Error(w, "Refresh token has expired", http.StatusUnauthorized)
		return
	}

	accessToken, err := utils.GenerateJWT(userID, email)
	if err != nil {
		log.Printf("Error generating access token: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

