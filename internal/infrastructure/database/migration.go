package database

import (
	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Ticket{},
	)

	if err != nil {
		return err
	}

	return nil
}
