package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/geekzy/gspend-app/apps/financial-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FamilyCategoryData represents the seed data for family-oriented categories
type FamilyCategoryData struct {
	Name      string                `json:"name"`
	Type      domain.CategoryType   `json:"type"`
	Icon      string                `json:"icon"`
	Color     string                `json:"color"`
	SortOrder int                   `json:"sortOrder"`
}

// getFamilyCategories returns the predefined family-oriented categories
func getFamilyCategories() []FamilyCategoryData {
	return []FamilyCategoryData{
		// Housing & Utilities (1-10)
		{Name: "Rent/Mortgage", Type: domain.CategoryTypeExpense, Icon: "🏠", Color: "#3B82F6", SortOrder: 1},
		{Name: "Utilities", Type: domain.CategoryTypeExpense, Icon: "💡", Color: "#EF4444", SortOrder: 2},
		{Name: "Home Maintenance", Type: domain.CategoryTypeExpense, Icon: "🔧", Color: "#6B7280", SortOrder: 3},
		
		// Food & Groceries (11-20)
		{Name: "Groceries", Type: domain.CategoryTypeExpense, Icon: "🛒", Color: "#10B981", SortOrder: 11},
		{Name: "Dining Out", Type: domain.CategoryTypeExpense, Icon: "🍽️", Color: "#F59E0B", SortOrder: 12},
		{Name: "Coffee & Snacks", Type: domain.CategoryTypeExpense, Icon: "☕", Color: "#92400E", SortOrder: 13},
		
		// Children & Family (21-40)
		{Name: "Childcare", Type: domain.CategoryTypeExpense, Icon: "👶", Color: "#8B5CF6", SortOrder: 21},
		{Name: "School Expenses", Type: domain.CategoryTypeExpense, Icon: "📚", Color: "#06B6D4", SortOrder: 22},
		{Name: "Kids Activities", Type: domain.CategoryTypeExpense, Icon: "🎨", Color: "#EC4899", SortOrder: 23},
		{Name: "Kids Clothing", Type: domain.CategoryTypeExpense, Icon: "👕", Color: "#84CC16", SortOrder: 24},
		{Name: "Toys & Games", Type: domain.CategoryTypeExpense, Icon: "🧸", Color: "#F472B6", SortOrder: 25},
		{Name: "Baby Supplies", Type: domain.CategoryTypeExpense, Icon: "🍼", Color: "#A78BFA", SortOrder: 26},
		
		// Transportation (41-50)
		{Name: "Car Payment", Type: domain.CategoryTypeExpense, Icon: "🚗", Color: "#6366F1", SortOrder: 41},
		{Name: "Gas", Type: domain.CategoryTypeExpense, Icon: "⛽", Color: "#EF4444", SortOrder: 42},
		{Name: "Car Maintenance", Type: domain.CategoryTypeExpense, Icon: "🔧", Color: "#374151", SortOrder: 43},
		{Name: "Public Transport", Type: domain.CategoryTypeExpense, Icon: "🚌", Color: "#059669", SortOrder: 44},
		
		// Healthcare (51-60)
		{Name: "Medical", Type: domain.CategoryTypeExpense, Icon: "🏥", Color: "#DC2626", SortOrder: 51},
		{Name: "Insurance", Type: domain.CategoryTypeExpense, Icon: "🛡️", Color: "#059669", SortOrder: 52},
		{Name: "Pharmacy", Type: domain.CategoryTypeExpense, Icon: "💊", Color: "#7C2D12", SortOrder: 53},
		{Name: "Dental", Type: domain.CategoryTypeExpense, Icon: "🦷", Color: "#1E40AF", SortOrder: 54},
		
		// Personal & Clothing (61-70)
		{Name: "Clothing", Type: domain.CategoryTypeExpense, Icon: "👔", Color: "#7C3AED", SortOrder: 61},
		{Name: "Personal Care", Type: domain.CategoryTypeExpense, Icon: "🧴", Color: "#BE185D", SortOrder: 62},
		{Name: "Haircuts", Type: domain.CategoryTypeExpense, Icon: "✂️", Color: "#9333EA", SortOrder: 63},
		
		// Entertainment & Recreation (71-80)
		{Name: "Entertainment", Type: domain.CategoryTypeExpense, Icon: "🎬", Color: "#7C2D12", SortOrder: 71},
		{Name: "Subscriptions", Type: domain.CategoryTypeExpense, Icon: "📺", Color: "#1F2937", SortOrder: 72},
		{Name: "Hobbies", Type: domain.CategoryTypeExpense, Icon: "🎯", Color: "#0F766E", SortOrder: 73},
		{Name: "Vacation", Type: domain.CategoryTypeExpense, Icon: "✈️", Color: "#0369A1", SortOrder: 74},
		
		// Miscellaneous (81-90)
		{Name: "Gifts", Type: domain.CategoryTypeExpense, Icon: "🎁", Color: "#BE123C", SortOrder: 81},
		{Name: "Charity", Type: domain.CategoryTypeExpense, Icon: "❤️", Color: "#DC2626", SortOrder: 82},
		{Name: "Pet Care", Type: domain.CategoryTypeExpense, Icon: "🐕", Color: "#92400E", SortOrder: 83},
		{Name: "Other", Type: domain.CategoryTypeExpense, Icon: "📦", Color: "#6B7280", SortOrder: 89},
		
		// Income Categories (91-100)
		{Name: "Salary", Type: domain.CategoryTypeIncome, Icon: "💼", Color: "#10B981", SortOrder: 91},
		{Name: "Freelance", Type: domain.CategoryTypeIncome, Icon: "💻", Color: "#3B82F6", SortOrder: 92},
		{Name: "Side Business", Type: domain.CategoryTypeIncome, Icon: "🏪", Color: "#059669", SortOrder: 93},
		{Name: "Investment", Type: domain.CategoryTypeIncome, Icon: "📈", Color: "#0D9488", SortOrder: 94},
		{Name: "Gift Money", Type: domain.CategoryTypeIncome, Icon: "💝", Color: "#BE123C", SortOrder: 95},
		{Name: "Other Income", Type: domain.CategoryTypeIncome, Icon: "💰", Color: "#047857", SortOrder: 99},
	}
}

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
	collection := db.Collection("categories")

	fmt.Println("Seeding family-oriented categories...")

	// Check if system categories already exist
	count, err := collection.CountDocuments(context.Background(), bson.M{"isSystem": true})
	if err != nil {
		log.Fatalf("Failed to check existing categories: %v", err)
	}

	if count > 0 {
		fmt.Printf("Found %d existing system categories. Skipping seed to avoid duplicates.\n", count)
		fmt.Println("To re-seed, first remove existing system categories.")
		return
	}

	// Prepare categories for insertion
	familyCategories := getFamilyCategories()
	var categories []interface{}
	now := time.Now()

	for _, catData := range familyCategories {
		category := domain.Category{
			ID:        primitive.NewObjectID(),
			UserID:    nil, // System categories have no user
			Name:      catData.Name,
			Type:      catData.Type,
			Icon:      catData.Icon,
			Color:     catData.Color,
			IsSystem:  true,
			SortOrder: catData.SortOrder,
			CreatedAt: now,
			UpdatedAt: now,
		}
		categories = append(categories, category)
	}

	// Insert all categories
	result, err := collection.InsertMany(context.Background(), categories)
	if err != nil {
		log.Fatalf("Failed to insert categories: %v", err)
	}

	fmt.Printf("✓ Successfully seeded %d family-oriented categories!\n", len(result.InsertedIDs))
	
	// Print summary by type
	expenseCount := 0
	incomeCount := 0
	for _, catData := range familyCategories {
		if catData.Type == domain.CategoryTypeExpense {
			expenseCount++
		} else {
			incomeCount++
		}
	}
	
	fmt.Printf("  - %d expense categories\n", expenseCount)
	fmt.Printf("  - %d income categories\n", incomeCount)
	fmt.Println("\nCategories are organized by:")
	fmt.Println("  • Housing & Utilities")
	fmt.Println("  • Food & Groceries") 
	fmt.Println("  • Children & Family")
	fmt.Println("  • Transportation")
	fmt.Println("  • Healthcare")
	fmt.Println("  • Personal & Clothing")
	fmt.Println("  • Entertainment & Recreation")
	fmt.Println("  • Income Sources")
	fmt.Println("\nFamily financial management is now ready to use!")
}