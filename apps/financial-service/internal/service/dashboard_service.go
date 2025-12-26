package service

import (
	"context"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardService struct {
	dashboardRepo domain.DashboardRepository
}

func NewDashboardService(dashboardRepo domain.DashboardRepository) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
	}
}

func (s *DashboardService) GetDashboardSummary(ctx context.Context, userID string) (*domain.DashboardSummary, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	// Get total balance
	totalBalance, err := s.dashboardRepo.GetTotalBalance(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	// Get monthly income
	monthlyIncome, err := s.dashboardRepo.GetMonthlyIncome(ctx, userObjectID, now)
	if err != nil {
		return nil, err
	}

	// Get monthly expenses
	monthlyExpenses, err := s.dashboardRepo.GetMonthlyExpenses(ctx, userObjectID, now)
	if err != nil {
		return nil, err
	}

	// Get budget progress
	budgetProgress, err := s.dashboardRepo.GetCurrentMonthBudgetProgress(ctx, userObjectID)
	if err != nil {
		return nil, err
	}

	// Get top spending categories (top 3)
	topCategories, err := s.dashboardRepo.GetTopSpendingCategories(ctx, userObjectID, now, 3)
	if err != nil {
		return nil, err
	}

	// Get recent transactions (last 5)
	recentTransactions, err := s.dashboardRepo.GetRecentTransactions(ctx, userObjectID, 5)
	if err != nil {
		return nil, err
	}

	return &domain.DashboardSummary{
		TotalBalance:       totalBalance,
		MonthlyIncome:      monthlyIncome,
		MonthlyExpenses:    monthlyExpenses,
		BudgetProgress:     budgetProgress,
		TopCategories:      topCategories,
		RecentTransactions: recentTransactions,
	}, nil
}

func (s *DashboardService) GetRecentTransactions(ctx context.Context, userID string, limit int) ([]*domain.Transaction, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	return s.dashboardRepo.GetRecentTransactions(ctx, userObjectID, limit)
}

func (s *DashboardService) GetTopSpendingCategories(ctx context.Context, userID string, month time.Time) ([]*domain.CategorySpending, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	return s.dashboardRepo.GetTopSpendingCategories(ctx, userObjectID, month, 10)
}

func (s *DashboardService) GetCurrentMonthBudgetProgress(ctx context.Context, userID string) (*domain.BudgetProgress, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	return s.dashboardRepo.GetCurrentMonthBudgetProgress(ctx, userObjectID)
}