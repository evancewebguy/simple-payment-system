package user

import (
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	FullName    string `json:"full_name" gorm:"size:255;not null"`
	Email       string `json:"email" gorm:"size:100;uniqueIndex;not null"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"-" gorm:"size:255;not null"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`
	IsVerified  bool   `json:"is_verified" gorm:"default:false"`
}
