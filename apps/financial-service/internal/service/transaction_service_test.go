package service

import (
	"context"
	"testing"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockTransactionRepository is a manual mock for domain.TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	args := m.Called(ctx, transaction)
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

func (m *MockTransactionRepository) Update(ctx context.Context, transaction *domain.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockBudgetRepository is a manual mock for domain.BudgetRepository
type MockBudgetRepository struct {
	mock.Mock
}

func (m *MockBudgetRepository) Create(ctx context.Context, budget *domain.Budget) error {
	args := m.Called(ctx, budget)
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

func (m *MockBudgetRepository) Update(ctx context.Context, budget *domain.Budget) error {
	args := m.Called(ctx, budget)
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

func TestTransactionService_CreateTransaction(t *testing.T) {
	ctx := context.Background()
	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()
	budgetID := primitive.NewObjectID()

	t.Run("Create Expense with Budget - Succeeds and updates budget", func(t *testing.T) {
		txRepo := new(MockTransactionRepository)
		budgetRepo := new(MockBudgetRepository)
		service := NewTransactionService(txRepo, budgetRepo)

		transaction := &domain.Transaction{
			UserID:     userID,
			CategoryID: categoryID,
			BudgetID:   &budgetID,
			Type:       domain.TransactionTypeExpense,
			Amount:     100.0,
		}

		txRepo.On("Create", ctx, transaction).Return(nil)
		budgetRepo.On("UpdateSpentAmount", ctx, budgetID, categoryID, 100.0).Return(nil)

		err := service.CreateTransaction(ctx, transaction)

		assert.NoError(t, err)
		txRepo.AssertExpectations(t)
		budgetRepo.AssertExpectations(t)
	})

	t.Run("Create Income - Succeeds and does not update budget", func(t *testing.T) {
		txRepo := new(MockTransactionRepository)
		budgetRepo := new(MockBudgetRepository)
		service := NewTransactionService(txRepo, budgetRepo)

		transaction := &domain.Transaction{
			UserID:     userID,
			CategoryID: categoryID,
			Type:       domain.TransactionTypeIncome,
			Amount:     500.0,
		}

		txRepo.On("Create", ctx, transaction).Return(nil)

		err := service.CreateTransaction(ctx, transaction)

		assert.NoError(t, err)
		txRepo.AssertExpectations(t)
		budgetRepo.AssertNotCalled(t, "UpdateSpentAmount", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestTransactionService_DeleteTransaction(t *testing.T) {
	ctx := context.Background()
	txID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()
	budgetID := primitive.NewObjectID()

	t.Run("Delete Expense with Budget - Succeeds and reverses budget update", func(t *testing.T) {
		txRepo := new(MockTransactionRepository)
		budgetRepo := new(MockBudgetRepository)
		service := NewTransactionService(txRepo, budgetRepo)

		transaction := &domain.Transaction{
			ID:         txID,
			UserID:     userID,
			CategoryID: categoryID,
			BudgetID:   &budgetID,
			Type:       domain.TransactionTypeExpense,
			Amount:     100.0,
		}

		txRepo.On("GetByID", ctx, txID).Return(transaction, nil)
		txRepo.On("Delete", ctx, txID).Return(nil)
		budgetRepo.On("UpdateSpentAmount", ctx, budgetID, categoryID, -100.0).Return(nil)

		err := service.DeleteTransaction(ctx, txID.Hex())

		assert.NoError(t, err)
		txRepo.AssertExpectations(t)
		budgetRepo.AssertExpectations(t)
	})
}

func TestTransactionService_GetTransaction(t *testing.T) {
	ctx := context.Background()
	txID := primitive.NewObjectID()

	t.Run("Get Transaction - Success", func(t *testing.T) {
		txRepo := new(MockTransactionRepository)
		budgetRepo := new(MockBudgetRepository)
		service := NewTransactionService(txRepo, budgetRepo)

		expected := &domain.Transaction{ID: txID}
		txRepo.On("GetByID", ctx, txID).Return(expected, nil)

		res, err := service.GetTransaction(ctx, txID.Hex())
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("Get Transaction - Invalid ID", func(t *testing.T) {
		txRepo := new(MockTransactionRepository)
		budgetRepo := new(MockBudgetRepository)
		service := NewTransactionService(txRepo, budgetRepo)

		res, err := service.GetTransaction(ctx, "invalid")
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestTransactionService_ListUserTransactions(t *testing.T) {
	ctx := context.Background()
	userID := primitive.NewObjectID()

	t.Run("List User Transactions - Success", func(t *testing.T) {
		txRepo := new(MockTransactionRepository)
		budgetRepo := new(MockBudgetRepository)
		service := NewTransactionService(txRepo, budgetRepo)

		expected := []*domain.Transaction{{ID: primitive.NewObjectID()}}
		txRepo.On("ListByUserID", ctx, userID, mock.Anything).Return(expected, nil)

		res, err := service.ListUserTransactions(ctx, userID.Hex(), nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestTransactionService_UpdateTransaction(t *testing.T) {
	ctx := context.Background()
	txRepo := new(MockTransactionRepository)
	budgetRepo := new(MockBudgetRepository)
	service := NewTransactionService(txRepo, budgetRepo)

	t.Run("Update Transaction - Success", func(t *testing.T) {
		transaction := &domain.Transaction{ID: primitive.NewObjectID()}
		txRepo.On("Update", ctx, transaction).Return(nil)

		err := service.UpdateTransaction(ctx, transaction)
		assert.NoError(t, err)
	})
}
