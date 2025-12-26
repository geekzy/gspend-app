#!/bin/bash

# Test script for profile management endpoints
# This script tests the new profile management functionality

echo "Testing Profile Management Endpoints"
echo "===================================="

# Test data
EMAIL="testuser@example.com"
PASSWORD="TestPassword123"
NEW_EMAIL="updated@example.com"
NEW_PASSWORD="NewPassword456"

echo "1. Testing Registration with new password validation..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"fullName\": \"Test User\",
    \"familySize\": 2
  }")

echo "Register Response: $REGISTER_RESPONSE"

# Extract token from response (assuming jq is available)
if command -v jq &> /dev/null; then
    TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.accessToken // empty')
    if [ -n "$TOKEN" ]; then
        echo "✓ Registration successful, token obtained"
        
        echo "2. Testing Profile Update..."
        UPDATE_RESPONSE=$(curl -s -X PUT http://localhost:8081/api/v1/auth/me \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $TOKEN" \
          -d "{
            \"email\": \"$NEW_EMAIL\",
            \"fullName\": \"Updated User\",
            \"familySize\": 3
          }")
        
        echo "Update Response: $UPDATE_RESPONSE"
        
        echo "3. Testing Password Change..."
        PASSWORD_RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/auth/change-password \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $TOKEN" \
          -d "{
            \"currentPassword\": \"$PASSWORD\",
            \"newPassword\": \"$NEW_PASSWORD\"
          }")
        
        echo "Password Change Response: $PASSWORD_RESPONSE"
        
        echo "4. Testing Profile Retrieval..."
        PROFILE_RESPONSE=$(curl -s -X GET http://localhost:8081/api/v1/auth/me \
          -H "Authorization: Bearer $TOKEN")
        
        echo "Profile Response: $PROFILE_RESPONSE"
        
    else
        echo "✗ Registration failed or no token received"
    fi
else
    echo "jq not available, cannot extract token for further testing"
fi

echo "===================================="
echo "Profile Management Test Complete"