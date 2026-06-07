package usecase

import (
	"context"
	"errors"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/dto"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/repository"
)

type healthUseCase struct {
	ticketRepo repository.TicketRepositoy
	userRepo   repository.UserRepository
	fileRepo   repository.FileRepository
	cacheRepo  repository.CacheRepository
}

func NewHealthUseCase(ticketRepo repository.TicketRepositoy, userRepo repository.UserRepository, fileRepo repository.FileRepository, cacheRepo repository.CacheRepository) HealthUseCase {
	return &healthUseCase{
		ticketRepo: ticketRepo,
		userRepo:   userRepo,
		fileRepo:   fileRepo,
		cacheRepo:  cacheRepo,
	}
}

func (uc *healthUseCase) Readiness() (*dto.ServicesStatus, error) {
	ctx := context.Background()
	var msg error
	status := dto.ServicesStatus{
		CacheRepoStatus: "up",
		BucketExists:    "up",
		DatabaseStatus:  "up",
	}

	if err := uc.cacheRepo.Ping(ctx); err != nil {
		status.CacheRepoStatus = "down"
		msg = errors.New("Service_not_ready")
	}
	if err := uc.fileRepo.Ping(ctx); err != nil {
		status.BucketExists = "down"
		msg = errors.New("Service_not_ready")
	}
	if err := uc.ticketRepo.Ping(ctx); err != nil {
		status.DatabaseStatus = "down"
		msg = errors.New("Service_not_ready")
	}

	return &status, msg
}
