package utils

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDatabase() {
	connStr := os.Getenv("DATABASE_URL")
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Database connection established")

	ensureTablesExist()
}

func GetDB() *sql.DB {
	return db
}

func ensureTablesExist() {
	// Ensure users table exists
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		verified BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	} else {
		log.Println("Users table ensured")
	}

	// Ensure refresh_tokens table exists
	createRefreshTokensTable := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id SERIAL PRIMARY KEY,
		token TEXT NOT NULL UNIQUE,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createRefreshTokensTable)
	if err != nil {
		log.Fatalf("Failed to create refresh_tokens table: %v", err)
	} else {
		log.Println("Refresh tokens table ensured")
	}

	// Ensure email_confirmation_tokens table exists
	createConfirmationTokensTable := `
	CREATE TABLE IF NOT EXISTS email_confirmation_tokens (
		id SERIAL PRIMARY KEY,
		token TEXT NOT NULL UNIQUE,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createConfirmationTokensTable)
	if err != nil {
		log.Fatalf("Failed to create email_confirmation_tokens table: %v", err)
	} else {
		log.Println("Email confirmation tokens table ensured")
	}
}
