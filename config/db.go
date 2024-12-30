package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}


	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_CONNECTION") 

	mongoURI := host

	clientOptions := options.Client().ApplyURI(mongoURI).SetAuth(options.Credential{
		Username: username,
		Password: password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	MongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}

	if err = MongoClient.Ping(ctx, nil); err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}

	fmt.Println("Connected to MongoDB with authentication")
}
