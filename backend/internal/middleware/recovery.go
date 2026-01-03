package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery recovers from panics and returns a 500 Internal Server Error
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic and stack trace
				log.Printf("PANIC: %v\n%s", err, debug.Stack())

				// Return 500 to client
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
