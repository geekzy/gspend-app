# gSpend Testing Guide

This guide explains how to run tests for the gSpend application using Docker containers for CI/CD compatibility.

## Testing Strategy

We use a **multi-layered testing approach**:

1. **Unit Tests** - Fast, isolated tests for individual components
2. **Integration Tests** - Docker-based tests for complete backend workflows  
3. **End-to-End Tests** - Full application testing (future)

## Quick Start

### Unit Tests (Fast)
```bash
# Run all unit tests
make test

# Run with verbose output
make test-v

# Run with coverage
make test-coverage
```

### Integration Tests (Docker-based)
```bash
# Run complete integration test suite
make test-integration

# Clean up test environment
make test-clean
```

## Docker-Based Integration Testing

Our integration tests run in **isolated Docker containers** to ensure:
- ✅ **Consistent environment** across development and CI/CD
- ✅ **No local dependencies** (MongoDB, Redis, services)
- ✅ **Automatic seeding** of test data
- ✅ **Complete cleanup** after tests
- ✅ **CI/CD ready** with GitHub Actions

### What Gets Tested

**Database Layer:**
- ✅ MongoDB connection and operations
- ✅ Seeded categories (30+ family-oriented categories)
- ✅ Database indexes for performance

**Auth Service:**
- ✅ User registration with password validation
- ✅ Profile management (update, password change)
- ✅ JWT token generation and validation
- ✅ Email uniqueness enforcement

**Financial Service:**
- ✅ Service health and connectivity
- ✅ Categories API endpoints
- ✅ Dashboard API accessibility
- ✅ Database integration

### Test Environment

The Docker test environment includes:
- **MongoDB Test Instance** (port 27018)
- **Redis Test Instance** (port 6380)  
- **Auth Service** (port 8083)
- **Financial Service** (port 8084)
- **Test Runner** with automated test execution

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

The prerequisites script handles seeding automatically, but if you need manual intervention:

**Force re-seed categories:**
```bash
# Remove existing categories first
mongosh gspend --eval "db.categories.deleteMany({isSystem: true})"
# Then run the script again
./scripts/verify_test_prerequisites.sh
```

**Clean up test users:**
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

1. **Use the verification functions** - Import checks from the prerequisites script
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

## CI/CD Integration

The Docker-based tests are designed for **GitHub Actions** and other CI/CD systems:

```yaml
# .github/workflows/ci.yml
- name: Run Integration Tests
  run: make test-integration
```

**Benefits for CI/CD:**
- ✅ No external dependencies
- ✅ Consistent test environment  
- ✅ Automatic cleanup
- ✅ JSON test reports
- ✅ Proper exit codes for CI systems

## Legacy Local Testing

For local development, you can still use the original scripts:

```bash
# Verify local prerequisites
./scripts/verify_test_prerequisites.sh

# Run local tests (requires manual service startup)
cd apps/auth-service && ./test_profile_endpoints.sh
cd apps/financial-service && ./test_budget_items.sh
```

## Test Output

**Docker Integration Tests:**
- 🧪 Comprehensive test execution report
- 📊 JSON test results in `test-results/`
- ✅ Clear pass/fail indicators
- 📈 Success rate calculation

**Unit Tests:**
- 🔍 Individual component testing
- 📊 Code coverage reports
- ⚡ Fast execution (< 30 seconds)

## Environment Variables

Customize test behavior:

```bash
# For Docker tests
export MONGODB_DATABASE=custom_test_db
export AUTH_SERVICE_URL=http://custom-auth:8081

# For local tests  
export MONGODB_URI=mongodb://localhost:27017
export MONGODB_DATABASE=gspend_test
```

## Troubleshooting

### Docker Issues
```bash
# Clean up test environment
make test-clean

# Check Docker resources
docker system df
docker system prune -f
```

### Test Failures
```bash
# View detailed test results
cat test-results/integration-test-results.json | jq '.'

# Check container logs
docker-compose -f docker-compose.test.yml logs test-runner
```

## Contributing

When adding new tests:

1. **Unit Tests** - Add to existing `*_test.go` files
2. **Integration Tests** - Extend `scripts/integration-test.sh`
3. **Docker Tests** - Update `docker-compose.test.yml` if needed
4. **CI/CD** - Tests run automatically on PR/push

## Next Steps

After successful testing:
- ✅ All backend services are working
- ✅ Database is properly seeded  
- ✅ API endpoints are functional
- ✅ Ready for frontend integration
- ✅ CI/CD pipeline validated