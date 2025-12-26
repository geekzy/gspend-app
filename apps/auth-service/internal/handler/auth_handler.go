package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := new(dto.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	resp, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(dto.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	resp, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)

	resp, err := h.authService.GetProfile(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)

	req := new(dto.UpdateProfileRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	user, err := h.authService.UpdateProfile(c.Request().Context(), userID, req)
	if err != nil {
		if err.Error() == "email already in use" {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		if err.Error() == "family size must be between 0 and 5 children" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ProfileResponse{
		Success: true,
		Message: "Profile updated successfully",
		User:    *user,
	})
}

func (h *AuthHandler) ChangePassword(c echo.Context) error {
	userID := c.Get("user_id").(string)

	req := new(dto.ChangePasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	err := h.authService.ChangePassword(c.Request().Context(), userID, req)
	if err != nil {
		if err.Error() == "current password is incorrect" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		if err.Error() == "new password must be at least 8 characters with uppercase, lowercase, and number" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ProfileResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}
