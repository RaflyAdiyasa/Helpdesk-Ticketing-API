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
	UserID      string       `gorm:"type:char(60)" json:"user_id"`
	Title       string       `gorm:"size:30" json:"title"`
	Description string       `gorm:"size:300" json:"description"`
	Status      TicketStatus `gorm:"size:10" json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletetAt   time.Time    `json:"deleted_at"`

	//Relationship
	Owner User `gorm:"foreignKey:UserID" json:"owner"`

	//TO-DO : Images
}
