package v1

import (
	"backend/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func RegisterWithLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decoding registration request:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	db := utils.GetDB()
	var userID int
	err = db.QueryRow(
		`INSERT INTO users (email, password, verified) VALUES ($1, $2, FALSE) RETURNING id`,
		req.Email, hashedPassword,
	).Scan(&userID)
	if err != nil {
		log.Printf("Error inserting user: %v\n", err)
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}

	confirmationToken, err := utils.GenerateConfirmationToken()
	if err != nil {
		log.Printf("Error generating confirmation token: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = utils.StoreConfirmationToken(userID, confirmationToken)
	if err != nil {
		log.Printf("Error storing confirmation token: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	templatePath := "./templates/email_address_confirmation.html"

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	confirmationLink := fmt.Sprintf("%s/auth/email-confirmation?token=%s", baseURL, confirmationToken)

	data := map[string]interface{}{
		"ConfirmationLink": confirmationLink,
	}

	emailBody, err := utils.ParseHtmlTemplate(templatePath, data)
	if err != nil {
		log.Printf("Error parsing email template: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	emailSender, err := utils.CreateEmailSender()
	if err != nil {
		log.Printf("Error initializing email sender: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = emailSender.SendEmail(req.Email, "Confirm Your Email", emailBody, true)
	if err != nil {
		log.Printf("Error sending confirmation email: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Registration successful. Please check your email to confirm your account."}`))
}
