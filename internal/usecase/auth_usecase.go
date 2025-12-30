package usecase

import (
	"errors"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/repository"
	pkg "github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/pkg/jwt"
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/pkg/utils"
)

type authUseCase struct {
	userRepo   repository.UserRepository
	jwtService pkg.JWTservice
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtService pkg.JWTservice) AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (uc *authUseCase) Register(username, email, password string) (*entity.User, error) {
	existingUser, _ := uc.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("User sudah ada")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		UserID:   utils.GenerateUserID(string(entity.RoleUser)),
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Role:     entity.RoleUser,
	}
	return uc.userRepo.Create(user)
}

func (uc *authUseCase) Login(username, password string) (string, error) {
	user, err := uc.userRepo.FIndByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials: username tidak ditemukan")
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials: email tidak ditemukan")
	}

	token, err := uc.jwtService.GenerateToken(user.UserID, string(user.Role))
	if err != nil {
		return "", err
	}
	return token, nil
}
