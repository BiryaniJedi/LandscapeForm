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

**❌ Not Implemented:**
- HTTP API handlers (only `/health` exists)
- Authentication/authorization system
- Frontend UI (just Next.js boilerplate)
- Search/filter endpoints
- Admin functionality
- PDF upload/export features

---

## Implementation Phases

### Phase 1: Backend API Layer
**Goal:** Expose the repository layer as REST endpoints

**What to Build:**
- HTTP handlers connecting repository methods to routes
- Request/response DTOs (data transfer objects)
- Input validation for all endpoints
- Proper error response formatting
- Middleware for logging, CORS, and panic recovery
- Wire up chi router (already in dependencies)

**Key Files to Create:**
- `backend/internal/handlers/handlers.go` - Handler struct with dependencies
- `backend/internal/handlers/forms.go` - Form CRUD endpoints
- `backend/internal/middleware/` - Logging, CORS, recovery middleware
- Update `backend/cmd/api/main.go` - Wire routes to handlers

**Endpoints to Implement:**
```
POST   /api/forms/shrub          Create shrub form
POST   /api/forms/pesticide      Create pesticide form
GET    /api/forms                List user's forms (with sort params)
GET    /api/forms/{id}           Get single form
PUT    /api/forms/{id}           Update form
DELETE /api/forms/{id}           Delete form
```

**Key Considerations:**
- Extract user ID from context (set by auth middleware later)
- Return proper HTTP status codes (200, 201, 400, 404, 500)
- Validate all inputs before calling repository
- Handle `sql.ErrNoRows` → 404 responses
- Use JSON for all requests/responses

---

### Phase 2: Authentication & Authorization
**Goal:** Secure the API and identify users

**What to Build:**
- User registration endpoint (email + password)
- User login endpoint (returns JWT token)
- Password hashing with bcrypt
- JWT token generation and validation
- Auth middleware to extract current user from token
- Protected route wrapper

**Key Files to Create:**
- `backend/internal/auth/auth.go` - Password hashing, JWT helpers
- `backend/internal/auth/middleware.go` - JWT validation middleware
- `backend/internal/users/repository.go` - User CRUD operations
- `backend/internal/handlers/auth.go` - Registration/login handlers

**Database Changes:**
- Add `password_hash` column to users table
- Add `is_admin` boolean column to users table

**Endpoints to Implement:**
```
POST   /api/auth/register        Create new user account
POST   /api/auth/login           Login, return JWT token
GET    /api/auth/me              Get current user info (protected)
```

**Key Considerations:**
- Use bcrypt for password hashing (cost factor 12-14)
- JWT secret from environment variable
- Token expiration (24 hours recommended)
- Store user ID and is_admin flag in JWT claims
- Middleware extracts token from Authorization header
- All `/api/forms/*` routes should require authentication

**Dependencies to Add:**
```bash
go get golang.org/x/crypto/bcrypt
go get github.com/golang-jwt/jwt/v5
```

---

### Phase 3: Search, Filter & Pagination
**Goal:** Enhance form listing with query capabilities

**What to Build:**
- Query parameter parsing for search/filter/sort
- Enhanced repository method accepting filter params
- Pagination support (page number + page size)
- Total count for pagination metadata

**Repository Changes:**
- Update `ListFormsByUserId` to accept filter parameters:
  - `search` - Search in first_name/last_name (case-insensitive)
  - `form_type` - Filter by shrub/pesticide
  - `sort_by` - Column to sort (first_name, last_name, created_at)
  - `order` - ASC or DESC
  - `page` - Page number (1-indexed)
  - `page_size` - Results per page (default 20)

**API Changes:**
- `GET /api/forms` accepts query params:
  ```
  /api/forms?search=smith&form_type=shrub&sort_by=last_name&order=DESC&page=1&page_size=20
  ```
- Return pagination metadata in response:
  ```json
  {
    "forms": [...],
    "total": 150,
    "page": 1,
    "page_size": 20,
    "total_pages": 8
  }
  ```

**Key Considerations:**
- Use database indexes already created for efficient queries
- Validate sort_by column is allowed (prevent SQL injection)
- Default values for missing params
- Handle empty results gracefully

---

### Phase 4: Admin Functionality
**Goal:** Allow admins to view all forms across all users

**What to Build:**
- New repository method to list all forms (no user filter)
- Admin-only endpoint
- Middleware to check admin role
- Include user/creator info in admin responses

**Key Files to Modify:**
- `backend/internal/forms/forms.go` - Add `ListAllForms()` method
- `backend/internal/handlers/forms.go` - Add admin handler
- `backend/internal/middleware/admin.go` - Admin check middleware

**Endpoints to Implement:**
```
GET    /api/admin/forms          List all forms (admin only)
```

**Key Considerations:**
- Check `is_admin` flag from JWT claims
- Return 403 Forbidden if not admin
- Include created_by user email in response for context
- Support same search/filter/pagination as user endpoint

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

1. **Phase 1** (Backend API) - Get basic endpoints working
2. **Phase 2** (Auth) - Secure the API
3. **Phase 5** (Frontend Auth) - Login/register UI
4. **Phase 6** (Frontend Forms) - Main functionality
5. **Phase 3** (Search/Filter) - Enhanced querying
6. **Phase 4** (Admin) - Admin features
7. **Phase 7** (Admin UI) - Admin dashboard
8. **Testing** - Add test coverage
9. **Phase 8** (PDF) - Optional advanced feature
10. **Deployment** - Ship to production

---

## Estimated Complexity

| Phase | Complexity | Estimated Effort |
|-------|------------|------------------|
| 1. Backend API | Medium | Core functionality |
| 2. Authentication | Medium-High | Critical for security |
| 3. Search/Filter | Low | Leverages existing DB indexes |
| 4. Admin API | Low | Extension of existing patterns |
| 5. Frontend Auth | Medium | Standard auth flow |
| 6. Frontend Forms | Medium-High | Main UI bulk |
| 7. Admin UI | Low | Reuse dashboard components |
| 8. PDF Features | High | Complex, optional |

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
