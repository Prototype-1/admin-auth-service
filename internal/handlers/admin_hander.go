package handlers

import (
	"context"
	"github.com/Prototype-1/admin-auth-service/internal/usecase"
	pb "github.com/Prototype-1/admin-auth-service/proto"
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

func (s *AdminServer) BlockUser(ctx context.Context, req *pb.UserRequest) (*pb.StatusResponse, error) {
	err := s.usecase.BlockUser(uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.StatusResponse{Message: "User blocked successfully"}, nil
}

func (s *AdminServer) UnblockUser(ctx context.Context, req *pb.UserRequest) (*pb.StatusResponse, error) {
	err := s.usecase.UnblockUser(uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.StatusResponse{Message: "User unblocked successfully"}, nil
}

func (s *AdminServer) SuspendUser(ctx context.Context, req *pb.UserRequest) (*pb.StatusResponse, error) {
	err := s.usecase.SuspendUser(uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.StatusResponse{Message: "User suspended successfully"}, nil
}

func (s *AdminServer) GetAllUsers(ctx context.Context, req *pb.Empty) (*pb.UserList, error) {
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