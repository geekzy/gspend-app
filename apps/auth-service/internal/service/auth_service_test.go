package service

import (
	"context"
	"errors"
	"testing"

	"github.com/geekzy/gspend-app/apps/auth-service/internal/config"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository is a manual mock of domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		RefreshSecret: "test-refresh-secret",
	}
	authService := NewAuthService(mockRepo, cfg)

	ctx := context.Background()
	email := "test@example.com"
	password := "Password123"
	hashedPassword, _ := util.HashPassword(password)
	userID := primitive.NewObjectID()

	user := &domain.User{
		ID:           userID,
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     "Test User",
		FamilySize:   4,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, email).Return(user, nil).Once()

		req := &dto.LoginRequest{
			Email:    email,
			Password: password,
		}

		resp, err := authService.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, email, resp.User.Email)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Credentials - User Not Found", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, "wrong@example.com").Return(nil, nil).Once()

		req := &dto.LoginRequest{
			Email:    "wrong@example.com",
			Password: password,
		}

		resp, err := authService.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Credentials - Wrong Password", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, email).Return(user, nil).Once()

		req := &dto.LoginRequest{
			Email:    email,
			Password: "wrongpassword",
		}

		resp, err := authService.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo.On("GetByEmail", ctx, email).Return(nil, errors.New("db error")).Once()

		req := &dto.LoginRequest{
			Email:    email,
			Password: password,
		}

		resp, err := authService.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{
		JWTSecret:     "test-secret",
		RefreshSecret: "test-refresh-secret",
	}
	authService := NewAuthService(mockRepo, cfg)

	ctx := context.Background()
	req := &dto.RegisterRequest{
		Email:      "new@example.com",
		Password:   "Password123",
		FullName:   "New User",
		FamilySize: 2,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Exists", ctx, req.Email).Return(false, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		resp, err := authService.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, req.Email, resp.User.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Already Exists", func(t *testing.T) {
		mockRepo.On("Exists", ctx, req.Email).Return(true, nil).Once()

		resp, err := authService.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "user already exists", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_GetProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, &config.Config{})

	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()

	t.Run("Success", func(t *testing.T) {
		user := &domain.User{Email: "test@example.com"}
		mockRepo.On("GetByID", ctx, userID).Return(user, nil).Once()

		res, err := authService.GetProfile(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, user.Email, res.Email)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_CheckUserExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, &config.Config{})

	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()

	t.Run("Exists", func(t *testing.T) {
		user := &domain.User{Email: "test@example.com"}
		mockRepo.On("GetByID", ctx, userID).Return(user, nil).Once()

		exists, err := authService.CheckUserExists(ctx, userID)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Exists", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, userID).Return(nil, nil).Once()

		exists, err := authService.CheckUserExists(ctx, userID)

		assert.NoError(t, err)
		assert.False(t, exists)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_UpdateProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, &config.Config{})

	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()

	t.Run("Success", func(t *testing.T) {
		user := &domain.User{
			ID:         primitive.NewObjectID(),
			Email:      "old@example.com",
			FullName:   "Old Name",
			FamilySize: 2,
		}
		req := &dto.UpdateProfileRequest{
			Email:      "new@example.com",
			FullName:   "New Name",
			FamilySize: 3,
		}

		mockRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
		mockRepo.On("Exists", ctx, req.Email).Return(false, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		result, err := authService.UpdateProfile(ctx, userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Email, result.Email)
		assert.Equal(t, req.FullName, result.FullName)
		assert.Equal(t, req.FamilySize, result.FamilySize)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Email Already In Use", func(t *testing.T) {
		user := &domain.User{
			ID:         primitive.NewObjectID(),
			Email:      "old@example.com",
			FullName:   "Old Name",
			FamilySize: 2,
		}
		req := &dto.UpdateProfileRequest{
			Email:      "existing@example.com",
			FullName:   "New Name",
			FamilySize: 3,
		}

		mockRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
		mockRepo.On("Exists", ctx, req.Email).Return(true, nil).Once()

		result, err := authService.UpdateProfile(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "email already in use", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Family Size", func(t *testing.T) {
		user := &domain.User{
			ID:         primitive.NewObjectID(),
			Email:      "old@example.com",
			FullName:   "Old Name",
			FamilySize: 2,
		}
		req := &dto.UpdateProfileRequest{
			Email:      "new@example.com",
			FullName:   "New Name",
			FamilySize: 6, // Invalid - more than 5
		}

		mockRepo.On("GetByID", ctx, userID).Return(user, nil).Once()

		result, err := authService.UpdateProfile(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "family size must be between 0 and 5 children", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_ChangePassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := NewAuthService(mockRepo, &config.Config{})

	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	currentPassword := "OldPassword123"
	hashedCurrentPassword, _ := util.HashPassword(currentPassword)

	t.Run("Success", func(t *testing.T) {
		user := &domain.User{
			ID:           primitive.NewObjectID(),
			Email:        "test@example.com",
			PasswordHash: hashedCurrentPassword,
		}
		req := &dto.ChangePasswordRequest{
			CurrentPassword: currentPassword,
			NewPassword:     "NewPassword123",
		}

		mockRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		err := authService.ChangePassword(ctx, userID, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Incorrect Current Password", func(t *testing.T) {
		user := &domain.User{
			ID:           primitive.NewObjectID(),
			Email:        "test@example.com",
			PasswordHash: hashedCurrentPassword,
		}
		req := &dto.ChangePasswordRequest{
			CurrentPassword: "WrongPassword123",
			NewPassword:     "NewPassword123",
		}

		mockRepo.On("GetByID", ctx, userID).Return(user, nil).Once()

		err := authService.ChangePassword(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, "current password is incorrect", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid New Password", func(t *testing.T) {
		user := &domain.User{
			ID:           primitive.NewObjectID(),
			Email:        "test@example.com",
			PasswordHash: hashedCurrentPassword,
		}
		req := &dto.ChangePasswordRequest{
			CurrentPassword: currentPassword,
			NewPassword:     "weak", // Invalid password
		}

		mockRepo.On("GetByID", ctx, userID).Return(user, nil).Once()

		err := authService.ChangePassword(ctx, userID, req)

		assert.Error(t, err)
		assert.Equal(t, "new password must be at least 8 characters with uppercase, lowercase, and number", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
