package auth

import (
	"auth-application/internal/domain/models"
	"auth-application/internal/storage"
	"context"
	ath "github.com/FreylGit/protoModel/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Auth interface {
	RegistrationNewUser(ctx context.Context, newUser models.NewUser) (*models.UserResponse, error)
	RefreshToken(ctx context.Context, userId int64, refreshToken string) (id int64, rToken string, atoken string, err error)
	Login(ctx context.Context, email string, password string) (*models.UserResponse, error)
}

type serverApi struct {
	ath.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ath.RegisterAuthServer(gRPC,
		&serverApi{auth: auth},
	)
}

func (s *serverApi) Register(ctx context.Context, req *ath.RegisterRequest) (*ath.RegisterResponse, error) {

	newUser := models.NewUser{
		Email:       req.GetEmail(),
		Password:    req.GetPassword(),
		Name:        req.GetName(),
		DateOfBirth: time.Date(2001, 7, 9, 0, 0, 0, 0, time.UTC),
	}
	rspModel, err := s.auth.RegistrationNewUser(ctx, newUser)

	if err != nil {

	}

	return &ath.RegisterResponse{
		AccessToken:  rspModel.AccessToken,
		RefreshToken: rspModel.RefreshToken,
		UserId:       rspModel.UserId,
	}, nil
}

func (s *serverApi) Login(ctx context.Context, req *ath.LoginRequest) (*ath.LoginResponse, error) {
	res, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "error auth")
	}
	return &ath.LoginResponse{
		UserId:       res.UserId,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken}, nil
}

func (s *serverApi) Refresh(ctx context.Context, req *ath.RefreshRequest) (*ath.RefreshResponse, error) {
	id, rtoken, atoken, err := s.auth.RefreshToken(ctx, req.GetUserId(), req.GetRefreshToken())
	switch err {
	case nil:
		// pass
	case storage.ErrorScan:
		return nil, status.Errorf(codes.InvalidArgument, "Invalid arguments")
	case storage.ErrNotFound:
		return nil, status.Errorf(codes.InvalidArgument, "Refresh token not found or expired")
	case storage.ErrorUpdate:
		return nil, status.Errorf(codes.InvalidArgument, "Failed to update refresh token")

	default:
		return nil, status.Errorf(codes.Internal, "Internal error")
	}
	_ = id

	return &ath.RefreshResponse{
		AccessToken:  atoken,
		RefreshToken: rtoken,
	}, nil
}
