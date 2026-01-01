# Database Schema

## MongoDB Collections

gSpend uses MongoDB 7 with the following collections:

---

## users

Stores user account information.

```javascript
{
  _id: ObjectId,
  email: String,           // unique
  password_hash: String,
  full_name: String,
  family_size: Number,     // 1-5
  created_at: ISODate,
  updated_at: ISODate,
  deleted_at: ISODate      // null for active users
}
```

**Indexes:**
- `{ email: 1 }` (unique)
- `{ deleted_at: 1 }`

---

## categories

Stores expense and income categories.

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,       // null for system categories
  name: String,
  type: String,            // "expense" | "income"
  icon: String,
  color: String,
  is_system: Boolean,
  sort_order: Number,
  created_at: ISODate,
  updated_at: ISODate
}
```

**Indexes:**
- `{ user_id: 1, type: 1 }`
- `{ is_system: 1 }`
- `{ user_id: 1, sort_order: 1 }`

---

## incomes

Stores income records.

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  source: String,
  amount: Decimal128,
  frequency: String,       // "one-time" | "daily" | "weekly" | "monthly" | "yearly"
  effective_date: ISODate,
  created_at: ISODate,
  updated_at: ISODate
}
```

**Indexes:**
- `{ user_id: 1, effective_date: -1 }`
- `{ user_id: 1, frequency: 1 }`

---

## budgets

Stores budget plans with embedded items.

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  name: String,
  period_type: String,     // "monthly" | "quarterly" | "yearly"
  start_date: ISODate,
  end_date: ISODate,
  total_amount: Decimal128,
  items: [
    {
      _id: ObjectId,
      category_id: ObjectId,
      category_name: String,    // denormalized
      planned_amount: Decimal128,
      spent_amount: Decimal128,
      notes: String
    }
  ],
  created_at: ISODate,
  updated_at: ISODate
}
```

**Indexes:**
- `{ user_id: 1, start_date: -1, end_date: -1 }`
- `{ user_id: 1, period_type: 1 }`
- `{ "items.category_id": 1 }`

---

## transactions

Stores all financial transactions.

```javascript
{
  _id: ObjectId,
  user_id: ObjectId,
  category_id: ObjectId,
  budget_id: ObjectId,     // nullable
  type: String,            // "income" | "expense"
  amount: Decimal128,
  description: String,
  transaction_date: ISODate,
  payment_method: String,
  notes: String,
  metadata: {
    category_name: String,  // denormalized
    budget_name: String     // denormalized
  },
  created_at: ISODate,
  updated_at: ISODate
}
```

**Indexes:**
- `{ user_id: 1, transaction_date: -1 }`
- `{ user_id: 1, category_id: 1 }`
- `{ user_id: 1, type: 1, transaction_date: -1 }`
- `{ budget_id: 1 }`
- `{ transaction_date: 1 }`
- Compound: `{ user_id: 1, transaction_date: -1, type: 1, category_id: 1 }`

---

## Design Patterns

| Pattern | Usage |
|---------|-------|
| Embedded Documents | Budget items within budgets |
| Denormalization | Category/budget names in transactions |
| Decimal128 | Monetary precision |
| Soft Deletes | User deletion recovery |
| Compound Indexes | Query optimization |

---

## Entity Relationships

```
users ──┬── incomes
        ├── budgets ─── budget_items
        ├── transactions
        └── categories (custom)

categories (system) ──┬── budget_items
                      └── transactions
```
