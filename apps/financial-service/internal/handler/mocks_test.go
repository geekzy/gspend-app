package handler

import (
	"context"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockTransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) ListByUserID(ctx context.Context, userID primitive.ObjectID, filter map[string]interface{}) ([]*domain.Transaction, error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Enhanced methods for EnhancedTransactionRepository interface
func (m *MockTransactionRepository) FindWithFilters(ctx context.Context, userID primitive.ObjectID, filters domain.TransactionFilters) ([]*domain.Transaction, int, error) {
	args := m.Called(ctx, userID, filters)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Transaction), args.Int(1), args.Error(2)
}

func (m *MockTransactionRepository) GetSpendingByCategory(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*domain.CategorySpending, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CategorySpending), args.Error(1)
}

func (m *MockTransactionRepository) GetMonthlySpendingTrends(ctx context.Context, userID primitive.ObjectID, months int) ([]*domain.MonthlySpending, error) {
	args := m.Called(ctx, userID, months)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.MonthlySpending), args.Error(1)
}

// MockBudgetRepository
type MockBudgetRepository struct {
	mock.Mock
}

func (m *MockBudgetRepository) Create(ctx context.Context, b *domain.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBudgetRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Budget, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) ListByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Budget, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) GetActiveByDate(ctx context.Context, userID primitive.ObjectID, date time.Time) (*domain.Budget, error) {
	args := m.Called(ctx, userID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) Update(ctx context.Context, b *domain.Budget) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBudgetRepository) UpdateSpentAmount(ctx context.Context, budgetID primitive.ObjectID, categoryID primitive.ObjectID, amount float64) error {
	args := m.Called(ctx, budgetID, categoryID, amount)
	return args.Error(0)
}

func (m *MockBudgetRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockIncomeRepository
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
