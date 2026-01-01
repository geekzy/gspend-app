# Auth Service API Specification

Base URL: `/api/v1/auth`

## Endpoints

### Authentication

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/register` | Register new user | No |
| POST | `/login` | User login | No |
| POST | `/refresh` | Refresh JWT token | Yes (Refresh) |
| POST | `/logout` | Logout and invalidate token | Yes |
| GET | `/me` | Get current user profile | Yes |
| PUT | `/profile` | Update user profile | Yes |
| PUT | `/password` | Change password | Yes |
| GET | `/health` | Health check | No |

---

## Request/Response Examples

### POST /register

**Request:**
```json
{
  "email": "family@example.com",
  "password": "SecurePass123",
  "full_name": "John Doe",
  "family_size": 5
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "64f1a2b3c4d5e6f7a8b9c0d1",
    "email": "family@example.com",
    "full_name": "John Doe",
    "family_size": 5,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### POST /login

**Request:**
```json
{
  "email": "family@example.com",
  "password": "SecurePass123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

### GET /me

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "64f1a2b3c4d5e6f7a8b9c0d1",
    "email": "family@example.com",
    "full_name": "John Doe",
    "family_size": 5,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### PUT /profile

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request:**
```json
{
  "full_name": "John Smith",
  "family_size": 4
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "64f1a2b3c4d5e6f7a8b9c0d1",
    "email": "family@example.com",
    "full_name": "John Smith",
    "family_size": 4,
    "updated_at": "2024-01-16T14:20:00Z"
  }
}
```

### PUT /password

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request:**
```json
{
  "current_password": "OldPass123",
  "new_password": "NewSecurePass456"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password updated successfully"
}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  }
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or expired token"
  }
}
```

### 409 Conflict
```json
{
  "success": false,
  "error": {
    "code": "EMAIL_EXISTS",
    "message": "Email already registered"
  }
}
```

---

## Password Requirements

- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
