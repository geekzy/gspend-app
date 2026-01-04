package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the user entity in the system
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"passwordHash" json:"-"`
	FullName     string             `bson:"fullName" json:"fullName"`
	FamilySize   int                `bson:"familySize" json:"familySize"`
	// Email verification
	EmailVerified       bool       `bson:"emailVerified" json:"emailVerified"`
	VerificationToken   string     `bson:"verificationToken,omitempty" json:"-"`
	VerificationExpiry  *time.Time `bson:"verificationExpiry,omitempty" json:"-"`
	// Password reset
	ResetToken          string     `bson:"resetToken,omitempty" json:"-"`
	ResetTokenExpiry    *time.Time `bson:"resetTokenExpiry,omitempty" json:"-"`
	// Timestamps
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, email string) (bool, error)
	// Token operations
	GetByVerificationToken(ctx context.Context, token string) (*User, error)
	GetByResetToken(ctx context.Context, token string) (*User, error)
}
