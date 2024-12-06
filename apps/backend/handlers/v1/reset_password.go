package v1

import (
	"backend/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	var userID string
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		log.Printf("Reset password requested for non-existent email: %s", req.Email)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "If this email exists, a reset link has been sent."}`))
		return
	} else if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"userID": userID,
		"email":  req.Email,
		"jti":    jti,
	}

	tokenLifetime := 15*time.Minute

	token, err := utils.GenerateGenericJWT(claims, tokenLifetime)
	if err != nil {
		log.Printf("Failed to generate reset token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	redisClient := utils.GetRedisClient()
	err = redisClient.Set(utils.GetRedisContext(), jti, "true", tokenLifetime).Err()
	if err != nil {
		log.Printf("Failed to store JTI in Redis: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	emailSender, err := utils.CreateEmailSender()
	if err != nil {
		log.Printf("Failed to initialize email sender: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	templatePath := "./templates/reset_password.html"

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	// THIS NEEDS THE FRONTEND URL FOR PW RESET
	resetLink := fmt.Sprintf("%s/api/v1/auth/reset-password?token=%s", baseURL, token)

	data := map[string]interface{}{
		"ResetLink": resetLink,
	}
	
	emailBody, err := utils.ParseHtmlTemplate(templatePath, data)
	if err != nil {
			log.Printf("Failed to parse email template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
	}


	err = emailSender.SendEmail(req.Email, "Password Reset Request", emailBody, true)
	if err != nil {
		log.Printf("Failed to send reset email to %s: %v", req.Email, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "If this email exists, a reset link has been sent."}`))
}
