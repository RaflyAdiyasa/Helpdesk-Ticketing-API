package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/dto"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/repository"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/pkg/utils"
)

type ticketUseCase struct {
	ticketRepo repository.TicketRepositoy
	userRepo   repository.UserRepository
	fileRepo   repository.FileRepository
	cacheRepo  repository.CacheRepository
}

func NewTicketUseCase(ticketRepo repository.TicketRepositoy, userRepo repository.UserRepository, fileRepo repository.FileRepository, cacheRepo repository.CacheRepository) TicketUseCase {
	return &ticketUseCase{
		ticketRepo: ticketRepo,
		userRepo:   userRepo,
		fileRepo:   fileRepo,
		cacheRepo:  cacheRepo,
	}
}

func (uc *ticketUseCase) CreateTicket(userID, title, description, image string, imageFile multipart.File, imageFilename string, imageSize int64, imageContentType string) (*entity.Ticket, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	if description == "" {
		return nil, errors.New("descripton is empty")
	}

	ticketID := utils.GenerateTicketID()

	if imageFile != nil {
		if uc.fileRepo == nil {
			return nil, errors.New("file repository is not configured")
		}

		objectName := buildTicketImageObjectName(ticketID, imageFilename)
		uploadedImage, err := uc.fileRepo.Upload(imageFile, objectName, imageSize, imageContentType)
		if err != nil {
			return nil, fmt.Errorf("failed to upload image: %w", err)
		}
		image = uploadedImage
	}

	ticket := &entity.Ticket{
		TicketID:    ticketID,
		Title:       title,
		Description: description,
		Status:      entity.StatusOpen,
		UserID:      userID,
		Image:       image,
	}
	ctx := context.Background()
	_ = uc.cacheRepo.Delete(ctx, "dashboard:overview")
	return uc.ticketRepo.Create(ticket)

}

func buildTicketImageObjectName(ticketID, filename string) string {
	normalizedFilename := strings.ReplaceAll(filename, "\\", "/")
	ext := strings.ToLower(path.Ext(normalizedFilename))
	if ext == "" {
		ext = ".bin"
	}

	return fmt.Sprintf("tickets/%s/%s%s", ticketID, utils.GeneratePrefixedUUID("image"), ext)
}

func (uc *ticketUseCase) GetUserTickets(userID string) ([]*entity.Ticket, error) {
	tickets, err := uc.ticketRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, ticket := range tickets {
		if ticket.Image != "" {
			imageURL, err := uc.fileRepo.GetPresignedURL(context.Background(), ticket.Image)
			if err != nil {
				return nil, err
			}

			ticket.Image = imageURL
		}
	}
	return tickets, nil
}

func (uc *ticketUseCase) GetAllTicket() ([]*entity.Ticket, error) {
	tickets, err := uc.ticketRepo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, ticket := range tickets {
		if ticket.Image != "" {
			imageURL, err := uc.fileRepo.GetPresignedURL(context.Background(), ticket.Image)
			if err != nil {
				return nil, err
			}

			ticket.Image = imageURL
		}
	}

	return tickets, nil
}

func (uc *ticketUseCase) UpdateTicketStatus(ticketID, updatedBy string, status entity.TicketStatus) (*entity.Ticket, error) {
	_, err := uc.ticketRepo.FindByID(ticketID)
	if err != nil {
		return nil, errors.New("Ticket tidak ditemukan")
	}

	user, err := uc.userRepo.FindByID(updatedBy)
	if err != nil || user.Role != entity.RoleAdmin {
		return nil, errors.New("unauthorized")
	}

	if err := uc.ticketRepo.UpdateStatus(ticketID, status); err != nil {
		return nil, err
	}
	ctx := context.Background()
	_ = uc.cacheRepo.Delete(ctx, "dashboard:overview")

	return uc.ticketRepo.FindByID(ticketID)

}

func (uc *ticketUseCase) GetSummaryStat() (*dto.DashboardSummary, error) {
	ctx := context.Background()

	raw, err := uc.cacheRepo.Get(
		ctx,
		"dashboard:overview",
	)

	if err == nil {
		var stat dto.DashboardSummary

		if err := json.Unmarshal(
			[]byte(raw),
			&stat,
		); err == nil {
			return &stat, nil
		}
	}

	stat, err := uc.ticketRepo.GetSumary()
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(stat)
	if err == nil {
		_ = uc.cacheRepo.Set(
			ctx,
			"dashboard:overview",
			string(data),
			5*time.Minute,
		)
	}

	return stat, nil
}
