# Integration Testing Guide

This guide covers the integration testing setup and execution for the GSpend application.

## Overview

The integration testing framework provides comprehensive end-to-end testing of the GSpend backend services in a containerized environment that closely mirrors production.

## Quick Start

### Run Full Integration Tests
```bash
make test-integration
```

### Run Quick Integration Tests (Reuse Containers)
```bash
make test-integration-quick
```

### Setup Test Environment Only
```bash
make test-integration-setup
```

### Teardown Test Environment
```bash
make test-integration-teardown
```

## Test Architecture

### Test Environment Components
- **MongoDB Test Instance** - Isolated database on port 27018
- **Redis Test Instance** - Isolated cache on port 6380  
- **Auth Service Test** - Authentication service on port 8083
- **Financial Service Test** - Financial service on port 8084
- **Test Runner Container** - Executes integration test suite

### Test Categories

#### 1. Service Health Tests
- Verify all services start successfully
- Check health endpoints respond correctly
- Validate service dependencies

#### 2. Authentication Tests
- User registration flow
- Login/logout functionality
- JWT token validation
- Profile management (view/update)
- Password change functionality

#### 3. Financial Service Tests
- Categories endpoint accessibility
- Dashboard endpoint connectivity
- Basic API response validation

#### 4. Database Tests
- MongoDB connectivity
- Seeded data validation
- Index creation verification

## Test Data Management

### Automatic Data Seeding
The test environment automatically seeds:
- **System Categories** - Default expense and income categories
- **Database Indexes** - Performance optimization indexes
- **Test Users** - Generated during test execution

### Data Isolation
- Each test run uses a fresh database
- Test data is automatically cleaned up
- No interference with development data

## Test Execution Workflow

### Full Integration Test (`make test-integration`)
1. **Cleanup** - Remove any existing test containers
2. **Build** - Build latest service images
3. **Start Services** - Launch test environment
4. **Wait for Health** - Ensure all services are ready
5. **Run Tests** - Execute comprehensive test suite
6. **Generate Report** - Create JSON test results
7. **Cleanup** - Remove test containers and volumes

### Quick Integration Test (`make test-integration-quick`)
1. **Check Existing** - Look for running test containers
2. **Reuse or Rebuild** - Use existing containers if available
3. **Run Tests** - Execute test suite against existing environment
4. **Generate Report** - Create JSON test results

## Test Results

### JSON Report Format
```json
{
    "timestamp": "2024-01-15T10:30:00Z",
    "total_tests": 12,
    "tests_passed": 11,
    "tests_failed": 1,
    "success_rate": 92,
    "status": "FAILED",
    "environment": {
        "auth_service_url": "http://auth-service-test:8081",
        "financial_service_url": "http://financial-service-test:8082",
        "mongodb_uri": "mongodb://mongodb-test:27017",
        "mongodb_database": "gspend_test"
    }
}
```

### Report Location
- **File**: `test-results/integration-test-results.json`
- **View**: `make test-show-results`

## Debugging and Troubleshooting

### View Test Logs
```bash
# All services
make test-integration-logs

# Follow logs in real-time
make test-integration-logs-follow

# Individual service logs
docker-compose -f docker-compose.test.yml logs auth-service-test
docker-compose -f docker-compose.test.yml logs financial-service-test
```

### Manual Test Environment
```bash
# Setup environment and keep it running
make test-integration-setup

# Run manual tests
curl http://localhost:8083/api/v1/auth/health
curl http://localhost:8084/api/v1/health

# Cleanup when done
make test-integration-teardown
```

### Common Issues

#### Services Not Starting
```bash
# Check container status
docker-compose -f docker-compose.test.yml ps

# Check specific service logs
docker-compose -f docker-compose.test.yml logs mongodb-test
docker-compose -f docker-compose.test.yml logs auth-service-test
```

#### Port Conflicts
The test environment uses these ports:
- **27018** - MongoDB Test
- **6380** - Redis Test  
- **8083** - Auth Service Test
- **8084** - Financial Service Test

Ensure these ports are available before running tests.

#### Database Connection Issues
```bash
# Test MongoDB connectivity
docker-compose -f docker-compose.test.yml exec mongodb-test mongosh gspend_test --eval "db.runCommand('ping')"

# Check seeded data
docker-compose -f docker-compose.test.yml exec mongodb-test mongosh gspend_test --eval "db.categories.countDocuments()"
```

## CI/CD Integration

### GitHub Actions Example
```yaml
name: Integration Tests
on: [push, pull_request]

jobs:
  integration-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Integration Tests
        run: make test-integration
      - name: Upload Test Results
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: integration-test-results
          path: test-results/
```

### Local Development Workflow
```bash
# During development
make test-integration-setup  # Setup once
make test-integration-quick  # Run tests quickly
make test-integration-quick  # Run again after changes
make test-integration-teardown  # Cleanup when done

# Before committing
make test-integration  # Full test run
```

## Performance Considerations

### Test Execution Times
- **Full Integration Test**: ~3-5 minutes (includes build time)
- **Quick Integration Test**: ~30-60 seconds (reuses containers)
- **Setup Only**: ~2-3 minutes
- **Teardown Only**: ~10-20 seconds

### Resource Usage
- **Memory**: ~2GB RAM recommended
- **Disk**: ~1GB for images and volumes
- **CPU**: Moderate usage during test execution

### Optimization Tips
1. Use `test-integration-quick` for rapid iteration
2. Keep test environment running during development
3. Use `test-integration-setup` once, then multiple `test-integration-quick`
4. Clean up regularly with `test-integration-teardown`

## Extending Tests

### Adding New Test Cases
Edit `scripts/integration-test.sh`:

```bash
# Add new test function
test_new_feature() {
    local base_url="$FINANCIAL_SERVICE_URL/api/v1"
    
    run_test "New Feature Test" "curl -f -s '$base_url/new-endpoint' > /dev/null"
}

# Call from main function
main() {
    # ... existing tests ...
    test_new_feature
    # ... rest of main ...
}
```

### Custom Test Environment
Create custom docker-compose file:
```yaml
# docker-compose.custom-test.yml
services:
  # ... copy from docker-compose.test.yml ...
  # Add custom services or configurations
```

Run with custom environment:
```bash
docker-compose -f docker-compose.custom-test.yml up -d
# Run custom tests
docker-compose -f docker-compose.custom-test.yml down
```

## Best Practices

### Test Development
1. **Isolation** - Each test should be independent
2. **Cleanup** - Tests should clean up their own data
3. **Idempotency** - Tests should produce same results when run multiple times
4. **Fast Feedback** - Use quick tests for rapid development cycles

### Environment Management
1. **Fresh Start** - Use full integration tests for CI/CD
2. **Development Speed** - Use quick tests during development
3. **Resource Management** - Clean up test environments regularly
4. **Monitoring** - Check test logs for performance issues

### Debugging
1. **Incremental Testing** - Test individual components first
2. **Log Analysis** - Use detailed logging for troubleshooting
3. **Manual Verification** - Use manual endpoints to verify behavior
4. **Environment Inspection** - Examine database state during failures