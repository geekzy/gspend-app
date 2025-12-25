package service

import (
	"context"

	"github.com/imam/gspend-app/apps/financial-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IncomeService struct {
	incomeRepo domain.IncomeRepository
}

func NewIncomeService(incomeRepo domain.IncomeRepository) *IncomeService {
	return &IncomeService{
		incomeRepo: incomeRepo,
	}
}

func (s *IncomeService) CreateIncome(ctx context.Context, income *domain.Income) error {
	return s.incomeRepo.Create(ctx, income)
}

func (s *IncomeService) GetIncome(ctx context.Context, id string) (*domain.Income, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.incomeRepo.GetByID(ctx, objectID)
}

func (s *IncomeService) ListUserIncomes(ctx context.Context, userID string) ([]*domain.Income, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return s.incomeRepo.ListByUserID(ctx, objectID)
}

func (s *IncomeService) UpdateIncome(ctx context.Context, income *domain.Income) error {
	return s.incomeRepo.Update(ctx, income)
}

func (s *IncomeService) DeleteIncome(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.incomeRepo.Delete(ctx, objectID)
}
