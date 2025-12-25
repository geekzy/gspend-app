package service

import (
	"context"
	"errors"
	"time"

	"github.com/imam/gspend-app/apps/auth-service/internal/config"
	"github.com/imam/gspend-app/apps/auth-service/internal/domain"
	"github.com/imam/gspend-app/apps/auth-service/internal/dto"
	"github.com/imam/gspend-app/apps/auth-service/internal/util"
)

type AuthService struct {
	userRepo domain.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo domain.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	exists, err := s.userRepo.Exists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		FamilySize:   req.FamilySize,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !util.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*dto.UserDTO, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserDTO{
		ID:         user.ID.Hex(),
		Email:      user.Email,
		FullName:   user.FullName,
		FamilySize: user.FamilySize,
	}, nil
}

func (s *AuthService) generateAuthResponse(user *domain.User) (*dto.AuthResponse, error) {
	// Generate Access Token
	accessToken, err := util.GenerateToken(
		user.ID.Hex(),
		user.Email,
		s.config.JWTSecret,
		15*time.Minute,
	)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshToken, err := util.GenerateToken(
		user.ID.Hex(),
		user.Email,
		s.config.RefreshSecret,
		7*24*time.Hour,
	)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserDTO{
			ID:         user.ID.Hex(),
			Email:      user.Email,
			FullName:   user.FullName,
			FamilySize: user.FamilySize,
		},
	}, nil
}
