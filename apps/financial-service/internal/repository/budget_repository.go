package repository

import (
	"context"
	"errors"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoBudgetRepository struct {
	collection *mongo.Collection
}

func NewMongoBudgetRepository(db *mongo.Database) domain.BudgetRepository {
	return &mongoBudgetRepository{
		collection: db.Collection("budgets"),
	}
}

func (r *mongoBudgetRepository) Create(ctx context.Context, budget *domain.Budget) error {
	budget.ID = primitive.NewObjectID()
	budget.CreatedAt = time.Now()
	budget.UpdatedAt = time.Now()
	// Ensure items have IDs
	for i := range budget.Items {
		if budget.Items[i].ID.IsZero() {
			budget.Items[i].ID = primitive.NewObjectID()
		}
	}

	_, err := r.collection.InsertOne(ctx, budget)
	return err
}

func (r *mongoBudgetRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Budget, error) {
	var budget domain.Budget
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&budget)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &budget, nil
}

func (r *mongoBudgetRepository) ListByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Budget, error) {
	filter := bson.M{"userId": userID}
	opts := options.Find().SetSort(bson.M{"startDate": -1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var budgets []*domain.Budget
	if err := cursor.All(ctx, &budgets); err != nil {
		return nil, err
	}
	return budgets, nil
}

func (r *mongoBudgetRepository) GetActiveByDate(ctx context.Context, userID primitive.ObjectID, date time.Time) (*domain.Budget, error) {
	filter := bson.M{
		"userId":    userID,
		"startDate": bson.M{"$lte": date},
		"endDate":   bson.M{"$gte": date},
	}

	var budget domain.Budget
	err := r.collection.FindOne(ctx, filter).Decode(&budget)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &budget, nil
}

func (r *mongoBudgetRepository) Update(ctx context.Context, budget *domain.Budget) error {
	budget.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": budget.ID}, budget)
	return err
}

func (r *mongoBudgetRepository) UpdateSpentAmount(ctx context.Context, budgetID primitive.ObjectID, categoryID primitive.ObjectID, amount float64) error {
	filter := bson.M{
		"_id":               budgetID,
		"items.categoryId": categoryID,
	}
	update := bson.M{
		"$inc": bson.M{"items.$.spentAmount": amount},
		"$set": bson.M{"updatedAt": time.Now()},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoBudgetRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
