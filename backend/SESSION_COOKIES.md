# Session Cookie Authentication Guide

## Current Implementation vs Session Cookies

### What We Have Now
Currently, the backend returns the JWT token in the **response body**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { ... }
}
```

The frontend would need to:
1. Store the token in localStorage/sessionStorage
2. Manually add it to every request header: `Authorization: Bearer <token>`

### What Session Cookies Provide
With session cookies, the backend sends the JWT token as an **HttpOnly cookie**:
- Browser automatically stores it securely
- Browser automatically sends it with every request to the same domain
- JavaScript **cannot** access it (prevents XSS attacks)
- Cookie can be set to expire automatically

## Backend Changes Needed

### 1. Update Login Handler to Set Cookie

**File: `internal/handlers/auth.go`**

Add this after generating the token:

```go
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    // ... existing login logic ...

    // Generate JWT token
    token, err := auth.GenerateToken(user.ID, user.Role)
    if err != nil {
        respondError(w, http.StatusInternalServerError, "Failed to generate token")
        return
    }

    // Set token as HttpOnly cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    token,
        Path:     "/",
        HttpOnly: true,  // Prevents JavaScript access (XSS protection)
        Secure:   true,  // Only send over HTTPS (set to false for localhost dev)
        SameSite: http.SameSiteStrictMode, // CSRF protection
        MaxAge:   86400, // 24 hours (matches JWT expiration)
    })

    // OPTIONAL: Still return token in response for mobile apps
    respondJSON(w, http.StatusOK, LoginResponse{
        Token: token, // Can omit this if only using cookies
        User:  userResponse,
    })
}
```

### 2. Update Auth Middleware to Check Cookie

**File: `internal/middleware/auth.go`**

Modify `AuthMiddleware` to check cookies first, then fall back to Authorization header:

```go
func AuthMiddleware(usersRepo *users.UsersRepository) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            var token string

            // Try to get token from cookie first
            cookie, err := r.Cookie("auth_token")
            if err == nil {
                token = cookie.Value
            } else {
                // Fall back to Authorization header
                authHeader := r.Header.Get("Authorization")
                if authHeader == "" {
                    w.Header().Set("Content-Type", "application/json")
                    http.Error(w, `{"error":"Unauthorized","message":"Missing authorization"}`, http.StatusUnauthorized)
                    return
                }

                parts := strings.Split(authHeader, " ")
                if len(parts) != 2 || parts[0] != "Bearer" {
                    w.Header().Set("Content-Type", "application/json")
                    http.Error(w, `{"error":"Unauthorized","message":"Invalid authorization header"}`, http.StatusUnauthorized)
                    return
                }
                token = parts[1]
            }

            // ... rest of existing validation ...
        })
    }
}
```

### 3. Add Logout Endpoint

**File: `internal/handlers/auth.go`**

```go
// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    // Clear the auth cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        MaxAge:   -1, // Immediately expire
    })

    respondJSON(w, http.StatusOK, map[string]string{
        "message": "Logged out successfully",
    })
}
```

### 4. Update Router in main.go

```go
// Add logout route
r.Post("/api/auth/logout", authHandler.Logout)
```

### 5. CORS Configuration for Cookies

**File: `internal/middleware/cors.go`**

Update CORS to allow credentials:

```go
func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Allow credentials (cookies)
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        // Set allowed origin (MUST be specific when using credentials, not *)
        origin := r.Header.Get("Origin")
        if origin == "http://localhost:3000" || origin == "https://yourdomain.com" {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }

        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

## Frontend Usage (Next.js)

### 1. Login Request

```typescript
// app/login/page.tsx
async function handleLogin(username: string, password: string) {
  const response = await fetch('http://localhost:8080/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include', // IMPORTANT: Send/receive cookies
    body: JSON.stringify({ username, password })
  });

  const data = await response.json();
  // Cookie is automatically stored by browser
  // Redirect to dashboard
  router.push('/dashboard');
}
```

### 2. Protected API Requests

```typescript
// All subsequent requests automatically include cookie
async function fetchForms() {
  const response = await fetch('http://localhost:8080/api/forms', {
    credentials: 'include', // IMPORTANT: Include cookies
  });

  return response.json();
}
```

### 3. Logout Request

```typescript
async function handleLogout() {
  await fetch('http://localhost:8080/api/auth/logout', {
    method: 'POST',
    credentials: 'include',
  });

  // Cookie is automatically cleared
  router.push('/login');
}
```

## Cookie vs LocalStorage Comparison

| Feature | HttpOnly Cookie | LocalStorage |
|---------|----------------|--------------|
| **XSS Protection** | ✅ Yes - JS can't access | ❌ No - JS can read |
| **CSRF Protection** | ⚠️ Need SameSite | ✅ Not vulnerable |
| **Auto-sent** | ✅ Yes | ❌ Manual |
| **Works with subdomains** | ✅ Yes | ❌ No |
| **Mobile app support** | ❌ Limited | ✅ Yes |
| **Browser compatibility** | ✅ Universal | ✅ Universal |

## Security Best Practices

1. **Use HTTPS in production** - Set `Secure: true` on cookies
2. **Set SameSite=Strict** - Prevents CSRF attacks
3. **HttpOnly=true** - Prevents XSS attacks
4. **Short expiration** - 24 hours or less
5. **Refresh token pattern** - For longer sessions (advanced)

## Development vs Production

### Development (localhost)
```go
Secure:   false,  // Allow over HTTP
SameSite: http.SameSiteLaxMode,
```

### Production
```go
Secure:   true,   // Require HTTPS
SameSite: http.SameSiteStrictMode,
```

## Testing Cookies

```bash
# Login and save cookies
curl -c cookies.txt -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Use saved cookies
curl -b cookies.txt http://localhost:8080/api/users
```

## Summary

**Do you need to implement cookies?**

It depends on your use case:

- **Web app only** → Session cookies (recommended)
- **Mobile app + Web** → Keep current approach (Authorization header)
- **Both needs** → Support both (cookie + header)

The current implementation already works! Session cookies are an enhancement for better security and UX in web browsers.
