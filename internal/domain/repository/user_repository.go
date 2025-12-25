package repository

import (
	"github.com/RaflyAdiyasa/Helpdest-Ticketing-API/internal/domain/entity"
)

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id string) error
	FindAll() (*[]entity.User, error)
}
