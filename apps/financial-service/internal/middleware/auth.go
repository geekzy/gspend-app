package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/grpc/client"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware validates the JWT token using the Auth Service gRPC client
func AuthMiddleware(authClient *client.AuthGRPCClient) echo.MiddlewareFunc {
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
			
			// Validate token via gRPC
			resp, err := authClient.ValidateToken(context.Background(), tokenString)
			if err != nil || !resp.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}

			// Store user info in context
			c.Set("user_id", resp.UserId)
			c.Set("email", resp.Email)

			return next(c)
		}
	}
}
