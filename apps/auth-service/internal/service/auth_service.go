package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/geekzy/gspend-app/apps/auth-service/internal/config"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/util"
)

type AuthService struct {
	userRepo     domain.UserRepository
	config       *config.Config
	emailService *EmailService
}

func NewAuthService(userRepo domain.UserRepository, cfg *config.Config, emailService *EmailService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		config:       cfg,
		emailService: emailService,
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

	// Validate password strength
	if !util.ValidatePassword(req.Password) {
		return nil, errors.New("password must be at least 8 characters with uppercase, lowercase, and number")
	}

	// Validate family size
	if req.FamilySize < 0 || req.FamilySize > 5 {
		return nil, errors.New("family size must be between 0 and 5 children")
	}

	// Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Generate verification token
	verificationToken, err := GenerateToken()
	if err != nil {
		return nil, err
	}
	verificationExpiry := time.Now().Add(24 * time.Hour)

	// Create user
	user := &domain.User{
		Email:              req.Email,
		PasswordHash:       hashedPassword,
		FullName:           req.FullName,
		FamilySize:         req.FamilySize,
		EmailVerified:      false,
		VerificationToken:  verificationToken,
		VerificationExpiry: &verificationExpiry,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Send verification email (async, don't block registration)
	go func() {
		if err := s.emailService.SendVerificationEmail(user.Email, user.FullName, verificationToken); err != nil {
			log.Printf("[AUTH] Failed to send verification email: %v", err)
		}
	}()

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

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	user, err := s.userRepo.GetByVerificationToken(ctx, token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid verification token")
	}

	// Check if token is expired
	if user.VerificationExpiry != nil && time.Now().After(*user.VerificationExpiry) {
		return errors.New("verification token has expired")
	}

	// Mark email as verified
	user.EmailVerified = true
	user.VerificationToken = ""
	user.VerificationExpiry = nil

	return s.userRepo.Update(ctx, user)
}

func (s *AuthService) ResendVerification(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	// Don't reveal if user exists - return success anyway for security
	if user == nil {
		return nil
	}

	// Already verified
	if user.EmailVerified {
		return nil
	}

	// Generate new verification token
	verificationToken, err := GenerateToken()
	if err != nil {
		return err
	}
	verificationExpiry := time.Now().Add(24 * time.Hour)

	user.VerificationToken = verificationToken
	user.VerificationExpiry = &verificationExpiry

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Send verification email
	return s.emailService.SendVerificationEmail(user.Email, user.FullName, verificationToken)
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	// Don't reveal if user exists - return success anyway for security
	if user == nil {
		return nil
	}

	// Generate reset token
	resetToken, err := GenerateToken()
	if err != nil {
		return err
	}
	resetExpiry := time.Now().Add(1 * time.Hour)

	user.ResetToken = resetToken
	user.ResetTokenExpiry = &resetExpiry

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Send password reset email
	return s.emailService.SendPasswordResetEmail(user.Email, user.FullName, resetToken)
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	user, err := s.userRepo.GetByResetToken(ctx, token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid reset token")
	}

	// Check if token is expired
	if user.ResetTokenExpiry != nil && time.Now().After(*user.ResetTokenExpiry) {
		return errors.New("reset token has expired")
	}

	// Validate new password strength
	if !util.ValidatePassword(newPassword) {
		return errors.New("password must be at least 8 characters with uppercase, lowercase, and number")
	}

	// Hash new password
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password and clear reset token
	user.PasswordHash = hashedPassword
	user.ResetToken = ""
	user.ResetTokenExpiry = nil

	return s.userRepo.Update(ctx, user)
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
		ID:            user.ID.Hex(),
		Email:         user.Email,
		FullName:      user.FullName,
		FamilySize:    user.FamilySize,
		EmailVerified: user.EmailVerified,
	}, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) (*dto.UserDTO, error) {
	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Validate family size (0-5 children) first
	if req.FamilySize < 0 || req.FamilySize > 5 {
		return nil, errors.New("family size must be between 0 and 5 children")
	}

	// Check if email is being changed and if it's already in use by another user
	if req.Email != user.Email {
		exists, err := s.userRepo.Exists(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already in use")
		}
	}

	// Update user fields
	user.FullName = req.FullName
	user.FamilySize = req.FamilySize
	user.Email = req.Email

	// Save to database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserDTO{
		ID:            user.ID.Hex(),
		Email:         user.Email,
		FullName:      user.FullName,
		FamilySize:    user.FamilySize,
		EmailVerified: user.EmailVerified,
	}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error {
	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify current password
	if !util.CheckPasswordHash(req.CurrentPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	// Validate new password strength
	if !util.ValidatePassword(req.NewPassword) {
		return errors.New("new password must be at least 8 characters with uppercase, lowercase, and number")
	}

	// Hash new password
	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = hashedPassword

	// Save to database
	return s.userRepo.Update(ctx, user)
}

func (s *AuthService) ValidateJWT(token string) (*util.JWTClaims, error) {
	return util.ValidateToken(token, s.config.JWT.Secret)
}

func (s *AuthService) CheckUserExists(ctx context.Context, userID string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (s *AuthService) generateTokens(userID, email string) (string, string, error) {
	// Generate Access Token
	accessToken, err := util.GenerateToken(
		userID,
		email,
		s.config.JWT.Secret,
		15*time.Minute,
	)
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshToken, err := util.GenerateToken(
		userID,
		email,
		s.config.JWT.RefreshSecret,
		7*24*time.Hour,
	)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) generateAuthResponse(user *domain.User) (*dto.AuthResponse, error) {
	// Generate Access Token and Refresh Token
	accessToken, refreshToken, err := s.generateTokens(user.ID.Hex(), user.Email)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserDTO{
			ID:            user.ID.Hex(),
			Email:         user.Email,
			FullName:      user.FullName,
			FamilySize:    user.FamilySize,
			EmailVerified: user.EmailVerified,
		},
	}, nil
}
