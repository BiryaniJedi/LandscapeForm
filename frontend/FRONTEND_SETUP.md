# Frontend Setup Guide

## Authentication Pages - Setup Complete! ✅

The boilerplate login and register pages have been created with full integration to your Go backend.

## What's Been Created

### 1. **API Client** (`lib/api.ts`)
- Singleton API client for backend communication
- Automatic JWT token management (stored in localStorage)
- Auto-redirect to login on 401 responses
- Type-safe request/response handling

### 2. **Auth Context** (`lib/auth.tsx`)
- React Context for app-wide authentication state
- `useAuth()` hook for accessing auth in any component
- Automatic token persistence across page reloads

### 3. **Login Page** (`app/login/page.tsx`)
- Clean, modern UI with Tailwind CSS
- Form validation
- Error handling and display
- Link to registration page
- Fully functional with backend `/api/auth/login`

### 4. **Register Page** (`app/register/page.tsx`)
- Multi-field registration form
- Password confirmation validation
- Date of birth picker
- Shows "pending approval" message after registration
- Fully functional with backend `/api/auth/register`

### 5. **Dashboard Page** (`app/dashboard/page.tsx`)
- Protected route (redirects to login if not authenticated)
- Displays user info and status
- Shows "pending approval" warning if applicable
- Admin badge for admin users
- Logout functionality
- Placeholder for forms management

### 6. **Home Page** (`app/page.tsx`)
- Landing page with login/register links
- Auto-redirects to dashboard if already logged in

## Environment Setup

The `.env.local` file has been created with:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

## Running the Frontend

1. **Install dependencies** (if not already done):
   ```bash
   cd frontend
   npm install
   ```

2. **Make sure your backend is running**:
   ```bash
   cd backend
   go run ./cmd/api
   ```

3. **Start the Next.js dev server**:
   ```bash
   cd frontend
   npm run dev
   ```

4. **Open your browser** to `http://localhost:3000`

## User Flow

1. **Landing Page** (`/`) - Shows login/register buttons
2. **Register** (`/register`) - Create a new account
   - New users are created with `pending: true`
   - They get a JWT token but see a "pending approval" message
3. **Login** (`/login`) - Sign in to existing account
4. **Dashboard** (`/dashboard`) - Protected page showing user info
   - Pending users see a warning message
   - Approved users can access all features (when implemented)

## API Integration

All pages are fully integrated with your backend:

- **POST /api/auth/register** - Creates user, returns JWT + user data
- **POST /api/auth/login** - Validates credentials, returns JWT + user data
- JWT tokens are automatically included in subsequent requests
- 401 responses trigger auto-logout and redirect to login

## Features Included

✅ JWT token management (localStorage)
✅ Auto-redirect on auth state changes
✅ Protected routes
✅ Error handling and display
✅ Loading states
✅ Dark mode support
✅ Responsive design
✅ Form validation
✅ Password confirmation
✅ Pending approval workflow support
✅ Admin role detection

## What's Next?

The authentication flow is complete! Next steps from PLAN.md:

1. **Phase 6: Frontend Form Management**
   - Forms list page
   - Create form page (shrub/pesticide)
   - Edit form page
   - Delete functionality

2. **Phase 7: Admin Dashboard**
   - View all users
   - Approve pending users
   - View all forms (across users)

## Tech Stack

- **Next.js 16** - React framework with App Router
- **React 19** - UI library
- **TypeScript** - Type safety
- **Tailwind CSS 4** - Styling
- **fetch API** - HTTP requests (built-in)

## Notes

- The AuthProvider wraps the entire app in `app/layout.tsx`
- All auth state is managed through React Context
- Token is persisted in localStorage (browser only)
- Server-side rendering is disabled for auth pages (client components)
