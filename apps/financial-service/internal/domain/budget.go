package domain

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PeriodType string

const (
	PeriodTypeMonthly   PeriodType = "monthly"
	PeriodTypeQuarterly PeriodType = "quarterly"
	PeriodTypeYearly    PeriodType = "yearly"
)

type BudgetItem struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CategoryID    primitive.ObjectID `bson:"categoryId" json:"categoryId"`
	CategoryName  string             `bson:"categoryName" json:"categoryName"`
	PlannedAmount float64            `bson:"plannedAmount" json:"plannedAmount"`
	SpentAmount   float64            `bson:"spentAmount" json:"spentAmount"`
	Notes         string             `bson:"notes" json:"notes"`
}

type Budget struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	Name        string             `bson:"name" json:"name"`
	PeriodType  PeriodType         `bson:"periodType" json:"periodType"`
	StartDate   time.Time          `bson:"startDate" json:"startDate"`
	EndDate     time.Time          `bson:"endDate" json:"endDate"`
	TotalAmount float64            `bson:"totalAmount" json:"totalAmount"`
	Items       []BudgetItem      `bson:"items" json:"items"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type BudgetRepository interface {
	Create(ctx context.Context, budget *Budget) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Budget, error)
	ListByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Budget, error)
	GetActiveByDate(ctx context.Context, userID primitive.ObjectID, date time.Time) (*Budget, error)
	Update(ctx context.Context, budget *Budget) error
	UpdateSpentAmount(ctx context.Context, budgetID primitive.ObjectID, categoryID primitive.ObjectID, amount float64) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// Budget errors
var (
	ErrBudgetNotFound     = errors.New("budget not found")
	ErrBudgetItemNotFound = errors.New("budget item not found")
)
