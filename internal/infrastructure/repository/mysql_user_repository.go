package repository

import (
	"context"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/domain/entity"
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
	return r.db.
		Delete(&entity.User{}, "user_id = ?", userID).
		Error
}
func (r *MySQLUserRepository) FindAll() ([]*entity.User, error) {
	var users []*entity.User
	if err := r.db.Model(&users).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *MySQLUserRepository) FIndByUsername(userName string) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, "username = ?", userName).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MySQLUserRepository) Ping(ctx context.Context) error {
	sqlDB, err := r.db.WithContext(ctx).DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}

func (r *MySQLUserRepository) UpdateProfile(user *entity.User) error {
	return r.db.Model(&entity.User{}).
		Where("user_id = ?", user.UserID).
		Updates(map[string]interface{}{
			"username":   user.Username,
			"email":      user.Email,
			"department": user.Department,
			"is_remote":  user.IsRemote,
		}).Error
}

func (r *MySQLUserRepository) UpdateProfilePicture(userID string, objectKey string) error {
	return r.db.Model(&entity.User{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"profile_pict": objectKey,
		}).Error
}
