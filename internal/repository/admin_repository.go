package repository

import (
	"errors"

	"github.com/Prototype-1/admin-auth-service/internal/models"
	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(admin *models.Admin) error
	GetAdminByEmail(email string) (*models.Admin, error)
	BlockUser(userID uint) error
	UnblockUser(userID uint) error
	SuspendUser(userID uint) error
	GetAllUsers() ([]*models.User, error)
}

type adminRepositoryImpl struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepositoryImpl{db: db}
}

func (r *adminRepositoryImpl) CreateAdmin(admin *models.Admin) error {
	return r.db.Create(admin).Error
}

func (r *adminRepositoryImpl) GetAdminByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &admin, err
}

func (r *adminRepositoryImpl) BlockUser(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("blocked_status", true).Error
}

func (r *adminRepositoryImpl) UnblockUser(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("blocked_status", false).Error
}

func (r *adminRepositoryImpl) SuspendUser(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("inactive_status", true).Error
}

func (r *adminRepositoryImpl) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
