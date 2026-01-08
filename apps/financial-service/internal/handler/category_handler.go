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
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

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

	// Get existing category first
	existingCategory, err := h.categoryService.GetCategory(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if existingCategory == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "category not found"})
	}

	req := new(dto.UpdateCategoryRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// Update only the provided fields
	category := &domain.Category{
		ID:        objectID,
		UserID:    existingCategory.UserID,
		Name:      req.Name,
		Type:      existingCategory.Type, // Keep existing type
		Icon:      req.Icon,
		Color:     req.Color,
		SortOrder: req.SortOrder,
		IsSystem:  existingCategory.IsSystem,
	}

	if err := h.categoryService.UpdateCategory(c.Request().Context(), category); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}

	// Check if category is system protected
	if err := h.categoryService.ValidateSystemCategoryProtection(c.Request().Context(), id); err != nil {
		if err == domain.ErrSystemCategoryProtected {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "System categories cannot be deleted"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err := h.categoryService.DeleteCategory(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
