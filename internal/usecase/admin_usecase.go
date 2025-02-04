package usecase

import (
	"errors"
	"fmt"
	"os"

	"github.com/Prototype-1/admin-auth-service/internal/models"
	"github.com/Prototype-1/admin-auth-service/internal/repository"
	utils "github.com/Prototype-1/admin-auth-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"github.com/joho/godotenv"
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

// ✅ Load .env when package initializes
func init() {
	err := godotenv.Load("config/.env") 
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		fmt.Println("⚠️ Warning: JWT_SECRET_KEY is not set in .env file")
	} else {
		fmt.Println("✅ JWT_SECRET_KEY loaded successfully")
	}
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

    secretKey := os.Getenv("JWT_SECRET_KEY") 
    fmt.Println("USECASE: JWT_SECRET_KEY =", secretKey) 

    if secretKey == "" {
        return "", errors.New("server error: missing JWT_SECRET_KEY")
    }

    token, err := utils.GenerateJWT(int(admin.ID), secretKey)
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
