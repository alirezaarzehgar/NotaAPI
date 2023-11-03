package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	USERS_ROLE_ADMIN = "admin"
	USERS_ROLE_USER  = "user"

	STORY_TYPE_EXPLORE = "explore"
	STORY_TYPE_NORMAL  = "normal"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null" json:"user"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"pass,omitempty"`
	Role     string `gorm:"default:user" json:"role"`
	Blocked  bool   `gorm:"default:false" json:"blocked"`
	Verified bool   `gorm:"default:true" json:"verified"`
}

type Token struct {
	gorm.Model
	UserID          uint      `json:"user_id"`
	JwtToken        string    `gorm:"not null" json:"token"`
	Blocked         bool      `gorm:"default:false" json:"blocked"`
	ScreenHeight    uint      `json:"screen_height"`
	ScreenWidth     uint      `json:"screen_width"`
	Resolution      uint      `json:"resolution"`
	DeviceType      string    `json:"device_type"`
	Version         uint      `json:"version"`
	LastRequestTime time.Time `gorm:"default:current_timestamp" json:"last_request_time"`

	Notification bool   `gorm:"default:false" json:"notification"`
	GCMToken     string `json:"gcm_token"`

	SavedStories []Story `gorm:"many2many:token_stories" json:"saved_stories"`
}

type Story struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint           `json:"user_id"`
	Code      string         `gorm:"not null; unique; primaryKey" json:"code"`
	Type      string         `gorm:"default:normal" json:"type"`
	Name      string         `gorm:"not null" json:"name"`
	From      time.Time      `json:"from"`
	To        time.Time      `json:"to"`

	FinalImageUrl     string `gorm:"not null" json:"final_image"`
	BackgroundUrl     string `gorm:"not null" json:"background_url"`
	MainBackgroundUrl string `json:"main_background_url"`
	BackgroundColor   string `gorm:"not null" json:"background_color"`

	LogoUrl       string `json:"logo_url"`
	LogoHeight    uint   `json:"logo_height"`
	LogoWidth     uint   `json:"logo_width"`
	LogoXLocation uint   `json:"logo_x_location"`
	LogoYLocation uint   `json:"logo_y_location"`
	LogoYScale    uint   `json:"logo_scale"`

	Text         string `gorm:"not null" json:"text"`
	TextPosition string `gorm:"not null" json:"text_position"`
	TextSize     uint   `gorm:"not null" json:"text_size"`
	TextColor    uint   `gorm:"not null" json:"text_color"`
	IsShared     bool   `gorm:"not null" json:"is_shared"`

	AttachedWebpage string `json:"attached_webpage"`
	AttachedFileUrl string `json:"attached_file_url"`

	IsPublic bool `gorm:"default:false" json:"is_public"`

	Tokens []Token `gorm:"many2many:token_stories" json:"-"`
}
