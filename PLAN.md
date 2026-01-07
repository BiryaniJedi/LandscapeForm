# LandscapeForm - Implementation Plan

## Current Status Summary

**✅ Completed:**
- PostgreSQL database schema (users, forms, shrubs, pesticides)
- Complete CRUD repository layer with ownership enforcement
- Domain models (Form, ShrubForm, PesticideForm, FormView)
- Database connection with pooling
- Test infrastructure and helpers
- Transaction-safe operations
- Efficient indexing for queries
- **HTTP API handlers - ALL endpoints implemented**
- **Authentication/authorization system - COMPLETE**
- **JWT token generation and validation**
- **All middleware (Auth, AdminOnly, RequireApproved, CORS, Logger, Recovery)**
- **Search/filter/pagination endpoints**
- **Admin functionality (ListAllForms, user management)**
- **Comprehensive API test script**

**❌ Not Implemented:**
- Input validation in handlers (TODOs in code)
- Better error handling (duplicate username, etc.)
- Backend handler/middleware/integration tests
- Frontend UI (just Next.js boilerplate)
- PDF upload/export features

---

## 🎯 What's Next?

### Immediate Backend Tasks (Optional Polish):
1. **Add Input Validation** - All handlers have TODO comments for:
   - Required field validation
   - Phone number format validation
   - Password strength requirements
   - Date of birth validation
   - Duplicate username error handling

2. **Add Backend Tests**:
   - Handler tests using httptest
   - Middleware tests
   - Integration tests

3. **Error Handling Enhancements**:
   - Better error messages for duplicate usernames
   - More specific error responses

### Main Development Path (Frontend):
**The backend is fully functional!** The next major phase is to build the frontend:
- Phase 5: Frontend Authentication (login/register pages)
- Phase 6: Frontend Form Management (dashboard, CRUD UI)
- Phase 7: Admin Dashboard

---

## Implementation Phases

### Phase 1: Backend API Layer ✅ **COMPLETE**
**Goal:** Expose the repository layer as REST endpoints

**What Was Built:**
- ✅ HTTP handlers connecting repository methods to routes
- ✅ Request/response DTOs (data transfer objects)
- ⚠️ Input validation for all endpoints (TODOs in code - still needed)
- ✅ Proper error response formatting
- ✅ Middleware for logging, CORS, and panic recovery
- ✅ Chi router wired up

**Key Files Created:**
- ✅ `backend/internal/handlers/types.go` - Request/response types
- ✅ `backend/internal/handlers/responses.go` - Response helpers
- ✅ `backend/internal/handlers/forms.go` - Form CRUD endpoints
- ✅ `backend/internal/middleware/logger.go` - Request logging
- ✅ `backend/internal/middleware/cors.go` - CORS headers
- ✅ `backend/internal/middleware/recovery.go` - Panic recovery
- ✅ `backend/cmd/api/main.go` - Full router setup

**Endpoints Implemented:**
```
✅ POST   /api/forms/shrub          Create shrub form
✅ POST   /api/forms/pesticide      Create pesticide form
✅ GET    /api/forms                List user's forms (with sort/filter params)
✅ GET    /api/forms/{id}           Get single form
✅ PUT    /api/forms/{id}           Update form
✅ DELETE /api/forms/{id}           Delete form
```

**Implementation Notes:**
- ✅ Extracts user ID from context (set by auth middleware)
- ✅ Returns proper HTTP status codes (200, 201, 400, 404, 500)
- ⚠️ Input validation TODOs remain in code
- ✅ Handles `sql.ErrNoRows` → 404 responses
- ✅ Uses JSON for all requests/responses

---

### Phase 2: Authentication & Authorization ✅ **COMPLETE**
**Goal:** Secure the API and identify users

**What Was Built:**
- ✅ User registration endpoint (username + password)
- ✅ User login endpoint (returns JWT token + user data)
- ✅ Password hashing with bcrypt
- ✅ JWT token generation and validation
- ✅ Auth middleware to extract current user from token AND load from DB
- ✅ Protected route wrapper (RequireApproved middleware)
- ✅ Admin-only route wrapper (AdminOnly middleware)
- ✅ User approval system (pending flag)

**Key Files Created:**
- ✅ `backend/internal/auth/jwt.go` - JWT generation and validation
- ✅ `backend/internal/middleware/auth.go` - Auth middleware (JWT validation + DB user lookup)
- ✅ `backend/internal/users/users.go` - User CRUD operations repository
- ✅ `backend/internal/users/models.go` - User domain models
- ✅ `backend/internal/handlers/auth.go` - Registration/login handlers
- ✅ `backend/internal/handlers/users.go` - User management handlers

**Database Schema:**
- ✅ `password_hash` column in users table
- ✅ `role` column with 'employee' or 'admin' values (not just boolean)
- ✅ `pending` boolean for approval workflow

**Endpoints Implemented:**
```
✅ POST   /api/auth/register        Create new user account (returns token + user)
✅ POST   /api/auth/login           Login, return JWT token + user data
✅ GET    /api/users/{id}           Get user by ID (protected)
✅ PUT    /api/users/{id}           Update user (protected)
✅ GET    /api/users                List all users (admin only)
✅ DELETE /api/users/{id}           Delete user (admin only)
✅ POST   /api/users/{id}/approve   Approve pending user (admin only)
```

**Implementation Details:**
- ✅ Uses bcrypt for password hashing (default cost)
- ✅ JWT secret from environment variable (with fallback for dev)
- ✅ Token expiration (24 hours)
- ✅ Stores user ID and role in JWT claims
- ✅ Middleware extracts token from Authorization header
- ✅ All `/api/forms/*` routes require authentication + approved account
- ✅ Auth middleware loads current user from DB (doesn't trust JWT role)

**Dependencies Added:**
- ✅ golang.org/x/crypto/bcrypt
- ✅ github.com/golang-jwt/jwt/v5

---

### Phase 3: Search, Filter & Pagination ✅ **COMPLETE**
**Goal:** Enhance form listing with query capabilities

**What Was Built:**
- ✅ Query parameter parsing for search/filter/sort
- ✅ Enhanced repository method accepting filter params
- ✅ Pagination support (limit + offset OR page number)
- ⚠️ Total count for pagination metadata (not implemented - could be added)

**Repository Implementation:**
- ✅ `ListFormsByUserId` accepts filter parameters:
  - ✅ `search` - Search in first_name/last_name (case-insensitive with ILIKE)
  - ✅ `form_type` - Filter by shrub/pesticide
  - ✅ `sort_by` - Column to sort (first_name, last_name, created_at)
  - ✅ `order` - ASC or DESC
  - ✅ `limit` - Number of results to return
  - ✅ `offset` - Number of results to skip
  - ✅ `page` - Page number (converted to offset internally)

**API Implementation:**
- ✅ `GET /api/forms` accepts query params:
  ```
  /api/forms?search=smith&type=shrub&sort_by=last_name&order=DESC&limit=20&offset=0
  /api/forms?search=smith&type=shrub&sort_by=last_name&order=DESC&page=1&limit=20
  ```
- Current response format:
  ```json
  {
    "forms": [...],
    "count": 15
  }
  ```
- ⚠️ Note: Does not return total count or total_pages (could be enhanced)

**Implementation Details:**
- ✅ Uses database indexes for efficient queries
- ✅ Validates sort_by column against allowed list (prevents SQL injection)
- ✅ Default values for missing params
- ✅ Handles empty results gracefully

---

### Phase 4: Admin Functionality ✅ **COMPLETE**
**Goal:** Allow admins to view all forms across all users

**What Was Built:**
- ✅ New repository method to list all forms (no user filter)
- ✅ Admin-only endpoint
- ✅ Middleware to check admin role
- ✅ Admin user management endpoints
- ✅ User approval workflow

**Key Files Modified:**
- ✅ `backend/internal/forms/forms.go` - Added `ListAllForms()` method
- ✅ `backend/internal/handlers/forms.go` - Added admin handler
- ✅ `backend/internal/middleware/auth.go` - Added `AdminOnly` middleware
- ✅ `backend/internal/handlers/users.go` - User management handlers (all admin-only)

**Endpoints Implemented:**
```
✅ GET    /api/admin/forms              List all forms (admin only)
✅ GET    /api/users                    List all users (admin only)
✅ DELETE /api/users/{id}               Delete user (admin only)
✅ POST   /api/users/{id}/approve       Approve pending user (admin only)
```

**Implementation Details:**
- ✅ Checks `role` from user loaded from DB (not from JWT - more secure)
- ✅ Returns 403 Forbidden if not admin
- ✅ Forms include created_by UUID in response
- ✅ Supports same search/filter/pagination as user endpoint
- ✅ User approval workflow (new users are pending until approved)

---

### Phase 5: Frontend Authentication
**Goal:** Build login/register pages and auth state management

**What to Build:**
- Login page (`/login`)
- Registration page (`/register`)
- Auth context provider for app-wide state
- API client utilities
- Protected route wrapper component
- Token storage in localStorage
- Auto-login on page load if token exists

**Key Files to Create:**
- `frontend/app/login/page.tsx` - Login form
- `frontend/app/register/page.tsx` - Registration form
- `frontend/lib/auth.tsx` - Auth context and hooks
- `frontend/lib/api.ts` - API client with token handling
- `frontend/components/ProtectedRoute.tsx` - Auth guard

**Key Features:**
- Form validation on client side
- Error message display
- Loading states during API calls
- Redirect to dashboard after successful login
- Logout functionality (clear token)
- Auto-redirect to login if unauthorized (401 response)

**API Client Pattern:**
```typescript
// Automatically include token in all requests
// Handle 401 responses globally (redirect to login)
// Provide typed methods for all endpoints
```

---

### Phase 6: Frontend Form Management
**Goal:** Build the main UI for viewing, creating, editing, and deleting forms

**What to Build:**

**6.1 Dashboard (Forms List)**
- Main page showing user's forms in table/card layout
- Search bar with filter controls
- Sort controls (name, date)
- Form type filter (all/shrub/pesticide)
- Pagination controls
- Delete button with confirmation modal

**6.2 Form Creation**
- Page with form type selector (shrub vs pesticide)
- Conditional form fields based on type
- Client info fields (first name, last name, phone)
- Type-specific fields (num_shrubs or pesticide_name)
- Form validation
- Submit to API and redirect to dashboard

**6.3 Form View/Edit**
- Form detail page showing all fields
- Edit mode toggle
- Update functionality
- Delete button
- Back to dashboard link

**6.4 Navigation & Layout**
- Navbar with links (Dashboard, New Form, Logout)
- Show current user email
- Highlight admin users
- Responsive design with Tailwind

**Key Files to Create:**
- `frontend/app/dashboard/page.tsx` - Forms list
- `frontend/app/forms/new/page.tsx` - Create form
- `frontend/app/forms/[id]/page.tsx` - View form
- `frontend/app/forms/[id]/edit/page.tsx` - Edit form
- `frontend/components/Navbar.tsx` - Navigation bar
- `frontend/components/FormList.tsx` - Table/grid of forms
- `frontend/components/FormCard.tsx` - Individual form display
- `frontend/components/SearchBar.tsx` - Search/filter UI
- `frontend/components/ShrubFormFields.tsx` - Shrub-specific inputs
- `frontend/components/PesticideFormFields.tsx` - Pesticide inputs

**Key Features:**
- Debounced search input
- Optimistic UI updates for delete
- Loading spinners during API calls
- Error toast notifications
- Empty state when no forms exist
- Responsive table/card toggle

---

### Phase 7: Admin Dashboard
**Goal:** Admin view to see all users' forms

**What to Build:**
- Admin dashboard page (accessible only to admins)
- Show all forms with user email column
- Same search/filter/pagination as user dashboard
- Visual distinction (different color scheme)
- Admin badge in navbar

**Key Files to Create:**
- `frontend/app/admin/page.tsx` - Admin forms list
- `frontend/lib/api.ts` - Add admin API methods

**Key Features:**
- Check `is_admin` from auth context
- Redirect non-admins to regular dashboard
- Show "Created By" column with user email
- Filter by user (optional enhancement)

---

### Phase 8: PDF Upload & Export (Optional)
**Goal:** Allow PDF import/export of forms

**What to Build:**

**8.1 PDF Upload**
- File upload component (drag & drop)
- Multipart form data handling
- PDF parsing to extract form data
- Preview before saving

**8.2 PDF Export**
- Generate PDF from form data
- Download endpoint
- PDF template design

**Key Files to Create:**
- `backend/internal/pdf/parser.go` - PDF text extraction
- `backend/internal/pdf/generator.go` - PDF creation
- `frontend/components/FileUpload.tsx` - Upload UI

**Dependencies Needed:**
- Backend: PDF library (e.g., `github.com/jung-kurt/gofpdf` or `github.com/unidoc/unipdf`)
- Frontend: File upload (built-in HTML5)

**Endpoints to Implement:**
```
POST   /api/forms/import/pdf     Upload and parse PDF
GET    /api/forms/{id}/pdf       Download form as PDF
```

**Key Considerations:**
- File size limits (5-10MB)
- File type validation (only PDF)
- Secure file storage
- OCR may be needed for scanned PDFs (complex)
- Consider cloud services (AWS Textract, Google Document AI)

---

## Testing Strategy

### Backend Testing
**What to Test:**
- Handler tests (using httptest)
- Auth flow tests (register, login, protected routes)
- Repository tests (already started)
- Middleware tests
- Integration tests (full request → database → response)

**Files to Create:**
- `backend/internal/handlers/handlers_test.go`
- `backend/internal/auth/auth_test.go`
- `backend/internal/middleware/middleware_test.go`

### Frontend Testing
**What to Test:**
- Component unit tests (React Testing Library)
- Auth flow tests
- Form submission tests
- API integration tests (mocked)

**Setup Required:**
- Jest + React Testing Library
- MSW (Mock Service Worker) for API mocking

---

## Deployment

### Backend Deployment
**Steps:**
1. Create Dockerfile for Go application
2. Build production binary
3. Set environment variables (DB credentials, JWT secret)
4. Deploy to cloud provider:
   - Railway (recommended, has PostgreSQL add-on)
   - Fly.io
   - Render
   - AWS ECS/Fargate

### Frontend Deployment
**Steps:**
1. Set `NEXT_PUBLIC_API_URL` environment variable
2. Deploy to Vercel (easiest for Next.js):
   ```bash
   vercel --prod
   ```
3. Alternative: Netlify, Cloudflare Pages

### Database
- Use hosted PostgreSQL (Railway, Supabase, AWS RDS)
- Run migrations/schema on production DB
- Create admin user manually via SQL

---

## Suggested Implementation Order

1. ✅ **Phase 1** (Backend API) - Get basic endpoints working - **DONE**
2. ✅ **Phase 2** (Auth) - Secure the API - **DONE**
3. ✅ **Phase 3** (Search/Filter) - Enhanced querying - **DONE**
4. ✅ **Phase 4** (Admin) - Admin features - **DONE**
5. **→ Phase 5** (Frontend Auth) - Login/register UI - **NEXT STEP**
6. **Phase 6** (Frontend Forms) - Main functionality
7. **Phase 7** (Admin UI) - Admin dashboard
8. **Testing** - Add test coverage (backend + frontend)
9. **Phase 8** (PDF) - Optional advanced feature
10. **Deployment** - Ship to production

---

## Estimated Complexity

| Phase | Complexity | Status |
|-------|------------|--------|
| 1. Backend API | Medium | ✅ **COMPLETE** |
| 2. Authentication | Medium-High | ✅ **COMPLETE** |
| 3. Search/Filter | Low | ✅ **COMPLETE** |
| 4. Admin API | Low | ✅ **COMPLETE** |
| 5. Frontend Auth | Medium | ❌ **TODO** |
| 6. Frontend Forms | Medium-High | ❌ **TODO** |
| 7. Admin UI | Low | ❌ **TODO** |
| 8. PDF Features | High | ❌ **TODO** (Optional) |

---

## Quick Reference: Architecture

```
┌─────────────────────────────────────────────────┐
│                   Frontend                       │
│  (Next.js, React, TypeScript, Tailwind)         │
│                                                  │
│  Pages: Login, Register, Dashboard, Forms       │
│  Components: FormList, SearchBar, Navbar        │
│  State: Auth Context, React Hooks              │
└─────────────────┬───────────────────────────────┘
                  │
                  │ HTTP/JSON (REST API)
                  │ JWT Token in Authorization header
                  │
┌─────────────────▼───────────────────────────────┐
│              Backend API (Go)                    │
│                                                  │
│  ┌──────────────────────────────────────────┐  │
│  │  Middleware                               │  │
│  │  - CORS, Logging, Recovery, JWT Auth     │  │
│  └──────────────┬───────────────────────────┘  │
│                 │                                │
│  ┌──────────────▼───────────────────────────┐  │
│  │  HTTP Handlers                            │  │
│  │  - Auth (register, login)                 │  │
│  │  - Forms (CRUD, list, search)            │  │
│  │  - Admin (list all)                       │  │
│  └──────────────┬───────────────────────────┘  │
│                 │                                │
│  ┌──────────────▼───────────────────────────┐  │
│  │  Repository Layer (forms, users)          │  │
│  │  - Ownership enforcement                  │  │
│  │  - Transaction management                 │  │
│  │  - SQL query builders                     │  │
│  └──────────────┬───────────────────────────┘  │
└─────────────────┼───────────────────────────────┘
                  │
                  │ SQL Queries
                  │
┌─────────────────▼───────────────────────────────┐
│          PostgreSQL Database                     │
│                                                  │
│  Tables: users, forms, shrubs, pesticides       │
│  Triggers: type validation, updated_at          │
│  Constraints: FK, CHECK, UNIQUE                 │
└─────────────────────────────────────────────────┘
```

---

## Environment Variables Needed

### Backend `.env`
```bash
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=landscaping
JWT_SECRET=<random-secret-key>  # Generate with: openssl rand -hex 32
PORT=8080
```

### Frontend `.env.local`
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

---

## Key Design Decisions Made

1. **Polymorphic Forms**: Using form_type + separate tables (shrubs, pesticides) rather than JSON blob
2. **Repository Pattern**: Data access layer separated from HTTP handlers
3. **JWT Authentication**: Stateless tokens rather than sessions
4. **UUID Primary Keys**: Better for distributed systems, no auto-increment leakage
5. **Database Triggers**: Enforce business rules at DB level (type validation, timestamps)
6. **Ownership at SQL Layer**: All queries filter by created_by for security

---

## Notes

- The foundation (data layer) is production-ready
- All phases build incrementally on each other
- PDF features (Phase 8) are optional and complex - recommend doing last
- Consider adding email verification in future (not in current plan)
- File upload for PDF might need cloud storage (S3) for production
- Current schema is minimal - may need more fields based on real form requirements
