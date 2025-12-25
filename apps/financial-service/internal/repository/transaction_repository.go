package repository

import (
	"context"
	"errors"
	"time"

	"github.com/imam/gspend-app/apps/financial-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoTransactionRepository struct {
	collection *mongo.Collection
}

func NewMongoTransactionRepository(db *mongo.Database) domain.TransactionRepository {
	return &mongoTransactionRepository{
		collection: db.Collection("transactions"),
	}
}

func (r *mongoTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	transaction.ID = primitive.NewObjectID()
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, transaction)
	return err
}

func (r *mongoTransactionRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *mongoTransactionRepository) ListByUserID(ctx context.Context, userID primitive.ObjectID, filter map[string]interface{}) ([]*domain.Transaction, error) {
	mongoFilter := bson.M{"userId": userID}
	
	for k, v := range filter {
		mongoFilter[k] = v
	}

	opts := options.Find().SetSort(bson.M{"transactionDate": -1})
	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*domain.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *mongoTransactionRepository) Update(ctx context.Context, transaction *domain.Transaction) error {
	transaction.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": transaction.ID}, transaction)
	return err
}

func (r *mongoTransactionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
