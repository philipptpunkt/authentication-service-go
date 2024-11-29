package v1

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type HealthResponse struct {
	Backend  string `json:"backend"`
	Database string `json:"database"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	dbStatus := "Not Connected"
	if checkDatabaseConnection() {
		dbStatus = "Ok"
	}

	response := HealthResponse{
		Backend:  "Ok",
		Database: dbStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func checkDatabaseConnection() bool {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Database connection error:", err)
		return false
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Println("Database ping error:", err)
		return false
	}
	return true
}
