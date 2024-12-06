package v1

import (
	"backend/utils"
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
	Redis    string `json:"redis"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	dbStatus := "Not Connected"
	redisStatus := "Not Connected"

	if checkDatabaseConnection() {
		dbStatus = "Ok"
	}

	if checkRedisConnection() {
		redisStatus = "Ok"
	}

	response := HealthResponse{
		Backend:  "Ok",
		Database: dbStatus,
		Redis: redisStatus,
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

func checkRedisConnection() bool {
	redisClient := utils.GetRedisClient()
	if redisClient == nil {
		log.Println("Redis client is not initialized")
		return false
	}

	if _, err := redisClient.Ping(utils.GetRedisContext()).Result(); err != nil {
		log.Println("Redis ping error:", err)
		return false
	}
	return true
}
