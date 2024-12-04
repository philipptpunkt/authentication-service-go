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

	http.HandleFunc("/api/v1/auth/register/link", v1.RegisterWithLinkHandler)
	http.HandleFunc("/api/v1/auth/register/code", v1.RegisterHandler)
	http.HandleFunc("/api/v1/auth/login", v1.LoginHandler)
	http.HandleFunc("/api/v1/auth/refresh", v1.RefreshTokenHandler)
	http.HandleFunc("/api/v1/auth/email-confirmation", v1.ConfirmEmailHandler)

	// Routes with Auth Middleware
	http.HandleFunc("/api/v1/auth/delete", v1.AuthMiddleware(v1.DeleteAccountHandler))

	log.Println("Starting server on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
