package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockCategoryRepository for testing
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) ListByUserID(ctx context.Context, userID primitive.ObjectID, categoryType domain.CategoryType) ([]*domain.Category, error) {
	args := m.Called(ctx, userID, categoryType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) ListSystem(ctx context.Context, categoryType domain.CategoryType) ([]*domain.Category, error) {
	args := m.Called(ctx, categoryType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCategoryHandler_Create(t *testing.T) {
	e := echo.New()
	categoryRepo := new(MockCategoryRepository)
	svc := service.NewCategoryService(categoryRepo)
	h := NewCategoryHandler(svc)

	userID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(dto.CreateCategoryRequest{
			Name:      "Groceries",
			Type:      "expense",
			Icon:      "🛒",
			Color:     "#10B981",
			SortOrder: 1,
		})

		categoryRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// The response is the category object directly, not wrapped in success/data
		var response domain.Category
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Groceries", response.Name)
		assert.Equal(t, domain.CategoryTypeExpense, response.Type)

		categoryRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid request body")
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		reqBody, _ := json.Marshal(dto.CreateCategoryRequest{
			Name: "Groceries",
			Type: "expense",
			Icon: "🛒",
			Color: "#10B981",
		})

		categoryRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(assert.AnError).Once()

		req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		categoryRepo.AssertExpectations(t)
	})
}

func TestCategoryHandler_List(t *testing.T) {
	e := echo.New()
	categoryRepo := new(MockCategoryRepository)
	svc := service.NewCategoryService(categoryRepo)
	h := NewCategoryHandler(svc)

	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	t.Run("Success - User Categories", func(t *testing.T) {
		categories := []*domain.Category{
			{
				ID:       categoryID,
				UserID:   &userID,
				Name:     "Groceries",
				Type:     domain.CategoryTypeExpense,
				Icon:     "🛒",
				Color:    "#10B981",
				IsSystem: false,
			},
		}

		categoryRepo.On("ListByUserID", mock.Anything, userID, domain.CategoryType("")).Return(categories, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.List(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []*domain.Category
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)

		categoryRepo.AssertExpectations(t)
	})

	t.Run("Success - System Categories", func(t *testing.T) {
		categories := []*domain.Category{
			{
				ID:       categoryID,
				UserID:   nil,
				Name:     "Rent",
				Type:     domain.CategoryTypeExpense,
				Icon:     "🏠",
				Color:    "#3B82F6",
				IsSystem: true,
			},
		}

		// The handler calls ListUserCategories with "system" type, which should still use ListByUserID
		categoryRepo.On("ListByUserID", mock.Anything, userID, domain.CategoryType("system")).Return(categories, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/categories?type=system", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.List(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		categoryRepo.AssertExpectations(t)
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		categoryRepo.On("ListByUserID", mock.Anything, userID, domain.CategoryType("")).Return(nil, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.List(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		categoryRepo.AssertExpectations(t)
	})
}

func TestCategoryHandler_Update(t *testing.T) {
	e := echo.New()
	categoryRepo := new(MockCategoryRepository)
	svc := service.NewCategoryService(categoryRepo)
	h := NewCategoryHandler(svc)

	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		existingCategory := &domain.Category{
			ID:     categoryID,
			UserID: &userID,
			Name:   "Groceries",
			Type:   domain.CategoryTypeExpense,
		}

		reqBody, _ := json.Marshal(dto.UpdateCategoryRequest{
			Name:      "Updated Groceries",
			Icon:      "🛒",
			Color:     "#10B981",
			SortOrder: 2,
		})

		categoryRepo.On("GetByID", mock.Anything, categoryID).Return(existingCategory, nil).Once()
		categoryRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Category")).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/categories/"+categoryID.Hex(), strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(categoryID.Hex())
		c.Set("user_id", userID.Hex())

		err := h.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response domain.Category
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Groceries", response.Name)

		categoryRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/categories/invalid-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid-id")
		c.Set("user_id", userID.Hex())

		err := h.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Error - Category Not Found", func(t *testing.T) {
		categoryRepo.On("GetByID", mock.Anything, categoryID).Return(nil, nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/categories/"+categoryID.Hex(), strings.NewReader("{}"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(categoryID.Hex())
		c.Set("user_id", userID.Hex())

		err := h.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		categoryRepo.AssertExpectations(t)
	})
}

func TestCategoryHandler_Delete(t *testing.T) {
	e := echo.New()
	categoryRepo := new(MockCategoryRepository)
	svc := service.NewCategoryService(categoryRepo)
	h := NewCategoryHandler(svc)

	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		userCategory := &domain.Category{
			ID:       categoryID,
			UserID:   &userID,
			Name:     "Custom Category",
			Type:     domain.CategoryTypeExpense,
			IsSystem: false,
		}

		categoryRepo.On("GetByID", mock.Anything, categoryID).Return(userCategory, nil).Once()
		categoryRepo.On("Delete", mock.Anything, categoryID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/categories/"+categoryID.Hex(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(categoryID.Hex())
		c.Set("user_id", userID.Hex())

		err := h.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		categoryRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/categories/invalid-id", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid-id")
		c.Set("user_id", userID.Hex())

		err := h.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		userCategory := &domain.Category{
			ID:       categoryID,
			UserID:   &userID,
			Name:     "Custom Category",
			Type:     domain.CategoryTypeExpense,
			IsSystem: false,
		}

		categoryRepo.On("GetByID", mock.Anything, categoryID).Return(userCategory, nil).Once()
		categoryRepo.On("Delete", mock.Anything, categoryID).Return(assert.AnError).Once()

		req := httptest.NewRequest(http.MethodDelete, "/categories/"+categoryID.Hex(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(categoryID.Hex())
		c.Set("user_id", userID.Hex())

		err := h.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		categoryRepo.AssertExpectations(t)
	})
}
func TestCategoryHandler_Delete_SystemProtection(t *testing.T) {
	e := echo.New()
	categoryRepo := new(MockCategoryRepository)
	svc := service.NewCategoryService(categoryRepo)
	h := NewCategoryHandler(svc)

	userID := primitive.NewObjectID()
	systemCategoryID := primitive.NewObjectID()

	t.Run("Error - System Category Protection", func(t *testing.T) {
		systemCategory := &domain.Category{
			ID:       systemCategoryID,
			UserID:   nil, // System category
			Name:     "Groceries",
			Type:     domain.CategoryTypeExpense,
			IsSystem: true,
		}

		categoryRepo.On("GetByID", mock.Anything, systemCategoryID).Return(systemCategory, nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/categories/"+systemCategoryID.Hex(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(systemCategoryID.Hex())
		c.Set("user_id", userID.Hex())

		err := h.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "System categories cannot be deleted")

		categoryRepo.AssertExpectations(t)
	})

	t.Run("Success - User Category Deletion", func(t *testing.T) {
		userCategoryID := primitive.NewObjectID()
		userCategory := &domain.Category{
			ID:       userCategoryID,
			UserID:   &userID,
			Name:     "Custom Category",
			Type:     domain.CategoryTypeExpense,
			IsSystem: false,
		}

		categoryRepo.On("GetByID", mock.Anything, userCategoryID).Return(userCategory, nil).Once()
		categoryRepo.On("Delete", mock.Anything, userCategoryID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/categories/"+userCategoryID.Hex(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(userCategoryID.Hex())
		c.Set("user_id", userID.Hex())

		err := h.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)

		categoryRepo.AssertExpectations(t)
	})
}