package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/Prototype-1/admin-auth-service/internal/models"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load("./config/.env")
	if err != nil {
		Log.Fatal("Error loading .env file", zap.Error(err))
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)

	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = DB.AutoMigrate(&models.Admin{})
	if err != nil {
		Log.Fatal("Error migrating database: ", zap.Error(err))
	}

	log.Println("Successfully connected to the database and migrated the schema")
}
