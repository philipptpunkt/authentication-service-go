package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

type TokenSubject struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

func GenerateJWT(userID int, email string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	claims := jwt.MapClaims{
		"sub": TokenSubject{
			UserID: userID,
			Email:  email,
		},
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenString string) (TokenSubject, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return TokenSubject{}, fmt.Errorf("JWT_SECRET is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return TokenSubject{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return TokenSubject{}, fmt.Errorf("invalid token claims")
	}

	sub, ok := claims["sub"].(map[string]interface{})
	if !ok {
		return TokenSubject{}, fmt.Errorf("missing or invalid 'sub' claim")
	}

	userID := int(sub["user_id"].(float64))
	email := sub["email"].(string)

	return TokenSubject{
		UserID: userID,
		Email:  email,
	}, nil
}
