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
	grpcServer.Serve(listener)
}

