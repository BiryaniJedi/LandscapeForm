package middleware

import (
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware is a placeholder for JWT authentication
// TODO: Implement actual JWT validation
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// TODO: Validate JWT token here
		// - Parse the token
		// - Verify signature
		// - Check expiration
		// - Extract user ID from claims

		// PLACEHOLDER: For now, just check token is not empty
		if token == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// TODO: Replace this with actual user ID from JWT claims
		// For development, you can hardcode a test user ID
		userID := "placeholder-user-id"

		// Add userID to request context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Optional: Helper to extract userID from context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}
