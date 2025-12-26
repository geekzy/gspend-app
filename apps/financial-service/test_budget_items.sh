#!/bin/bash

# Simple integration test for budget item management
# This script tests the budget item CRUD operations

echo "Testing Budget Item Management API..."
echo "===================================="

# Configuration
BASE_URL="http://localhost:8082"
MONGO_URI=${MONGODB_URI:-"mongodb://localhost:27017"}
MONGO_DB=${MONGODB_DATABASE:-"gspend"}

# Test data
BUDGET_ID="507f1f77bcf86cd799439011"
CATEGORY_ID="507f1f77bcf86cd799439012"

# Function to check if MongoDB is running
check_mongodb() {
    echo "Checking MongoDB connection..."
    if command -v mongosh &> /dev/null; then
        mongosh --quiet --eval "db.runCommand('ping')" "$MONGO_URI/$MONGO_DB" > /dev/null 2>&1
        return $?
    elif command -v mongo &> /dev/null; then
        mongo --quiet --eval "db.runCommand('ping')" "$MONGO_URI/$MONGO_DB" > /dev/null 2>&1
        return $?
    else
        echo "❌ Neither mongosh nor mongo client found"
        return 1
    fi
}

# Function to check if categories are seeded
check_seeded_categories() {
    echo "Checking for seeded categories..."
    if command -v mongosh &> /dev/null; then
        CATEGORY_COUNT=$(mongosh --quiet --eval "db.categories.countDocuments({isSystem: true})" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
    elif command -v mongo &> /dev/null; then
        CATEGORY_COUNT=$(mongo --quiet --eval "db.categories.countDocuments({isSystem: true})" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
    else
        echo "❌ Cannot check categories - no MongoDB client found"
        return 1
    fi
    
    if [ "$CATEGORY_COUNT" -gt 0 ]; then
        echo "✓ Found $CATEGORY_COUNT seeded categories"
        return 0
    else
        echo "❌ No seeded categories found"
        echo "Please run: go run scripts/seed_categories.go"
        return 1
    fi
}

# Function to get a real category ID from seeded data
get_real_category_id() {
    echo "Getting real category ID from seeded data..."
    if command -v mongosh &> /dev/null; then
        REAL_CATEGORY_ID=$(mongosh --quiet --eval "db.categories.findOne({isSystem: true, type: 'expense'})._id.toString()" "$MONGO_URI/$MONGO_DB" 2>/dev/null | grep -o '[a-f0-9]\{24\}')
    elif command -v mongo &> /dev/null; then
        REAL_CATEGORY_ID=$(mongo --quiet --eval "db.categories.findOne({isSystem: true, type: 'expense'})._id.toString()" "$MONGO_URI/$MONGO_DB" 2>/dev/null | grep -o '[a-f0-9]\{24\}')
    fi
    
    if [ -n "$REAL_CATEGORY_ID" ]; then
        echo "✓ Using real category ID: $REAL_CATEGORY_ID"
        CATEGORY_ID="$REAL_CATEGORY_ID"
        return 0
    else
        echo "⚠️  Could not get real category ID, using test ID: $CATEGORY_ID"
        return 0
    fi
}

# Function to check if financial service is running
check_service() {
    echo "Checking if financial service is running..."
    if curl -s "$BASE_URL/api/v1/health" > /dev/null 2>&1; then
        echo "✓ Financial service is running"
        return 0
    else
        echo "❌ Financial service is not running"
        echo "Please start the service: ./bin/financial-service"
        return 1
    fi
}

# Pre-flight checks
echo "Running pre-flight checks..."
echo "----------------------------"

if ! check_mongodb; then
    echo "❌ MongoDB is not running or not accessible"
    echo "Please start MongoDB and ensure it's accessible at: $MONGO_URI"
    exit 1
fi

if ! check_seeded_categories; then
    echo "❌ Seeded categories not found"
    echo "Please run the category seeding script first:"
    echo "  cd apps/financial-service"
    echo "  go run scripts/seed_categories.go"
    exit 1
fi

get_real_category_id

if ! check_service; then
    echo "❌ Financial service is not running"
    echo "Please start the service first:"
    echo "  cd apps/financial-service"
    echo "  go run cmd/server/main.go"
    exit 1
fi

echo "✓ All pre-flight checks passed!"
echo ""

# Test 1: Add budget item
echo "Test 1: Adding budget item..."
ADD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/budgets/$BUDGET_ID/items" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d "{
    \"categoryId\": \"$CATEGORY_ID\",
    \"categoryName\": \"Test Category\",
    \"plannedAmount\": 500.0,
    \"notes\": \"Test budget item\"
  }" \
  -w "\nHTTP_STATUS:%{http_code}")

HTTP_STATUS=$(echo "$ADD_RESPONSE" | grep "HTTP_STATUS:" | cut -d: -f2)
RESPONSE_BODY=$(echo "$ADD_RESPONSE" | sed '/HTTP_STATUS:/d')

echo "Response: $RESPONSE_BODY"
echo "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "201" ] || [ "$HTTP_STATUS" = "200" ]; then
    echo "✓ Budget item added successfully"
    # Try to extract item ID from response
    if command -v jq &> /dev/null; then
        ITEM_ID=$(echo "$RESPONSE_BODY" | jq -r '.data.id // .id // empty' 2>/dev/null)
        if [ -n "$ITEM_ID" ]; then
            echo "✓ Extracted item ID: $ITEM_ID"
        fi
    fi
else
    echo "❌ Failed to add budget item"
fi

echo ""

# Test 2: Get budget item (use extracted ID or fallback)
if [ -z "$ITEM_ID" ]; then
    ITEM_ID="507f1f77bcf86cd799439013"
    echo "Using fallback item ID for remaining tests: $ITEM_ID"
fi

echo "Test 2: Getting budget item..."
GET_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/budgets/$BUDGET_ID/items/$ITEM_ID" \
  -H "Authorization: Bearer test-token" \
  -w "\nHTTP_STATUS:%{http_code}")

HTTP_STATUS=$(echo "$GET_RESPONSE" | grep "HTTP_STATUS:" | cut -d: -f2)
RESPONSE_BODY=$(echo "$GET_RESPONSE" | sed '/HTTP_STATUS:/d')

echo "Response: $RESPONSE_BODY"
echo "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "200" ]; then
    echo "✓ Budget item retrieved successfully"
else
    echo "❌ Failed to get budget item"
fi

echo ""

# Test 3: Update budget item
echo "Test 3: Updating budget item..."
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/v1/budgets/$BUDGET_ID/items/$ITEM_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d "{
    \"categoryId\": \"$CATEGORY_ID\",
    \"categoryName\": \"Updated Category\",
    \"plannedAmount\": 600.0,
    \"notes\": \"Updated budget item\"
  }" \
  -w "\nHTTP_STATUS:%{http_code}")

HTTP_STATUS=$(echo "$UPDATE_RESPONSE" | grep "HTTP_STATUS:" | cut -d: -f2)
RESPONSE_BODY=$(echo "$UPDATE_RESPONSE" | sed '/HTTP_STATUS:/d')

echo "Response: $RESPONSE_BODY"
echo "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "200" ]; then
    echo "✓ Budget item updated successfully"
else
    echo "❌ Failed to update budget item"
fi

echo ""

# Test 4: Delete budget item
echo "Test 4: Deleting budget item..."
DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/v1/budgets/$BUDGET_ID/items/$ITEM_ID" \
  -H "Authorization: Bearer test-token" \
  -w "\nHTTP_STATUS:%{http_code}")

HTTP_STATUS=$(echo "$DELETE_RESPONSE" | grep "HTTP_STATUS:" | cut -d: -f2)
RESPONSE_BODY=$(echo "$DELETE_RESPONSE" | sed '/HTTP_STATUS:/d')

echo "Response: $RESPONSE_BODY"
echo "Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "200" ] || [ "$HTTP_STATUS" = "204" ]; then
    echo "✓ Budget item deleted successfully"
else
    echo "❌ Failed to delete budget item"
fi

echo ""
echo "===================================="
echo "Budget Item Management Tests Complete!"
echo ""
echo "Prerequisites verified:"
echo "✓ MongoDB connection"
echo "✓ Seeded categories ($CATEGORY_COUNT found)"
echo "✓ Financial service running"
echo ""
echo "Note: For production testing, ensure:"
echo "1. Valid authentication tokens"
echo "2. Existing budget with ID: $BUDGET_ID"
echo "3. Proper user permissions"