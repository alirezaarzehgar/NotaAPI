package models

const (
	USERS_ROLE_ADMIN = "admin"
	USERS_ROLE_USER  = "user"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"not null" json:"user"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"password,omitempty"`
	Role     string `gorm:"default:user" json:"role"`
	Blocked  bool   `gorm:"default:false" json:"blocked"`
}

type Token struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	JwtToken string `gorm:"not null" json:"token"`
	Blocked  bool   `gorm:"default:false" json:"blocked"`
}

type Story struct {
	Code string `gorm:"primaryKey" json:"code"`
}

type Guest struct {
	ID uint `gorm:"primaryKey" json:"id"`
}
