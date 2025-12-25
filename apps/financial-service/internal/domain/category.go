package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type Category struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    *primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"` // null for system categories
	Name      string             `bson:"name" json:"name"`
	Type      CategoryType       `bson:"type" json:"type"`
	Icon      string             `bson:"icon" json:"icon"`
	Color     string             `bson:"color" json:"color"`
	IsSystem  bool               `bson:"isSystem" json:"isSystem"`
	SortOrder int                `bson:"sortOrder" json:"sortOrder"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Category, error)
	ListByUserID(ctx context.Context, userID primitive.ObjectID, categoryType CategoryType) ([]*Category, error)
	ListSystem(ctx context.Context, categoryType CategoryType) ([]*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
