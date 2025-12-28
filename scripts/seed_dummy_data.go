package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type Category struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	UserID    *primitive.ObjectID `bson:"userId,omitempty"`
	Name      string              `bson:"name"`
	Type      string              `bson:"type"`
	Icon      string              `bson:"icon"`
	Color     string              `bson:"color"`
	IsSystem  bool                `bson:"isSystem"`
	SortOrder int                 `bson:"sortOrder"`
	CreatedAt time.Time           `bson:"createdAt"`
	UpdatedAt time.Time           `bson:"updatedAt"`
}

type Transaction struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty"`
	UserID          primitive.ObjectID  `bson:"userId"`
	CategoryID      primitive.ObjectID  `bson:"categoryId"`
	BudgetID        *primitive.ObjectID `bson:"budgetId,omitempty"`
	Type            string              `bson:"type"`
	Amount          float64             `bson:"amount"`
	Description     string              `bson:"description"`
	TransactionDate time.Time           `bson:"transactionDate"`
	PaymentMethod   string              `bson:"paymentMethod"`
	Notes           string              `bson:"notes"`
	Metadata        map[string]string   `bson:"metadata"`
	CreatedAt       time.Time           `bson:"createdAt"`
	UpdatedAt       time.Time           `bson:"updatedAt"`
}

type BudgetItem struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	CategoryID    primitive.ObjectID `bson:"categoryId"`
	CategoryName  string             `bson:"categoryName"`
	PlannedAmount float64            `bson:"plannedAmount"`
	SpentAmount   float64            `bson:"spentAmount"`
	Notes         string             `bson:"notes"`
}

type Budget struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"userId"`
	Name        string             `bson:"name"`
	PeriodType  string             `bson:"periodType"`
	StartDate   time.Time          `bson:"startDate"`
	EndDate     time.Time          `bson:"endDate"`
	TotalAmount float64            `bson:"totalAmount"`
	Items       []BudgetItem       `bson:"items"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

type Income struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        primitive.ObjectID `bson:"userId"`
	Source        string             `bson:"source"`
	Amount        float64            `bson:"amount"`
	Frequency     string             `bson:"frequency"`
	EffectiveDate time.Time          `bson:"effectiveDate"`
	CreatedAt     time.Time          `bson:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt"`
}

func main() {
	// Get MongoDB URI from environment or use default
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "gspend"
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(dbName)

	// Clear existing data
	fmt.Println("Clearing existing data...")
	collections := []string{"users", "categories", "transactions", "budgets", "incomes"}
	for _, collection := range collections {
		db.Collection(collection).Drop(context.Background())
	}

	// Create demo user
	fmt.Println("Creating demo user...")
	userID := primitive.NewObjectID()
	user := User{
		ID:        userID,
		Email:     "demo@gspend.com",
		Username:  "demo",
		Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password: "password"
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = db.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal("Failed to create user:", err)
	}

	// Create system categories
	fmt.Println("Creating categories...")
	expenseCategories := []Category{
		{ID: primitive.NewObjectID(), Name: "Food & Dining", Type: "expense", Icon: "🍽️", Color: "#FF6B6B", IsSystem: true, SortOrder: 1},
		{ID: primitive.NewObjectID(), Name: "Transportation", Type: "expense", Icon: "🚗", Color: "#4ECDC4", IsSystem: true, SortOrder: 2},
		{ID: primitive.NewObjectID(), Name: "Shopping", Type: "expense", Icon: "🛍️", Color: "#45B7D1", IsSystem: true, SortOrder: 3},
		{ID: primitive.NewObjectID(), Name: "Entertainment", Type: "expense", Icon: "🎬", Color: "#96CEB4", IsSystem: true, SortOrder: 4},
		{ID: primitive.NewObjectID(), Name: "Bills & Utilities", Type: "expense", Icon: "💡", Color: "#FFEAA7", IsSystem: true, SortOrder: 5},
		{ID: primitive.NewObjectID(), Name: "Healthcare", Type: "expense", Icon: "🏥", Color: "#DDA0DD", IsSystem: true, SortOrder: 6},
		{ID: primitive.NewObjectID(), Name: "Education", Type: "expense", Icon: "📚", Color: "#98D8C8", IsSystem: true, SortOrder: 7},
		{ID: primitive.NewObjectID(), Name: "Travel", Type: "expense", Icon: "✈️", Color: "#F7DC6F", IsSystem: true, SortOrder: 8},
	}

	incomeCategories := []Category{
		{ID: primitive.NewObjectID(), Name: "Salary", Type: "income", Icon: "💰", Color: "#2ECC71", IsSystem: true, SortOrder: 1},
		{ID: primitive.NewObjectID(), Name: "Freelance", Type: "income", Icon: "💻", Color: "#3498DB", IsSystem: true, SortOrder: 2},
		{ID: primitive.NewObjectID(), Name: "Investment", Type: "income", Icon: "📈", Color: "#9B59B6", IsSystem: true, SortOrder: 3},
		{ID: primitive.NewObjectID(), Name: "Other Income", Type: "income", Icon: "💵", Color: "#1ABC9C", IsSystem: true, SortOrder: 4},
	}

	// Add timestamps to categories
	now := time.Now()
	allCategories := append(expenseCategories, incomeCategories...)
	for i := range allCategories {
		allCategories[i].CreatedAt = now
		allCategories[i].UpdatedAt = now
	}

	// Insert categories
	var categoryDocs []interface{}
	for _, cat := range allCategories {
		categoryDocs = append(categoryDocs, cat)
	}
	_, err = db.Collection("categories").InsertMany(context.Background(), categoryDocs)
	if err != nil {
		log.Fatal("Failed to create categories:", err)
	}

	// Create income records
	fmt.Println("Creating income records...")
	incomes := []Income{
		{
			ID:            primitive.NewObjectID(),
			UserID:        userID,
			Source:        "Software Engineer Salary",
			Amount:        5000.00,
			Frequency:     "monthly",
			EffectiveDate: time.Now().AddDate(0, -6, 0),
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            primitive.NewObjectID(),
			UserID:        userID,
			Source:        "Freelance Project",
			Amount:        1200.00,
			Frequency:     "one-time",
			EffectiveDate: time.Now().AddDate(0, -2, 0),
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            primitive.NewObjectID(),
			UserID:        userID,
			Source:        "Investment Returns",
			Amount:        300.00,
			Frequency:     "monthly",
			EffectiveDate: time.Now().AddDate(0, -3, 0),
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	var incomeDocs []interface{}
	for _, income := range incomes {
		incomeDocs = append(incomeDocs, income)
	}
	_, err = db.Collection("incomes").InsertMany(context.Background(), incomeDocs)
	if err != nil {
		log.Fatal("Failed to create incomes:", err)
	}

	// Create budget
	fmt.Println("Creating budget...")
	startDate := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	budgetItems := []BudgetItem{
		{ID: primitive.NewObjectID(), CategoryID: expenseCategories[0].ID, CategoryName: "Food & Dining", PlannedAmount: 800.00, SpentAmount: 650.00},
		{ID: primitive.NewObjectID(), CategoryID: expenseCategories[1].ID, CategoryName: "Transportation", PlannedAmount: 300.00, SpentAmount: 280.00},
		{ID: primitive.NewObjectID(), CategoryID: expenseCategories[2].ID, CategoryName: "Shopping", PlannedAmount: 500.00, SpentAmount: 420.00},
		{ID: primitive.NewObjectID(), CategoryID: expenseCategories[3].ID, CategoryName: "Entertainment", PlannedAmount: 200.00, SpentAmount: 150.00},
		{ID: primitive.NewObjectID(), CategoryID: expenseCategories[4].ID, CategoryName: "Bills & Utilities", PlannedAmount: 400.00, SpentAmount: 380.00},
	}

	budget := Budget{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Name:        fmt.Sprintf("Monthly Budget - %s", startDate.Format("January 2006")),
		PeriodType:  "monthly",
		StartDate:   startDate,
		EndDate:     endDate,
		TotalAmount: 2200.00,
		Items:       budgetItems,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err = db.Collection("budgets").InsertOne(context.Background(), budget)
	if err != nil {
		log.Fatal("Failed to create budget:", err)
	}

	// Create transactions
	fmt.Println("Creating transactions...")
	rand.Seed(time.Now().UnixNano())

	var transactions []Transaction
	paymentMethods := []string{"Credit Card", "Debit Card", "Cash", "Bank Transfer", "Digital Wallet"}

	// Generate transactions for the last 3 months
	for month := 0; month < 3; month++ {
		monthStart := time.Now().AddDate(0, -month, 0)
		monthStart = time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, time.UTC)

		// Income transactions
		for _, incomeCategory := range incomeCategories {
			if incomeCategory.Name == "Salary" {
				transactions = append(transactions, Transaction{
					ID:              primitive.NewObjectID(),
					UserID:          userID,
					CategoryID:      incomeCategory.ID,
					Type:            "income",
					Amount:          5000.00,
					Description:     "Monthly Salary",
					TransactionDate: monthStart.AddDate(0, 0, 1),
					PaymentMethod:   "Bank Transfer",
					Notes:           "Regular monthly salary",
					Metadata:        map[string]string{"categoryName": incomeCategory.Name},
					CreatedAt:       now,
					UpdatedAt:       now,
				})
			}
		}

		// Expense transactions
		for day := 1; day <= 28; day++ {
			if rand.Float32() < 0.7 { // 70% chance of transaction per day
				categoryIndex := rand.Intn(len(expenseCategories))
				category := expenseCategories[categoryIndex]

				var amount float64
				var description string

				switch category.Name {
				case "Food & Dining":
					amount = 15.00 + rand.Float64()*50.00
					descriptions := []string{"Lunch at cafe", "Grocery shopping", "Dinner out", "Coffee", "Fast food"}
					description = descriptions[rand.Intn(len(descriptions))]
				case "Transportation":
					amount = 5.00 + rand.Float64()*30.00
					descriptions := []string{"Gas station", "Uber ride", "Bus fare", "Parking fee", "Car maintenance"}
					description = descriptions[rand.Intn(len(descriptions))]
				case "Shopping":
					amount = 20.00 + rand.Float64()*100.00
					descriptions := []string{"Clothing store", "Electronics", "Home goods", "Online shopping", "Pharmacy"}
					description = descriptions[rand.Intn(len(descriptions))]
				case "Entertainment":
					amount = 10.00 + rand.Float64()*60.00
					descriptions := []string{"Movie tickets", "Concert", "Streaming service", "Games", "Books"}
					description = descriptions[rand.Intn(len(descriptions))]
				case "Bills & Utilities":
					amount = 50.00 + rand.Float64()*150.00
					descriptions := []string{"Electricity bill", "Internet bill", "Phone bill", "Water bill", "Insurance"}
					description = descriptions[rand.Intn(len(descriptions))]
				default:
					amount = 10.00 + rand.Float64()*80.00
					description = fmt.Sprintf("%s expense", category.Name)
				}

				transactionDate := monthStart.AddDate(0, 0, day-1)
				transactions = append(transactions, Transaction{
					ID:              primitive.NewObjectID(),
					UserID:          userID,
					CategoryID:      category.ID,
					Type:            "expense",
					Amount:          amount,
					Description:     description,
					TransactionDate: transactionDate,
					PaymentMethod:   paymentMethods[rand.Intn(len(paymentMethods))],
					Notes:           "",
					Metadata:        map[string]string{"categoryName": category.Name},
					CreatedAt:       now,
					UpdatedAt:       now,
				})
			}
		}
	}

	// Insert transactions in batches
	batchSize := 100
	for i := 0; i < len(transactions); i += batchSize {
		end := i + batchSize
		if end > len(transactions) {
			end = len(transactions)
		}

		var batch []interface{}
		for j := i; j < end; j++ {
			batch = append(batch, transactions[j])
		}

		_, err = db.Collection("transactions").InsertMany(context.Background(), batch)
		if err != nil {
			log.Fatal("Failed to create transactions batch:", err)
		}
	}

	fmt.Printf("Successfully created dummy data:\n")
	fmt.Printf("- 1 demo user (email: demo@gspend.com, password: password)\n")
	fmt.Printf("- %d categories\n", len(allCategories))
	fmt.Printf("- %d income records\n", len(incomes))
	fmt.Printf("- 1 budget with %d items\n", len(budgetItems))
	fmt.Printf("- %d transactions\n", len(transactions))
	fmt.Println("\nYou can now log in with:")
	fmt.Println("Email: demo@gspend.com")
	fmt.Println("Password: password")
}