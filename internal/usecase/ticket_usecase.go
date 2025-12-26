package usecase

import (
	"errors"

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

func (uc *ticketUseCase) CreateTicket(userID, title, description string) (*entity.Ticket, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	if description == "" {
		return nil, errors.New("descripton is empty")
	}

	ticket := &entity.Ticket{
		Title:       title,
		Description: description,
		Status:      entity.StatusOpen,
		UserID:      userID,
	}
	return uc.ticketRepo.Create(ticket)

}

func (uc *ticketUseCase) GetUserTickets(userID string) ([]*entity.Ticket, error) {
	return uc.ticketRepo.FindByUserID(userID)
}

func (uc *ticketUseCase) GetAllTicket() ([]*entity.Ticket, error) {
	return uc.ticketRepo.FindAll()
}

func (uc *ticketUseCase) UpdateTicketStatus(ticketID, status, updatedBy string) (*entity.Ticket, error) {
	_, err := uc.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("Ticket tidak ditemukan")
	}

	user, err := uc.userRepo.FIndByUsername(updatedBy)
	if err != nil || user.Role != entity.RoleAdmin {
		return nil, errors.New("unauthorized")
	}

	if err := uc.ticketRepo.UpdateStatus(ticketID, status); err != nil {
		return nil, err
	}

	return uc.ticketRepo.FindByID(ticketID)

}
