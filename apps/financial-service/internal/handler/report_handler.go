package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/service"
	"github.com/labstack/echo/v4"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GetBudgetVsActual handles GET /api/v1/reports/budget-vs-actual
func (h *ReportHandler) GetBudgetVsActual(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Parse month parameter (default to current month)
	monthStr := c.QueryParam("month")
	month := time.Now()
	if monthStr != "" {
		if parsedMonth, err := time.Parse("2006-01", monthStr); err == nil {
			month = parsedMonth
		}
	}

	report, err := h.reportService.GetBudgetVsActualReport(c.Request().Context(), userID, month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "REPORT_ERROR",
				"message": "Failed to generate budget vs actual report",
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    report,
	})
}

// GetSpendingByCategory handles GET /api/v1/reports/spending-by-category
func (h *ReportHandler) GetSpendingByCategory(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Parse date range parameters
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": map[string]string{
					"code":    "VALIDATION_ERROR",
					"message": "Invalid start_date format. Use YYYY-MM-DD",
				},
			})
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": map[string]string{
					"code":    "VALIDATION_ERROR",
					"message": "Invalid end_date format. Use YYYY-MM-DD",
				},
			})
		}
	}

	// Validate date range
	if !startDate.IsZero() && !endDate.IsZero() && startDate.After(endDate) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "VALIDATION_ERROR",
				"message": "Start date must be before end date",
			},
		})
	}

	report, err := h.reportService.GetSpendingByCategoryReport(c.Request().Context(), userID, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "REPORT_ERROR",
				"message": "Failed to generate spending by category report",
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    report,
	})
}

// GetMonthlyTrends handles GET /api/v1/reports/monthly-trends
func (h *ReportHandler) GetMonthlyTrends(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Parse months parameter (default to 3)
	monthsStr := c.QueryParam("months")
	months := 3 // default
	if monthsStr != "" {
		if parsedMonths, err := strconv.Atoi(monthsStr); err == nil && parsedMonths > 0 && parsedMonths <= 12 {
			months = parsedMonths
		}
	}

	report, err := h.reportService.GetMonthlyTrendsReport(c.Request().Context(), userID, months)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "REPORT_ERROR",
				"message": "Failed to generate monthly trends report",
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    report,
	})
}