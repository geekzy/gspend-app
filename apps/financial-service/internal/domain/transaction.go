package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type TransactionMetadata struct {
	CategoryName string `bson:"categoryName" json:"categoryName"`
	BudgetName   string `bson:"budgetName,omitempty" json:"budgetName,omitempty"`
}

type Transaction struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID  `bson:"userId" json:"userId"`
	CategoryID      primitive.ObjectID  `bson:"categoryId" json:"categoryId"`
	BudgetID        *primitive.ObjectID `bson:"budgetId,omitempty" json:"budgetId,omitempty"`
	Type            TransactionType     `bson:"type" json:"type"`
	Amount          float64             `bson:"amount" json:"amount"`
	Description     string              `bson:"description" json:"description"`
	TransactionDate time.Time           `bson:"transactionDate" json:"transactionDate"`
	PaymentMethod   string              `bson:"paymentMethod" json:"paymentMethod"`
	Notes           string              `bson:"notes" json:"notes"`
	Metadata        TransactionMetadata `bson:"metadata" json:"metadata"`
	CreatedAt       time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Transaction, error)
	ListByUserID(ctx context.Context, userID primitive.ObjectID, filter map[string]interface{}) ([]*Transaction, error)
	Update(ctx context.Context, transaction *Transaction) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
