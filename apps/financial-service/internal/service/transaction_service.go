package service

import (
	"context"
	"math"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionService struct {
	transactionRepo domain.EnhancedTransactionRepository
	budgetRepo      domain.BudgetRepository
}

func NewTransactionService(transactionRepo domain.EnhancedTransactionRepository, budgetRepo domain.BudgetRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		budgetRepo:      budgetRepo,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return err
	}

	// If transaction is an expense and has a budget, update the budget spent amount
	if transaction.Type == domain.TransactionTypeExpense && transaction.BudgetID != nil {
		return s.budgetRepo.UpdateSpentAmount(ctx, *transaction.BudgetID, transaction.CategoryID, transaction.Amount)
	}

	return nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, id string) (*domain.Transaction, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.transactionRepo.GetByID(ctx, objectID)
}

func (s *TransactionService) ListUserTransactions(ctx context.Context, userID string, filter map[string]interface{}) ([]*domain.Transaction, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	return s.transactionRepo.ListByUserID(ctx, objectID, filter)
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	// For simplicity, we assume the amount/budget/category don't change in a way that requires
	// reverting and re-applying budget updates. In a production app, we would handle this.
	return s.transactionRepo.Update(ctx, transaction)
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	transaction, err := s.transactionRepo.GetByID(ctx, objectID)
	if err != nil {
		return err
	}
	if transaction == nil {
		return nil
	}

	if err := s.transactionRepo.Delete(ctx, objectID); err != nil {
		return err
	}

	// If transaction was an expense and had a budget, update the budget spent amount (decrease)
	if transaction.Type == domain.TransactionTypeExpense && transaction.BudgetID != nil {
		return s.budgetRepo.UpdateSpentAmount(ctx, *transaction.BudgetID, transaction.CategoryID, -transaction.Amount)
	}

	return nil
}
func (s *TransactionService) ListUserTransactionsWithFilters(ctx context.Context, userID string, filters domain.TransactionFilters) (*domain.PaginatedTransactions, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	transactions, total, err := s.transactionRepo.FindWithFilters(ctx, userObjectID, filters)
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	page := filters.Page
	if page < 1 {
		page = 1
	}
	perPage := filters.PerPage
	if perPage < 1 {
		perPage = 20
	}
	
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &domain.PaginatedTransactions{
		Transactions: transactions,
		Pagination: domain.Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
		FiltersApplied: filters,
	}, nil
}

func (s *TransactionService) GetSpendingByCategory(ctx context.Context, userID string, startDate, endDate time.Time) ([]*domain.CategorySpending, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	return s.transactionRepo.GetSpendingByCategory(ctx, userObjectID, startDate, endDate)
}

func (s *TransactionService) GetMonthlySpendingTrends(ctx context.Context, userID string, months int) ([]*domain.MonthlySpending, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	return s.transactionRepo.GetMonthlySpendingTrends(ctx, userObjectID, months)
}