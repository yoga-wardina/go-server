package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Status string

const (
	StatusOnline       Status = "Online"
	StatusBusy         Status = "Busy"
	StatusInvisible    Status = "Invisible"
	StatusDoNotDisturb Status = "DoNotDisturb"
)

type Subscription struct {
	PrivilegeLevel int       `bson:"privilegeLevel"` // Corrected spelling
	SubName        string    `bson:"subName"`
	ExpiryDate     time.Time `bson:"expiryDate"`
	StartDate      time.Time `bson:"startDate"`
	IsActive       bool      `bson:"isActive"`
}

type User struct {
	ID         		primitive.ObjectID	`bson:"_id,omitempty"`
	Email      		string            	`bson:"email"`
	Password   		string            	`bson:"password"`
	Name       		string            	`bson:"name"`
	ProfilePic 		string            	`bson:"profilePic"`
	BannerPic 		string            	`bson:"bannerPic"`
	Alias      		string            	`bson:"alias"`
	Status     		Status            	`bson:"status"`
	Subscription 	Subscription		`bson:"subscription"`
	IsAdmin    		bool              	`bson:"isAdmin"`
	CreatedAt  		time.Time         	`bson:"createdAt"`
	UpdatedAt  		time.Time          	`bson:"updatedAt"`
	LastLogin  		time.Time          	`bson:"lastLogin"`
}


func IsValidStatus(status Status) bool {
	switch status {
	case StatusOnline, StatusBusy, StatusInvisible, StatusDoNotDisturb:
		return true
	default:
		return false
	}
}