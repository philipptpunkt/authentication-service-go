package v1

import (
	"backend/backend/generated/auth"
	"backend/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request into Protobuf-generated type
	var req *auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, statusCode, err := LoginHandlerLogic(req)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func LoginHandlerLogic(req *auth.LoginRequest) (*auth.LoginResponse, int, error) {
	db := utils.GetDB()

	var userID int
	var hashedPassword string
	err := db.QueryRow("SELECT id, password FROM users WHERE email = $1", req.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, http.StatusUnauthorized, errors.New("invalid email or password")
		}
		log.Printf("Error querying user: %v\n", err)
		return nil, http.StatusInternalServerError, errors.New("server error")
	}

	if err := utils.ComparePasswords(hashedPassword, req.Password); err != nil {
		return nil, http.StatusUnauthorized, errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(userID, req.Email)
	if err != nil {
		log.Printf("Error generating token: %v\n", err)
		return nil, http.StatusInternalServerError, errors.New("server error")
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		log.Printf("Error generating refresh token: %v\n", err)
		return nil, http.StatusInternalServerError, errors.New("server error")
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days valid
	_, err = db.Exec("INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)",
		refreshToken, userID, expiresAt)
	if err != nil {
		log.Printf("Error storing refresh token: %v\n", err)
		return nil, http.StatusInternalServerError, errors.New("server error")
	}

	return &auth.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, http.StatusOK, nil
}
