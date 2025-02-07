package handlers

import (
	"context"
	"log"

	"github.com/Prototype-1/admin-auth-service/internal/usecase"
	"github.com/Prototype-1/admin-auth-service/internal/utils"
	pb "github.com/Prototype-1/admin-auth-service/proto/admin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	userpb "github.com/Prototype-1/admin-auth-service/proto/user"
	routepb  "github.com/Prototype-1/admin-auth-service/proto/routes"
)

type AdminServer struct {
	pb.UnimplementedAdminServiceServer
	usecase   usecase.AdminUsecase
	userClient userpb.UserServiceClient
	routeClient routepb.RouteServiceClient
}

func NewAdminServer(usecase usecase.AdminUsecase, userClient userpb.UserServiceClient, routeClient routepb.RouteServiceClient) *AdminServer {
	return &AdminServer{
		usecase: usecase, 
		userClient: userClient,
		routeClient: routeClient,
	}

}

func (s *AdminServer) AdminSignup(ctx context.Context, req *pb.AdminSignupRequest) (*pb.AuthResponse, error) {
    err := s.usecase.Signup(req.Email, req.Password)
    if err != nil {
        return nil, err
    }

    return &pb.AuthResponse{
        Message: "Signup successful", 
    }, nil
}


func (s *AdminServer) AdminLogin(ctx context.Context, req *pb.AdminLoginRequest) (*pb.AuthResponse, error) {
	token, err := s.usecase.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		AccessToken: token,
		Message:     "Login successful",
	}, nil
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

func (s *AdminServer) BlockUser(ctx context.Context, req *userpb.UserRequest) (*userpb.StatusResponse, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is blocking user ID %d", adminID, req.UserId)

	_, err = s.userClient.BlockUser(ctx, &userpb.UserRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	return &userpb.StatusResponse{Message: "User blocked successfully"}, nil
}

func (s *AdminServer) UnblockUser(ctx context.Context, req *userpb.UserRequest) (*userpb.StatusResponse, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is unblocking user ID %d", adminID, req.UserId)

	_, err = s.userClient.UnblockUser(ctx, &userpb.UserRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	return &userpb.StatusResponse{Message: "User unblocked successfully"}, nil
}

func (s *AdminServer) SuspendUser(ctx context.Context, req *userpb.UserRequest) (*userpb.StatusResponse, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is suspending user ID %d", adminID, req.UserId)

	_, err = s.userClient.SuspendUser(ctx, &userpb.UserRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	return &userpb.StatusResponse{Message: "User suspended successfully"}, nil
}

func (s *AdminServer) GetAllUsers(ctx context.Context, req *userpb.Empty) (*userpb.UserList, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is retrieving all users", adminID)

	userList, err := s.userClient.GetAllUsers(ctx, &userpb.Empty{})
	if err != nil {
		return nil, err
	}

	var userResponses []*userpb.User
	for _, user := range userList.Users {
		userResponses = append(userResponses, &userpb.User{
			Id:             uint32(user.Id),
			Name:           user.Name,
			Email:          user.Email,
			BlockedStatus:  user.BlockedStatus,
			InactiveStatus: user.InactiveStatus,
		})
	}

	return &userpb.UserList{Users: userResponses}, nil
}

func (s *AdminServer) AddRoute(ctx context.Context, req *routepb.AddRouteRequest) (*routepb.AddRouteResponse, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is adding a new route: %s", adminID, req.RouteName)

	response, err := s.routeClient.AddRoute(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to add route: %v", err)
	}

	return response, nil
}

func (s *AdminServer) UpdateRoute(ctx context.Context, req *routepb.UpdateRouteRequest) (*routepb.UpdateRouteResponse, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is updating route ID %d", adminID, req.RouteId)

	response, err := s.routeClient.UpdateRoute(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update route: %v", err)
	}

	return response, nil
}

func (s *AdminServer) DeleteRoute(ctx context.Context, req *routepb.DeleteRouteRequest) (*routepb.DeleteRouteResponse, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is deleting route ID %d", adminID, req.RouteId)

	response, err := s.routeClient.DeleteRoute(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete route: %v", err)
	}

	return response, nil
}

func (s *AdminServer) GetAllRoutes(ctx context.Context, req *routepb.GetAllRoutesRequest) (*routepb.GetAllRoutesResponse, error) {
	adminID, err := s.authenticateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Admin ID %d is retrieving all routes", adminID)

	response, err := s.routeClient.GetAllRoutes(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch routes: %v", err)
	}

	return response, nil
}
