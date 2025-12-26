package entity

import (
	"time"
)

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	UserID       string    `gorm:"primaryKey;type:char(60)" json:"user_id"`
	Username     string    `gorm:"uniqueIndex;size:50" json:"user_name"`
	Email        string    `gotm:"size:40" json:"email"`
	Password     string    `gotm:"size:60" json:"-"`
	Role         Role      `gotm:"size:10" json:"role"`
	ProfilePict  string    `gorm:"size:255" json:"profile_pict"`
	Department   string    `gorm:"size:40" json:"department"`
	AccessToken  string    `gorm:"size:500" json:"-"`
	RefreshToken string    `gorm:"size:500" json:"-"`
	IsRemote     bool      `json:"remote"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	//Relationship
	Tickets []Ticket `gorm:"foreignKey:UserID" json:"tickets,omitempty"`
}
