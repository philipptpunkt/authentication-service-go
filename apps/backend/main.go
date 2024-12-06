package main

import (
	"backend/backend/generated/auth"
	v1 "backend/handlers/v1"
	"backend/utils"
	"log"
	"net"
	"net/http"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type authServer struct {
	auth.UnimplementedAuthServiceServer
}

func main() {
	_ = godotenv.Load(".env")

	// Initialize the database
	utils.InitDatabase()
	// Initialize Redis
	utils.InitRedis()

	// Routes and Request Handlers
	http.HandleFunc("/api/v1/health", v1.HealthHandler)

	http.HandleFunc("/api/v1/auth/register/link", v1.RegisterWithLinkHandler)
	http.HandleFunc("/api/v1/auth/register/code", v1.RegisterHandler)
	http.HandleFunc("/api/v1/auth/login", v1.LoginHandler)
	http.HandleFunc("/api/v1/auth/refresh", v1.RefreshTokenHandler)
	http.HandleFunc("/api/v1/auth/email-confirmation", v1.ConfirmEmailHandler)
	http.HandleFunc("/api/v1/auth/reset-password", v1.ResetPasswordHandler)
	http.HandleFunc("/api/v1/auth/reset-password/confirm", v1.ResetPasswordConfirmHandler)

	// Routes with Auth Middleware
	http.HandleFunc("/api/v1/auth/delete", v1.AuthMiddleware(v1.DeleteAccountHandler))
	http.HandleFunc("/api/v1/auth/change-password", v1.AuthMiddleware(v1.ChangePasswordHandler))

	// Run gRPC server in a goroutine
	go func() {
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen on port 50051: %v", err)
		}

		grpcServer := grpc.NewServer()
		auth.RegisterAuthServiceServer(grpcServer, &authServer{})

		log.Println("Starting gRPC server on port 50051...")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	log.Println("Starting server on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
