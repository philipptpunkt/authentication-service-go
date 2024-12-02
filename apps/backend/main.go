package main

import (
	v1 "backend/handlers/v1"
	"backend/utils"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	// Initialize the database
	utils.InitDatabase()

	// Routes and Request Handlers
	http.HandleFunc("/api/v1/health", v1.HealthHandler)
	http.HandleFunc("/api/v1/auth/register", v1.RegisterHandler)
	http.HandleFunc("/api/v1/auth/login", v1.LoginHandler)

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
