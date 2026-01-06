package repository

import (
	"time"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
	"gorm.io/gorm"
)

type MySQLTicketRepository struct {
	db *gorm.DB
}

func NewMySQLTicketRepository(db *gorm.DB) *MySQLTicketRepository {
	return &MySQLTicketRepository{db: db}
}

func (r *MySQLTicketRepository) Create(ticket *entity.Ticket) (*entity.Ticket, error) {
	if err := r.db.Create(ticket).Error; err != nil {
		return nil, err
	}
	return ticket, nil
}
func (r *MySQLTicketRepository) FindByID(id string) (*entity.Ticket, error) {
	var ticket entity.Ticket
	if err := r.db.Preload("Owner").First(&ticket, "ticket_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}
func (r *MySQLTicketRepository) FindByUserID(userID string) ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket
	if err := r.db.Model(&tickets).Find(&tickets, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}
func (r *MySQLTicketRepository) FindAll() ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket
	if err := r.db.Model(&tickets).Preload("Owner").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}
func (r *MySQLTicketRepository) Update(ticket *entity.Ticket) error {
	ticket.UpdatedAt = time.Now()
	if err := r.db.Save(ticket).Error; err != nil {
		return err
	}
	return nil
}
func (r *MySQLTicketRepository) UpdateStatus(ticketID, status string) error {
	return r.db.Model(&entity.Ticket{}).Where("ticke_id = ?", ticketID).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
}
func (r *MySQLTicketRepository) Delete(ticketID string) error {
	return r.db.Delete(&entity.Ticket{}, ticketID).Error
}
