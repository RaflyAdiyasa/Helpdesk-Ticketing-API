package repository

import (
	"github.com/RaflyAdiyasa/Helpdest-Ticketing-API/internal/domain/entity"
	"gorm.io/gorm"
)

type MySQLUserRepository struct {
	db *gorm.DB
}

func NewMySQLUserRepository(db *gorm.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Create(user *entity.User) error {
	return nil
}
func (r *MySQLUserRepository) FindByID(id string) (entity.User, error) {
	return entity.User{}, nil
}
func (r *MySQLUserRepository) FindByEmail(email string) (entity.User, error) {
	return entity.User{}, nil
}
func (r *MySQLUserRepository) Update(user *entity.User) error {
	return nil
}
func (r *MySQLUserRepository) Delete(id string) error {
	return nil
}
func (r *MySQLUserRepository) FindAll() ([]entity.User, error) {
	return nil, nil
}
