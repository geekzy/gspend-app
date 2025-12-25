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

type mongoIncomeRepository struct {
	collection *mongo.Collection
}

func NewMongoIncomeRepository(db *mongo.Database) domain.IncomeRepository {
	return &mongoIncomeRepository{
		collection: db.Collection("incomes"),
	}
}

func (r *mongoIncomeRepository) Create(ctx context.Context, income *domain.Income) error {
	income.ID = primitive.NewObjectID()
	income.CreatedAt = time.Now()
	income.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, income)
	return err
}

func (r *mongoIncomeRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Income, error) {
	var income domain.Income
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&income)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &income, nil
}

func (r *mongoIncomeRepository) ListByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Income, error) {
	filter := bson.M{"userId": userID}
	opts := options.Find().SetSort(bson.M{"effectiveDate": -1})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var incomes []*domain.Income
	if err := cursor.All(ctx, &incomes); err != nil {
		return nil, err
	}
	return incomes, nil
}

func (r *mongoIncomeRepository) Update(ctx context.Context, income *domain.Income) error {
	income.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": income.ID}, income)
	return err
}

func (r *mongoIncomeRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
