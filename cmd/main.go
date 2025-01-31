package main

import (
	"github.com/Prototype-1/admin-auth-service/internal/utils/logger" 
	"fmt"
)

func main() {
	// Initialize the logger
	logger.InitLogger()
	logger.Log.Info("Admin Auth Service started successfully")

	fmt.Println("Server is running...")
}
