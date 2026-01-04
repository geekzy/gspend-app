package middleware

import (
	"net/http"
	"strings"

	"github.com/geekzy/gspend-app/apps/auth-service/internal/config"
	"github.com/geekzy/gspend-app/apps/auth-service/internal/util"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware validates the JWT token in the Authorization header
func AuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization header format"})
			}

			tokenString := parts[1]
			claims, err := util.ValidateToken(tokenString, cfg.JWT.Secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}

			// Store user info in context
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}
