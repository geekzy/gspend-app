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

// Budget Item Management Methods

func (s *BudgetService) AddBudgetItem(ctx context.Context, budgetID string, item *domain.BudgetItem) error {
	budgetObjectID, err := primitive.ObjectIDFromHex(budgetID)
	if err != nil {
		return err
	}

	// Get the budget first
	budget, err := s.budgetRepo.GetByID(ctx, budgetObjectID)
	if err != nil {
		return err
	}
	if budget == nil {
		return domain.ErrBudgetNotFound
	}

	// Generate ID for new item
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}

	// Add item to budget
	budget.Items = append(budget.Items, *item)
	
	// Recalculate total amount
	s.recalculateBudgetTotal(budget)

	return s.budgetRepo.Update(ctx, budget)
}

func (s *BudgetService) UpdateBudgetItem(ctx context.Context, budgetID string, itemID string, updatedItem *domain.BudgetItem) error {
	budgetObjectID, err := primitive.ObjectIDFromHex(budgetID)
	if err != nil {
		return err
	}

	itemObjectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return err
	}

	// Get the budget first
	budget, err := s.budgetRepo.GetByID(ctx, budgetObjectID)
	if err != nil {
		return err
	}
	if budget == nil {
		return domain.ErrBudgetNotFound
	}

	// Find and update the item
	found := false
	for i, item := range budget.Items {
		if item.ID == itemObjectID {
			// Preserve ID and spent amount, update other fields
			budget.Items[i].CategoryID = updatedItem.CategoryID
			budget.Items[i].CategoryName = updatedItem.CategoryName
			budget.Items[i].PlannedAmount = updatedItem.PlannedAmount
			budget.Items[i].Notes = updatedItem.Notes
			found = true
			break
		}
	}

	if !found {
		return domain.ErrBudgetItemNotFound
	}

	// Recalculate total amount
	s.recalculateBudgetTotal(budget)

	return s.budgetRepo.Update(ctx, budget)
}

func (s *BudgetService) DeleteBudgetItem(ctx context.Context, budgetID string, itemID string) error {
	budgetObjectID, err := primitive.ObjectIDFromHex(budgetID)
	if err != nil {
		return err
	}

	itemObjectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return err
	}

	// Get the budget first
	budget, err := s.budgetRepo.GetByID(ctx, budgetObjectID)
	if err != nil {
		return err
	}
	if budget == nil {
		return domain.ErrBudgetNotFound
	}

	// Find and remove the item
	found := false
	for i, item := range budget.Items {
		if item.ID == itemObjectID {
			// Remove item from slice
			budget.Items = append(budget.Items[:i], budget.Items[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return domain.ErrBudgetItemNotFound
	}

	// Recalculate total amount
	s.recalculateBudgetTotal(budget)

	return s.budgetRepo.Update(ctx, budget)
}

func (s *BudgetService) GetBudgetItem(ctx context.Context, budgetID string, itemID string) (*domain.BudgetItem, error) {
	budgetObjectID, err := primitive.ObjectIDFromHex(budgetID)
	if err != nil {
		return nil, err
	}

	itemObjectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		return nil, err
	}

	// Get the budget first
	budget, err := s.budgetRepo.GetByID(ctx, budgetObjectID)
	if err != nil {
		return nil, err
	}
	if budget == nil {
		return nil, domain.ErrBudgetNotFound
	}

	// Find the item
	for _, item := range budget.Items {
		if item.ID == itemObjectID {
			return &item, nil
		}
	}

	return nil, domain.ErrBudgetItemNotFound
}

// recalculateBudgetTotal updates the budget's total amount based on all items
func (s *BudgetService) recalculateBudgetTotal(budget *domain.Budget) {
	total := 0.0
	for _, item := range budget.Items {
		total += item.PlannedAmount
	}
	budget.TotalAmount = total
}
