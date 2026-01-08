package handler

import (
	"net/http"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BudgetHandler struct {
	budgetService service.BudgetServiceInterface
	validate      *validator.Validate
}

func NewBudgetHandler(budgetService service.BudgetServiceInterface) *BudgetHandler {
	return &BudgetHandler{
		budgetService: budgetService,
		validate:      validator.New(),
	}
}

func (h *BudgetHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	req := new(dto.CreateBudgetRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	var items []domain.BudgetItem
	for _, itemReq := range req.Items {
		categoryObjectID, _ := primitive.ObjectIDFromHex(itemReq.CategoryID)
		items = append(items, domain.BudgetItem{
			CategoryID:    categoryObjectID,
			CategoryName:  itemReq.CategoryName,
			PlannedAmount: itemReq.PlannedAmount,
			Notes:         itemReq.Notes,
		})
	}

	budget := &domain.Budget{
		UserID:      userObjectID,
		Name:        req.Name,
		PeriodType:  domain.PeriodType(req.PeriodType),
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		TotalAmount: req.TotalAmount,
		Items:       items,
	}

	if err := h.budgetService.CreateBudget(c.Request().Context(), budget); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, budget)
}

func (h *BudgetHandler) List(c echo.Context) error {
	userID := c.Get("user_id").(string)

	budgets, err := h.budgetService.ListUserBudgets(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, budgets)
}

func (h *BudgetHandler) GetActive(c echo.Context) error {
	userID := c.Get("user_id").(string)
	dateStr := c.QueryParam("date")
	
	var date time.Time
	if dateStr != "" {
		var err error
		date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date format"})
		}
	} else {
		date = time.Now()
	}

	budget, err := h.budgetService.GetActiveBudget(c.Request().Context(), userID, date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if budget == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "no active budget found"})
	}

	return c.JSON(http.StatusOK, budget)
}

func (h *BudgetHandler) Update(c echo.Context) error {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}

	req := new(dto.CreateBudgetRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	var items []domain.BudgetItem
	for _, itemReq := range req.Items {
		categoryObjectID, _ := primitive.ObjectIDFromHex(itemReq.CategoryID)
		items = append(items, domain.BudgetItem{
			CategoryID:    categoryObjectID,
			CategoryName:  itemReq.CategoryName,
			PlannedAmount: itemReq.PlannedAmount,
			Notes:         itemReq.Notes,
		})
	}

	budget := &domain.Budget{
		ID:          objectID,
		Name:        req.Name,
		PeriodType:  domain.PeriodType(req.PeriodType),
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		TotalAmount: req.TotalAmount,
		Items:       items,
	}

	if err := h.budgetService.UpdateBudget(c.Request().Context(), budget); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, budget)
}

func (h *BudgetHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.budgetService.DeleteBudget(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Budget Item Management Endpoints

func (h *BudgetHandler) AddBudgetItem(c echo.Context) error {
	budgetID := c.Param("id")
	
	req := new(dto.CreateBudgetItemRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	categoryObjectID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category ID"})
	}

	item := &domain.BudgetItem{
		CategoryID:    categoryObjectID,
		CategoryName:  req.CategoryName,
		PlannedAmount: req.PlannedAmount,
		Notes:         req.Notes,
	}

	if err := h.budgetService.AddBudgetItem(c.Request().Context(), budgetID, item); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "budget item added successfully",
		"itemId":  item.ID.Hex(),
	})
}

func (h *BudgetHandler) UpdateBudgetItem(c echo.Context) error {
	budgetID := c.Param("id")
	itemID := c.Param("itemId")
	
	req := new(dto.UpdateBudgetItemRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	categoryObjectID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category ID"})
	}

	item := &domain.BudgetItem{
		CategoryID:    categoryObjectID,
		CategoryName:  req.CategoryName,
		PlannedAmount: req.PlannedAmount,
		Notes:         req.Notes,
	}

	if err := h.budgetService.UpdateBudgetItem(c.Request().Context(), budgetID, itemID, item); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "budget item updated successfully"})
}

func (h *BudgetHandler) DeleteBudgetItem(c echo.Context) error {
	budgetID := c.Param("id")
	itemID := c.Param("itemId")

	if err := h.budgetService.DeleteBudgetItem(c.Request().Context(), budgetID, itemID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *BudgetHandler) GetBudgetItem(c echo.Context) error {
	budgetID := c.Param("id")
	itemID := c.Param("itemId")

	item, err := h.budgetService.GetBudgetItem(c.Request().Context(), budgetID, itemID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if item == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "budget item not found"})
	}

	itemDTO := dto.BudgetItemDTO{
		ID:            item.ID.Hex(),
		CategoryID:    item.CategoryID.Hex(),
		CategoryName:  item.CategoryName,
		PlannedAmount: item.PlannedAmount,
		SpentAmount:   item.SpentAmount,
		Notes:         item.Notes,
	}

	return c.JSON(http.StatusOK, itemDTO)
}
