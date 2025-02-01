package usecase

import (
	"errors"

	"github.com/Prototype-1/admin-auth-service/internal/models"
	"github.com/Prototype-1/admin-auth-service/internal/repository"
	utils "github.com/Prototype-1/admin-auth-service/internal/utils/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase interface {
	Signup(email, password string) error
	Login(email, password string) (string, error)
	BlockUser(userID uint) error
	UnblockUser(userID uint) error
	SuspendUser(userID uint) error
	GetAllUsers() ([]*models.User, error)
}

type adminUsecaseImpl struct {
	repo repository.AdminRepository
}

func NewAdminUsecase(repo repository.AdminRepository) AdminUsecase {
	return &adminUsecaseImpl{repo: repo}
}

func (u *adminUsecaseImpl) Signup(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &models.Admin{
		Email:    email,
		Password: string(hashedPassword),
	}

	return u.repo.CreateAdmin(admin)
}

func (u *adminUsecaseImpl) Login(email, password string) (string, error) {
	admin, err := u.repo.GetAdminByEmail(email)
	if err != nil || admin == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(admin.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *adminUsecaseImpl) BlockUser(userID uint) error {
	return u.repo.BlockUser(userID)
}

func (u *adminUsecaseImpl) UnblockUser(userID uint) error {
	return u.repo.UnblockUser(userID)
}

func (u *adminUsecaseImpl) SuspendUser(userID uint) error {
	return u.repo.SuspendUser(userID)
}

func (u *adminUsecaseImpl) GetAllUsers() ([]*models.User, error) {
	return u.repo.GetAllUsers()
}
