package service

import (
	"context"
	"errors"
	"testing"

	"github.com/imam/gspend-app/apps/auth-service/internal/config"
	"github.com/imam/gspend-app/apps/auth-service/internal/domain"
	"github.com/imam/gspend-app/apps/auth-service/internal/dto"
	"github.com/imam/gspend-app/apps/auth-service/internal/util"
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
	password := "password123"
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
