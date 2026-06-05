package usecase

import (
	"mime/multipart"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/dto"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
)

type AuthUseCase interface {
	Register(username, email, password, department string, role entity.Role, isRemote bool) (*entity.User, error)
	Login(username, password string) (string, error)
}

type TicketUseCase interface {
	CreateTicket(userID, title, description, image string, imageFile multipart.File, imageFilename string, imageSize int64, imageContentType string) (*entity.Ticket, error)
	GetUserTickets(userID string) ([]*entity.Ticket, error)
	GetAllTicket() ([]*entity.Ticket, error)
	UpdateTicketStatus(ticketID, updatedBy string, status entity.TicketStatus) (*entity.Ticket, error)
	GetSummaryStat() (*dto.DashboardSummary, error)
}
