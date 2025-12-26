package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetSummary(c echo.Context) error {
	userID := c.Get("user_id").(string)

	summary, err := h.dashboardService.GetDashboardSummary(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "DASHBOARD_ERROR",
				"message": "Failed to load dashboard data",
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    summary,
	})
}

func (h *DashboardHandler) GetRecentTransactions(c echo.Context) error {
	userID := c.Get("user_id").(string)
	
	limitStr := c.QueryParam("limit")
	limit := 5 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	transactions, err := h.dashboardService.GetRecentTransactions(c.Request().Context(), userID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "DASHBOARD_ERROR",
				"message": "Failed to load recent transactions",
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    transactions,
	})
}

func (h *DashboardHandler) GetTopCategories(c echo.Context) error {
	userID := c.Get("user_id").(string)
	
	// Parse month parameter (default to current month)
	monthStr := c.QueryParam("month")
	month := time.Now()
	if monthStr != "" {
		if parsedMonth, err := time.Parse("2006-01", monthStr); err == nil {
			month = parsedMonth
		}
	}

	categories, err := h.dashboardService.GetTopSpendingCategories(c.Request().Context(), userID, month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "DASHBOARD_ERROR",
				"message": "Failed to load spending categories",
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    categories,
	})
}

func (h *DashboardHandler) GetBudgetProgress(c echo.Context) error {
	userID := c.Get("user_id").(string)

	progress, err := h.dashboardService.GetCurrentMonthBudgetProgress(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "DASHBOARD_ERROR",
				"message": "Failed to load budget progress",
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    progress,
	})
}