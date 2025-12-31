package database

import (
	"fmt"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {

	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&entity.User{},
		&entity.Ticket{},
	)

	if err != nil {
		return err
	}

	return nil
}
