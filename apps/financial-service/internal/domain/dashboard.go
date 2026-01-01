package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DashboardSummary represents the main dashboard data
type DashboardSummary struct {
	TotalBalance       float64                `json:"totalBalance"`
	MonthlyIncome      float64                `json:"monthlyIncome"`
	MonthlyExpenses    float64                `json:"monthlyExpenses"`
	BudgetProgress     *BudgetProgress        `json:"budgetProgress"`
	TopCategories      []*CategorySpending    `json:"topCategories"`
	RecentTransactions []*Transaction         `json:"recentTransactions"`
}

// BudgetProgress represents budget vs spending progress
type BudgetProgress struct {
	TotalBudget     float64 `json:"totalBudget"`
	TotalSpent      float64 `json:"totalSpent"`
	RemainingBudget float64 `json:"remainingBudget"`
	PercentageUsed  float64 `json:"percentageUsed"`
}

// CategorySpending represents spending by category
type CategorySpending struct {
	CategoryID   primitive.ObjectID `json:"categoryId"`
	CategoryName string             `json:"categoryName"`
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
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	CategoryID *primitive.ObjectID `json:"categoryId"`
	Type      *TransactionType `json:"type"`
	Page      int        `json:"page"`
	PerPage   int        `json:"perPage"`
	SortBy    string     `json:"sortBy"`    // "transaction_date" or "amount"
	SortOrder string     `json:"sortOrder"` // "asc" or "desc"
}

// PaginatedTransactions represents paginated transaction results
type PaginatedTransactions struct {
	Transactions []*Transaction `json:"transactions"`
	Pagination   Pagination     `json:"pagination"`
	FiltersApplied TransactionFilters `json:"filtersApplied"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
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