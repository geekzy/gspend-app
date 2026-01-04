package handler

import (
	"net/http"

	"github.com/geekzy/gspend-app/apps/auth-service/internal/dto"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/service"
	"github.com/go-playground/validator/v10"
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

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	req := new(dto.VerifyEmailRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	err := h.authService.VerifyEmail(c.Request().Context(), req.Token)
	if err != nil {
		if err.Error() == "invalid verification token" || err.Error() == "verification token has expired" {
			return c.JSON(http.StatusBadRequest, dto.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.MessageResponse{
			Success: false,
			Message: "An error occurred while verifying your email",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Success: true,
		Message: "Email verified successfully! You can now log in.",
	})
}

func (h *AuthHandler) ResendVerification(c echo.Context) error {
	req := new(dto.ResendVerificationRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	// Always return success for security (don't reveal if email exists)
	_ = h.authService.ResendVerification(c.Request().Context(), req.Email)

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Success: true,
		Message: "If your email is registered, a verification link has been sent.",
	})
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	req := new(dto.ForgotPasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	// Always return success for security (don't reveal if email exists)
	_ = h.authService.ForgotPassword(c.Request().Context(), req.Email)

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Success: true,
		Message: "If your email is registered, a password reset link has been sent.",
	})
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	req := new(dto.ResetPasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "validation failed", "details": err.Error()})
	}

	err := h.authService.ResetPassword(c.Request().Context(), req.Token, req.NewPassword)
	if err != nil {
		if err.Error() == "invalid reset token" || err.Error() == "reset token has expired" {
			return c.JSON(http.StatusBadRequest, dto.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}
		if err.Error() == "password must be at least 8 characters with uppercase, lowercase, and number" {
			return c.JSON(http.StatusBadRequest, dto.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.MessageResponse{
			Success: false,
			Message: "An error occurred while resetting your password",
		})
	}

	return c.JSON(http.StatusOK, dto.MessageResponse{
		Success: true,
		Message: "Password reset successfully! You can now log in with your new password.",
	})
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

