package main

import (
	"fmt"
	"log"
	"net"

	"os"

	"github.com/Prototype-1/admin-auth-service/internal/handlers"
	"github.com/Prototype-1/admin-auth-service/internal/repository"
	"github.com/Prototype-1/admin-auth-service/internal/usecase"
	"github.com/Prototype-1/admin-auth-service/internal/utils"
	pb "github.com/Prototype-1/admin-auth-service/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func init() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if len(jwtSecretKey) == 0 {
		log.Fatal("JWT_SECRET_KEY environment variable not set")
	}
	fmt.Println(jwtSecretKey)
}

func main() {
	utils.InitDB()
	repo := repository.NewAdminRepository(utils.DB)
	usecase := usecase.NewAdminUsecase(repo)

	server := handlers.NewAdminServer(usecase)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAdminServiceServer(grpcServer, server)

	log.Println("Admin service running on port 50051...")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
