package v1

import (
	"backend/utils"
	"log"
	"net/http"
)

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user info from the context (set by Middleware)
	subject, ok := r.Context().Value(AuthKey).(utils.TokenSubject)
	if !ok {
		http.Error(w, "Failed to retrieve user info from context", http.StatusUnauthorized)
		return
	}

	db := utils.GetDB()
	_, err := db.Exec("DELETE FROM users WHERE id = $1 AND email = $2", subject.UserID, subject.Email)
	if err != nil {
		log.Printf("Error deleting user: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Account deleted successfully"}`))
}

