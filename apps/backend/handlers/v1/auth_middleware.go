package v1

import (
	"backend/utils"
	"context"
	"net/http"
)

type AuthContextKey string

const AuthKey AuthContextKey = "auth_subject"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		token := authHeader[7:]// Strip "Bearer " prefix

		subject, err := utils.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AuthKey, subject)
		next(w, r.WithContext(ctx))
	}
}
