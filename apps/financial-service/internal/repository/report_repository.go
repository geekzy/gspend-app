package repository

import (
	"context"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoReportRepository struct {
	transactionCollection *mongo.Collection
	budgetCollection      *mongo.Collection
	categoryCollection    *mongo.Collection
}

func NewMongoReportRepository(db *mongo.Database) domain.ReportRepository {
	return &mongoReportRepository{
		transactionCollection: db.Collection("transactions"),
		budgetCollection:      db.Collection("budgets"),
		categoryCollection:    db.Collection("categories"),
	}
}

func (r *mongoReportRepository) GetBudgetVsActualData(ctx context.Context, userID primitive.ObjectID, month time.Time) (*domain.BudgetVsActualReport, error) {
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	// Find active budget for the month
	var budget domain.Budget
	err := r.budgetCollection.FindOne(ctx, bson.M{
		"userId": userID,
		"startDate": bson.M{"$lte": startOfMonth},
		"endDate":   bson.M{"$gte": endOfMonth.Add(-time.Second)},
	}).Decode(&budget)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No budget found, return empty report
			return &domain.BudgetVsActualReport{
				Month:           startOfMonth,
				Categories:      []*domain.BudgetVsActualItem{},
				TotalBudgeted:   0,
				TotalSpent:      0,
				OverallVariance: 0,
			}, nil
		}
		return nil, err
	}

	// Get actual spending by category for the month
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

	// Create map of actual spending by category
	actualSpending := make(map[primitive.ObjectID]float64)
	categoryNames := make(map[primitive.ObjectID]string)

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
		actualSpending[result.ID] = result.Total
		categoryNames[result.ID] = result.Category.Name
	}

	// Build budget vs actual items
	var categories []*domain.BudgetVsActualItem
	totalBudgeted := 0.0
	totalSpent := 0.0

	for _, item := range budget.Items {
		actual := actualSpending[item.CategoryID]
		totalBudgeted += item.PlannedAmount
		totalSpent += actual

		categories = append(categories, &domain.BudgetVsActualItem{
			CategoryID:     item.CategoryID,
			CategoryName:   item.CategoryName,
			Budgeted:       item.PlannedAmount,
			Actual:         actual,
			Variance:       actual - item.PlannedAmount,
			PercentageUsed: 0, // Will be calculated in service
		})
	}

	return &domain.BudgetVsActualReport{
		Month:           startOfMonth,
		Categories:      categories,
		TotalBudgeted:   totalBudgeted,
		TotalSpent:      totalSpent,
		OverallVariance: 0, // Will be calculated in service
	}, nil
}

func (r *mongoReportRepository) GetSpendingByCategoryData(ctx context.Context, userID primitive.ObjectID, startDate, endDate time.Time) (*domain.SpendingByCategoryReport, error) {
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
		{"$sort": bson.M{"total": -1}},
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

	var categories []*domain.CategorySpending
	totalSpent := 0.0

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

		totalSpent += result.Total
		categories = append(categories, &domain.CategorySpending{
			CategoryID:   result.ID,
			CategoryName: result.Category.Name,
			Amount:       result.Total,
			Percentage:   0, // Will be calculated in service
		})
	}

	return &domain.SpendingByCategoryReport{
		StartDate:  startDate,
		EndDate:    endDate,
		Categories: categories,
		TotalSpent: totalSpent,
	}, nil
}

func (r *mongoReportRepository) GetMonthlyTrendsData(ctx context.Context, userID primitive.ObjectID, months int) (*domain.MonthlyTrendsReport, error) {
	// Calculate start date (months ago from now)
	now := time.Now()
	startDate := now.AddDate(0, -months, 0)
	startOfMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())

	pipeline := []bson.M{
		{"$match": bson.M{
			"userId": userID,
			"type":   domain.TransactionTypeExpense,
			"transactionDate": bson.M{
				"$gte": startOfMonth,
				"$lte": now,
			},
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

	cursor, err := r.transactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var monthlyData []*domain.MonthlySpending
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

		monthTime := time.Date(result.ID.Year, time.Month(result.ID.Month), 1, 0, 0, 0, 0, time.UTC)
		monthlyData = append(monthlyData, &domain.MonthlySpending{
			Month:  monthTime,
			Amount: result.Total,
		})
	}

	return &domain.MonthlyTrendsReport{
		Months:          months,
		MonthlyData:     monthlyData,
		AverageSpending: 0, // Will be calculated in service
		TrendDirection:  "", // Will be calculated in service
	}, nil
}