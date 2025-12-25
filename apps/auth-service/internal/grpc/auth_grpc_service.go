package grpc

import (
	"context"

	"github.com/imam/gspend-app/apps/auth-service/internal/service"
	authv1 "github.com/imam/gspend-app/apps/auth-service/pkg/proto/auth/v1"
)

type AuthGRPCService struct {
	authv1.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthGRPCService(authService *service.AuthService) *AuthGRPCService {
	return &AuthGRPCService{
		authService: authService,
	}
}

func (s *AuthGRPCService) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	claims, err := s.authService.ValidateJWT(req.Token)
	if err != nil {
		return &authv1.ValidateTokenResponse{
			Valid:        false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &authv1.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Email:  claims.Email,
	}, nil
}

func (s *AuthGRPCService) GetUserProfile(ctx context.Context, req *authv1.GetUserProfileRequest) (*authv1.GetUserProfileResponse, error) {
	user, err := s.authService.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &authv1.GetUserProfileResponse{
		User: &authv1.User{
			Id:         user.ID,
			Email:      user.Email,
			FullName:   user.FullName,
			FamilySize: int32(user.FamilySize),
		},
	}, nil
}

func (s *AuthGRPCService) CheckUserExists(ctx context.Context, req *authv1.CheckUserExistsRequest) (*authv1.CheckUserExistsResponse, error) {
	exists, err := s.authService.CheckUserExists(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &authv1.CheckUserExistsResponse{
		Exists: exists,
	}, nil
}
