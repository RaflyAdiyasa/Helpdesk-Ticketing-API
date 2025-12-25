package usecase

import (
	"github.com/RaflyAdiyasa/Helpdest-Ticketing-API/internal/domain/entity"
	"github.com/RaflyAdiyasa/Helpdest-Ticketing-API/internal/domain/repository"
	pkg "github.com/RaflyAdiyasa/Helpdest-Ticketing-API/pkg/jwt"
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
	// nanti

	return nil, nil
}

func (uc *authUseCase) Login(username, password string) (string, error) {
	/// nanti
	return "huan", nil
}
