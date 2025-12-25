package service

import (
	"context"
	"testing"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockIncomeRepository is a manual mock for domain.IncomeRepository
type MockIncomeRepository struct {
	mock.Mock
}

func (m *MockIncomeRepository) Create(ctx context.Context, income *domain.Income) error {
	args := m.Called(ctx, income)
	return args.Error(0)
}

func (m *MockIncomeRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Income, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Income), args.Error(1)
}

func (m *MockIncomeRepository) ListByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Income, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Income), args.Error(1)
}

func (m *MockIncomeRepository) Update(ctx context.Context, income *domain.Income) error {
	args := m.Called(ctx, income)
	return args.Error(0)
}

func (m *MockIncomeRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestIncomeService(t *testing.T) {
	ctx := context.Background()
	userID := primitive.NewObjectID()
	incomeID := primitive.NewObjectID()

	t.Run("Create Income", func(t *testing.T) {
		repo := new(MockIncomeRepository)
		service := NewIncomeService(repo)
		income := &domain.Income{Source: "Salary", Amount: 5000}

		repo.On("Create", ctx, income).Return(nil)
		err := service.CreateIncome(ctx, income)
		assert.NoError(t, err)
	})

	t.Run("Get Income", func(t *testing.T) {
		repo := new(MockIncomeRepository)
		service := NewIncomeService(repo)
		expected := &domain.Income{ID: incomeID, Source: "Salary"}

		repo.On("GetByID", ctx, incomeID).Return(expected, nil)
		res, err := service.GetIncome(ctx, incomeID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("List User Incomes", func(t *testing.T) {
		repo := new(MockIncomeRepository)
		service := NewIncomeService(repo)
		expected := []*domain.Income{{ID: incomeID}}

		repo.On("ListByUserID", ctx, userID).Return(expected, nil)
		res, err := service.ListUserIncomes(ctx, userID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("Update Income", func(t *testing.T) {
		repo := new(MockIncomeRepository)
		service := NewIncomeService(repo)
		income := &domain.Income{ID: incomeID, Source: "Bonus"}

		repo.On("Update", ctx, income).Return(nil)
		err := service.UpdateIncome(ctx, income)
		assert.NoError(t, err)
	})

	t.Run("Delete Income", func(t *testing.T) {
		repo := new(MockIncomeRepository)
		service := NewIncomeService(repo)

		repo.On("Delete", ctx, incomeID).Return(nil)
		err := service.DeleteIncome(ctx, incomeID.Hex())
		assert.NoError(t, err)
	})
}
