package handlers

import (
    "context"
    "github.com/Prototype-1/admin-auth-service/internal/repository"
    pb "github.com/Prototype-1/admin-auth-service/proto" 
)

type AdminServer struct {
    pb.UnimplementedAdminServiceServer
    adminRepo repository.AdminRepository
}

func NewAdminServer(adminRepo repository.AdminRepository) *AdminServer {
    return &AdminServer{adminRepo: adminRepo}
}

func (s *AdminServer) GetAllUsers(ctx context.Context, req *pb.Empty) (*pb.UserList, error) {
    users, err := s.adminRepo.GetAllUsers()
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

func (s *AdminServer) BlockUser(ctx context.Context, req *pb.UserRequest) (*pb.StatusResponse, error) {
    err := s.adminRepo.BlockUser(uint(req.UserId))
    if err != nil {
        return nil, err
    }
    return &pb.StatusResponse{Message: "User blocked successfully"}, nil
}

func (s *AdminServer) UnblockUser(ctx context.Context, req *pb.UserRequest) (*pb.StatusResponse, error) {
    err := s.adminRepo.UnblockUser(uint(req.UserId))
    if err != nil {
        return nil, err
    }
    return &pb.StatusResponse{Message: "User unblocked successfully"}, nil
}

func (s *AdminServer) SuspendUser(ctx context.Context, req *pb.UserRequest) (*pb.StatusResponse, error) {
    err := s.adminRepo.SuspendUser(uint(req.UserId))
    if err != nil {
        return nil, err
    }
    return &pb.StatusResponse{Message: "User suspended successfully"}, nil
}