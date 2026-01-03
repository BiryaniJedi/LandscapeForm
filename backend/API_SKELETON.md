# API Skeleton - Implementation Guide

This document explains the API skeleton structure and what you need to fill in.

## ✅ What's Implemented

### Structure
- **Handlers** (`internal/handlers/`) - HTTP request/response handlers
- **Middleware** (`internal/middleware/`) - Auth, logging, CORS, recovery
- **Router Setup** (`cmd/api/main.go`) - Chi router with all endpoints configured
- **Request/Response Types** (`internal/handlers/types.go`) - JSON schemas

### Available Endpoints

All endpoints are configured and ready to use (with auth placeholder):

```
GET    /health                    - Health check (no auth required)
GET    /api/forms                 - List all forms for authenticated user
POST   /api/forms/shrub           - Create a shrub form
POST   /api/forms/pesticide       - Create a pesticide form
GET    /api/forms/{id}            - Get a specific form
PUT    /api/forms/{id}            - Update a form
DELETE /api/forms/{id}            - Delete a form
```

### Middleware Stack

Every request goes through:
1. **Recovery** - Catches panics and returns 500
2. **Logger** - Logs method, path, status, duration
3. **CORS** - Allows cross-origin requests from frontend
4. **RequestID** - Adds unique ID to each request
5. **RealIP** - Extracts real client IP

Protected routes also use:
6. **AuthMiddleware** - Validates JWT token (placeholder)

## ❌ What You Need to Implement

### 1. Authentication (HIGH PRIORITY)

**File:** `backend/internal/middleware/auth.go:17-27`

Current placeholder:
```go
// TODO: Validate JWT token here
// - Parse the token
// - Verify signature
// - Check expiration
// - Extract user ID from claims

// PLACEHOLDER: For now, just check token is not empty
userID := "placeholder-user-id"
```

**What to add:**
- JWT parsing and validation
- Token signature verification
- Expiration check
- Extract user ID from JWT claims
- Return proper errors for invalid tokens

**Recommended libraries:**
- `github.com/golang-jwt/jwt/v5` - Industry standard JWT library

### 2. Request Validation (MEDIUM PRIORITY)

All handlers have `// TODO: Add validation` comments. You need to validate:

**Shrub Forms:**
- `first_name` - Required, non-empty string
- `last_name` - Required, non-empty string
- `home_phone` - Required, valid phone format
- `num_shrubs` - Required, integer > 0

**Pesticide Forms:**
- `first_name` - Required, non-empty string
- `last_name` - Required, non-empty string
- `home_phone` - Required, valid phone format
- `pesticide_name` - Required, non-empty string

**Recommended approach:**
- Create a `validators` package
- Use a validation library like `github.com/go-playground/validator/v10`
- Return proper 400 errors with field-specific messages

### 3. CORS Configuration (MEDIUM PRIORITY)

**File:** `backend/internal/middleware/cors.go:8-10`

Current: Allows all origins (`*`)

**What to change:**
```go
// Development
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

// Production
allowedOrigin := os.Getenv("ALLOWED_ORIGIN") // from .env
w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
```

### 4. Pagination (LOW PRIORITY)

**File:** `backend/internal/handlers/forms.go:99`

Current: Returns all forms

**What to add:**
- Parse `limit` and `offset` query params
- Update repository method to support pagination
- Return pagination metadata in response

### 5. Authentication Endpoints (LOW PRIORITY)

Create new handler for user registration and login:

**New file:** `backend/internal/handlers/auth.go`

Endpoints to implement:
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login (returns JWT)
- `POST /api/auth/refresh` - Refresh JWT token

### 6. Admin Endpoints (LOW PRIORITY)

**New file:** `backend/internal/handlers/admin.go`

Endpoints to implement:
- `GET /api/admin/forms` - List all forms (all users)
- Admin-only middleware

## 📝 Testing the API

### Without Authentication (Quick Test)

Temporarily comment out the auth middleware to test endpoints:

```go
// In cmd/api/main.go, line 56
// r.Use(middleware.AuthMiddleware) // Comment this out
```

### Example Requests

**Health Check:**
```bash
curl http://localhost:8080/health
```

**Create Shrub Form:**
```bash
curl -X POST http://localhost:8080/api/forms/shrub \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token-here" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "home_phone": "555-1234",
    "num_shrubs": 5
  }'
```

**List Forms:**
```bash
curl http://localhost:8080/api/forms?sort_by=created_at&order=DESC \
  -H "Authorization: Bearer your-token-here"
```

**Get Form by ID:**
```bash
curl http://localhost:8080/api/forms/{form-id} \
  -H "Authorization: Bearer your-token-here"
```

**Update Form:**
```bash
curl -X PUT http://localhost:8080/api/forms/{form-id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token-here" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith",
    "home_phone": "555-5678",
    "num_shrubs": 10
  }'
```

**Delete Form:**
```bash
curl -X DELETE http://localhost:8080/api/forms/{form-id} \
  -H "Authorization: Bearer your-token-here"
```

## 🏗️ Architecture

```
HTTP Request
    ↓
Middleware Stack (Recovery, Logger, CORS, Auth)
    ↓
Chi Router (route matching)
    ↓
Handler (forms.go)
    ↓
Repository (internal/forms/forms.go)
    ↓
Database (PostgreSQL)
```

## 📁 File Structure

```
backend/
├── cmd/api/
│   └── main.go                    # ✅ Router setup and server start
├── internal/
│   ├── handlers/
│   │   ├── forms.go               # ✅ Form CRUD handlers
│   │   ├── types.go               # ✅ Request/response types
│   │   └── responses.go           # ✅ Response helpers
│   ├── middleware/
│   │   ├── auth.go                # ⚠️ JWT validation needed
│   │   ├── logger.go              # ✅ Request logging
│   │   ├── recovery.go            # ✅ Panic recovery
│   │   └── cors.go                # ⚠️ Configure allowed origins
│   ├── forms/
│   │   ├── forms.go               # ✅ Repository (already complete)
│   │   └── models.go              # ✅ Domain models
│   └── db/
│       └── db.go                  # ✅ Database connection
```

## 🔄 Next Steps

1. **Implement JWT authentication** in `middleware/auth.go`
2. **Add request validation** to all handlers
3. **Configure CORS** for your frontend domain
4. **Test all endpoints** with real auth tokens
5. **Add pagination** to list endpoint
6. **Create auth endpoints** for login/register
7. **Add admin routes** for cross-user access

## 💡 Tips

- **Start with auth** - Everything depends on having a real user ID
- **Test incrementally** - Test each endpoint after implementing validation
- **Use a REST client** - Postman, Insomnia, or Thunder Client in VS Code
- **Check logs** - The logger middleware shows all requests and responses
- **Database is ready** - Repository layer is complete and tested
