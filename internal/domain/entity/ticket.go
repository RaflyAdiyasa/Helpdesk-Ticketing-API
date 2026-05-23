package entity

import (
	"time"
)

type TicketStatus string

const (
	StatusOpen       TicketStatus = "OPEN"
	StatusInProgress TicketStatus = "IN_PROGRESS"
	StatusDone       TicketStatus = "DONE"
)

type Ticket struct {
	TicketID    string       `gorm:"primaryKey;type:char(60)" json:"ticket_id"`
	UserID      string       `gorm:"type:char(60);not null;index" json:"user_id"`
	Title       string       `gorm:"size:80" json:"title"`
	Description string       `gorm:"type:text" json:"description"`
	Image       string       `gorm:"size:255" json:"image"`
	Status      TicketStatus `gorm:"size:20" json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at"`

	// Relationship
	Owner *User `gorm:"foreignKey:UserID;references:UserID" json:"owner"`
}

func (Ticket) TableOptions() string {
	return "ENGINE=InnoDB"
}
