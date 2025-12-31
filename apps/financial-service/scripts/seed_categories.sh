#!/bin/bash
set -e

# Default values
MONGODB_URI="${MONGODB_URI:-mongodb://localhost:27017}"
MONGODB_DATABASE="${MONGODB_DATABASE:-gspend}"

echo "Starting category seeding for database: $MONGODB_DATABASE"

# Run the seeding logic using mongosh
mongosh "$MONGODB_URI/$MONGODB_DATABASE" --quiet --eval '
try {
    const targetDb = db.getSiblingDB("'$MONGODB_DATABASE'");
    const collection = targetDb.getCollection("categories");

    // Check if system categories already exist
    const count = collection.countDocuments({ isSystem: true });

    if (count > 0) {
        print(`Found ${count} existing system categories. Skipping seed to avoid duplicates.`);
        print("To re-seed, first remove existing system categories.");
    } else {
        print("Seeding family-oriented categories...");
        
        const now = new Date();
        const categories = [
            // Housing & Utilities (1-10)
            {name: "Rent/Mortgage", type: "expense", icon: "🏠", color: "#3B82F6", sortOrder: 1},
            {name: "Utilities", type: "expense", icon: "💡", color: "#EF4444", sortOrder: 2},
            {name: "Home Maintenance", type: "expense", icon: "🔧", color: "#6B7280", sortOrder: 3},

            // Food & Groceries (11-20)
            {name: "Groceries", type: "expense", icon: "🛒", color: "#10B981", sortOrder: 11},
            {name: "Dining Out", type: "expense", icon: "🍽️", color: "#F59E0B", sortOrder: 12},
            {name: "Coffee & Snacks", type: "expense", icon: "☕", color: "#92400E", sortOrder: 13},

            // Children & Family (21-40)
            {name: "Childcare", type: "expense", icon: "👶", color: "#8B5CF6", sortOrder: 21},
            {name: "School Expenses", type: "expense", icon: "📚", color: "#06B6D4", sortOrder: 22},
            {name: "Kids Activities", type: "expense", icon: "🎨", color: "#EC4899", sortOrder: 23},
            {name: "Kids Clothing", type: "expense", icon: "👕", color: "#84CC16", sortOrder: 24},
            {name: "Toys & Games", type: "expense", icon: "🧸", color: "#F472B6", sortOrder: 25},
            {name: "Baby Supplies", type: "expense", icon: "🍼", color: "#A78BFA", sortOrder: 26},

            // Transportation (41-50)
            {name: "Car Payment", type: "expense", icon: "🚗", color: "#6366F1", sortOrder: 41},
            {name: "Gas", type: "expense", icon: "⛽", color: "#EF4444", sortOrder: 42},
            {name: "Car Maintenance", type: "expense", icon: "🔧", color: "#374151", sortOrder: 43},
            {name: "Public Transport", type: "expense", icon: "🚌", color: "#059669", sortOrder: 44},

            // Healthcare (51-60)
            {name: "Medical", type: "expense", icon: "🏥", color: "#DC2626", sortOrder: 51},
            {name: "Insurance", type: "expense", icon: "🛡️", color: "#059669", sortOrder: 52},
            {name: "Pharmacy", type: "expense", icon: "💊", color: "#7C2D12", sortOrder: 53},
            {name: "Dental", type: "expense", icon: "🦷", color: "#1E40AF", sortOrder: 54},

            // Personal & Clothing (61-70)
            {name: "Clothing", type: "expense", icon: "👔", color: "#7C3AED", sortOrder: 61},
            {name: "Personal Care", type: "expense", icon: "🧴", color: "#BE185D", sortOrder: 62},
            {name: "Haircuts", type: "expense", icon: "✂️", color: "#9333EA", sortOrder: 63},

            // Entertainment & Recreation (71-80)
            {name: "Entertainment", type: "expense", icon: "🎬", color: "#7C2D12", sortOrder: 71},
            {name: "Subscriptions", type: "expense", icon: "📺", color: "#1F2937", sortOrder: 72},
            {name: "Hobbies", type: "expense", icon: "🎯", color: "#0F766E", sortOrder: 73},
            {name: "Vacation", type: "expense", icon: "✈️", color: "#0369A1", sortOrder: 74},

            // Miscellaneous (81-90)
            {name: "Gifts", type: "expense", icon: "🎁", color: "#BE123C", sortOrder: 81},
            {name: "Charity", type: "expense", icon: "❤️", color: "#DC2626", sortOrder: 82},
            {name: "Pet Care", type: "expense", icon: "🐕", color: "#92400E", sortOrder: 83},
            {name: "Other", type: "expense", icon: "📦", color: "#6B7280", sortOrder: 89},

            // Income Categories (91-100)
            {name: "Salary", type: "income", icon: "💼", color: "#10B981", sortOrder: 91},
            {name: "Freelance", type: "income", icon: "💻", color: "#3B82F6", sortOrder: 92},
            {name: "Side Business", type: "income", icon: "🏪", color: "#059669", sortOrder: 93},
            {name: "Investment", type: "income", icon: "📈", color: "#0D9488", sortOrder: 94},
            {name: "Gift Money", type: "income", icon: "💝", color: "#BE123C", sortOrder: 95},
            {name: "Other Income", type: "income", icon: "💰", color: "#047857", sortOrder: 99}
        ];

        const documents = categories.map(cat => ({
            userId: null,
            name: cat.name,
            type: cat.type,
            icon: cat.icon,
            color: cat.color,
            isSystem: true,
            sortOrder: cat.sortOrder,
            createdAt: now,
            updatedAt: now
        }));

        const result = collection.insertMany(documents);
        print(`✓ Successfully seeded ${result.insertedIds.length} family-oriented categories!`);
        
        // Print Summary logic similar to Go script
        let expenseCount = 0;
        let incomeCount = 0;
        categories.forEach(c => {
            if(c.type === "expense") expenseCount++;
            else incomeCount++;
        });
        
        print(`  - ${expenseCount} expense categories`);
        print(`  - ${incomeCount} income categories`);
        print("\nFamily financial management is now ready to use!");
    }
} catch (err) {
    print("Error seeding categories: " + err);
    quit(1);
}
'
