#!/bin/bash

# Comprehensive integration test script for CI/CD
# This script runs all backend integration tests in a containerized environment

set -e  # Exit on any error

echo "🚀 gSpend Backend Integration Tests"
echo "=================================="

# Configuration
AUTH_SERVICE_URL=${AUTH_SERVICE_URL:-"http://auth-service-test:8081"}
FINANCIAL_SERVICE_URL=${FINANCIAL_SERVICE_URL:-"http://financial-service-test:8082"}
MONGODB_URI=${MONGODB_URI:-"mongodb://mongodb-test:27017"}
MONGODB_DATABASE=${MONGODB_DATABASE:-"gspend_test"}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results tracking
TESTS_PASSED=0
TESTS_FAILED=0
TEST_RESULTS_FILE="/app/test-results/integration-test-results.json"

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "success") echo -e "${GREEN}✅${NC} $message" ;;
        "error") echo -e "${RED}❌${NC} $message" ;;
        "warning") echo -e "${YELLOW}⚠️${NC} $message" ;;
        "info") echo -e "${BLUE}ℹ️${NC} $message" ;;
    esac
}

# Function to wait for service to be ready
wait_for_service() {
    local service_url=$1
    local service_name=$2
    local max_attempts=30
    local attempt=1
    
    print_status "info" "Waiting for $service_name to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$service_url" > /dev/null 2>&1; then
            print_status "success" "$service_name is ready"
            return 0
        fi
        
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    print_status "error" "$service_name failed to start within $((max_attempts * 2)) seconds"
    return 1
}

# Function to seed test data
seed_test_data() {
    print_status "info" "Seeding test data..."
    
    # Seed categories
    if /app/bin/seed-categories; then
        print_status "success" "Categories seeded successfully"
    else
        print_status "error" "Failed to seed categories"
        return 1
    fi
    
    # Create indexes
    if /app/bin/create-indexes; then
        print_status "success" "Database indexes created"
    else
        print_status "warning" "Index creation may have failed (might already exist)"
    fi
    
    return 0
}

# Function to run a test and track results
run_test() {
    local test_name=$1
    local test_command=$2
    
    print_status "info" "Running test: $test_name"
    
    if eval "$test_command"; then
        print_status "success" "✅ $test_name PASSED"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        print_status "error" "❌ $test_name FAILED"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Function to test auth service endpoints
test_auth_service() {
    local base_url="$AUTH_SERVICE_URL/api/v1/auth"
    
    print_status "info" "Testing Auth Service..."
    
    # Test 1: Health check
    run_test "Auth Health Check" "curl -f -s '$base_url/health' > /dev/null"
    
    # Test 2: User registration
    local email="test-$(date +%s)@example.com"
    local password="TestPassword123"
    
    local register_response=$(curl -s -X POST "$base_url/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"$email\",
            \"password\": \"$password\",
            \"fullName\": \"Integration Test User\",
            \"familySize\": 2
        }")
    
    if echo "$register_response" | jq -e '.accessToken' > /dev/null 2>&1; then
        print_status "success" "✅ User Registration PASSED"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        
        # Extract token for further tests
        local token=$(echo "$register_response" | jq -r '.accessToken')
        
        # Test 3: Profile retrieval
        if curl -f -s -H "Authorization: Bearer $token" "$base_url/me" > /dev/null; then
            print_status "success" "✅ Profile Retrieval PASSED"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            print_status "error" "❌ Profile Retrieval FAILED"
            TESTS_FAILED=$((TESTS_FAILED + 1))
        fi
        
        # Test 4: Profile update
        local update_response=$(curl -s -X PUT "$base_url/me" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -d "{
                \"email\": \"updated-$email\",
                \"fullName\": \"Updated Test User\",
                \"familySize\": 3
            }")
        
        if echo "$update_response" | jq -e '.success' > /dev/null 2>&1; then
            print_status "success" "✅ Profile Update PASSED"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            print_status "error" "❌ Profile Update FAILED"
            TESTS_FAILED=$((TESTS_FAILED + 1))
        fi
        
        # Test 5: Password change
        local password_response=$(curl -s -X POST "$base_url/change-password" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -d "{
                \"currentPassword\": \"$password\",
                \"newPassword\": \"NewPassword456\"
            }")
        
        if echo "$password_response" | jq -e '.success' > /dev/null 2>&1; then
            print_status "success" "✅ Password Change PASSED"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            print_status "error" "❌ Password Change FAILED"
            TESTS_FAILED=$((TESTS_FAILED + 1))
        fi
        
    else
        print_status "error" "❌ User Registration FAILED"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Function to test financial service endpoints
test_financial_service() {
    local base_url="$FINANCIAL_SERVICE_URL/api/v1"
    
    print_status "info" "Testing Financial Service..."
    
    # Test 1: Health check
    run_test "Financial Health Check" "curl -f -s '$base_url/health' > /dev/null"
    
    # Test 2: Categories endpoint
    run_test "Categories List" "curl -f -s '$base_url/categories' > /dev/null"
    
    # Test 3: Dashboard endpoint (might need auth, but test basic connectivity)
    if curl -s "$base_url/dashboard" | grep -q "error\|unauthorized"; then
        print_status "success" "✅ Dashboard Endpoint Accessible PASSED"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        print_status "error" "❌ Dashboard Endpoint FAILED"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Function to test database connectivity and seeded data
test_database() {
    print_status "info" "Testing Database..."
    
    # Test MongoDB connection
    if mongosh --quiet --eval "db.runCommand('ping')" "$MONGODB_URI/$MONGODB_DATABASE" > /dev/null 2>&1; then
        print_status "success" "✅ MongoDB Connection PASSED"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        print_status "error" "❌ MongoDB Connection FAILED"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
    
    # Test seeded categories
    local category_count=$(mongosh --quiet --eval "db.categories.countDocuments({isSystem: true})" "$MONGODB_URI/$MONGODB_DATABASE" 2>/dev/null)
    
    if [ "$category_count" -gt 30 ]; then
        print_status "success" "✅ Seeded Categories ($category_count found) PASSED"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        print_status "error" "❌ Seeded Categories FAILED (only $category_count found)"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Function to generate test report
generate_test_report() {
    local total_tests=$((TESTS_PASSED + TESTS_FAILED))
    local success_rate=0
    
    if [ $total_tests -gt 0 ]; then
        success_rate=$(( (TESTS_PASSED * 100) / total_tests ))
    fi
    
    # Create test results directory
    mkdir -p /app/test-results
    
    # Generate JSON report
    cat > "$TEST_RESULTS_FILE" << EOF
{
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "total_tests": $total_tests,
    "tests_passed": $TESTS_PASSED,
    "tests_failed": $TESTS_FAILED,
    "success_rate": $success_rate,
    "status": "$([ $TESTS_FAILED -eq 0 ] && echo "PASSED" || echo "FAILED")",
    "environment": {
        "auth_service_url": "$AUTH_SERVICE_URL",
        "financial_service_url": "$FINANCIAL_SERVICE_URL",
        "mongodb_uri": "$MONGODB_URI",
        "mongodb_database": "$MONGODB_DATABASE"
    }
}
EOF
    
    # Print summary
    echo ""
    echo "=================================="
    echo "🧪 Integration Test Results"
    echo "=================================="
    echo "Total Tests: $total_tests"
    echo "Passed: $TESTS_PASSED"
    echo "Failed: $TESTS_FAILED"
    echo "Success Rate: $success_rate%"
    echo ""
    
    if [ $TESTS_FAILED -eq 0 ]; then
        print_status "success" "🎉 ALL TESTS PASSED!"
        echo "Test report saved to: $TEST_RESULTS_FILE"
        return 0
    else
        print_status "error" "💥 SOME TESTS FAILED!"
        echo "Test report saved to: $TEST_RESULTS_FILE"
        return 1
    fi
}

# Main execution
main() {
    print_status "info" "Starting integration tests..."
    
    # Wait for services to be ready
    wait_for_service "$AUTH_SERVICE_URL/api/v1/auth/health" "Auth Service" || exit 1
    wait_for_service "$FINANCIAL_SERVICE_URL/api/v1/health" "Financial Service" || exit 1
    
    # Seed test data
    seed_test_data || exit 1
    
    # Run tests
    test_database
    test_auth_service
    test_financial_service
    
    # Generate report and exit with appropriate code
    generate_test_report
}

# Run main function
main "$@"