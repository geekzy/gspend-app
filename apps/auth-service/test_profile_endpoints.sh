#!/bin/bash

# Test script for profile management endpoints
# This script tests the new profile management functionality

echo "Testing Profile Management Endpoints"
echo "===================================="

# Configuration
AUTH_BASE_URL="http://localhost:8081"
MONGO_URI=${MONGODB_URI:-"mongodb://localhost:27017"}
MONGO_DB=${MONGODB_DATABASE:-"gspend"}

# Test data
EMAIL="testuser@example.com"
PASSWORD="TestPassword123"
NEW_EMAIL="updated@example.com"
NEW_PASSWORD="NewPassword456"

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

# Function to check if auth service is running
check_auth_service() {
    echo "Checking if auth service is running..."
    if curl -s "$AUTH_BASE_URL/api/v1/auth/health" > /dev/null 2>&1; then
        echo "✓ Auth service is running"
        return 0
    else
        echo "❌ Auth service is not running"
        echo "Please start the service: ./bin/auth-service"
        return 1
    fi
}

# Function to clean up test user (to avoid conflicts)
cleanup_test_user() {
    echo "Cleaning up any existing test user..."
    if command -v mongosh &> /dev/null; then
        mongosh --quiet --eval "db.users.deleteOne({email: '$EMAIL'})" "$MONGO_URI/$MONGO_DB" > /dev/null 2>&1
        mongosh --quiet --eval "db.users.deleteOne({email: '$NEW_EMAIL'})" "$MONGO_URI/$MONGO_DB" > /dev/null 2>&1
    elif command -v mongo &> /dev/null; then
        mongo --quiet --eval "db.users.deleteOne({email: '$EMAIL'})" "$MONGO_URI/$MONGO_DB" > /dev/null 2>&1
        mongo --quiet --eval "db.users.deleteOne({email: '$NEW_EMAIL'})" "$MONGO_URI/$MONGO_DB" > /dev/null 2>&1
    fi
    echo "✓ Test user cleanup completed"
}

# Function to validate response
validate_response() {
    local response="$1"
    local expected_status="$2"
    local test_name="$3"
    
    if command -v jq &> /dev/null; then
        local success=$(echo "$response" | jq -r '.success // empty' 2>/dev/null)
        local error=$(echo "$response" | jq -r '.error // empty' 2>/dev/null)
        
        if [ "$success" = "true" ] || [ -n "$(echo "$response" | jq -r '.accessToken // empty' 2>/dev/null)" ]; then
            echo "✓ $test_name successful"
            return 0
        elif [ -n "$error" ]; then
            echo "❌ $test_name failed: $error"
            return 1
        else
            echo "⚠️  $test_name response unclear: $response"
            return 1
        fi
    else
        # Fallback without jq
        if [[ "$response" == *"success\":true"* ]] || [[ "$response" == *"accessToken"* ]]; then
            echo "✓ $test_name successful"
            return 0
        elif [[ "$response" == *"error"* ]]; then
            echo "❌ $test_name failed"
            return 1
        else
            echo "⚠️  $test_name response unclear"
            return 1
        fi
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

if ! check_auth_service; then
    echo "❌ Auth service is not running"
    echo "Please start the service first:"
    echo "  cd apps/auth-service"
    echo "  go run cmd/server/main.go"
    exit 1
fi

cleanup_test_user

echo "✓ All pre-flight checks passed!"
echo ""

# Test 1: Registration with new password validation
echo "1. Testing Registration with new password validation..."
REGISTER_RESPONSE=$(curl -s -X POST "$AUTH_BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"fullName\": \"Test User\",
    \"familySize\": 2
  }")

echo "Register Response: $REGISTER_RESPONSE"

if validate_response "$REGISTER_RESPONSE" "201" "Registration"; then
    # Extract token from response
    if command -v jq &> /dev/null; then
        TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.accessToken // empty')
        if [ -n "$TOKEN" ]; then
            echo "✓ Token obtained successfully"
            
            echo ""
            echo "2. Testing Profile Update..."
            UPDATE_RESPONSE=$(curl -s -X PUT "$AUTH_BASE_URL/api/v1/auth/me" \
              -H "Content-Type: application/json" \
              -H "Authorization: Bearer $TOKEN" \
              -d "{
                \"email\": \"$NEW_EMAIL\",
                \"fullName\": \"Updated User\",
                \"familySize\": 3
              }")
            
            echo "Update Response: $UPDATE_RESPONSE"
            validate_response "$UPDATE_RESPONSE" "200" "Profile Update"
            
            echo ""
            echo "3. Testing Password Change..."
            PASSWORD_RESPONSE=$(curl -s -X POST "$AUTH_BASE_URL/api/v1/auth/change-password" \
              -H "Content-Type: application/json" \
              -H "Authorization: Bearer $TOKEN" \
              -d "{
                \"currentPassword\": \"$PASSWORD\",
                \"newPassword\": \"$NEW_PASSWORD\"
              }")
            
            echo "Password Change Response: $PASSWORD_RESPONSE"
            validate_response "$PASSWORD_RESPONSE" "200" "Password Change"
            
            echo ""
            echo "4. Testing Profile Retrieval..."
            PROFILE_RESPONSE=$(curl -s -X GET "$AUTH_BASE_URL/api/v1/auth/me" \
              -H "Authorization: Bearer $TOKEN")
            
            echo "Profile Response: $PROFILE_RESPONSE"
            validate_response "$PROFILE_RESPONSE" "200" "Profile Retrieval"
            
            echo ""
            echo "5. Testing Login with New Password..."
            LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_BASE_URL/api/v1/auth/login" \
              -H "Content-Type: application/json" \
              -d "{
                \"email\": \"$NEW_EMAIL\",
                \"password\": \"$NEW_PASSWORD\"
              }")
            
            echo "Login Response: $LOGIN_RESPONSE"
            validate_response "$LOGIN_RESPONSE" "200" "Login with New Password"
            
        else
            echo "❌ No token received from registration"
        fi
    else
        echo "⚠️  jq not available, cannot extract token for further testing"
        echo "Install jq for complete testing: brew install jq (macOS) or apt-get install jq (Ubuntu)"
    fi
else
    echo "❌ Registration failed, cannot continue with other tests"
fi

# Cleanup
echo ""
echo "Cleaning up test data..."
cleanup_test_user

echo ""
echo "===================================="
echo "Profile Management Test Complete"
echo ""
echo "Prerequisites verified:"
echo "✓ MongoDB connection"
echo "✓ Auth service running"
echo "✓ Test data cleanup"
echo ""
echo "Tests performed:"
echo "• User registration with password validation"
echo "• Profile update (name, email, family size)"
echo "• Password change with current password verification"
echo "• Profile retrieval"
echo "• Login with updated credentials"