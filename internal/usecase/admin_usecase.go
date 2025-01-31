package usecase

import (
	"github.com/Prototype-1/admin-auth-service/internal/models"
	"github.com/Prototype-1/admin-auth-service/internal/repository"
	"errors"
	"time"
)

type AdminUsecase interface {
	GetAllUsers() ([]models.User, error)
	BlockUser(userID uint) error
	UnblockUser(userID uint) error
	SuspendUser(userID uint) error
}

type adminUsecaseImpl struct {
	adminRepo repository.AdminRepository
}

func NewAdminUsecase(adminRepo repository.AdminRepository) AdminUsecase {
	return &adminUsecaseImpl{adminRepo: adminRepo}
}

// GetAllUsers fetches all normal users
func (u *adminUsecaseImpl) GetAllUsers() ([]models.User, error) {
	return u.adminRepo.GetAllUsers()
}

// BlockUser blocks a user by setting `BlockedStatus` to true
func (u *adminUsecaseImpl) BlockUser(userID uint) error {
	user, err := u.adminRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if user.BlockedStatus {
		return errors.New("user is already blocked")
	}

	user.BlockedStatus = true
	return u.adminRepo.UpdateUser(user)
}

func (u *adminUsecaseImpl) UnblockUser(userID uint) error {
	user, err := u.adminRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if !user.BlockedStatus {
		return errors.New("user is not blocked")
	}

	user.BlockedStatus = false
	return u.adminRepo.UpdateUser(user)
}

func (u *adminUsecaseImpl) SuspendUser(userID uint) error {
	user, err := u.adminRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if user.InactiveStatus {
		return errors.New("user is already suspended")
	}

	now := time.Now()
	user.InactiveStatus = true
	user.SuspendedAt = &now

	return u.adminRepo.UpdateUser(user)
}
