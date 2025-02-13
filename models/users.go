package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Status ENUM
type Status string

const (
	StatusOnline       Status = "Online"
	StatusBusy         Status = "Busy"
	StatusInvisible    Status = "Invisible"
	StatusDoNotDisturb Status = "DoNotDisturb"
)

// Subscription Model
type Subscription struct {
	ID             uint      `gorm:"primaryKey"`
	PrivilegeLevel int       `gorm:"column:privilege_level"`
	SubName        string    `gorm:"column:sub_name"`
	ExpiryDate     time.Time `gorm:"column:expiry_date"`
	StartDate      time.Time `gorm:"column:start_date"`
	IsActive       bool      `gorm:"column:is_active"`
}

// User Model
type User struct {
	ID           uint         `gorm:"primaryKey"`
	UID           string       `gorm:"type:char(36);unique;not null"` 
	Email        string       `gorm:"unique;not null"`
	Password     string       `gorm:"not null"`
	Name         string       `gorm:"not null"`
	ProfilePic   string       `gorm:"column:profile_pic"`
	BannerPic    string       `gorm:"column:banner_pic"`
	Alias        string       `gorm:"column:alias"`
	Status       Status       `gorm:"type:status_enum;default:'Online'"`
	Subscription Subscription `gorm:"foreignKey:SubscriptionID"`
	SubscriptionID uint       
	IsAdmin      bool         `gorm:"default:false"`
	CreatedAt    time.Time    `gorm:"autoCreateTime"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime"`
	LastLogin    time.Time
}

// IsValidStatus checks if a status is valid
func IsValidStatus(status Status) bool {
	switch status {
	case StatusOnline, StatusBusy, StatusInvisible, StatusDoNotDisturb:
		return true
	default:
		return false
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UID = uuid.New().String()
	return nil
}

func MigrateDB(db *gorm.DB) {
	// Create ENUM type for status
	db.Exec("CREATE TYPE status_enum AS ENUM ('Online', 'Busy', 'Invisible', 'DoNotDisturb');")

	// AutoMigrate tables
	db.AutoMigrate(&Subscription{}, &User{})
}
