# Testing Guide

## Overview

gSpend uses multiple testing strategies to ensure code quality.

---

## Backend Testing (Go)

### Running Tests

```bash
# Run all tests
make test

# Run with verbose output
make test-v

# Run with coverage
make test-coverage

# Test specific service
make test-auth
make test-finance
```

### Unit Tests

Located alongside source files with `_test.go` suffix.

```go
// auth_service_test.go
func TestAuthService_ValidateToken(t *testing.T) {
    // Arrange
    mockRepo := NewMockUserRepository()
    service := NewAuthService(mockRepo)
    
    // Act
    result, err := service.ValidateToken("valid-token")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Test Structure
```
internal/
├── service/
│   ├── auth_service.go
│   └── auth_service_test.go
├── handler/
│   ├── auth_handler.go
│   └── auth_handler_test.go
```

### Coverage Target

**Goal: 80%+ coverage**

```bash
# View coverage report
make test-auth-coverage
make test-finance-coverage
```

---

## Frontend Testing (TypeScript)

### Running Tests

```bash
cd apps/frontend

# Run tests
npm run test

# Run with coverage
npm run test:coverage

# Run in watch mode
npm run test:watch
```

### Component Testing

```typescript
// TransactionList.spec.ts
import { mount } from '@vue/test-utils'
import TransactionList from './TransactionList.vue'

describe('TransactionList', () => {
  it('renders transactions', () => {
    const wrapper = mount(TransactionList, {
      props: {
        transactions: [
          { id: '1', amount: 100, description: 'Test' }
        ]
      }
    })
    
    expect(wrapper.text()).toContain('Test')
  })
})
```

---

## Integration Testing

### Running Integration Tests

```bash
# Full integration test suite
make test-integration

# Quick run (reuse containers)
make test-integration-quick

# Setup test environment
make test-integration-setup

# Teardown
make test-integration-teardown
```

### Test Environment

Uses `docker-compose.test.yml` with isolated:
- MongoDB (port 27018)
- Redis (port 6380)
- Auth Service (port 8083)
- Financial Service (port 8084)

---

## Test Data

### Factories
Create test data using factory functions:

```go
func CreateTestUser(t *testing.T) *domain.User {
    return &domain.User{
        ID:         primitive.NewObjectID(),
        Email:      "test@example.com",
        FullName:   "Test User",
        FamilySize: 3,
    }
}
```

### Demo Data
Use demo environment for manual testing:

```bash
make demo-start
# Login: demo@gspend.com / password
```

---

## Best Practices

1. **Test naming**: `Test<Function>_<Scenario>_<Expected>`
2. **Arrange-Act-Assert** pattern
3. **Mock external dependencies**
4. **Test edge cases**
5. **Keep tests fast and isolated**

---

## CI/CD Integration

Tests run automatically on:
- Pull requests
- Push to `main` or `develop`

```yaml
# .github/workflows/ci.yml
- run: make test
- run: make test-coverage
```
