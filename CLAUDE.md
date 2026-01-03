# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a monorepo for a landscaping company forms management system with:
- **Backend**: Go REST API with PostgreSQL
- **Frontend**: Next.js 16 with React 19 and TypeScript

Employees create, view, and manage two types of forms: **shrub forms** and **pesticide forms**. Admins can view all forms across all users.

## Development Commands

### Backend (Go)

**Run server:**
```bash
cd backend
go run ./cmd/api
```

**Run all tests:**
```bash
cd backend
go test ./...
```

**Run specific package tests:**
```bash
cd backend
go test ./internal/forms -v
```

**Run single test:**
```bash
cd backend
go test ./internal/forms -run TestCreateAndGetShrubForm -v
```

### Frontend (Next.js)

**Install dependencies:**
```bash
cd frontend
npm install
```

**Run dev server:**
```bash
cd frontend
npm run dev
```

**Build:**
```bash
cd frontend
npm run build
```

**Lint:**
```bash
cd frontend
npm run lint
```

### Database

**Start PostgreSQL (Docker):**
```bash
cd backend/db
docker-compose up -d
```

**Apply schema:**
```bash
cd backend/db
./setup_db_from_schema.sh
```

The database runs on `localhost:5433` (not default 5432). Test database is `landscaping_test`.

## Architecture

### Backend Architecture

**Layered Architecture:**
```
HTTP Layer (not yet implemented)
    ↓
Repository Layer (backend/internal/forms/forms.go)
    ↓
Database (PostgreSQL)
```

**Repository Pattern:**
- `FormsRepository` provides all data access methods
- All methods enforce **ownership at SQL level** (queries filter by `created_by`)
- Returns `sql.ErrNoRows` for not found OR unauthorized access (no leaking of existence)
- Uses transactions for multi-table operations (form + shrub/pesticide)

**Key Files:**
- `backend/internal/forms/forms.go` - Repository with CRUD methods
- `backend/internal/forms/models.go` - Domain models
- `backend/internal/forms/testdb.go` - Test database setup helper
- `backend/cmd/api/main.go` - Application entry point (minimal, only `/health` endpoint)

### Database Schema Design

**Polymorphic Forms Pattern:**
- Base `forms` table with `form_type` CHECK constraint ('shrub' | 'pesticide')
- Type-specific tables (`shrubs`, `pesticides`) with `form_id` as PK and FK to forms
- ON DELETE CASCADE ensures cleanup
- Database triggers enforce type safety (shrub forms can't have pesticide records)

**Key Constraints:**
- UUID primary keys (not auto-increment)
- `updated_at` automatically updated by trigger
- Type enforcement triggers prevent form type mismatches
- Indices optimized for user queries and name searches

**Authorization Model:**
- All queries include `WHERE created_by = $userID` in repository layer
- No user can see/modify forms they don't own
- Returns same error for not-found vs unauthorized to prevent enumeration

### Domain Models

**Core Types:**
- `Form` - Base form with common fields
- `ShrubForm` - Embeds `Form` + `ShrubDetails`
- `PesticideForm` - Embeds `Form` + `PesticideDetails`
- `FormView` - Polymorphic wrapper with `Shrub *ShrubForm` | `Pesticide *PesticideForm`

**Input Types:**
- `CreateFormInput` - Common fields for creation
- `UpdateFormInput` - Fields that can be updated
- Type-specific details passed separately to create/update methods

### Frontend Architecture

**Status:** Currently just Next.js boilerplate (not implemented)

**Planned Structure:**
- Next.js 13+ App Router (not Pages Router)
- React Server Components
- Tailwind CSS 4 for styling
- TypeScript strict mode

## Implementation Status

### ✅ Completed

- PostgreSQL database schema with polymorphic forms
- Complete CRUD repository layer
- Transaction-safe operations
- Ownership enforcement at SQL layer
- Database triggers for type safety and timestamps
- Comprehensive test suite (19 tests)
- Test database setup with automatic schema loading

### ❌ Not Implemented

- HTTP API handlers (chi router imported but not used)
- Authentication/authorization system
- Request validation and error handling
- All frontend pages and components
- Search/filter API endpoints
- Admin functionality
- PDF import/export features

See `PLAN.md` for detailed implementation roadmap.

## Testing

**Test Database:**
- Uses separate `landscaping_test` database
- `testdb.go` helper resets schema before each test
- `createTestUser()` helper for test isolation

**Test Pattern:**
```go
func TestSomething(t *testing.T) {
    ctx := context.Background()
    db := testDB(t) // Sets up fresh schema
    repo := NewFormsRepository(db)
    userID := createTestUser(t, db)

    // Test code...
}
```

**Coverage:**
- All CRUD operations (happy paths)
- Authorization checks (wrong user access)
- Validation (nil details, both details provided)
- Edge cases (empty lists, non-existent forms)
- Sorting (first_name, last_name, created_at)
- Database constraints (cascade deletes)

## Key Patterns & Conventions

### Repository Methods

All repository methods follow this pattern:
1. Accept `context.Context` as first parameter
2. Accept `userID` to enforce ownership
3. Use transactions for multi-table operations
4. Return `sql.ErrNoRows` for not found OR unauthorized
5. Use parameterized queries ($1, $2) to prevent SQL injection

**Example:**
```go
func (r *FormsRepository) GetFormById(ctx context.Context, formID string, userID string) (*FormView, error)
```

### FormView Access

`FormView` is polymorphic. Access fields based on type:

```go
if formView.Shrub != nil {
    firstName := formView.Shrub.Form.FirstName
    numShrubs := formView.Shrub.ShrubDetails.NumShrubs
}

if formView.Pesticide != nil {
    firstName := formView.Pesticide.Form.FirstName
    pesticideName := formView.Pesticide.PesticideDetails.PesticideName
}
```

### Sorting

`ListFormsByUserId` supports:
- `sortBy`: "first_name", "last_name", "created_at" (defaults to "created_at")
- `order`: "ASC" or "DESC" (defaults to "DESC")

Invalid values fall back to defaults (no error).

## Environment Variables

**Backend (.env):**
```
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=landscaping
```

**Testing (.env.testing):**
```
DB_NAME=landscaping_test
# Other DB vars same as .env
```

**Frontend (.env.local - not yet created):**
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

## Database Connection

Connection pooling configured in `internal/db/db.go`:
- Max open connections: 10
- Max idle connections: 5
- Connection max lifetime: 1 hour

## Important Notes

1. **No HTTP layer yet** - Repository is complete but not exposed via API
2. **Chi router imported** - `go.mod` has chi v5 but it's not used in `main.go`
3. **UUID everywhere** - All IDs are UUIDs, not integers
4. **Form type is immutable** - Cannot change form type after creation
5. **Test isolation** - Each test gets fresh schema via `testDB(t)`
6. **Cascade deletes** - Deleting a form auto-deletes shrub/pesticide records
7. **Timestamp triggers** - `updated_at` automatically updated by database

## Future Work

See `PLAN.md` for the complete implementation plan. Key phases:
1. Backend API Layer (HTTP handlers, chi router)
2. Authentication & Authorization (JWT, user registration/login)
3. Search & Filter (query params, pagination)
4. Frontend (Auth pages, form management UI)
5. Admin features (view all forms)
6. PDF import/export (optional)
