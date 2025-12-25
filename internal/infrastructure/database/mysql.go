package database

import (
	"fmt"
	"log"

	"github.com/RaflyAdiyasa/Helpdest-Ticketing-API/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLConnection(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}
