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

type mongoCategoryRepository struct {
	collection *mongo.Collection
}

func NewMongoCategoryRepository(db *mongo.Database) domain.CategoryRepository {
	return &mongoCategoryRepository{
		collection: db.Collection("categories"),
	}
}

func (r *mongoCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, category)
	return err
}

func (r *mongoCategoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error) {
	var category domain.Category
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *mongoCategoryRepository) ListByUserID(ctx context.Context, userID primitive.ObjectID, categoryType domain.CategoryType) ([]*domain.Category, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"userId": userID},
			{"isSystem": true},
		},
	}
	if categoryType != "" {
		filter["type"] = categoryType
	}

	opts := options.Find().SetSort(bson.M{"sortOrder": 1, "name": 1})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*domain.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *mongoCategoryRepository) ListSystem(ctx context.Context, categoryType domain.CategoryType) ([]*domain.Category, error) {
	filter := bson.M{"isSystem": true}
	if categoryType != "" {
		filter["type"] = categoryType
	}

	opts := options.Find().SetSort(bson.M{"sortOrder": 1, "name": 1})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*domain.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *mongoCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	category.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": category.ID}, category)
	return err
}

func (r *mongoCategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id, "isSystem": false})
	return err
}
