package repository

import (
	"context"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDashboardRepository struct {
	transactionCollection *mongo.Collection
	budgetCollection      *mongo.Collection
	categoryCollection    *mongo.Collection
}

func NewMongoDashboardRepository(db *mongo.Database) domain.DashboardRepository {
	return &mongoDashboardRepository{
		transactionCollection: db.Collection("transactions"),
		budgetCollection:      db.Collection("budgets"),
		categoryCollection:    db.Collection("categories"),
	}
}

func (r *mongoDashboardRepository) GetTotalBalance(ctx context.Context, userID primitive.ObjectID) (float64, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"userId": userID}},
		{"$group": bson.M{
			"_id": "$type",
			"total": bson.M{"$sum": "$amount"},
		}},
	}

	cursor, err := r.transactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var income, expenses float64
	for cursor.Next(ctx) {
		var result struct {
			ID    string  `bson:"_id"`
			Total float64 `bson:"total"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		
		if result.ID == string(domain.TransactionTypeIncome) {
			income = result.Total
		} else if result.ID == string(domain.TransactionTypeExpense) {
			expenses = result.Total
		}
	}

	return income - expenses, nil
}

func (r *mongoDashboardRepository) GetMonthlyIncome(ctx context.Context, userID primitive.ObjectID, month time.Time) (float64, error) {
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	pipeline := []bson.M{
		{"$match": bson.M{
			"userId": userID,
			"type":   domain.TransactionTypeIncome,
			"transactionDate": bson.M{
				"$gte": startOfMonth,
				"$lt":  endOfMonth,
			},
		}},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}},
	}

	cursor, err := r.transactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Total, nil
	}

	return 0, nil
}

func (r *mongoDashboardRepository) GetMonthlyExpenses(ctx context.Context, userID primitive.ObjectID, month time.Time) (float64, error) {
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	pipeline := []bson.M{
		{"$match": bson.M{
			"userId": userID,
			"type":   domain.TransactionTypeExpense,
			"transactionDate": bson.M{
				"$gte": startOfMonth,
				"$lt":  endOfMonth,
			},
		}},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}},
	}

	cursor, err := r.transactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Total, nil
	}

	return 0, nil
}

func (r *mongoDashboardRepository) GetTopSpendingCategories(ctx context.Context, userID primitive.ObjectID, month time.Time, limit int) ([]*domain.CategorySpending, error) {
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	pipeline := []bson.M{
		{"$match": bson.M{
			"userId": userID,
			"type":   domain.TransactionTypeExpense,
			"transactionDate": bson.M{
				"$gte": startOfMonth,
				"$lt":  endOfMonth,
			},
		}},
		{"$group": bson.M{
			"_id":   "$categoryId",
			"total": bson.M{"$sum": "$amount"},
		}},
		{"$sort": bson.M{"total": -1}},
		{"$limit": limit},
		{"$lookup": bson.M{
			"from":         "categories",
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": "$category"},
	}

	cursor, err := r.transactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Get total expenses for percentage calculation
	totalExpenses, err := r.GetMonthlyExpenses(ctx, userID, month)
	if err != nil {
		return nil, err
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

func (r *mongoDashboardRepository) GetRecentTransactions(ctx context.Context, userID primitive.ObjectID, limit int) ([]*domain.Transaction, error) {
	opts := options.Find().
		SetSort(bson.M{"transactionDate": -1}).
		SetLimit(int64(limit))

	cursor, err := r.transactionCollection.Find(ctx, bson.M{"userId": userID}, opts)
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

func (r *mongoDashboardRepository) GetCurrentMonthBudgetProgress(ctx context.Context, userID primitive.ObjectID) (*domain.BudgetProgress, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	// Find active budget for current month
	var budget struct {
		TotalAmount float64 `bson:"totalAmount"`
	}
	
	err := r.budgetCollection.FindOne(ctx, bson.M{
		"userId": userID,
		"startDate": bson.M{"$lte": startOfMonth},
		"endDate":   bson.M{"$gte": endOfMonth.Add(-time.Second)},
	}).Decode(&budget)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &domain.BudgetProgress{}, nil
		}
		return nil, err
	}

	// Get total spent this month
	totalSpent, err := r.GetMonthlyExpenses(ctx, userID, now)
	if err != nil {
		return nil, err
	}

	remaining := budget.TotalAmount - totalSpent
	percentageUsed := float64(0)
	if budget.TotalAmount > 0 {
		percentageUsed = (totalSpent / budget.TotalAmount) * 100
	}

	return &domain.BudgetProgress{
		TotalBudget:     budget.TotalAmount,
		TotalSpent:      totalSpent,
		RemainingBudget: remaining,
		PercentageUsed:  percentageUsed,
	}, nil
}