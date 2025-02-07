package usecase

import (
    "context"
    "errors"
    "fmt"
    "os"

    "github.com/Prototype-1/admin-auth-service/internal/models"
    "github.com/Prototype-1/admin-auth-service/internal/repository"
    "github.com/Prototype-1/admin-auth-service/internal/utils"
    userpb "github.com/Prototype-1/admin-auth-service/proto/user"
    routepb "github.com/Prototype-1/admin-auth-service/proto/routes"
    "github.com/joho/godotenv"
    "golang.org/x/crypto/bcrypt"
)

type AdminUsecase interface {
    Signup(email, password string) error
    Login(email, password string) (string, error)

    BlockUser(userID uint) error
    UnblockUser(userID uint) error
    SuspendUser(userID uint) error
    GetAllUsers() ([]*models.User, error)

    AddRoute(routeName string, startStopID, endStopID, categoryID int) error
    UpdateRoute(routeID int, routeName string, startStopID, endStopID, categoryID int) error
    DeleteRoute(routeID int) error
    GetAllRoutes() ([]*models.Route, error)
}

type adminUsecaseImpl struct {
    repo        repository.AdminRepository
    userService userpb.UserServiceClient
    routeService routepb.RouteServiceClient
}

func init() {
    err := godotenv.Load("config/.env")
    if err != nil {
        fmt.Println("Error loading .env file:", err)
    }
    secretKey := os.Getenv("JWT_SECRET_KEY")
    if secretKey == "" {
        fmt.Println("Warning: JWT_SECRET_KEY is not set in .env file")
    } else {
        fmt.Println("JWT_SECRET_KEY loaded successfully")
    }
}

func NewAdminUsecase(repo repository.AdminRepository, userClient userpb.UserServiceClient, routeClient routepb.RouteServiceClient) AdminUsecase {
    return &adminUsecaseImpl{
        repo:        repo,
        userService: userClient,
        routeService: routeClient,
    }
}

func (u *adminUsecaseImpl) Signup(email, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    admin := &models.Admin{
        Email:    email,
        Password: string(hashedPassword),
    }
    return u.repo.CreateAdmin(admin)
}

func (u *adminUsecaseImpl) Login(email, password string) (string, error) {
    admin, err := u.repo.GetAdminByEmail(email)
    if err != nil || admin == nil {
        return "", errors.New("invalid credentials")
    }
    if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }
    secretKey := os.Getenv("JWT_SECRET_KEY")
    fmt.Println("USECASE: JWT_SECRET_KEY =", secretKey)
    if secretKey == "" {
        return "", errors.New("server error: missing JWT_SECRET_KEY")
    }
    token, err := utils.GenerateJWT(int(admin.ID), secretKey)
    if err != nil {
        return "", err
    }
    return token, nil
}

func (u *adminUsecaseImpl) BlockUser(userID uint) error {
    _, err := u.userService.BlockUser(context.Background(), &userpb.UserRequest{
        UserId: uint32(userID),
    })
    return err
}

func (u *adminUsecaseImpl) UnblockUser(userID uint) error {
    _, err := u.userService.UnblockUser(context.Background(), &userpb.UserRequest{
        UserId: uint32(userID),
    })
    return err
}

func (u *adminUsecaseImpl) SuspendUser(userID uint) error {
    _, err := u.userService.SuspendUser(context.Background(), &userpb.UserRequest{
        UserId: uint32(userID),
    })
    return err
}

func (u *adminUsecaseImpl) GetAllUsers() ([]*models.User, error) {
    res, err := u.userService.GetAllUsers(context.Background(), &userpb.Empty{})
    if err != nil {
        return nil, err
    }
    var users []*models.User
    for _, u := range res.Users {
        users = append(users, &models.User{
            ID:             uint(u.Id),
            Email:          u.Email,
            Name:           u.Name,
            BlockedStatus:  u.BlockedStatus,
            InactiveStatus: u.InactiveStatus,
        })
    }
    return users, nil
}

func (u *adminUsecaseImpl) AddRoute(routeName string, startStopID, endStopID, categoryID int) error {
    _, err := u.routeService.AddRoute(context.Background(), &routepb.AddRouteRequest{
        RouteName:   routeName,
        StartStopId: int32(startStopID),
        EndStopId:   int32(endStopID),
        CategoryId:  int32(categoryID),
    })
    return err
}

func (u *adminUsecaseImpl) UpdateRoute(routeID int, routeName string, startStopID, endStopID, categoryID int) error {
    _, err := u.routeService.UpdateRoute(context.Background(), &routepb.UpdateRouteRequest{
        RouteId:     int32(routeID),
        RouteName:   routeName,
        StartStopId: int32(startStopID),
        EndStopId:   int32(endStopID),
        CategoryId:  int32(categoryID),
    })
    return err
}

func (u *adminUsecaseImpl) DeleteRoute(routeID int) error {
    _, err := u.routeService.DeleteRoute(context.Background(), &routepb.DeleteRouteRequest{
        RouteId: int32(routeID),
    })
    return err
}

func (u *adminUsecaseImpl) GetAllRoutes() ([]*models.Route, error) {
    res, err := u.routeService.GetAllRoutes(context.Background(), &routepb.GetAllRoutesRequest{})
    if err != nil {
        return nil, err
    }

    var routes []*models.Route
    for _, r := range res.Routes {
        routes = append(routes, &models.Route{
            RouteID:    int(r.RouteId),
            RouteName:  r.RouteName,
            StartStopID: int(r.StartStopId),
            EndStopID:   int(r.EndStopId),
            CategoryID: int(r.CategoryId),
        })
    }
    return routes, nil
}