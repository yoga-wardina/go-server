package config

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateToken(email string) (string, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		return "", err
	}

	// Get JWT secret from .env
	jwtSecret := os.Getenv("JWT_SECRET")

	// Define token claims
	claims := jwt.MapClaims{
		"_id": email,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
