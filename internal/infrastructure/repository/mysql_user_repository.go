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

func (r *MySQLUserRepository) Create(user *entity.User) (*entity.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (r *MySQLUserRepository) FindByID(userID string) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *MySQLUserRepository) FindByEmail(email string) (*entity.User, error) {

	var user entity.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *MySQLUserRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}
func (r *MySQLUserRepository) Delete(userID string) error {
	return r.db.Delete(&entity.User{}, userID).Error
}
func (r *MySQLUserRepository) FindAll() ([]*entity.User, error) {
	var users []*entity.User
	if err := r.db.Model(&users).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
