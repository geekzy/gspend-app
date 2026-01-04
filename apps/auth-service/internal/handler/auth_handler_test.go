package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/geekzy/gspend-app/apps/auth-service/internal/config"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/service"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository is duplicated here for handler tests
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

func (m *MockUserRepository) GetByResetToken(ctx context.Context, token string) (*domain.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByVerificationToken(ctx context.Context, token string) (*domain.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateResetToken(ctx context.Context, userID, token string, expiry time.Time) error {
	args := m.Called(ctx, userID, token, expiry)
	return args.Error(0)
}

func TestAuthHandler_Login(t *testing.T) {
	e := echo.New()
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{
		JWT: config.JWTConfig{Secret: "test", RefreshSecret: "test-refresh"},
	}
	emailService := service.NewEmailService(cfg)
	authService := service.NewAuthService(mockRepo, cfg, emailService)
	h := NewAuthHandler(authService)

	t.Run("Success", func(t *testing.T) {
		userID := primitive.NewObjectID()
		password := "Password123"
		hashed, _ := util.HashPassword(password)
		user := &domain.User{
			ID:           userID,
			Email:        "test@example.com",
			PasswordHash: hashed,
		}

		mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil).Once()

		reqBody, _ := json.Marshal(dto.LoginRequest{
			Email:    "test@example.com",
			Password: password,
		})
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp dto.AuthResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", resp.User.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		mockRepo.On("GetByEmail", mock.Anything, "wrong@example.com").Return(nil, nil).Once()

		reqBody, _ := json.Marshal(dto.LoginRequest{
			Email:    "wrong@example.com",
			Password: "Password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthHandler_GetProfile(t *testing.T) {
	e := echo.New()
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{}
	emailService := service.NewEmailService(cfg)
	authService := service.NewAuthService(mockRepo, cfg, emailService)
	h := NewAuthHandler(authService)

	t.Run("Success", func(t *testing.T) {
		userID := primitive.NewObjectID()
		user := &domain.User{ID: userID, Email: "test@example.com"}
		mockRepo.On("GetByID", mock.Anything, userID.Hex()).Return(user, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetProfile(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp dto.UserDTO
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", resp.Email)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthHandler_Register(t *testing.T) {
	e := echo.New()
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{
		JWT: config.JWTConfig{Secret: "test", RefreshSecret: "test-refresh"},
	}
	emailService := service.NewEmailService(cfg)
	authService := service.NewAuthService(mockRepo, cfg, emailService)
	h := NewAuthHandler(authService)

	t.Run("Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(dto.RegisterRequest{
			Email:      "new@example.com",
			Password:   "Password123",
			FullName:   "New User",
			FamilySize: 2,
		})

		mockRepo.On("Exists", mock.Anything, "new@example.com").Return(false, nil).Once()
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthHandler_UpdateProfile(t *testing.T) {
	e := echo.New()
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{}
	emailService := service.NewEmailService(cfg)
	authService := service.NewAuthService(mockRepo, cfg, emailService)
	h := NewAuthHandler(authService)

	t.Run("Success", func(t *testing.T) {
		userID := primitive.NewObjectID()
		user := &domain.User{
			ID:         userID,
			Email:      "old@example.com",
			FullName:   "Old Name",
			FamilySize: 2,
		}

		reqBody, _ := json.Marshal(dto.UpdateProfileRequest{
			Email:      "new@example.com",
			FullName:   "New Name",
			FamilySize: 3,
		})

		mockRepo.On("GetByID", mock.Anything, userID.Hex()).Return(user, nil).Once()
		mockRepo.On("Exists", mock.Anything, "new@example.com").Return(false, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/me", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.UpdateProfile(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp dto.ProfileResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "Profile updated successfully", resp.Message)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthHandler_ChangePassword(t *testing.T) {
	e := echo.New()
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{}
	emailService := service.NewEmailService(cfg)
	authService := service.NewAuthService(mockRepo, cfg, emailService)
	h := NewAuthHandler(authService)

	t.Run("Success", func(t *testing.T) {
		userID := primitive.NewObjectID()
		currentPassword := "OldPassword123"
		hashedPassword, _ := util.HashPassword(currentPassword)
		user := &domain.User{
			ID:           userID,
			Email:        "test@example.com",
			PasswordHash: hashedPassword,
		}

		reqBody, _ := json.Marshal(dto.ChangePasswordRequest{
			CurrentPassword: currentPassword,
			NewPassword:     "NewPassword123",
		})

		mockRepo.On("GetByID", mock.Anything, userID.Hex()).Return(user, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/change-password", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.ChangePassword(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp dto.ProfileResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "Password changed successfully", resp.Message)
		mockRepo.AssertExpectations(t)
	})
}
