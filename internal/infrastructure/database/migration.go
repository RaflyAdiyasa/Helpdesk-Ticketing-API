package database

import (
	"fmt"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&entity.User{},
		&entity.Ticket{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
