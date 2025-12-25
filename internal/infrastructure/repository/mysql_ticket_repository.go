package repository

import (
	"github.com/RaflyAdiyasa/Helpdest-Ticketing-API/internal/domain/entity"
	"gorm.io/gorm"
)

type MySQLTicketRepository struct {
	db *gorm.DB
}

func NewMySQLTicketRepository(db *gorm.DB) *MySQLTicketRepository {
	return &MySQLTicketRepository{db: db}
}

func (r *MySQLTicketRepository) Create(ticket *entity.Ticket) error {
	return nil
}
func (r *MySQLTicketRepository) FindByID(id string) (entity.Ticket, error) {
	return entity.Ticket{}, nil
}
func (r *MySQLTicketRepository) FindByUserID(id string) ([]entity.Ticket, error) {
	return nil, nil
}
func (r *MySQLTicketRepository) FindAll() ([]entity.Ticket, error) {
	return nil, nil
}
func (r *MySQLTicketRepository) Update(ticket *entity.Ticket) error {
	return nil
}
func (r *MySQLTicketRepository) UpdateStatus(ticket *entity.Ticket, status string) error {
	return nil
}
func (r *MySQLTicketRepository) Delete(id string) error {
	return nil
}
