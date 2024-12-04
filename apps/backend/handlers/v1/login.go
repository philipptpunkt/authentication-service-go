package v1

import (
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error decoding login request:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var userID int
	var hashedPassword string
	db := utils.GetDB()
	err := db.QueryRow("SELECT id, password FROM users WHERE email = $1", req.Email).Scan(&userID, &hashedPassword)
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

	token, err := utils.GenerateJWT(userID, req.Email)
	if err != nil {
		log.Printf("Error generating token: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
