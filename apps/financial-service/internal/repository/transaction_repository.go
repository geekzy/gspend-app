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

type mongoTransactionRepository struct {
	collection *mongo.Collection
}

func NewMongoTransactionRepository(db *mongo.Database) domain.EnhancedTransactionRepository {
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
// FindWithFilters implements enhanced filtering with pagination
func (r *mongoTransactionRepository) FindWithFilters(ctx context.Context, userID primitive.ObjectID, filters domain.TransactionFilters) ([]*domain.Transaction, int, error) {
	// Build filter
	mongoFilter := bson.M{"userId": userID}
	
	if filters.StartDate != nil && filters.EndDate != nil {
		mongoFilter["transactionDate"] = bson.M{
			"$gte": *filters.StartDate,
			"$lte": *filters.EndDate,
		}
	} else if filters.StartDate != nil {
		mongoFilter["transactionDate"] = bson.M{"$gte": *filters.StartDate}
	} else if filters.EndDate != nil {
		mongoFilter["transactionDate"] = bson.M{"$lte": *filters.EndDate}
	}
	
	if filters.CategoryID != nil {
		mongoFilter["categoryId"] = *filters.CategoryID
	}
	
	if filters.Type != nil {
		mongoFilter["type"] = *filters.Type
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	// Set defaults for pagination
	page := filters.Page
	if page < 1 {
		page = 1
	}
	perPage := filters.PerPage
	if perPage < 1 {
		perPage = 20
	}

	// Set defaults for sorting
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "transactionDate"
	}
	sortOrder := filters.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Build sort
	sortDirection := -1
	if sortOrder == "asc" {
		sortDirection = 1
	}
	sort := bson.M{sortBy: sortDirection}

	// Build options
	skip := (page - 1) * perPage
	opts := options.Find().
		SetSort(sort).
		SetSkip(int64(skip)).
		SetLimit(int64(perPage))

	// Execute query
	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var transactions []*domain.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, 0, err
	}

	return transactions, int(total), nil
}

// GetSpendingByCategory aggregates spending by category for a date range
func (r *mongoTransactionRepository) GetSpendingByCategory(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) ([]*domain.CategorySpending, error) {
	pipeline := []bson.M{
		{"$match": bson.M{
			"userId": userID,
			"type":   domain.TransactionTypeExpense,
			"transactionDate": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}},
		{"$group": bson.M{
			"_id":   "$categoryId",
			"total": bson.M{"$sum": "$amount"},
		}},
		{"$lookup": bson.M{
			"from":         "categories",
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": "$category"},
		{"$sort": bson.M{"total": -1}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Get total expenses for percentage calculation
	totalPipeline := []bson.M{
		{"$match": bson.M{
			"userId": userID,
			"type":   domain.TransactionTypeExpense,
			"transactionDate": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}},
	}

	totalCursor, err := r.collection.Aggregate(ctx, totalPipeline)
	if err != nil {
		return nil, err
	}
	defer totalCursor.Close(ctx)

	var totalExpenses float64
	if totalCursor.Next(ctx) {
		var totalResult struct {
			Total float64 `bson:"total"`
		}
		if err := totalCursor.Decode(&totalResult); err == nil {
			totalExpenses = totalResult.Total
		}
	}

	var categories []*domain.CategorySpending
	for cursor.Next(ctx) {
		var result struct {
			ID       primitive.ObjectID `bson:"_id"`
			Total    float64            `bson:"total"`
			Category struct {
				Name string `bson:"name"`
			} `bson:"category"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}

		percentage := float64(0)
		if totalExpenses > 0 {
			percentage = (result.Total / totalExpenses) * 100
		}

		categories = append(categories, &domain.CategorySpending{
			CategoryID:   result.ID,
			CategoryName: result.Category.Name,
			Amount:       result.Total,
			Percentage:   percentage,
		})
	}

	return categories, nil
}

// GetMonthlySpendingTrends gets spending trends for the last N months
func (r *mongoTransactionRepository) GetMonthlySpendingTrends(ctx context.Context, userID primitive.ObjectID, months int) ([]*domain.MonthlySpending, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -months+1, 0)

	pipeline := []bson.M{
		{"$match": bson.M{
			"userId": userID,
			"type":   domain.TransactionTypeExpense,
			"transactionDate": bson.M{"$gte": startDate},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"year":  bson.M{"$year": "$transactionDate"},
				"month": bson.M{"$month": "$transactionDate"},
			},
			"total": bson.M{"$sum": "$amount"},
		}},
		{"$sort": bson.M{"_id.year": 1, "_id.month": 1}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trends []*domain.MonthlySpending
	for cursor.Next(ctx) {
		var result struct {
			ID struct {
				Year  int `bson:"year"`
				Month int `bson:"month"`
			} `bson:"_id"`
			Total float64 `bson:"total"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}

		month := time.Date(result.ID.Year, time.Month(result.ID.Month), 1, 0, 0, 0, 0, time.UTC)
		trends = append(trends, &domain.MonthlySpending{
			Month:  month,
			Amount: result.Total,
		})
	}

	return trends, nil
}