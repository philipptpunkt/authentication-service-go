package v1

import (
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
)

type ConfirmEmailRequest struct {
	Token string `json:"token"`
}

func ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decoding email confirmation request:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	// Validate the token and fetch the associated user
	var userID int
	err := db.QueryRow(`
		SELECT user_id 
		FROM email_confirmation_tokens 
		WHERE token = $1 AND expires_at > NOW()`, req.Token).Scan(&userID)
	if err != nil {
		log.Printf("Error validating token: %v\n", err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Mark the user as verified
	_, err = db.Exec(`UPDATE users SET verified = TRUE WHERE id = $1`, userID)
	if err != nil {
		log.Printf("Error verifying user: %v\n", err)
		http.Error(w, "Failed to verify account", http.StatusInternalServerError)
		return
	}

	// Delete the used token
	_, err = db.Exec(`DELETE FROM email_confirmation_tokens WHERE token = $1`, req.Token)
	if err != nil {
		log.Printf("Error deleting token: %v\n", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Email confirmed successfully"}`))
}
