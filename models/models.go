package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	USERS_ROLE_ADMIN = "admin"
	USERS_ROLE_USER  = "user"
)

type User struct {
	gorm.Model
	Username        string    `gorm:"not null" json:"user"`
	Email           string    `gorm:"not null;unique" json:"email"`
	Password        string    `gorm:"not null" json:"pass,omitempty"`
	Role            string    `gorm:"default:user" json:"role"`
	Blocked         bool      `gorm:"default:false" json:"blocked"`
	Verified        bool      `gorm:"default:true" json:"verified"`
	LastRequestTime time.Time `gorm:"default:current_timestamp" json:"last_request_time"`
}

type Token struct {
	gorm.Model
	JwtToken string `gorm:"not null" json:"token"`
	Blocked  bool   `gorm:"default:false" json:"blocked"`
}

type Story struct {
	gorm.Model
	Code string `gorm:"not null; unique" json:"code"`
}

type Guest struct {
	gorm.Model
	LastRequestTime time.Time `gorm:"not null" json:"last_request_time"`
}
