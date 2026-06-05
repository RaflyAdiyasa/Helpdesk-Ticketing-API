package repository

import (
	"time"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/dto"
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
	if err := r.db.Preload("Owner").Find(&tickets, "user_id = ?", userID).Error; err != nil {
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
func (r *MySQLTicketRepository) UpdateStatus(ticketID string, status entity.TicketStatus) error {
	return r.db.Model(&entity.Ticket{}).Where("ticket_id = ?", ticketID).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
}
func (r *MySQLTicketRepository) Delete(ticketID string) error {
	return r.db.Delete(&entity.Ticket{}, ticketID).Error
}

func (r *MySQLTicketRepository) CountAllTicket() (int64, error) {
	var count int64

	err := r.db.Model(&entity.Ticket{}).Count(&count).Error
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *MySQLTicketRepository) CountTicketByUserID(userID string) (int64, error) {
	var count int64

	err := r.db.Model(&entity.Ticket{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *MySQLTicketRepository) GetSumary() (*dto.DashboardSummary, error) {
	var res dto.DashboardSummary
	err := r.db.Model(&entity.Ticket{}).Count(&res.TotalTicket).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&entity.Ticket{}).Where("status = ?", entity.StatusOpen).Count(&res.OpenTicket).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&entity.Ticket{}).Where("status = ?", entity.StatusInProgress).Count(&res.InProgressTicket).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&entity.Ticket{}).Where("status = ?", entity.StatusDone).Count(&res.ClosedTicket).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&entity.User{}).Count(&res.TotalUser).Error
	if err != nil {
		return nil, err
	}

	return &res, err
}
