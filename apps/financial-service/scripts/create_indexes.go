package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Get MongoDB URI from environment or use default
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoDatabase := os.Getenv("MONGODB_DATABASE")
	if mongoDatabase == "" {
		mongoDatabase = "gspend"
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(mongoDatabase)

	fmt.Println("Creating MongoDB indexes for optimal dashboard performance...")

	// Compound index for transaction filtering
	transactionCollection := db.Collection("transactions")
	_, err = transactionCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "transactionDate", Value: -1},
			{Key: "type", Value: 1},
			{Key: "categoryId", Value: 1},
		},
		Options: options.Index().SetName("transaction_filtering_idx"),
	})
	if err != nil {
		log.Printf("Warning: Failed to create transaction filtering index: %v", err)
	} else {
		fmt.Println("✓ Created transaction filtering index")
	}

	// Index for dashboard aggregations
	_, err = transactionCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "transactionDate", Value: -1},
		},
		Options: options.Index().SetName("dashboard_aggregation_idx"),
	})
	if err != nil {
		log.Printf("Warning: Failed to create dashboard aggregation index: %v", err)
	} else {
		fmt.Println("✓ Created dashboard aggregation index")
	}

	// Index for budget tracking
	budgetCollection := db.Collection("budgets")
	_, err = budgetCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "startDate", Value: -1},
			{Key: "endDate", Value: -1},
		},
		Options: options.Index().SetName("budget_tracking_idx"),
	})
	if err != nil {
		log.Printf("Warning: Failed to create budget tracking index: %v", err)
	} else {
		fmt.Println("✓ Created budget tracking index")
	}

	// Index for category lookups
	categoryCollection := db.Collection("categories")
	_, err = categoryCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "name", Value: 1},
		},
		Options: options.Index().SetName("category_lookup_idx"),
	})
	if err != nil {
		log.Printf("Warning: Failed to create category lookup index: %v", err)
	} else {
		fmt.Println("✓ Created category lookup index")
	}

	fmt.Println("MongoDB indexes created successfully!")
}