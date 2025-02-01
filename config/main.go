package main

import (
	"github.com/Prototype-1/admin-auth-service/internal/handlers"
	"github.com/Prototype-1/admin-auth-service/internal/repository"
	"github.com/Prototype-1/admin-auth-service/internal/usecase"
	"github.com/Prototype-1/admin-auth-service/internal/utils"
	pb "github.com/Prototype-1/admin-auth-service/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// Initialize the database and perform auto-migration
	utils.InitDB()

	// Initialize the repository and use case layers
	repo := repository.NewAdminRepository(utils.DB)
	usecase := usecase.NewAdminUsecase(repo)

	// Initialize the gRPC server
	server := handlers.NewAdminServer(usecase)

	// Start listening on port 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterAdminServiceServer(grpcServer, server)

	// Log that the server is running
	log.Println("Admin service running on port 50051...")

	// Start the server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


