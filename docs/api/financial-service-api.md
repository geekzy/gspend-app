# Financial Service API Specification

Base URL: `/api/v1`

All endpoints require authentication via Bearer token unless specified otherwise.

---

## Income Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/incomes` | List all incomes |
| POST | `/incomes` | Create new income |
| GET | `/incomes/:id` | Get income details |
| PUT | `/incomes/:id` | Update income |
| DELETE | `/incomes/:id` | Delete income |
| GET | `/incomes/summary` | Get income summary |

### POST /incomes

**Request:**
```json
{
  "source": "Salary",
  "amount": 5000.00,
  "frequency": "monthly",
  "effective_date": "2024-01-01"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "64f1a2b3c4d5e6f7a8b9c0d1",
    "source": "Salary",
    "amount": 5000.00,
    "frequency": "monthly",
    "effective_date": "2024-01-01T00:00:00Z"
  }
}
```

---

## Budget Planning

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/budgets` | List all budgets |
| POST | `/budgets` | Create new budget |
| GET | `/budgets/:id` | Get budget with items |
| PUT | `/budgets/:id` | Update budget |
| DELETE | `/budgets/:id` | Delete budget |
| GET | `/budgets/current` | Get current active budget |
| POST | `/budgets/:id/items` | Add budget item |
| PUT | `/budgets/:id/items/:item_id` | Update budget item |
| DELETE | `/budgets/:id/items/:item_id` | Delete budget item |

### POST /budgets

**Request:**
```json
{
  "name": "January 2024 Budget",
  "period_type": "monthly",
  "start_date": "2024-01-01",
  "end_date": "2024-01-31",
  "total_amount": 4000.00,
  "items": [
    {
      "category_id": "64f1a2b3...",
      "planned_amount": 800.00
    }
  ]
}
```

---

## Transaction Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/transactions` | List transactions with filters |
| POST | `/transactions` | Create new transaction |
| GET | `/transactions/:id` | Get transaction details |
| PUT | `/transactions/:id` | Update transaction |
| DELETE | `/transactions/:id` | Delete transaction |
| GET | `/transactions/stats` | Get transaction statistics |

### Query Parameters for GET /transactions

| Parameter | Type | Description |
|-----------|------|-------------|
| `start_date` | string | Start date filter (YYYY-MM-DD) |
| `end_date` | string | End date filter (YYYY-MM-DD) |
| `category_id` | string | Filter by category |
| `type` | string | Filter by type (income/expense) |
| `page` | integer | Page number (default: 1) |
| `per_page` | integer | Items per page (default: 20, max: 100) |
| `sort` | string | Sort field (transaction_date, amount) |
| `order` | string | Sort order (asc/desc) |

### POST /transactions

**Request:**
```json
{
  "category_id": "64f1a2b3c4d5e6f7a8b9c0d1",
  "budget_id": "64f1a2b3c4d5e6f7a8b9c0d2",
  "type": "expense",
  "amount": 150.00,
  "description": "Weekly groceries",
  "transaction_date": "2024-01-15",
  "payment_method": "credit_card",
  "notes": "Including household items"
}
```

---

## Category Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/categories` | List all categories |
| POST | `/categories` | Create custom category |
| GET | `/categories/:id` | Get category details |
| PUT | `/categories/:id` | Update category |
| DELETE | `/categories/:id` | Delete category (custom only) |
| GET | `/categories/system` | Get system categories |

---

## Dashboard

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/dashboard/summary` | Get dashboard overview |

### GET /dashboard/summary

**Response:**
```json
{
  "success": true,
  "data": {
    "total_balance": 15000.00,
    "monthly_income": 6000.00,
    "monthly_expenses": 4500.00,
    "monthly_budget": 5000.00,
    "budget_used_percentage": 90.0,
    "recent_transactions": [...],
    "top_categories": [...]
  }
}
```

---

## Reports

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/reports/budget-vs-actual` | Budget vs actual report |
| GET | `/reports/spending-by-category` | Spending by category |
| GET | `/reports/monthly-trends` | Monthly trends report |

### GET /reports/budget-vs-actual

**Query Parameters:**
- `month`: Month filter (YYYY-MM, default: current month)

### GET /reports/spending-by-category

**Query Parameters:**
- `start_date`: Start date (YYYY-MM-DD)
- `end_date`: End date (YYYY-MM-DD)

### GET /reports/monthly-trends

**Query Parameters:**
- `months`: Number of months (1-12, default: 3)

---

## Health Check

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/health` | Service health check | No |

---

## Pagination Response

```json
{
  "success": true,
  "data": [...],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}
```
