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

// MockReportRepository for testing
type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) GetBudgetVsActualData(ctx context.Context, userID primitive.ObjectID, month time.Time) (*domain.BudgetVsActualReport, error) {
	args := m.Called(ctx, userID, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BudgetVsActualReport), args.Error(1)
}

func (m *MockReportRepository) GetSpendingByCategoryData(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (*domain.SpendingByCategoryReport, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SpendingByCategoryReport), args.Error(1)
}

func (m *MockReportRepository) GetMonthlyTrendsData(ctx context.Context, userID primitive.ObjectID, months int) (*domain.MonthlyTrendsReport, error) {
	args := m.Called(ctx, userID, months)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MonthlyTrendsReport), args.Error(1)
}

func TestReportHandler_GetBudgetVsActual(t *testing.T) {
	e := echo.New()
	reportRepo := new(MockReportRepository)
	txRepo := new(MockTransactionRepository)
	budgetRepo := new(MockBudgetRepository)
	svc := service.NewReportService(reportRepo, txRepo, budgetRepo)
	h := NewReportHandler(svc)

	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	t.Run("Success - Current Month", func(t *testing.T) {
		now := time.Now()
		expectedReport := &domain.BudgetVsActualReport{
			Month: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()),
			Categories: []*domain.BudgetVsActualItem{
				{
					CategoryID:     categoryID,
					CategoryName:   "Groceries",
					Budgeted:       500.0,
					Actual:         450.0,
					Variance:       -50.0,
					PercentageUsed: 90.0,
				},
			},
			TotalBudgeted:   500.0,
			TotalSpent:      450.0,
			OverallVariance: -50.0,
		}

		reportRepo.On("GetBudgetVsActualData", mock.Anything, userID, mock.AnythingOfType("time.Time")).Return(expectedReport, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/reports/budget-vs-actual", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetBudgetVsActual(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))
		
		data := response["data"].(map[string]interface{})
		assert.Equal(t, 500.0, data["total_budgeted"])
		assert.Equal(t, 450.0, data["total_spent"])
		assert.Equal(t, -50.0, data["overall_variance"])

		reportRepo.AssertExpectations(t)
	})

	t.Run("Success - Specific Month", func(t *testing.T) {
		month := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		expectedReport := &domain.BudgetVsActualReport{
			Month:           month,
			Categories:      []*domain.BudgetVsActualItem{},
			TotalBudgeted:   0.0,
			TotalSpent:      0.0,
			OverallVariance: 0.0,
		}

		reportRepo.On("GetBudgetVsActualData", mock.Anything, userID, month).Return(expectedReport, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/reports/budget-vs-actual?month=2024-01", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetBudgetVsActual(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		reportRepo.AssertExpectations(t)
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		reportRepo.On("GetBudgetVsActualData", mock.Anything, userID, mock.AnythingOfType("time.Time")).Return(nil, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodGet, "/reports/budget-vs-actual", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetBudgetVsActual(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))

		reportRepo.AssertExpectations(t)
	})
}

func TestReportHandler_GetSpendingByCategory(t *testing.T) {
	e := echo.New()
	reportRepo := new(MockReportRepository)
	txRepo := new(MockTransactionRepository)
	budgetRepo := new(MockBudgetRepository)
	svc := service.NewReportService(reportRepo, txRepo, budgetRepo)
	h := NewReportHandler(svc)

	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	t.Run("Success - Default Date Range", func(t *testing.T) {
		expectedReport := &domain.SpendingByCategoryReport{
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			Categories: []*domain.CategorySpending{
				{
					CategoryID:   categoryID,
					CategoryName: "Groceries",
					Amount:       300.0,
					Percentage:   60.0,
				},
			},
			TotalSpent: 500.0,
		}

		reportRepo.On("GetSpendingByCategoryData", mock.Anything, userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(expectedReport, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/reports/spending-by-category", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetSpendingByCategory(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, 500.0, data["total_spent"])

		reportRepo.AssertExpectations(t)
	})

	t.Run("Success - Custom Date Range", func(t *testing.T) {
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
		
		expectedReport := &domain.SpendingByCategoryReport{
			StartDate:  startDate,
			EndDate:    endDate,
			Categories: []*domain.CategorySpending{},
			TotalSpent: 0.0,
		}

		reportRepo.On("GetSpendingByCategoryData", mock.Anything, userID, startDate, endDate).Return(expectedReport, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/reports/spending-by-category?start_date=2024-01-01&end_date=2024-01-31", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetSpendingByCategory(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		reportRepo.AssertExpectations(t)
	})

	t.Run("Error - Invalid Start Date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reports/spending-by-category?start_date=invalid-date", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetSpendingByCategory(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(map[string]interface{})["message"], "Invalid start_date format")
	})

	t.Run("Error - Invalid End Date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reports/spending-by-category?end_date=invalid-date", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetSpendingByCategory(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Error - Start Date After End Date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reports/spending-by-category?start_date=2024-01-31&end_date=2024-01-01", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetSpendingByCategory(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"].(map[string]interface{})["message"], "Start date must be before end date")
	})
}

func TestReportHandler_GetMonthlyTrends(t *testing.T) {
	e := echo.New()
	reportRepo := new(MockReportRepository)
	txRepo := new(MockTransactionRepository)
	budgetRepo := new(MockBudgetRepository)
	svc := service.NewReportService(reportRepo, txRepo, budgetRepo)
	h := NewReportHandler(svc)

	userID := primitive.NewObjectID()

	t.Run("Success - Default 3 Months", func(t *testing.T) {
		expectedReport := &domain.MonthlyTrendsReport{
			Months: 3,
			MonthlyData: []*domain.MonthlySpending{
				{
					Month:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Amount: 550.0,
				},
				{
					Month:  time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
					Amount: 550.0,
				},
				{
					Month:  time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					Amount: 550.0,
				},
			},
			AverageSpending: 550.0,
			TrendDirection:  "stable",
		}

		reportRepo.On("GetMonthlyTrendsData", mock.Anything, userID, 3).Return(expectedReport, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/reports/monthly-trends", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetMonthlyTrends(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(3), data["months"])
		assert.Equal(t, 550.0, data["average_spending"])
		assert.Equal(t, "stable", data["trend_direction"])

		reportRepo.AssertExpectations(t)
	})

	t.Run("Success - Custom Months", func(t *testing.T) {
		expectedReport := &domain.MonthlyTrendsReport{
			Months:          6,
			MonthlyData:     []*domain.MonthlySpending{},
			AverageSpending: 0.0,
			TrendDirection:  "stable",
		}

		reportRepo.On("GetMonthlyTrendsData", mock.Anything, userID, 6).Return(expectedReport, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/reports/monthly-trends?months=6", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetMonthlyTrends(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		reportRepo.AssertExpectations(t)
	})

	t.Run("Success - Invalid Months Parameter Uses Default", func(t *testing.T) {
		expectedReport := &domain.MonthlyTrendsReport{
			Months:          3,
			MonthlyData:     []*domain.MonthlySpending{},
			AverageSpending: 0.0,
			TrendDirection:  "stable",
		}

		// Should use default of 3 months when invalid parameter is provided
		reportRepo.On("GetMonthlyTrendsData", mock.Anything, userID, 3).Return(expectedReport, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/reports/monthly-trends?months=invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetMonthlyTrends(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		reportRepo.AssertExpectations(t)
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		reportRepo.On("GetMonthlyTrendsData", mock.Anything, userID, 3).Return(nil, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodGet, "/reports/monthly-trends", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetMonthlyTrends(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))

		reportRepo.AssertExpectations(t)
	})
}