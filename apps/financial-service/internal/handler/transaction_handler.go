package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
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
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

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
