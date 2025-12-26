#!/bin/bash

# Simple integration test for budget item management
# This script tests the budget item CRUD operations

echo "Testing Budget Item Management API..."

# Test data
BUDGET_ID="507f1f77bcf86cd799439011"
CATEGORY_ID="507f1f77bcf86cd799439012"

# Start the server in background (assuming MongoDB is running)
echo "Starting financial service..."
./bin/financial-service &
SERVER_PID=$!

# Wait for server to start
sleep 3

# Test 1: Add budget item
echo "Test 1: Adding budget item..."
curl -X POST "http://localhost:8080/api/v1/budgets/$BUDGET_ID/items" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "categoryId": "'$CATEGORY_ID'",
    "categoryName": "Test Category",
    "plannedAmount": 500.0,
    "notes": "Test budget item"
  }' \
  -w "\nStatus: %{http_code}\n"

echo ""

# Test 2: Get budget item (would need actual item ID from response)
echo "Test 2: Getting budget item..."
ITEM_ID="507f1f77bcf86cd799439013"
curl -X GET "http://localhost:8080/api/v1/budgets/$BUDGET_ID/items/$ITEM_ID" \
  -H "Authorization: Bearer test-token" \
  -w "\nStatus: %{http_code}\n"

echo ""

# Test 3: Update budget item
echo "Test 3: Updating budget item..."
curl -X PUT "http://localhost:8080/api/v1/budgets/$BUDGET_ID/items/$ITEM_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "categoryId": "'$CATEGORY_ID'",
    "categoryName": "Updated Category",
    "plannedAmount": 600.0,
    "notes": "Updated budget item"
  }' \
  -w "\nStatus: %{http_code}\n"

echo ""

# Test 4: Delete budget item
echo "Test 4: Deleting budget item..."
curl -X DELETE "http://localhost:8080/api/v1/budgets/$BUDGET_ID/items/$ITEM_ID" \
  -H "Authorization: Bearer test-token" \
  -w "\nStatus: %{http_code}\n"

echo ""

# Clean up
echo "Stopping server..."
kill $SERVER_PID

echo "Budget item management tests completed!"
echo ""
echo "Note: These tests require:"
echo "1. MongoDB running on localhost:27017"
echo "2. Valid authentication (mock auth service or disable auth middleware)"
echo "3. Existing budget with ID: $BUDGET_ID"