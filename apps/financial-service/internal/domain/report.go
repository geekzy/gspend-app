package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BudgetVsActualReport represents budget vs actual spending comparison
type BudgetVsActualReport struct {
	Month           time.Time              `json:"month"`
	Categories      []*BudgetVsActualItem  `json:"categories"`
	TotalBudgeted   float64                `json:"total_budgeted"`
	TotalSpent      float64                `json:"total_spent"`
	OverallVariance float64                `json:"overall_variance"`
}

// BudgetVsActualItem represents individual category budget vs actual
type BudgetVsActualItem struct {
	CategoryID     primitive.ObjectID `json:"category_id"`
	CategoryName   string             `json:"category_name"`
	Budgeted       float64            `json:"budgeted"`
	Actual         float64            `json:"actual"`
	Variance       float64            `json:"variance"`
	PercentageUsed float64            `json:"percentage_used"`
}

// SpendingByCategoryReport represents spending breakdown by category
type SpendingByCategoryReport struct {
	StartDate  time.Time            `json:"start_date"`
	EndDate    time.Time            `json:"end_date"`
	Categories []*CategorySpending  `json:"categories"`
	TotalSpent float64              `json:"total_spent"`
}

// MonthlyTrendsReport represents spending trends over multiple months
type MonthlyTrendsReport struct {
	Months          int                 `json:"months"`
	MonthlyData     []*MonthlySpending  `json:"monthly_data"`
	AverageSpending float64             `json:"average_spending"`
	TrendDirection  string              `json:"trend_direction"` // "increasing", "decreasing", "stable"
}

// Enhanced MonthlySpending with more details
type MonthlySpendingDetail struct {
	Month         time.Time          `json:"month"`
	TotalIncome   float64            `json:"total_income"`
	TotalExpenses float64            `json:"total_expenses"`
	NetAmount     float64            `json:"net_amount"`
	TopCategory   *CategorySpending  `json:"top_category"`
}

// ReportRepository defines methods for generating reports
type ReportRepository interface {
	GetBudgetVsActualData(ctx context.Context, userID primitive.ObjectID, month time.Time) (*BudgetVsActualReport, error)
	GetSpendingByCategoryData(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (*SpendingByCategoryReport, error)
	GetMonthlyTrendsData(ctx context.Context, userID primitive.ObjectID, months int) (*MonthlyTrendsReport, error)
}