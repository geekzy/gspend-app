package service

import (
	"context"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BudgetService struct {
	budgetRepo domain.BudgetRepository
}

func NewBudgetService(budgetRepo domain.BudgetRepository) *BudgetService {
	return &BudgetService{
		budgetRepo: budgetRepo,
	}
}

func (s *BudgetService) CreateBudget(ctx context.Context, budget *domain.Budget) error {
	return s.budgetRepo.Create(ctx, budget)
}

func (s *BudgetService) GetBudget(ctx context.Context, id string) (*domain.Budget, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.budgetRepo.GetByID(ctx, objectID)
}

func (s *BudgetService) ListUserBudgets(ctx context.Context, userID string) ([]*domain.Budget, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return s.budgetRepo.ListByUserID(ctx, objectID)
}

func (s *BudgetService) GetActiveBudget(ctx context.Context, userID string, date time.Time) (*domain.Budget, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	if date.IsZero() {
		date = time.Now()
	}
	return s.budgetRepo.GetActiveByDate(ctx, objectID, date)
}

func (s *BudgetService) UpdateBudget(ctx context.Context, budget *domain.Budget) error {
	return s.budgetRepo.Update(ctx, budget)
}

func (s *BudgetService) DeleteBudget(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.budgetRepo.Delete(ctx, objectID)
}
