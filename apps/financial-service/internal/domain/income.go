package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Frequency string

const (
	FrequencyOneTime Frequency = "one-time"
	FrequencyDaily   Frequency = "daily"
	FrequencyWeekly  Frequency = "weekly"
	FrequencyMonthly Frequency = "monthly"
	FrequencyYearly  Frequency = "yearly"
)

type Income struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"userId" json:"userId"`
	Source        string             `bson:"source" json:"source"`
	Amount        float64            `bson:"amount" json:"amount"` // Using float64 for simplicity, architecture mentions Decimal128
	Frequency     Frequency          `bson:"frequency" json:"frequency"`
	EffectiveDate time.Time          `bson:"effectiveDate" json:"effectiveDate"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type IncomeRepository interface {
	Create(ctx context.Context, income *Income) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Income, error)
	ListByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Income, error)
	Update(ctx context.Context, income *Income) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
