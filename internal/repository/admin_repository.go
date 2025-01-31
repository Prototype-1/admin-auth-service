package repository

import (
	"errors"
	"time"
	"gorm.io/gorm"
	"github.com/Prototype-1/admin-auth-service/internal/models"
)

type AdminRepository interface {
	GetAllUsers() ([]models.User, error)
	BlockUser(userID uint) error
	UnblockUser(userID uint) error
	SuspendUser(userID uint) error
	GetUserByID(userID uint) (*models.User, error)
}

type adminRepositoryImpl struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepositoryImpl{db: db}
}

func (r *adminRepositoryImpl) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("role = ?", "user").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *adminRepositoryImpl) BlockUser(userID uint) error {
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Update("blocked_status", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *adminRepositoryImpl) UnblockUser(userID uint) error {
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Update("blocked_status", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *adminRepositoryImpl) SuspendUser(userID uint) error {
	now := time.Now()
	result := r.db.Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"inactive_status": true,
			"suspended_at":    now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *adminRepositoryImpl) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

