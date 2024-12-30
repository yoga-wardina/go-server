package Routes

import (
	"context"
	"encoding/json"
	"go-server/config"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	// Connect to the "users" collection
	collection := config.MongoClient.Database("goDB").Collection("users")

	// Fetch all documents
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var users []bson.M
	if err = cursor.All(ctx, &users); err != nil {
		http.Error(w, "Error decoding users", http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
