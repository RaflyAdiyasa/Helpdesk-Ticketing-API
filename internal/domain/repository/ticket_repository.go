package repository

import (
	"github.com/RaflyAdiyasa/Helpdest-Ticketing-API/internal/domain/entity"
)

type TicketRepositoy interface {
	Create(ticket *entity.Ticket) (*entity.Ticket, error)
	FindByID(id string) (*entity.Ticket, error)
	FindByUserID(id string) ([]*entity.Ticket, error)
	FindAll() ([]*entity.Ticket, error)
	Update(ticket *entity.Ticket) error
	UpdateStatus(ticketID, status string) error
	Delete(id string) error
}
