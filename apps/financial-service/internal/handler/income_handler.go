package handler

import (
	"net/http"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IncomeHandler struct {
	incomeService *service.IncomeService
	validate      *validator.Validate
}

func NewIncomeHandler(incomeService *service.IncomeService) *IncomeHandler {
	return &IncomeHandler{
		incomeService: incomeService,
		validate:      validator.New(),
	}
}

func (h *IncomeHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	req := new(dto.CreateIncomeRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	income := &domain.Income{
		UserID:        userObjectID,
		Source:        req.Source,
		Amount:        req.Amount,
		Frequency:     domain.Frequency(req.Frequency),
		EffectiveDate: req.EffectiveDate,
	}

	if err := h.incomeService.CreateIncome(c.Request().Context(), income); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, income)
}

func (h *IncomeHandler) List(c echo.Context) error {
	userID := c.Get("user_id").(string)

	incomes, err := h.incomeService.ListUserIncomes(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, incomes)
}

func (h *IncomeHandler) Update(c echo.Context) error {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}

	req := new(dto.CreateIncomeRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	income := &domain.Income{
		ID:            objectID,
		Source:        req.Source,
		Amount:        req.Amount,
		Frequency:     domain.Frequency(req.Frequency),
		EffectiveDate: req.EffectiveDate,
	}

	if err := h.incomeService.UpdateIncome(c.Request().Context(), income); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, income)
}

func (h *IncomeHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.incomeService.DeleteIncome(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
