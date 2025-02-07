package main

import (
	"log"
	"net"

	"os"

	"github.com/Prototype-1/admin-auth-service/internal/handlers"
	"github.com/Prototype-1/admin-auth-service/internal/repository"
	"github.com/Prototype-1/admin-auth-service/internal/usecase"
	"github.com/Prototype-1/admin-auth-service/internal/utils"
	pb "github.com/Prototype-1/admin-auth-service/proto/admin"
	userpb "github.com/Prototype-1/admin-auth-service/proto/user"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func init() {
	err := godotenv.Load("config/.env")
	if err != nil {
		utils.Log.Sugar().Fatalf("Error loading .env file: %v", err)
	}
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if len(jwtSecretKey) == 0 {
		log.Fatal("JWT_SECRET_KEY environment variable not set")
	}
}

func main() {
	utils.InitLogger() 
	utils.Log.Info("Logger initialized successfully")

	utils.InitDB()

conn, err := grpc.Dial(":50052", grpc.WithInsecure()) 
if err != nil {
    log.Fatalf("did not connect: %v", err)
}
defer conn.Close()

userClient := userpb.NewUserServiceClient(conn)


	repo := repository.NewAdminRepository(utils.DB)
	usecase := usecase.NewAdminUsecase(repo, userClient)

	server := handlers.NewAdminServer(usecase, userClient)

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
