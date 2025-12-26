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

func TestBudgetService_AddBudgetItem(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)
	ctx := context.Background()

	budgetID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	// Create a test budget
	budget := &domain.Budget{
		ID:          budgetID,
		UserID:      primitive.NewObjectID(),
		Name:        "Test Budget",
		PeriodType:  domain.PeriodTypeMonthly,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
		TotalAmount: 1000.0,
		Items: []domain.BudgetItem{
			{
				ID:            primitive.NewObjectID(),
				CategoryID:    primitive.NewObjectID(),
				CategoryName:  "Existing Category",
				PlannedAmount: 500.0,
				SpentAmount:   0.0,
			},
		},
	}

	newItem := &domain.BudgetItem{
		CategoryID:    categoryID,
		CategoryName:  "New Category",
		PlannedAmount: 300.0,
		Notes:         "Test notes",
	}

	// Mock repository calls
	mockRepo.On("GetByID", ctx, budgetID).Return(budget, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Budget")).Return(nil)

	// Execute
	err := service.AddBudgetItem(ctx, budgetID.Hex(), newItem)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, newItem.ID) // ID should be generated
	mockRepo.AssertExpectations(t)
}

func TestBudgetService_UpdateBudgetItem(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)
	ctx := context.Background()

	budgetID := primitive.NewObjectID()
	itemID := primitive.NewObjectID()
	categoryID := primitive.NewObjectID()

	// Create a test budget with an existing item
	budget := &domain.Budget{
		ID:          budgetID,
		UserID:      primitive.NewObjectID(),
		Name:        "Test Budget",
		PeriodType:  domain.PeriodTypeMonthly,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
		TotalAmount: 1000.0,
		Items: []domain.BudgetItem{
			{
				ID:            itemID,
				CategoryID:    primitive.NewObjectID(),
				CategoryName:  "Old Category",
				PlannedAmount: 500.0,
				SpentAmount:   100.0,
			},
		},
	}

	updatedItem := &domain.BudgetItem{
		CategoryID:    categoryID,
		CategoryName:  "Updated Category",
		PlannedAmount: 600.0,
		Notes:         "Updated notes",
	}

	// Mock repository calls
	mockRepo.On("GetByID", ctx, budgetID).Return(budget, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Budget")).Return(nil)

	// Execute
	err := service.UpdateBudgetItem(ctx, budgetID.Hex(), itemID.Hex(), updatedItem)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBudgetService_DeleteBudgetItem(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)
	ctx := context.Background()

	budgetID := primitive.NewObjectID()
	itemID := primitive.NewObjectID()

	// Create a test budget with items
	budget := &domain.Budget{
		ID:          budgetID,
		UserID:      primitive.NewObjectID(),
		Name:        "Test Budget",
		PeriodType:  domain.PeriodTypeMonthly,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
		TotalAmount: 1000.0,
		Items: []domain.BudgetItem{
			{
				ID:            itemID,
				CategoryID:    primitive.NewObjectID(),
				CategoryName:  "Category to Delete",
				PlannedAmount: 500.0,
				SpentAmount:   100.0,
			},
			{
				ID:            primitive.NewObjectID(),
				CategoryID:    primitive.NewObjectID(),
				CategoryName:  "Other Category",
				PlannedAmount: 300.0,
				SpentAmount:   50.0,
			},
		},
	}

	// Mock repository calls
	mockRepo.On("GetByID", ctx, budgetID).Return(budget, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Budget")).Return(nil)

	// Execute
	err := service.DeleteBudgetItem(ctx, budgetID.Hex(), itemID.Hex())

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBudgetService_GetBudgetItem(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)
	ctx := context.Background()

	budgetID := primitive.NewObjectID()
	itemID := primitive.NewObjectID()

	// Create a test budget with an item
	budget := &domain.Budget{
		ID:          budgetID,
		UserID:      primitive.NewObjectID(),
		Name:        "Test Budget",
		PeriodType:  domain.PeriodTypeMonthly,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
		TotalAmount: 1000.0,
		Items: []domain.BudgetItem{
			{
				ID:            itemID,
				CategoryID:    primitive.NewObjectID(),
				CategoryName:  "Test Category",
				PlannedAmount: 500.0,
				SpentAmount:   100.0,
				Notes:         "Test notes",
			},
		},
	}

	// Mock repository calls
	mockRepo.On("GetByID", ctx, budgetID).Return(budget, nil)

	// Execute
	item, err := service.GetBudgetItem(ctx, budgetID.Hex(), itemID.Hex())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, itemID, item.ID)
	assert.Equal(t, "Test Category", item.CategoryName)
	assert.Equal(t, 500.0, item.PlannedAmount)
	assert.Equal(t, 100.0, item.SpentAmount)
	assert.Equal(t, "Test notes", item.Notes)
	mockRepo.AssertExpectations(t)
}

func TestBudgetService_RecalculateBudgetTotal(t *testing.T) {
	service := NewBudgetService(nil)

	budget := &domain.Budget{
		TotalAmount: 0.0, // Will be recalculated
		Items: []domain.BudgetItem{
			{PlannedAmount: 500.0},
			{PlannedAmount: 300.0},
			{PlannedAmount: 200.0},
		},
	}

	// Execute
	service.recalculateBudgetTotal(budget)

	// Assert
	assert.Equal(t, 1000.0, budget.TotalAmount)
}

func TestBudgetService_BudgetItemNotFound(t *testing.T) {
	mockRepo := new(MockBudgetRepository)
	service := NewBudgetService(mockRepo)
	ctx := context.Background()

	budgetID := primitive.NewObjectID()
	nonExistentItemID := primitive.NewObjectID()

	// Create a test budget without the requested item
	budget := &domain.Budget{
		ID:          budgetID,
		UserID:      primitive.NewObjectID(),
		Name:        "Test Budget",
		PeriodType:  domain.PeriodTypeMonthly,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 1, 0),
		TotalAmount: 1000.0,
		Items:       []domain.BudgetItem{}, // Empty items
	}

	// Mock repository calls
	mockRepo.On("GetByID", ctx, budgetID).Return(budget, nil)

	// Execute
	item, err := service.GetBudgetItem(ctx, budgetID.Hex(), nonExistentItemID.Hex())

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrBudgetItemNotFound, err)
	assert.Nil(t, item)
	mockRepo.AssertExpectations(t)
}