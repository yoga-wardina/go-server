package Config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func init() {
    // Initialize MongoDB client
    var err error
    ctx := context.Background()
    MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        fmt.Println("Failed to create MongoDB client:", err)
        return
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    err = MongoClient.Ping(ctx, nil)
    if err != nil {
        fmt.Println("Failed to connect to MongoDB:", err)
        return
    }
    fmt.Println("Connected to MongoDB")
}
