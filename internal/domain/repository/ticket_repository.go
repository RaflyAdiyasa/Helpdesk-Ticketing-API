package repository

import (
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
)

type TicketRepositoy interface {
	Create(ticket *entity.Ticket) (*entity.Ticket, error)
	FindByID(id string) (*entity.Ticket, error)
	FindByUserID(id string) ([]*entity.Ticket, error)
	FindAll() ([]*entity.Ticket, error)
	Update(ticket *entity.Ticket) error
	UpdateStatus(ticketID string, status entity.TicketStatus) error
	Delete(id string) error
}
