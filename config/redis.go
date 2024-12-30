package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
)

var RedisClient *redis.Client

func init() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUser := os.Getenv("REDIS_USER")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_CACHE_DB"))

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Username: redisUser,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test Redis connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
	} else {
		fmt.Println("Connected to Redis")
	}
}
