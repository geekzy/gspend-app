# gSpend Testing Guide

This guide explains how to run functional tests for the gSpend application.

## Prerequisites

Before running any functional tests, ensure all prerequisites are met by running:

```bash
./scripts/verify_test_prerequisites.sh
```

This script will check:
- ✅ MongoDB connection
- ✅ Seeded system categories
- ✅ Database indexes
- ✅ Auth service running
- ✅ Financial service running
- ✅ Required tools (curl, jq, go)

## Automatic Setup

The verification script can automatically seed the database if needed:

```bash
# Run verification and follow prompts for automatic seeding
./scripts/verify_test_prerequisites.sh
```

## Manual Setup

If you prefer manual setup:

### 1. Start MongoDB
```bash
# Using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or using local MongoDB installation
mongod --dbpath /path/to/your/db
```

### 2. Seed Categories
```bash
cd apps/financial-service
go run scripts/seed_categories.go
```

### 3. Create Database Indexes
```bash
cd apps/financial-service
go run scripts/create_indexes.go
```

### 4. Start Services

**Auth Service:**
```bash
cd apps/auth-service
go run cmd/server/main.go
# Runs on http://localhost:8081
```

**Financial Service:**
```bash
cd apps/financial-service
go run cmd/server/main.go
# Runs on http://localhost:8082
```

## Running Tests

### Profile Management Tests
Tests user registration, profile updates, and password changes:

```bash
cd apps/auth-service
./test_profile_endpoints.sh
```

**What it tests:**
- ✅ User registration with password validation
- ✅ Profile update (name, email, family size)
- ✅ Password change with current password verification
- ✅ Profile retrieval
- ✅ Login with updated credentials

### Budget Item Management Tests
Tests budget item CRUD operations:

```bash
cd apps/financial-service
./test_budget_items.sh
```

**What it tests:**
- ✅ Add budget item with real category data
- ✅ Get budget item
- ✅ Update budget item
- ✅ Delete budget item

## Test Features

### Smart Prerequisites Checking
- ✅ Verifies MongoDB connection
- ✅ Checks for seeded categories (uses real category IDs)
- ✅ Validates service availability
- ✅ Cleans up test data automatically

### Robust Error Handling
- ✅ Clear error messages with solutions
- ✅ Graceful fallbacks when tools are missing
- ✅ Automatic test data cleanup

### Real Data Integration
- ✅ Uses actual seeded category IDs
- ✅ Validates responses with proper JSON parsing
- ✅ Tests complete workflows end-to-end

## Environment Variables

You can customize the test environment:

```bash
# MongoDB connection
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE="gspend"

# Service URLs (if different from defaults)
export AUTH_SERVICE_URL="http://localhost:8081"
export FINANCIAL_SERVICE_URL="http://localhost:8082"
```

## Troubleshooting

### Common Issues

**MongoDB not accessible:**
```bash
# Check if MongoDB is running
mongosh --eval "db.runCommand('ping')"

# Start MongoDB if needed
brew services start mongodb-community  # macOS
sudo systemctl start mongod            # Linux
```

**Services not starting:**
```bash
# Check if ports are in use
lsof -i :8081  # Auth service
lsof -i :8082  # Financial service

# Kill processes if needed
kill -9 <PID>
```

**Missing tools:**
```bash
# Install jq for better JSON parsing
brew install jq        # macOS
apt-get install jq     # Ubuntu

# Install curl if missing
brew install curl      # macOS
apt-get install curl   # Ubuntu
```

### Test Data Issues

**Categories not seeded:**
```bash
cd apps/financial-service
go run scripts/seed_categories.go
```

**Database indexes missing:**
```bash
cd apps/financial-service
go run scripts/create_indexes.go
```

**Test user conflicts:**
The test scripts automatically clean up test users, but you can manually clean up:
```bash
mongosh gspend --eval "db.users.deleteMany({email: /test/})"
```

## Test Output

Successful test output will show:
- ✅ Green checkmarks for passed tests
- ❌ Red X marks for failed tests  
- ⚠️ Yellow warnings for non-critical issues
- Clear response data and HTTP status codes

## Contributing

When adding new functional tests:

1. **Check prerequisites** - Use the verification functions
2. **Clean up data** - Remove test data after tests
3. **Use real data** - Integrate with seeded categories/users
4. **Validate responses** - Parse JSON and check success/error states
5. **Provide clear output** - Use colored status indicators

## Next Steps

After successful testing:
- ✅ All backend services are working
- ✅ Database is properly seeded
- ✅ API endpoints are functional
- ✅ Ready for frontend integration testing