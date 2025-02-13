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

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
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
    db := config.PostgresConn
    if db == nil {
        utils.JSONResponse(w, http.StatusInternalServerError, "Database connection is nil", nil)
        return
    }

    if r.Method != http.MethodPost {
        utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

    var request models.User
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        utils.JSONResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
        return
    }

    hashedPassword, err := utils.HashPassword(request.Password)
    if err != nil {
        utils.JSONResponse(w, http.StatusInternalServerError, "Error hashing password", nil)
        return
    }

    user := models.User{
        Email:     request.Email,
        Password:  hashedPassword,
        Name:      request.Name,
        IsAdmin:   false,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        LastLogin: time.Time{}, 
        Status:    models.StatusOnline,
        Subscription: models.Subscription{
            PrivilegeLevel: 1,
            IsActive:       false,
            StartDate:      time.Now(),
            ExpiryDate:     time.Now().AddDate(1, 0, 0), 
        },
    }

    if err := db.Create(&user).Error; err != nil {
        utils.JSONResponse(w, http.StatusInternalServerError, "Error creating user", nil)
        return
    }

    utils.JSONResponse(w, http.StatusCreated, "Sign up successful", user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db := config.PostgresConn

	if r.Method!= http.MethodPost {
        utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
        return
    }

	if r.Body == nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Input can't be empty", nil)
        return
	}

    var request models.User
    if err := json.NewDecoder(r.Body).Decode(&request); err!= nil {
        utils.JSONResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
        return
    }

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

    var user models.User
	result := db.WithContext(ctx).Table("users").Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
        utils.JSONResponse(w, http.StatusUnauthorized, "User not found", request)
        return
    }
	
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err!= nil {
        utils.JSONResponse(w, http.StatusUnauthorized, "Please check your email and passwords credentials", nil)
        return
    }
	
    initToken := uuid.New().String()

	token, err := config.GenerateToken(initToken)
	if err!= nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Error generating token", nil)
		return
	}

	userJSON, err := json.Marshal(user)
	if err!= nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Error marshaling user", nil)
        return
	}
	
    config.RedisClient.Set(ctx, fmt.Sprintf( "token:%s", initToken), userJSON, 24*time.Hour)
	utils.JSONResponse(w, http.StatusOK, "Login successful", map[string]string{"token": token})
}	
