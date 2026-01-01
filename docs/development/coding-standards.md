# Coding Standards

## General Principles

- **Readability** - Code should be self-documenting
- **Simplicity** - Prefer simple solutions  
- **Consistency** - Follow existing patterns
- **Testing** - Write tests for new code

---

## Go (Backend)

### Project Structure
```
internal/
├── config/      # Configuration
├── domain/      # Business entities
├── repository/  # Data access
├── service/     # Business logic
├── handler/     # HTTP handlers
├── middleware/  # Middleware
├── dto/         # Data transfer objects
└── util/        # Utilities
```

### Naming Conventions
```go
// Packages: lowercase, single word
package handler

// Interfaces: verb + "er" suffix
type UserRepository interface { ... }

// Structs: PascalCase
type AuthService struct { ... }

// Functions/Methods: PascalCase for exported
func (s *AuthService) ValidateToken(token string) error

// Variables: camelCase
userID := ctx.Get("user_id")

// Constants: PascalCase or SCREAMING_SNAKE_CASE
const MaxRetries = 3
const HTTP_TIMEOUT = 30
```

### Error Handling
```go
// Always check errors
user, err := s.repo.FindByID(ctx, id)
if err != nil {
    return nil, err
}

// Wrap errors with context
return nil, fmt.Errorf("failed to find user: %w", err)
```

### Logging
```go
// Use structured logging (Zap)
logger.Info("User logged in",
    zap.String("user_id", userID),
    zap.Duration("duration", duration),
)
```

---

## TypeScript (Frontend)

### Project Structure
```
src/
├── components/   # Vue components
├── composables/  # Composition API hooks
├── services/     # API services
├── stores/       # Pinia stores
├── views/        # Page components
├── router/       # Vue Router
└── utils/        # Utilities
```

### Naming Conventions
```typescript
// Files: PascalCase for components, camelCase for others
TransactionList.vue
useAuth.ts
auth.service.ts

// Components: PascalCase
<TransactionList />

// Functions: camelCase
function formatCurrency(amount: number): string

// Constants: SCREAMING_SNAKE_CASE
const API_BASE_URL = '/api/v1'

// Types/Interfaces: PascalCase
interface Transaction { ... }
type TransactionType = 'income' | 'expense'
```

### Vue Component Structure
```vue
<template>
  <!-- Template -->
</template>

<script setup lang="ts">
// 1. Imports
// 2. Props & Emits
// 3. Composables
// 4. State (ref/reactive)
// 5. Computed
// 6. Methods
// 7. Lifecycle hooks
</script>

<style scoped>
/* Component styles */
</style>
```

---

## API Design

### Response Format
```json
{
  "success": true,
  "data": { ... },
  "meta": { ... }
}
```

### Error Format
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

### URL Conventions
- Use plural nouns: `/api/v1/transactions`
- Use lowercase with hyphens: `/budget-vs-actual`
- Version prefix: `/api/v1/`

---

## Database

### MongoDB Collections
- Use snake_case: `budget_items`
- Singular names: `user`, `transaction`

### Field Names
- Use snake_case in storage
- Convert to camelCase in API responses
