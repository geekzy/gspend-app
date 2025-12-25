package service

import (
	"context"
	"testing"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBudgetService_GetActiveBudget(t *testing.T) {
	ctx := context.Background()
	budgetRepo := new(MockBudgetRepository)
	service := NewBudgetService(budgetRepo)

	userID := primitive.NewObjectID()
	now := time.Now()

	t.Run("Get Active Budget - Returns budget if exists", func(t *testing.T) {
		expectedBudget := &domain.Budget{
			ID:     primitive.NewObjectID(),
			UserID: userID,
			Name:   "Monthly Budget",
		}

		budgetRepo.On("GetActiveByDate", ctx, userID, now).Return(expectedBudget, nil)

		budget, err := service.GetActiveBudget(ctx, userID.Hex(), now)

		assert.NoError(t, err)
		assert.Equal(t, expectedBudget, budget)
		budgetRepo.AssertExpectations(t)
	})

	t.Run("Get Active Budget - Fails with invalid userID", func(t *testing.T) {
		budget, err := service.GetActiveBudget(ctx, "invalid-id", now)

		assert.Error(t, err)
		assert.Nil(t, budget)
	})
}

func TestBudgetService_CreateBudget(t *testing.T) {
	ctx := context.Background()
	budgetRepo := new(MockBudgetRepository)
	service := NewBudgetService(budgetRepo)

	t.Run("Create Budget - Succeeds", func(t *testing.T) {
		budget := &domain.Budget{
			Name: "New Budget",
		}

		budgetRepo.On("Create", ctx, budget).Return(nil)

		err := service.CreateBudget(ctx, budget)

		assert.NoError(t, err)
		budgetRepo.AssertExpectations(t)
	})
}
