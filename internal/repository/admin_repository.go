package repository

import (
	"errors"
	"log"
	"github.com/Prototype-1/admin-auth-service/internal/models"
	"gorm.io/gorm"
)

type AdminRepository interface {
	CreateAdmin(admin *models.Admin) error
	GetAdminByEmail(email string) (*models.Admin, error)
}

type adminRepositoryImpl struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepositoryImpl{db: db}
}

func (r *adminRepositoryImpl) CreateAdmin(admin *models.Admin) error {
    var existingAdmin models.Admin
    if err := r.db.Where("email = ?", admin.Email).First(&existingAdmin).Error; err == nil {
        return errors.New("admin with this email already exists")
    }
    return r.db.Create(admin).Error
}

func (r *adminRepositoryImpl) GetAdminByEmail(email string) (*models.Admin, error) {
	log.Println("Searching for admin with email:", email) 
	var admin models.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Admin not found in DB")
		return nil, nil
	}
	return &admin, err
}


