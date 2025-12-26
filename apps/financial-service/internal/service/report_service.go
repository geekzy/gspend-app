package service

import (
	"context"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportService struct {
	reportRepo      domain.ReportRepository
	transactionRepo domain.TransactionRepository
	budgetRepo      domain.BudgetRepository
}

func NewReportService(reportRepo domain.ReportRepository, transactionRepo domain.TransactionRepository, budgetRepo domain.BudgetRepository) *ReportService {
	return &ReportService{
		reportRepo:      reportRepo,
		transactionRepo: transactionRepo,
		budgetRepo:      budgetRepo,
	}
}

// GetBudgetVsActualReport generates a budget vs actual spending report for a specific month
func (s *ReportService) GetBudgetVsActualReport(ctx context.Context, userID string, month time.Time) (*domain.BudgetVsActualReport, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// If month is zero, use current month
	if month.IsZero() {
		month = time.Now()
	}

	// Get the report data from repository
	report, err := s.reportRepo.GetBudgetVsActualData(ctx, userObjectID, month)
	if err != nil {
		return nil, err
	}

	// Calculate overall variance
	if report.TotalBudgeted > 0 {
		report.OverallVariance = report.TotalSpent - report.TotalBudgeted
	}

	// Calculate percentage used and variance for each category
	for _, item := range report.Categories {
		if item.Budgeted > 0 {
			item.PercentageUsed = (item.Actual / item.Budgeted) * 100
			item.Variance = item.Actual - item.Budgeted
		}
	}

	return report, nil
}

// GetSpendingByCategoryReport generates a spending breakdown by category for a date range
func (s *ReportService) GetSpendingByCategoryReport(ctx context.Context, userID string, startDate, endDate time.Time) (*domain.SpendingByCategoryReport, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Default to current month if dates are not provided
	if startDate.IsZero() || endDate.IsZero() {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, -1) // Last day of the month
	}

	report, err := s.reportRepo.GetSpendingByCategoryData(ctx, userObjectID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate percentages for each category
	if report.TotalSpent > 0 {
		for _, category := range report.Categories {
			category.Percentage = (category.Amount / report.TotalSpent) * 100
		}
	}

	return report, nil
}

// GetMonthlyTrendsReport generates a monthly spending trends report
func (s *ReportService) GetMonthlyTrendsReport(ctx context.Context, userID string, months int) (*domain.MonthlyTrendsReport, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Default to 3 months if not specified
	if months <= 0 {
		months = 3
	}

	report, err := s.reportRepo.GetMonthlyTrendsData(ctx, userObjectID, months)
	if err != nil {
		return nil, err
	}

	// Calculate average spending
	if len(report.MonthlyData) > 0 {
		totalSpending := 0.0
		for _, monthData := range report.MonthlyData {
			totalSpending += monthData.Amount
		}
		report.AverageSpending = totalSpending / float64(len(report.MonthlyData))
	}

	// Determine trend direction
	report.TrendDirection = s.calculateTrendDirection(report.MonthlyData)

	return report, nil
}

// calculateTrendDirection analyzes the trend in monthly spending data
func (s *ReportService) calculateTrendDirection(monthlyData []*domain.MonthlySpending) string {
	if len(monthlyData) < 2 {
		return "stable"
	}

	// Simple trend analysis: compare first half with second half
	midPoint := len(monthlyData) / 2
	firstHalfAvg := 0.0
	secondHalfAvg := 0.0

	for i := 0; i < midPoint; i++ {
		firstHalfAvg += monthlyData[i].Amount
	}
	firstHalfAvg /= float64(midPoint)

	for i := midPoint; i < len(monthlyData); i++ {
		secondHalfAvg += monthlyData[i].Amount
	}
	secondHalfAvg /= float64(len(monthlyData) - midPoint)

	// Calculate percentage change
	if firstHalfAvg == 0 {
		return "stable"
	}

	percentageChange := ((secondHalfAvg - firstHalfAvg) / firstHalfAvg) * 100

	// Threshold for considering a trend significant (5%)
	if percentageChange > 5 {
		return "increasing"
	} else if percentageChange < -5 {
		return "decreasing"
	}

	return "stable"
}