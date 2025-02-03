package handlers

import (
	"context"
	"log"

	"github.com/Prototype-1/admin-auth-service/internal/usecase"
	"github.com/Prototype-1/admin-auth-service/internal/utils"
	pb "github.com/Prototype-1/admin-auth-service/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type AdminServer struct {
	pb.UnimplementedAdminServiceServer
	usecase usecase.AdminUsecase
}

func NewAdminServer(usecase usecase.AdminUsecase) *AdminServer {
	return &AdminServer{usecase: usecase}
}

func (s *AdminServer) AdminSignup(ctx context.Context, req *pb.AdminSignupRequest) (*pb.AuthResponse, error) {
	err := s.usecase.Signup(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{Message: "Signup successful"}, nil
}

func (s *AdminServer) AdminLogin(ctx context.Context, req *pb.AdminLoginRequest) (*pb.AuthResponse, error) {
	token, err := s.usecase.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{AccessToken: token, Message: "Login successful"}, nil
}

func (s *AdminServer) authenticateAdmin(ctx context.Context) (uint, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, "Missing metadata")
	}

	tokenList, exists := md["authorization"]
	if !exists || len(tokenList) == 0 {
		return 0, status.Errorf(codes.Unauthenticated, "Authorization token not provided")
	}

	tokenString := tokenList[0]
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	adminID, err := utils.ParseJWT(tokenString)
	if err != nil {
		return 0, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}

	return adminID, nil
}

func (s *AdminServer) BlockUser(ctx context.Context, req *pb.UserRequest) (*pb.StatusResponse, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is blocking user ID %d", adminID, req.UserId)
	err = s.usecase.BlockUser(uint(req.UserId))
	if err != nil {
		return nil, err
	}

	return &pb.StatusResponse{Message: "User blocked successfully"}, nil
}

func (s *AdminServer) UnblockUser(ctx context.Context, req *pb.UserRequest) (*pb.StatusResponse, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is unblocking user ID %d", adminID, req.UserId)
	err = s.usecase.UnblockUser(uint(req.UserId))
	if err != nil {
		return nil, err
	}

	return &pb.StatusResponse{Message: "User unblocked successfully"}, nil
}

func (s *AdminServer) SuspendUser(ctx context.Context, req *pb.UserRequest) (*pb.StatusResponse, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is suspending user ID %d", adminID, req.UserId)
	err = s.usecase.SuspendUser(uint(req.UserId))
	if err != nil {
		return nil, err
	}

	return &pb.StatusResponse{Message: "User suspended successfully"}, nil
}

func (s *AdminServer) GetAllUsers(ctx context.Context, req *pb.Empty) (*pb.UserList, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is retrieving all users", adminID)
	users, err := s.usecase.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var userResponses []*pb.User
	for _, user := range users {
		userResponses = append(userResponses, &pb.User{
			Id:             uint32(user.ID),
			Name:           user.FirstName + " " + user.LastName,
			Email:          user.Email,
			BlockedStatus:  user.BlockedStatus,
			InactiveStatus: user.InactiveStatus,
		})
	}

	return &pb.UserList{Users: userResponses}, nil
}
