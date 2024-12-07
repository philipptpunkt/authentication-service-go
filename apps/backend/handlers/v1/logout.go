package v1

import (
	"backend/utils"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	subject, ok := r.Context().Value(AuthKey).(utils.TokenSubject)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := subject.UserID

	db := utils.GetDB()
	_, err := db.Exec("DELETE FROM refresh_tokens WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logged out successfully"}`))
}
