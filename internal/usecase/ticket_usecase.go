package usecase

import (
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/repository"
)

type ticketUseCase struct {
	ticketRepo repository.TicketRepositoy
	userRepo   repository.UserRepository
}

func NewTicketUseCase(ticketRepo repository.TicketRepositoy, userRepo repository.UserRepository) TicketUseCase {
	return &ticketUseCase{
		ticketRepo: ticketRepo,
		userRepo:   userRepo,
	}
}

func (uc *ticketUseCase) CreateTicket(userID, title, descriptiom string) (*entity.Ticket, error) {
	return nil, nil
}

func (uc *ticketUseCase) GetUserTickets(userID string) ([]*entity.Ticket, error) {
	return nil, nil
}

func (uc *ticketUseCase) GetAllTicket() ([]*entity.Ticket, error) {
	return nil, nil
}

func (uc *ticketUseCase) UpdateTicketStatus(ticketID, status, updatedBy string) (*entity.Ticket, error) {
	return nil, nil
}
