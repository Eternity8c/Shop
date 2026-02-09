package authgrpc

import (
	authproto "auth-service/api/gen/go/api/proto"
	"auth-service/internal/services/auth"
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	RegisterNewUser(ctx context.Context, email string, password string, fullName string) (userID int64, err error)
	Login(ctx context.Context, email string, password string) (string, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	authproto.UnimplementedAuthServiceServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authproto.RegisterAuthServiceServer(gRPC, &serverAPI{auth: auth})
}

const (
	emptyValue = 0
)

func (s *serverAPI) Login(ctx context.Context, req *authproto.LoginRequest) (*authproto.LoginResponse, error) {
	if err := validatelogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authproto.LoginResponse{
		Token: token,
	}, nil

}

func (s *serverAPI) Register(ctx context.Context, req *authproto.RegisterRequest) (*authproto.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword(), req.GetFullName())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authproto.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authproto.IsAdminRequest) (*authproto.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authproto.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validatelogin(req *authproto.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateRegister(req *authproto.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateIsAdmin(req *authproto.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}
