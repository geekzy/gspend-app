package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockDashboardRepository for testing
type MockDashboardRepository struct {
	mock.Mock
}

func (m *MockDashboardRepository) GetTotalBalance(ctx context.Context, userID primitive.ObjectID) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockDashboardRepository) GetMonthlyIncome(ctx context.Context, userID primitive.ObjectID, month time.Time) (float64, error) {
	args := m.Called(ctx, userID, month)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockDashboardRepository) GetMonthlyExpenses(ctx context.Context, userID primitive.ObjectID, month time.Time) (float64, error) {
	args := m.Called(ctx, userID, month)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockDashboardRepository) GetTopSpendingCategories(ctx context.Context, userID primitive.ObjectID, month time.Time, limit int) ([]*domain.CategorySpending, error) {
	args := m.Called(ctx, userID, month, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CategorySpending), args.Error(1)
}

func (m *MockDashboardRepository) GetRecentTransactions(ctx context.Context, userID primitive.ObjectID, limit int) ([]*domain.Transaction, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockDashboardRepository) GetCurrentMonthBudgetProgress(ctx context.Context, userID primitive.ObjectID) (*domain.BudgetProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BudgetProgress), args.Error(1)
}

func TestDashboardHandler_GetSummary(t *testing.T) {
	e := echo.New()
	dashboardRepo := new(MockDashboardRepository)
	svc := service.NewDashboardService(dashboardRepo)
	h := NewDashboardHandler(svc)

	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()
	transactionID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		
		// Mock all the dashboard repository calls
		dashboardRepo.On("GetTotalBalance", mock.Anything, userID).Return(2500.0, nil).Once()
		dashboardRepo.On("GetMonthlyIncome", mock.Anything, userID, mock.AnythingOfType("time.Time")).Return(3000.0, nil).Once()
		dashboardRepo.On("GetMonthlyExpenses", mock.Anything, userID, mock.AnythingOfType("time.Time")).Return(1800.0, nil).Once()
		
		budgetProgress := &domain.BudgetProgress{
			TotalBudget:     2000.0,
			TotalSpent:      1800.0,
			RemainingBudget: 200.0,
			PercentageUsed:  90.0,
		}
		dashboardRepo.On("GetCurrentMonthBudgetProgress", mock.Anything, userID).Return(budgetProgress, nil).Once()
		
		topCategories := []*domain.CategorySpending{
			{
				CategoryID:   categoryID,
				CategoryName: "Groceries",
				Amount:       600.0,
				Percentage:   33.3,
			},
		}
		dashboardRepo.On("GetTopSpendingCategories", mock.Anything, userID, mock.AnythingOfType("time.Time"), 3).Return(topCategories, nil).Once()
		
		recentTransactions := []*domain.Transaction{
			{
				ID:              transactionID,
				UserID:          userID,
				CategoryID:      categoryID,
				Type:            domain.TransactionTypeExpense,
				Amount:          85.50,
				Description:     "Weekly groceries",
				TransactionDate: now,
				Metadata: domain.TransactionMetadata{
					CategoryName: "Groceries",
				},
			},
		}
		dashboardRepo.On("GetRecentTransactions", mock.Anything, userID, 5).Return(recentTransactions, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard/summary", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetSummary(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, 2500.0, data["totalBalance"])
		assert.Equal(t, 3000.0, data["monthlyIncome"])
		assert.Equal(t, 1800.0, data["monthlyExpenses"])

		budgetProgressData := data["budgetProgress"].(map[string]interface{})
		assert.Equal(t, 2000.0, budgetProgressData["totalBudget"])
		assert.Equal(t, 90.0, budgetProgressData["percentageUsed"])

		dashboardRepo.AssertExpectations(t)
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		dashboardRepo.On("GetTotalBalance", mock.Anything, userID).Return(0.0, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard/summary", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetSummary(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))
		assert.Equal(t, "DASHBOARD_ERROR", response["error"].(map[string]interface{})["code"])

		dashboardRepo.AssertExpectations(t)
	})
}

func TestDashboardHandler_GetRecentTransactions(t *testing.T) {
	e := echo.New()
	dashboardRepo := new(MockDashboardRepository)
	svc := service.NewDashboardService(dashboardRepo)
	h := NewDashboardHandler(svc)

	userID := primitive.NewObjectID()
	transactionID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	t.Run("Success - Default Limit", func(t *testing.T) {
		transactions := []*domain.Transaction{
			{
				ID:              transactionID,
				UserID:          userID,
				CategoryID:      categoryID,
				Type:            domain.TransactionTypeExpense,
				Amount:          85.50,
				Description:     "Weekly groceries",
				TransactionDate: time.Now(),
			},
		}

		dashboardRepo.On("GetRecentTransactions", mock.Anything, userID, 5).Return(transactions, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard/recent-transactions", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetRecentTransactions(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		data := response["data"].([]interface{})
		assert.Len(t, data, 1)

		dashboardRepo.AssertExpectations(t)
	})

	t.Run("Success - Custom Limit", func(t *testing.T) {
		transactions := []*domain.Transaction{}

		dashboardRepo.On("GetRecentTransactions", mock.Anything, userID, 10).Return(transactions, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard/recent-transactions?limit=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetRecentTransactions(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		dashboardRepo.AssertExpectations(t)
	})

	t.Run("Success - Invalid Limit Uses Default", func(t *testing.T) {
		transactions := []*domain.Transaction{}

		dashboardRepo.On("GetRecentTransactions", mock.Anything, userID, 5).Return(transactions, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard/recent-transactions?limit=invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetRecentTransactions(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		dashboardRepo.AssertExpectations(t)
	})
}

func TestDashboardHandler_GetTopCategories(t *testing.T) {
	e := echo.New()
	dashboardRepo := new(MockDashboardRepository)
	svc := service.NewDashboardService(dashboardRepo)
	h := NewDashboardHandler(svc)

	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	t.Run("Success - Current Month", func(t *testing.T) {
		categories := []*domain.CategorySpending{
			{
				CategoryID:   categoryID,
				CategoryName: "Groceries",
				Amount:       600.0,
				Percentage:   40.0,
			},
		}

		dashboardRepo.On("GetTopSpendingCategories", mock.Anything, userID, mock.AnythingOfType("time.Time"), 10).Return(categories, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard/top-categories", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetTopCategories(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		data := response["data"].([]interface{})
		assert.Len(t, data, 1)

		dashboardRepo.AssertExpectations(t)
	})

	t.Run("Success - Specific Month", func(t *testing.T) {
		categories := []*domain.CategorySpending{}

		dashboardRepo.On("GetTopSpendingCategories", mock.Anything, userID, mock.AnythingOfType("time.Time"), 10).Return(categories, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard/top-categories?month=2024-01", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetTopCategories(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		dashboardRepo.AssertExpectations(t)
	})
}

func TestDashboardHandler_GetBudgetProgress(t *testing.T) {
	e := echo.New()
	dashboardRepo := new(MockDashboardRepository)
	svc := service.NewDashboardService(dashboardRepo)
	h := NewDashboardHandler(svc)

	userID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		budgetProgress := &domain.BudgetProgress{
			TotalBudget:     2000.0,
			TotalSpent:      1500.0,
			RemainingBudget: 500.0,
			PercentageUsed:  75.0,
		}

		dashboardRepo.On("GetCurrentMonthBudgetProgress", mock.Anything, userID).Return(budgetProgress, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard/budget-progress", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetBudgetProgress(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, 2000.0, data["totalBudget"])
		assert.Equal(t, 1500.0, data["totalSpent"])
		assert.Equal(t, 75.0, data["percentageUsed"])

		dashboardRepo.AssertExpectations(t)
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		dashboardRepo.On("GetCurrentMonthBudgetProgress", mock.Anything, userID).Return(nil, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodGet, "/dashboard/budget-progress", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetBudgetProgress(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))

		dashboardRepo.AssertExpectations(t)
	})
}