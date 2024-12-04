package v1

import (
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type DeleteAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decoding delete account request:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	db := utils.GetDB()
	err := db.QueryRow("SELECT password FROM users WHERE email = $1", req.Email).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Error querying user: %v\n", err)
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE email = $1", req.Email)
	if err != nil {
		log.Printf("Error deleting user: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Account deleted successfully"}`))
}
