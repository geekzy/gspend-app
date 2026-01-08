package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
	validate           *validator.Validate
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		validate:           validator.New(),
	}
}

func (h *TransactionHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	req := new(dto.CreateTransactionRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	categoryObjectID, _ := primitive.ObjectIDFromHex(req.CategoryID)
	
	var budgetObjectID *primitive.ObjectID
	if req.BudgetID != "" {
		id, _ := primitive.ObjectIDFromHex(req.BudgetID)
		budgetObjectID = &id
	}

	transaction := &domain.Transaction{
		UserID:          userObjectID,
		CategoryID:      categoryObjectID,
		BudgetID:        budgetObjectID,
		Type:            domain.TransactionType(req.Type),
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
		PaymentMethod:   req.PaymentMethod,
		Notes:           req.Notes,
		// Metadata can be filled by service if needed, or by client
	}

	if err := h.transactionService.CreateTransaction(c.Request().Context(), transaction); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) List(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Basic filtering from query params
	filter := make(map[string]interface{})
	if t := c.QueryParam("type"); t != "" {
		filter["type"] = t
	}
	if catID := c.QueryParam("categoryId"); catID != "" {
		id, _ := primitive.ObjectIDFromHex(catID)
		filter["categoryId"] = id
	}
	if bID := c.QueryParam("budgetId"); bID != "" {
		id, _ := primitive.ObjectIDFromHex(bID)
		filter["budgetId"] = id
	}

	transactions, err := h.transactionService.ListUserTransactions(c.Request().Context(), userID, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	transaction, err := h.transactionService.GetTransaction(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if transaction == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "transaction not found"})
	}

	return c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) Update(c echo.Context) error {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}

	req := new(dto.CreateTransactionRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	categoryObjectID, _ := primitive.ObjectIDFromHex(req.CategoryID)
	
	var budgetObjectID *primitive.ObjectID
	if req.BudgetID != "" {
		id, _ := primitive.ObjectIDFromHex(req.BudgetID)
		budgetObjectID = &id
	}

	transaction := &domain.Transaction{
		ID:              objectID,
		CategoryID:      categoryObjectID,
		BudgetID:        budgetObjectID,
		Type:            domain.TransactionType(req.Type),
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
		PaymentMethod:   req.PaymentMethod,
		Notes:           req.Notes,
	}

	if err := h.transactionService.UpdateTransaction(c.Request().Context(), transaction); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.transactionService.DeleteTransaction(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
func (h *TransactionHandler) ListWithFilters(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Parse filters from query parameters
	filters := domain.TransactionFilters{}
	
	// Date range filters
	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = &startDate
		}
	}
	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Set to end of day
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filters.EndDate = &endDate
		}
	}
	
	// Category filter
	if categoryIDStr := c.QueryParam("category_id"); categoryIDStr != "" {
		if categoryID, err := primitive.ObjectIDFromHex(categoryIDStr); err == nil {
			filters.CategoryID = &categoryID
		}
	}
	
	// Type filter
	if typeStr := c.QueryParam("type"); typeStr != "" {
		transactionType := domain.TransactionType(typeStr)
		if transactionType == domain.TransactionTypeIncome || transactionType == domain.TransactionTypeExpense {
			filters.Type = &transactionType
		}
	}
	
	// Pagination
	if pageStr := c.QueryParam("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		}
	}
	if perPageStr := c.QueryParam("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 && perPage <= 100 {
			filters.PerPage = perPage
		}
	}
	
	// Sorting
	if sortBy := c.QueryParam("sort_by"); sortBy != "" {
		if sortBy == "transaction_date" || sortBy == "amount" {
			filters.SortBy = sortBy
		}
	}
	if sortOrder := c.QueryParam("sort_order"); sortOrder != "" {
		if sortOrder == "asc" || sortOrder == "desc" {
			filters.SortOrder = sortOrder
		}
	}

	result, err := h.transactionService.ListUserTransactionsWithFilters(c.Request().Context(), userID, filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "TRANSACTION_ERROR",
				"message": "Failed to load transactions",
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

func (h *TransactionHandler) GetSpendingByCategory(c echo.Context) error {
	userID := c.Get("user_id").(string)
	
	// Parse date range (default to current month)
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)
	
	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}
	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	categories, err := h.transactionService.GetSpendingByCategory(c.Request().Context(), userID, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "TRANSACTION_ERROR",
				"message": "Failed to load spending by category",
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    categories,
	})
}

func (h *TransactionHandler) GetMonthlyTrends(c echo.Context) error {
	userID := c.Get("user_id").(string)
	
	months := 3 // default to last 3 months
	if monthsStr := c.QueryParam("months"); monthsStr != "" {
		if parsed, err := strconv.Atoi(monthsStr); err == nil && parsed > 0 && parsed <= 12 {
			months = parsed
		}
	}

	trends, err := h.transactionService.GetMonthlySpendingTrends(c.Request().Context(), userID, months)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "TRANSACTION_ERROR",
				"message": "Failed to load monthly trends",
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    trends,
	})
}