package usecase

import (
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
)

type AuthUseCase interface {
	Register(username, email, password string) (*entity.User, error)
	Login(username, password string) (string, error)
}

type TicketUseCase interface {
	CreateTicket(userID, title, descriptiom string) (*entity.Ticket, error)
	GetUserTickets(userID string) ([]*entity.Ticket, error)
	GetAllTicket() ([]*entity.Ticket, error)
	UpdateTicketStatus(ticketID, status, updatedBy string) (*entity.Ticket, error)
}
