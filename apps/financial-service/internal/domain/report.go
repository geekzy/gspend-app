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
	TotalBudgeted   float64                `json:"totalBudgeted"`
	TotalSpent      float64                `json:"totalSpent"`
	OverallVariance float64                `json:"overallVariance"`
}

// BudgetVsActualItem represents individual category budget vs actual
type BudgetVsActualItem struct {
	CategoryID     primitive.ObjectID `json:"categoryId"`
	CategoryName   string             `json:"categoryName"`
	Budgeted       float64            `json:"budgeted"`
	Actual         float64            `json:"actual"`
	Variance       float64            `json:"variance"`
	PercentageUsed float64            `json:"percentageUsed"`
}

// SpendingByCategoryReport represents spending breakdown by category
type SpendingByCategoryReport struct {
	StartDate  time.Time            `json:"startDate"`
	EndDate    time.Time            `json:"endDate"`
	Categories []*CategorySpending  `json:"categories"`
	TotalSpent float64              `json:"totalSpent"`
}

// MonthlyTrendsReport represents spending trends over multiple months
type MonthlyTrendsReport struct {
	Months          int                 `json:"months"`
	MonthlyData     []*MonthlySpending  `json:"monthlyData"`
	AverageSpending float64             `json:"averageSpending"`
	TrendDirection  string              `json:"trendDirection"` // "increasing", "decreasing", "stable"
}

// Enhanced MonthlySpending with more details
type MonthlySpendingDetail struct {
	Month         time.Time          `json:"month"`
	TotalIncome   float64            `json:"totalIncome"`
	TotalExpenses float64            `json:"totalExpenses"`
	NetAmount     float64            `json:"netAmount"`
	TopCategory   *CategorySpending  `json:"topCategory"`
}

// ReportRepository defines methods for generating reports
type ReportRepository interface {
	GetBudgetVsActualData(ctx context.Context, userID primitive.ObjectID, month time.Time) (*BudgetVsActualReport, error)
	GetSpendingByCategoryData(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (*SpendingByCategoryReport, error)
	GetMonthlyTrendsData(ctx context.Context, userID primitive.ObjectID, months int) (*MonthlyTrendsReport, error)
}