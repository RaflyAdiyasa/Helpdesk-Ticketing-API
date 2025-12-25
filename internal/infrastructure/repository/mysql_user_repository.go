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

func (r *MySQLUserRepository) Create(ticket *entity.Ticket) error {
	return nil
}
func (r *MySQLUserRepository) FindByID(id string) (entity.Ticket, error) {
	return entity.Ticket{}, nil
}
func (r *MySQLUserRepository) FindByUserID(id string) ([]entity.Ticket, error) {
	return nil, nil
}
func (r *MySQLUserRepository) FindAll() ([]entity.Ticket, error) {
	return nil, nil
}
func (r *MySQLUserRepository) Update(ticket *entity.Ticket) error {
	return nil
}
func (r *MySQLUserRepository) UpdateStatus(ticket *entity.Ticket, status string) error {
	return nil
}
func (r *MySQLUserRepository) Delete(id string) error {
	return nil
}
