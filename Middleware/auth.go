package Middleware

import (
	"context"
	"fmt"
	"go-server/config"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func DecodeJWT(tokenString string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idRaw, exists := claims["_id"]
		if !exists {
			return "", fmt.Errorf("_id not found in token claims")
		}
		id, ok := idRaw.(string)
		if !ok {
			return "", fmt.Errorf("_id is not a string in token claims")
		}

		return id, nil
	}

	return "", fmt.Errorf("invalid token claims")
}


func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		id, err := DecodeJWT(token)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		redisKey := fmt.Sprintf("token:%s", id)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		val, err := config.RedisClient.Get(ctx, redisKey).Result()
		if err != nil {
			http.Error(w, "Unauthorized: token not found in Redis", http.StatusUnauthorized)
			return
		}
		fmt.Println("Token found in Redis:", val) 
		next.ServeHTTP(w, r)
	})
}
