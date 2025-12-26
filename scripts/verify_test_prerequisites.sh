#!/bin/bash

# Comprehensive test prerequisites verification script
# This script checks all prerequisites needed for functional testing

echo "gSpend Test Prerequisites Verification"
echo "======================================"

# Configuration
MONGO_URI=${MONGODB_URI:-"mongodb://localhost:27017"}
MONGO_DB=${MONGODB_DATABASE:-"gspend"}
AUTH_SERVICE_URL="http://localhost:8081"
FINANCIAL_SERVICE_URL="http://localhost:8082"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "success") echo -e "${GREEN}✓${NC} $message" ;;
        "error") echo -e "${RED}❌${NC} $message" ;;
        "warning") echo -e "${YELLOW}⚠️${NC} $message" ;;
    esac
}

# Function to check if MongoDB is running
check_mongodb() {
    echo "1. Checking MongoDB connection..."
    if command -v mongosh &> /dev/null; then
        if mongosh --quiet --eval "db.runCommand('ping')" "$MONGO_URI/$MONGO_DB" > /dev/null 2>&1; then
            print_status "success" "MongoDB is running and accessible"
            return 0
        else
            print_status "error" "MongoDB is not accessible at $MONGO_URI"
            return 1
        fi
    elif command -v mongo &> /dev/null; then
        if mongo --quiet --eval "db.runCommand('ping')" "$MONGO_URI/$MONGO_DB" > /dev/null 2>&1; then
            print_status "success" "MongoDB is running and accessible"
            return 0
        else
            print_status "error" "MongoDB is not accessible at $MONGO_URI"
            return 1
        fi
    else
        print_status "error" "Neither mongosh nor mongo client found"
        print_status "warning" "Please install MongoDB client tools"
        return 1
    fi
}

# Function to check seeded categories
check_seeded_categories() {
    echo "2. Checking seeded categories..."
    local category_count=0
    
    if command -v mongosh &> /dev/null; then
        category_count=$(mongosh --quiet --eval "db.categories.countDocuments({isSystem: true})" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
    elif command -v mongo &> /dev/null; then
        category_count=$(mongo --quiet --eval "db.categories.countDocuments({isSystem: true})" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
    fi
    
    if [ "$category_count" -gt 0 ]; then
        print_status "success" "Found $category_count seeded system categories"
        
        # Check for both income and expense categories
        if command -v mongosh &> /dev/null; then
            expense_count=$(mongosh --quiet --eval "db.categories.countDocuments({isSystem: true, type: 'expense'})" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
            income_count=$(mongosh --quiet --eval "db.categories.countDocuments({isSystem: true, type: 'income'})" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
        elif command -v mongo &> /dev/null; then
            expense_count=$(mongo --quiet --eval "db.categories.countDocuments({isSystem: true, type: 'expense'})" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
            income_count=$(mongo --quiet --eval "db.categories.countDocuments({isSystem: true, type: 'income'})" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
        fi
        
        print_status "success" "  - $expense_count expense categories"
        print_status "success" "  - $income_count income categories"
        return 0
    else
        print_status "error" "No seeded categories found"
        echo "  To seed categories, run:"
        echo "    cd apps/financial-service"
        echo "    go run scripts/seed_categories.go"
        return 1
    fi
}

# Function to check database indexes
check_database_indexes() {
    echo "3. Checking database indexes..."
    
    if command -v mongosh &> /dev/null; then
        # Check if transaction indexes exist
        transaction_indexes=$(mongosh --quiet --eval "db.transactions.getIndexes().length" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
        category_indexes=$(mongosh --quiet --eval "db.categories.getIndexes().length" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
        user_indexes=$(mongosh --quiet --eval "db.users.getIndexes().length" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
    elif command -v mongo &> /dev/null; then
        transaction_indexes=$(mongo --quiet --eval "db.transactions.getIndexes().length" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
        category_indexes=$(mongo --quiet --eval "db.categories.getIndexes().length" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
        user_indexes=$(mongo --quiet --eval "db.users.getIndexes().length" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
    fi
    
    if [ "$transaction_indexes" -gt 1 ] && [ "$category_indexes" -gt 1 ] && [ "$user_indexes" -gt 1 ]; then
        print_status "success" "Database indexes are present"
        print_status "success" "  - Transactions: $transaction_indexes indexes"
        print_status "success" "  - Categories: $category_indexes indexes"
        print_status "success" "  - Users: $user_indexes indexes"
        return 0
    else
        print_status "warning" "Some database indexes may be missing"
        echo "  To create indexes, run:"
        echo "    cd apps/financial-service"
        echo "    go run scripts/create_indexes.go"
        return 1
    fi
}

# Function to check auth service
check_auth_service() {
    echo "4. Checking Auth Service..."
    if curl -s "$AUTH_SERVICE_URL/api/v1/auth/health" > /dev/null 2>&1; then
        print_status "success" "Auth service is running at $AUTH_SERVICE_URL"
        return 0
    else
        print_status "error" "Auth service is not running"
        echo "  To start auth service:"
        echo "    cd apps/auth-service"
        echo "    go run cmd/server/main.go"
        return 1
    fi
}

# Function to check financial service
check_financial_service() {
    echo "5. Checking Financial Service..."
    if curl -s "$FINANCIAL_SERVICE_URL/api/v1/health" > /dev/null 2>&1; then
        print_status "success" "Financial service is running at $FINANCIAL_SERVICE_URL"
        return 0
    else
        print_status "error" "Financial service is not running"
        echo "  To start financial service:"
        echo "    cd apps/financial-service"
        echo "    go run cmd/server/main.go"
        return 1
    fi
}

# Function to check required tools
check_tools() {
    echo "6. Checking required tools..."
    local tools_ok=true
    
    if command -v curl &> /dev/null; then
        print_status "success" "curl is available"
    else
        print_status "error" "curl is not installed"
        tools_ok=false
    fi
    
    if command -v jq &> /dev/null; then
        print_status "success" "jq is available (for JSON parsing)"
    else
        print_status "warning" "jq is not installed (recommended for better test output)"
        echo "  Install with: brew install jq (macOS) or apt-get install jq (Ubuntu)"
    fi
    
    if command -v go &> /dev/null; then
        print_status "success" "Go is available"
    else
        print_status "error" "Go is not installed"
        tools_ok=false
    fi
    
    if [ "$tools_ok" = true ]; then
        return 0
    else
        return 1
    fi
}

# Function to run seed scripts if needed
run_seeding() {
    echo ""
    echo "Running automatic seeding..."
    echo "----------------------------"
    
    # Seed categories if not present
    if command -v mongosh &> /dev/null; then
        category_count=$(mongosh --quiet --eval "db.categories.countDocuments({isSystem: true})" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
    elif command -v mongo &> /dev/null; then
        category_count=$(mongo --quiet --eval "db.categories.countDocuments({isSystem: true})" "$MONGO_URI/$MONGO_DB" 2>/dev/null)
    fi
    
    if [ "$category_count" -eq 0 ]; then
        echo "Seeding categories..."
        if [ -f "apps/financial-service/scripts/seed_categories.go" ]; then
            cd apps/financial-service
            if go run scripts/seed_categories.go; then
                print_status "success" "Categories seeded successfully"
            else
                print_status "error" "Failed to seed categories"
            fi
            cd - > /dev/null
        else
            print_status "error" "Category seeding script not found"
        fi
    fi
    
    # Create indexes if needed
    echo "Creating database indexes..."
    if [ -f "apps/financial-service/scripts/create_indexes.go" ]; then
        cd apps/financial-service
        if go run scripts/create_indexes.go; then
            print_status "success" "Database indexes created successfully"
        else
            print_status "warning" "Index creation may have failed (might already exist)"
        fi
        cd - > /dev/null
    else
        print_status "warning" "Index creation script not found"
    fi
}

# Main execution
main() {
    local all_checks_passed=true
    
    # Run all checks
    check_tools || all_checks_passed=false
    check_mongodb || all_checks_passed=false
    check_seeded_categories || all_checks_passed=false
    check_database_indexes
    check_auth_service || all_checks_passed=false
    check_financial_service || all_checks_passed=false
    
    echo ""
    echo "======================================"
    
    if [ "$all_checks_passed" = true ]; then
        print_status "success" "All prerequisites are met! Ready for functional testing."
        echo ""
        echo "Available test scripts:"
        echo "  • apps/auth-service/test_profile_endpoints.sh"
        echo "  • apps/financial-service/test_budget_items.sh"
        return 0
    else
        print_status "error" "Some prerequisites are not met."
        echo ""
        echo "Would you like to run automatic seeding? (y/n)"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            run_seeding
            echo ""
            print_status "warning" "Please restart services and run this script again to verify."
        fi
        return 1
    fi
}

# Run main function
main "$@"