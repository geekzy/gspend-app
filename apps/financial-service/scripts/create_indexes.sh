#!/bin/bash
set -e

# Default values
MONGODB_URI="${MONGODB_URI:-mongodb://localhost:27017}"
MONGODB_DATABASE="${MONGODB_DATABASE:-gspend}"

echo "Starting index creation for database: $MONGODB_DATABASE"
echo "URI: $MONGODB_URI"

# Function to run mongo command
run_mongo_cmd() {
    local collection=$1
    local index_spec=$2
    local options=$3
    
    echo "Creating index on $collection: $index_spec"
    mongosh "$MONGODB_URI/$MONGODB_DATABASE" --eval "printjson(db.getCollection('$collection').createIndex($index_spec, $options))"
}

# Transaction Collection Indexes
echo "--- Transactions ---"
# Compound index for transaction filtering
run_mongo_cmd "transactions" \
    '{ "userId": 1, "transactionDate": -1, "type": 1, "categoryId": 1 }' \
    '{ "name": "transaction_filtering_idx" }'

# Index for dashboard aggregations
run_mongo_cmd "transactions" \
    '{ "userId": 1, "transactionDate": -1 }' \
    '{ "name": "dashboard_aggregation_idx" }'

# Budget Collection Indexes
echo "--- Budgets ---"
# Index for budget tracking
run_mongo_cmd "budgets" \
    '{ "userId": 1, "startDate": -1, "endDate": -1 }' \
    '{ "name": "budget_tracking_idx" }'

# Category Collection Indexes
echo "--- Categories ---"
# Index for category lookups
run_mongo_cmd "categories" \
    '{ "userId": 1, "name": 1 }' \
    '{ "name": "category_lookup_idx" }'

echo "Index creation complete!"
