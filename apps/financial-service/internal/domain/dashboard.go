package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DashboardSummary represents the main dashboard data
type DashboardSummary struct {
	TotalBalance       float64                `json:"total_balance"`
	MonthlyIncome      float64                `json:"monthly_income"`
	MonthlyExpenses    float64                `json:"monthly_expenses"`
	BudgetProgress     *BudgetProgress        `json:"budget_progress"`
	TopCategories      []*CategorySpending    `json:"top_categories"`
	RecentTransactions []*Transaction         `json:"recent_transactions"`
}

// BudgetProgress represents budget vs spending progress
type BudgetProgress struct {
	TotalBudget     float64 `json:"total_budget"`
	TotalSpent      float64 `json:"total_spent"`
	RemainingBudget float64 `json:"remaining_budget"`
	PercentageUsed  float64 `json:"percentage_used"`
}

// CategorySpending represents spending by category
type CategorySpending struct {
	CategoryID   primitive.ObjectID `json:"category_id"`
	CategoryName string             `json:"category_name"`
	Amount       float64            `json:"amount"`
	Percentage   float64            `json:"percentage"`
}

// MonthlySpending represents spending trends by month
type MonthlySpending struct {
	Month  time.Time `json:"month"`
	Amount float64   `json:"amount"`
}

// TransactionFilters represents filtering options for transactions
type TransactionFilters struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	CategoryID *primitive.ObjectID `json:"category_id"`
	Type      *TransactionType `json:"type"`
	Page      int        `json:"page"`
	PerPage   int        `json:"per_page"`
	SortBy    string     `json:"sort_by"`    // "transaction_date" or "amount"
	SortOrder string     `json:"sort_order"` // "asc" or "desc"
}

// PaginatedTransactions represents paginated transaction results
type PaginatedTransactions struct {
	Transactions []*Transaction `json:"transactions"`
	Pagination   Pagination     `json:"pagination"`
	FiltersApplied TransactionFilters `json:"filters_applied"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// DashboardRepository defines methods for dashboard data aggregation
type DashboardRepository interface {
	GetTotalBalance(ctx context.Context, userID primitive.ObjectID) (float64, error)
	GetMonthlyIncome(ctx context.Context, userID primitive.ObjectID, month time.Time) (float64, error)
	GetMonthlyExpenses(ctx context.Context, userID primitive.ObjectID, month time.Time) (float64, error)
	GetTopSpendingCategories(ctx context.Context, userID primitive.ObjectID, month time.Time, limit int) ([]*CategorySpending, error)
	GetRecentTransactions(ctx context.Context, userID primitive.ObjectID, limit int) ([]*Transaction, error)
	GetCurrentMonthBudgetProgress(ctx context.Context, userID primitive.ObjectID) (*BudgetProgress, error)
}

// Enhanced TransactionRepository interface with filtering
type EnhancedTransactionRepository interface {
	TransactionRepository
	FindWithFilters(ctx context.Context, userID primitive.ObjectID, filters TransactionFilters) ([]*Transaction, int, error)
	GetSpendingByCategory(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*CategorySpending, error)
	GetMonthlySpendingTrends(ctx context.Context, userID primitive.ObjectID, months int) ([]*MonthlySpending, error)
}