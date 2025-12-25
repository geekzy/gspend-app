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

type CategoryHandler struct {
	categoryService *service.CategoryService
	validate        *validator.Validate
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		validate:        validator.New(),
	}
}

func (h *CategoryHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(string)
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	req := new(dto.CreateCategoryRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	category := &domain.Category{
		UserID:    &userObjectID,
		Name:      req.Name,
		Type:      domain.CategoryType(req.Type),
		Icon:      req.Icon,
		Color:     req.Color,
		SortOrder: req.SortOrder,
		IsSystem:  false,
	}

	if err := h.categoryService.CreateCategory(c.Request().Context(), category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) List(c echo.Context) error {
	userID := c.Get("user_id").(string)
	categoryType := c.QueryParam("type")

	categories, err := h.categoryService.ListUserCategories(c.Request().Context(), userID, domain.CategoryType(categoryType))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) Update(c echo.Context) error {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}

	req := new(dto.CreateCategoryRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	category := &domain.Category{
		ID:        objectID,
		Name:      req.Name,
		Type:      domain.CategoryType(req.Type),
		Icon:      req.Icon,
		Color:     req.Color,
		SortOrder: req.SortOrder,
	}

	if err := h.categoryService.UpdateCategory(c.Request().Context(), category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.categoryService.DeleteCategory(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
