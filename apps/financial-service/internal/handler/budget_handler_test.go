package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockBudgetService is a mock implementation of service.BudgetServiceInterface
type MockBudgetService struct {
	mock.Mock
}

func (m *MockBudgetService) CreateBudget(ctx context.Context, budget *domain.Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func (m *MockBudgetService) GetBudget(ctx context.Context, id string) (*domain.Budget, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Budget), args.Error(1)
}

func (m *MockBudgetService) ListUserBudgets(ctx context.Context, userID string) ([]*domain.Budget, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Budget), args.Error(1)
}

func (m *MockBudgetService) GetActiveBudget(ctx context.Context, userID string, date time.Time) (*domain.Budget, error) {
	args := m.Called(ctx, userID, date)
	return args.Get(0).(*domain.Budget), args.Error(1)
}

func (m *MockBudgetService) UpdateBudget(ctx context.Context, budget *domain.Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func (m *MockBudgetService) DeleteBudget(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBudgetService) AddBudgetItem(ctx context.Context, budgetID string, item *domain.BudgetItem) error {
	args := m.Called(ctx, budgetID, item)
	return args.Error(0)
}

func (m *MockBudgetService) UpdateBudgetItem(ctx context.Context, budgetID string, itemID string, item *domain.BudgetItem) error {
	args := m.Called(ctx, budgetID, itemID, item)
	return args.Error(0)
}

func (m *MockBudgetService) DeleteBudgetItem(ctx context.Context, budgetID string, itemID string) error {
	args := m.Called(ctx, budgetID, itemID)
	return args.Error(0)
}

func (m *MockBudgetService) GetBudgetItem(ctx context.Context, budgetID string, itemID string) (*domain.BudgetItem, error) {
	args := m.Called(ctx, budgetID, itemID)
	return args.Get(0).(*domain.BudgetItem), args.Error(1)
}

func TestBudgetHandler_AddBudgetItem(t *testing.T) {
	mockService := new(MockBudgetService)
	handler := &BudgetHandler{
		budgetService: mockService,
		validate:      validator.New(),
	}

	// Test data
	budgetID := primitive.NewObjectID().Hex()
	categoryID := primitive.NewObjectID().Hex()
	
	requestBody := dto.CreateBudgetItemRequest{
		CategoryID:    categoryID,
		CategoryName:  "Test Category",
		PlannedAmount: 500.0,
		Notes:         "Test notes",
	}

	// Mock service call
	mockService.On("AddBudgetItem", mock.Anything, budgetID, mock.AnythingOfType("*domain.BudgetItem")).Return(nil)

	// Create request
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/budgets/"+budgetID+"/items", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(budgetID)

	// Execute
	err := handler.AddBudgetItem(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestBudgetHandler_UpdateBudgetItem(t *testing.T) {
	mockService := new(MockBudgetService)
	handler := &BudgetHandler{
		budgetService: mockService,
		validate:      validator.New(),
	}

	// Test data
	budgetID := primitive.NewObjectID().Hex()
	itemID := primitive.NewObjectID().Hex()
	categoryID := primitive.NewObjectID().Hex()
	
	requestBody := dto.UpdateBudgetItemRequest{
		CategoryID:    categoryID,
		CategoryName:  "Updated Category",
		PlannedAmount: 600.0,
		Notes:         "Updated notes",
	}

	// Mock service call
	mockService.On("UpdateBudgetItem", mock.Anything, budgetID, itemID, mock.AnythingOfType("*domain.BudgetItem")).Return(nil)

	// Create request
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/budgets/"+budgetID+"/items/"+itemID, bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id", "itemId")
	c.SetParamValues(budgetID, itemID)

	// Execute
	err := handler.UpdateBudgetItem(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestBudgetHandler_DeleteBudgetItem(t *testing.T) {
	mockService := new(MockBudgetService)
	handler := &BudgetHandler{
		budgetService: mockService,
		validate:      validator.New(),
	}

	// Test data
	budgetID := primitive.NewObjectID().Hex()
	itemID := primitive.NewObjectID().Hex()

	// Mock service call
	mockService.On("DeleteBudgetItem", mock.Anything, budgetID, itemID).Return(nil)

	// Create request
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/budgets/"+budgetID+"/items/"+itemID, nil)
	
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id", "itemId")
	c.SetParamValues(budgetID, itemID)

	// Execute
	err := handler.DeleteBudgetItem(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockService.AssertExpectations(t)
}

func TestBudgetHandler_GetBudgetItem(t *testing.T) {
	mockService := new(MockBudgetService)
	handler := &BudgetHandler{
		budgetService: mockService,
		validate:      validator.New(),
	}

	// Test data
	budgetID := primitive.NewObjectID().Hex()
	itemID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()
	
	budgetItem := &domain.BudgetItem{
		ID:            itemID,
		CategoryID:    categoryID,
		CategoryName:  "Test Category",
		PlannedAmount: 500.0,
		SpentAmount:   100.0,
		Notes:         "Test notes",
	}

	// Mock service call
	mockService.On("GetBudgetItem", mock.Anything, budgetID, itemID.Hex()).Return(budgetItem, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/budgets/"+budgetID+"/items/"+itemID.Hex(), nil)
	
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id", "itemId")
	c.SetParamValues(budgetID, itemID.Hex())

	// Execute
	err := handler.GetBudgetItem(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	
	var response dto.BudgetItemDTO
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, itemID.Hex(), response.ID)
	assert.Equal(t, categoryID.Hex(), response.CategoryID)
	assert.Equal(t, "Test Category", response.CategoryName)
	assert.Equal(t, 500.0, response.PlannedAmount)
	assert.Equal(t, 100.0, response.SpentAmount)
	assert.Equal(t, "Test notes", response.Notes)
	
	mockService.AssertExpectations(t)
}

func TestBudgetHandler_AddBudgetItem_ValidationError(t *testing.T) {
	mockService := new(MockBudgetService)
	handler := &BudgetHandler{
		budgetService: mockService,
		validate:      validator.New(),
	}

	// Test data with invalid request (missing required fields)
	budgetID := primitive.NewObjectID().Hex()
	
	requestBody := dto.CreateBudgetItemRequest{
		// Missing required fields
		PlannedAmount: -100.0, // Invalid negative amount
	}

	// Create request
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/budgets/"+budgetID+"/items", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(budgetID)

	// Execute
	err := handler.AddBudgetItem(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// Should not call service due to validation error
	mockService.AssertNotCalled(t, "AddBudgetItem")
}

func TestBudgetHandler_GetBudgetItem_NotFound(t *testing.T) {
	mockService := new(MockBudgetService)
	handler := &BudgetHandler{
		budgetService: mockService,
		validate:      validator.New(),
	}

	// Test data
	budgetID := primitive.NewObjectID().Hex()
	itemID := primitive.NewObjectID().Hex()

	// Mock service call returning nil (not found)
	mockService.On("GetBudgetItem", mock.Anything, budgetID, itemID).Return((*domain.BudgetItem)(nil), nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/budgets/"+budgetID+"/items/"+itemID, nil)
	
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id", "itemId")
	c.SetParamValues(budgetID, itemID)

	// Execute
	err := handler.GetBudgetItem(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockService.AssertExpectations(t)
}