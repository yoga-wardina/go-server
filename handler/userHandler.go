package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go-server/config"
	"go-server/models"
	"go-server/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	collection := config.MongoClient.Database("goDB").Collection("users")

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Error hashing password", nil)
		return
	}

	user.ID = primitive.NewObjectID()
	user.Password = hashedPassword
	user.IsAdmin = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LastLogin = time.Time{}
	user.Status = models.StatusOnline
	user.Subscription.PrivilegeLevel = 1
	user.Subscription.IsActive = false
	

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.MongoClient.Database("goDB").Collection("users")
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			utils.JSONResponse(w, http.StatusConflict, "Email already exists", nil)
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, "Error creating user", nil)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, "User created successfully", user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method!= http.MethodPost {
        utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err!= nil {
        utils.JSONResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
        return
    }
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	

    collection := config.MongoClient.Database("goDB").Collection("users")
	result := collection.FindOne(ctx, bson.D{{Key: "email", Value: user.Email}})
	if result.Err()!= nil {
        utils.JSONResponse(w, http.StatusNotFound, "User not found", nil)
        return
    }
	
	var foundUser models.User
	
	if err := result.Decode(&foundUser); err!= nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Error decoding user", nil)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err!= nil {
        utils.JSONResponse(w, http.StatusUnauthorized, "nvalid credentials", nil)
        return
    }
	
	token, err := config.GenerateToken(foundUser.Email)
	if err!= nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Error generating token", nil)
		return
	}

	userJSON, err := json.Marshal(foundUser)
	if err!= nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Error marshaling user", nil)
        return
	}
    config.RedisClient.Set(ctx, fmt.Sprintf( "token:%s", foundUser.Email), userJSON, 24*time.Hour)
	utils.JSONResponse(w, http.StatusOK, "Login successful", map[string]string{"token": token})
}	