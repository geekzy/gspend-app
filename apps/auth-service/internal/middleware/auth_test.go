package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/geekzy/gspend-app/apps/auth-service/internal/config"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
		},
	}

	handler := AuthMiddleware(cfg)(func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "missing authorization header")
	})

	t.Run("Invalid Header Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidFormat token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid authorization header format")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid or expired token")
	})

	t.Run("Success", func(t *testing.T) {
		token, _ := util.GenerateToken("user-123", "test@example.com", cfg.JWT.Secret, time.Hour)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "user-123", c.Get("user_id"))
		assert.Equal(t, "test@example.com", c.Get("email"))
	})
}
