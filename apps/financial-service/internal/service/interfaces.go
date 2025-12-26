package service

import (
	"context"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
)

// BudgetServiceInterface defines the interface for budget service operations
type BudgetServiceInterface interface {
	CreateBudget(ctx context.Context, budget *domain.Budget) error
	GetBudget(ctx context.Context, id string) (*domain.Budget, error)
	ListUserBudgets(ctx context.Context, userID string) ([]*domain.Budget, error)
	GetActiveBudget(ctx context.Context, userID string, date time.Time) (*domain.Budget, error)
	UpdateBudget(ctx context.Context, budget *domain.Budget) error
	DeleteBudget(ctx context.Context, id string) error
	AddBudgetItem(ctx context.Context, budgetID string, item *domain.BudgetItem) error
	UpdateBudgetItem(ctx context.Context, budgetID string, itemID string, item *domain.BudgetItem) error
	DeleteBudgetItem(ctx context.Context, budgetID string, itemID string) error
	GetBudgetItem(ctx context.Context, budgetID string, itemID string) (*domain.BudgetItem, error)
}