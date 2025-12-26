package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTransactionHandler_Create(t *testing.T) {
	e := echo.New()
	txRepo := new(MockTransactionRepository)
	budgetRepo := new(MockBudgetRepository)
	svc := service.NewTransactionService(txRepo, budgetRepo)
	h := NewTransactionHandler(svc)

	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(dto.CreateTransactionRequest{
			CategoryID:      categoryID.Hex(),
			Type:            "expense",
			Amount:          100.0,
			Description:     "Groceries",
			TransactionDate: time.Now(),
			PaymentMethod:   "Cash",
		})

		txRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		if assert.NoError(t, h.Create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("Validation Error", func(t *testing.T) {
		reqBody, _ := json.Marshal(dto.CreateTransactionRequest{
			Type: "invalid", // Invalid type
		})

		req := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.Create(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestTransactionHandler_List(t *testing.T) {
	e := echo.New()
	txRepo := new(MockTransactionRepository)
	budgetRepo := new(MockBudgetRepository)
	svc := service.NewTransactionService(txRepo, budgetRepo)
	h := NewTransactionHandler(svc)

	userID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		txRepo.On("ListByUserID", mock.Anything, userID, mock.Anything).Return([]*domain.Transaction{}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		if assert.NoError(t, h.List(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}

func TestTransactionHandler_Update(t *testing.T) {
	e := echo.New()
	txRepo := new(MockTransactionRepository)
	budgetRepo := new(MockBudgetRepository)
	svc := service.NewTransactionService(txRepo, budgetRepo)
	h := NewTransactionHandler(svc)

	txID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(dto.CreateTransactionRequest{
			CategoryID:      categoryID.Hex(),
			Type:            "expense",
			Amount:          150.0,
			Description:     "Dinner",
			TransactionDate: time.Now(),
			PaymentMethod:   "Credit Card",
		})

		txRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/transactions/"+txID.Hex(), strings.NewReader(string(reqBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(txID.Hex())

		if assert.NoError(t, h.Update(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}

func TestTransactionHandler_Delete(t *testing.T) {
	e := echo.New()
	txRepo := new(MockTransactionRepository)
	budgetRepo := new(MockBudgetRepository)
	svc := service.NewTransactionService(txRepo, budgetRepo)
	h := NewTransactionHandler(svc)

	txID := primitive.NewObjectID()

	t.Run("Success", func(t *testing.T) {
		transaction := &domain.Transaction{
			ID:     txID,
			Type:   domain.TransactionTypeIncome, // No budget update for income
			Amount: 100.0,
		}
		txRepo.On("GetByID", mock.Anything, txID).Return(transaction, nil).Once()
		txRepo.On("Delete", mock.Anything, txID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/transactions/"+txID.Hex(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(txID.Hex())

		if assert.NoError(t, h.Delete(c)) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}
	})

	t.Run("Success with Budget Reversal", func(t *testing.T) {
		budgetID := primitive.NewObjectID()
		categoryID := primitive.NewObjectID()
		transaction := &domain.Transaction{
			ID:         txID,
			Type:       domain.TransactionTypeExpense,
			Amount:     100.0,
			BudgetID:   &budgetID,
			CategoryID: categoryID,
		}
		txRepo.On("GetByID", mock.Anything, txID).Return(transaction, nil).Once()
		txRepo.On("Delete", mock.Anything, txID).Return(nil).Once()
		budgetRepo.On("UpdateSpentAmount", mock.Anything, budgetID, categoryID, -100.0).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/transactions/"+txID.Hex(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(txID.Hex())

		if assert.NoError(t, h.Delete(c)) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}
	})
}
func TestTransactionHandler_ListWithFilters(t *testing.T) {
	e := echo.New()
	txRepo := new(MockTransactionRepository)
	budgetRepo := new(MockBudgetRepository)
	svc := service.NewTransactionService(txRepo, budgetRepo)
	h := NewTransactionHandler(svc)

	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()
	transactionID := primitive.NewObjectID()

	t.Run("Success - Default Filters", func(t *testing.T) {
		expectedResult := &domain.PaginatedTransactions{
			Transactions: []*domain.Transaction{
				{
					ID:              transactionID,
					UserID:          userID,
					CategoryID:      categoryID,
					Type:            domain.TransactionTypeExpense,
					Amount:          100.0,
					Description:     "Groceries",
					TransactionDate: time.Now(),
				},
			},
			Pagination: domain.Pagination{
				Page:       1,
				PerPage:    20,
				Total:      1,
				TotalPages: 1,
			},
			FiltersApplied: domain.TransactionFilters{
				Page:    1,
				PerPage: 20,
			},
		}

		txRepo.On("FindWithFilters", mock.Anything, userID, mock.AnythingOfType("domain.TransactionFilters")).Return(expectedResult.Transactions, 1, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/filtered", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.ListWithFilters(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		pagination := data["pagination"].(map[string]interface{})
		assert.Equal(t, float64(1), pagination["page"])
		assert.Equal(t, float64(20), pagination["per_page"])
		assert.Equal(t, float64(1), pagination["total"])

		txRepo.AssertExpectations(t)
	})

	t.Run("Success - With Date Range Filter", func(t *testing.T) {
		expectedResult := &domain.PaginatedTransactions{
			Transactions: []*domain.Transaction{},
			Pagination: domain.Pagination{
				Page:       1,
				PerPage:    20,
				Total:      0,
				TotalPages: 0,
			},
			FiltersApplied: domain.TransactionFilters{
				Page:    1,
				PerPage: 20,
			},
		}

		txRepo.On("FindWithFilters", mock.Anything, userID, mock.AnythingOfType("domain.TransactionFilters")).Return(expectedResult.Transactions, 0, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/filtered?start_date=2024-01-01&end_date=2024-01-31", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.ListWithFilters(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		txRepo.AssertExpectations(t)
	})

	t.Run("Success - With Category and Type Filters", func(t *testing.T) {
		expectedResult := &domain.PaginatedTransactions{
			Transactions: []*domain.Transaction{},
			Pagination: domain.Pagination{
				Page:       1,
				PerPage:    20,
				Total:      0,
				TotalPages: 0,
			},
			FiltersApplied: domain.TransactionFilters{
				Page:    1,
				PerPage: 20,
			},
		}

		txRepo.On("FindWithFilters", mock.Anything, userID, mock.AnythingOfType("domain.TransactionFilters")).Return(expectedResult.Transactions, 0, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/filtered?category_id="+categoryID.Hex()+"&type=expense", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.ListWithFilters(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		txRepo.AssertExpectations(t)
	})

	t.Run("Success - With Pagination and Sorting", func(t *testing.T) {
		expectedResult := &domain.PaginatedTransactions{
			Transactions: []*domain.Transaction{},
			Pagination: domain.Pagination{
				Page:       2,
				PerPage:    10,
				Total:      25,
				TotalPages: 3,
			},
			FiltersApplied: domain.TransactionFilters{
				Page:      2,
				PerPage:   10,
				SortBy:    "amount",
				SortOrder: "desc",
			},
		}

		txRepo.On("FindWithFilters", mock.Anything, userID, mock.AnythingOfType("domain.TransactionFilters")).Return(expectedResult.Transactions, 25, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/filtered?page=2&per_page=10&sort_by=amount&sort_order=desc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.ListWithFilters(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		pagination := data["pagination"].(map[string]interface{})
		assert.Equal(t, float64(2), pagination["page"])
		assert.Equal(t, float64(10), pagination["per_page"])
		assert.Equal(t, float64(25), pagination["total"])
		assert.Equal(t, float64(3), pagination["total_pages"])

		txRepo.AssertExpectations(t)
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		txRepo.On("FindWithFilters", mock.Anything, userID, mock.AnythingOfType("domain.TransactionFilters")).Return(nil, 0, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/filtered", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.ListWithFilters(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))
		assert.Equal(t, "TRANSACTION_ERROR", response["error"].(map[string]interface{})["code"])

		txRepo.AssertExpectations(t)
	})
}

func TestTransactionHandler_GetSpendingByCategory(t *testing.T) {
	e := echo.New()
	txRepo := new(MockTransactionRepository)
	budgetRepo := new(MockBudgetRepository)
	svc := service.NewTransactionService(txRepo, budgetRepo)
	h := NewTransactionHandler(svc)

	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	t.Run("Success - Default Date Range", func(t *testing.T) {
		expectedCategories := []*domain.CategorySpending{
			{
				CategoryID:   categoryID,
				CategoryName: "Groceries",
				Amount:       300.0,
				Percentage:   60.0,
			},
		}

		txRepo.On("GetSpendingByCategory", mock.Anything, userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(expectedCategories, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/spending-by-category", nil)
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

		data := response["data"].([]interface{})
		assert.Len(t, data, 1)

		txRepo.AssertExpectations(t)
	})

	t.Run("Success - Custom Date Range", func(t *testing.T) {
		expectedCategories := []*domain.CategorySpending{}

		txRepo.On("GetSpendingByCategory", mock.Anything, userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(expectedCategories, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/spending-by-category?start_date=2024-01-01&end_date=2024-01-31", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetSpendingByCategory(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		txRepo.AssertExpectations(t)
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		txRepo.On("GetSpendingByCategory", mock.Anything, userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(nil, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/spending-by-category", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetSpendingByCategory(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response["success"].(bool))

		txRepo.AssertExpectations(t)
	})
}

func TestTransactionHandler_GetMonthlyTrends(t *testing.T) {
	e := echo.New()
	txRepo := new(MockTransactionRepository)
	budgetRepo := new(MockBudgetRepository)
	svc := service.NewTransactionService(txRepo, budgetRepo)
	h := NewTransactionHandler(svc)

	userID := primitive.NewObjectID()

	t.Run("Success - Default 3 Months", func(t *testing.T) {
		expectedTrends := []*domain.MonthlySpending{
			{
				Month:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Amount: 500.0,
			},
			{
				Month:  time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				Amount: 600.0,
			},
			{
				Month:  time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
				Amount: 550.0,
			},
		}

		txRepo.On("GetMonthlySpendingTrends", mock.Anything, userID, 3).Return(expectedTrends, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/monthly-trends", nil)
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

		data := response["data"].([]interface{})
		assert.Len(t, data, 3)

		txRepo.AssertExpectations(t)
	})

	t.Run("Success - Custom Months", func(t *testing.T) {
		expectedTrends := []*domain.MonthlySpending{}

		txRepo.On("GetMonthlySpendingTrends", mock.Anything, userID, 6).Return(expectedTrends, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/monthly-trends?months=6", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetMonthlyTrends(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		txRepo.AssertExpectations(t)
	})

	t.Run("Success - Invalid Months Uses Default", func(t *testing.T) {
		expectedTrends := []*domain.MonthlySpending{}

		// Should use default of 3 months when invalid parameter is provided
		txRepo.On("GetMonthlySpendingTrends", mock.Anything, userID, 3).Return(expectedTrends, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/monthly-trends?months=invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", userID.Hex())

		err := h.GetMonthlyTrends(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		txRepo.AssertExpectations(t)
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		txRepo.On("GetMonthlySpendingTrends", mock.Anything, userID, 3).Return(nil, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodGet, "/transactions/monthly-trends", nil)
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

		txRepo.AssertExpectations(t)
	})
}