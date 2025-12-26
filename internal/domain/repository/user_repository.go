package repository

import (
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
)

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	FindByID(userID string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FIndByUsername(userName string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(userID string) error
	FindAll() (*[]entity.User, error)
}
