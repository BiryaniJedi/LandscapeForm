package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/auth"
	"github.com/BiryaniJedi/LandscapeForm-backend/internal/users"
)

// AuthMiddleware validates JWT tokens and loads current user data from database
// This ensures we always have the latest user role and pending status
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

				// Expected format: "Bearer <token>"
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					w.Header().Set("Content-Type", "application/json")
					http.Error(w, `{"error":"Unauthorized","message":"Invalid authorization header format"}`, http.StatusUnauthorized)
					return
				}

				token = parts[1]
			}

			// Validate JWT token to get user ID
			claims, err := auth.ValidateToken(token)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"error":"Unauthorized","message":"Invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			// Query database to get current user role and status
			user, err := usersRepo.GetUserById(r.Context(), claims.UserID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					w.Header().Set("Content-Type", "application/json")
					http.Error(w, `{"error":"Unauthorized","message":"User not found"}`, http.StatusUnauthorized)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"error":"Internal Server Error","message":"Failed to verify user"}`, http.StatusInternalServerError)
				return
			}

			// Debug Info:
			//fmt.Printf("From auth middleware:\n\t- userID: %s\n\t- userRole: %s\n\t- userPending: %v\n", user.ID, user.Role, user.Pending)

			// Add user info to request context
			ctx := context.WithValue(r.Context(), "userID", user.ID)
			ctx = context.WithValue(ctx, "userRole", user.Role)
			ctx = context.WithValue(ctx, "userPending", user.Pending)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireApproved middleware ensures only non-pending users can access endpoints
// Must be used AFTER AuthMiddleware
func RequireApproved(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract pending status from context (set by AuthMiddleware)
		userPending, ok := r.Context().Value("userPending").(bool)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error":"Forbidden","message":"User status not found"}`, http.StatusForbidden)
			return
		}

		// Check if user is pending approval
		if userPending {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error":"Forbidden","message":"Account pending admin approval"}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AdminOnly middleware ensures only admin users can access the endpoint
// Must be used AFTER AuthMiddleware
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user role from context (set by AuthMiddleware)
		userRole, ok := r.Context().Value("userRole").(string)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error":"Forbidden","message":"User role not found"}`, http.StatusForbidden)
			return
		}

		// Check if user is admin
		if userRole != "admin" {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error":"Forbidden","message":"Admin access required"}`, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserID extracts userID from context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}

// GetUserRole extracts user role from context
func GetUserRole(ctx context.Context) (string, bool) {
	userRole, ok := ctx.Value("userRole").(string)
	return userRole, ok
}

// GetUserPending extracts user pending status from context
func GetUserPending(ctx context.Context) (bool, bool) {
	userPending, ok := ctx.Value("userPending").(bool)
	return userPending, ok
}
