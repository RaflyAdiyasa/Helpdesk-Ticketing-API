package repository

import (
	"github.com/RaflyAdiyasa/Helpdest-Ticketing-API/internal/domain/entity"
)

type TicketRepositoy interface {
	Create(ticket *entity.Ticket) error
	FindByID(id string) (entity.Ticket, error)
	FindByUserID(id string) ([]entity.Ticket, error)
	FindAll() ([]entity.Ticket, error)
	Update(ticket *entity.Ticket) error
	UpdateStatus(ticket *entity.Ticket, status string) error
	Delete(id string) error
}
