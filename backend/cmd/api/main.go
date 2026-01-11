package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/db"
	"github.com/BiryaniJedi/LandscapeForm-backend/internal/forms"
	"github.com/BiryaniJedi/LandscapeForm-backend/internal/handlers"
	"github.com/BiryaniJedi/LandscapeForm-backend/internal/middleware"
	"github.com/BiryaniJedi/LandscapeForm-backend/internal/users"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Status:  "success",
		Message: "Server is running",
		Code:    http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func setupRouter(formsHandler *handlers.FormsHandler, usersHandler *handlers.UsersHandler, authHandler *handlers.AuthHandler, usersRepo *users.UsersRepository) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Recovery)     // Recover from panics
	r.Use(middleware.Logger)       // Log all requests
	r.Use(middleware.CORS)         // Enable CORS
	r.Use(chimiddleware.RequestID) // Add request ID to each request
	r.Use(chimiddleware.RealIP)    // Get real client IP

	// Public routes (no auth required)
	r.Get("/health", healthHandler)

	// Authentication routes (public)
	r.Post("/api/auth/login", authHandler.Login)       // POST /api/auth/login
	r.Post("/api/auth/register", authHandler.Register) // POST /api/auth/register
	r.Post("/api/auth/logout", authHandler.Logout)     // POST /api/auth/logout
	//r.Get("/api/auth/me", authHandler.Me)

	/*// User registration (public)
	r.Post("/api/users", usersHandler.CreateUser) // POST /api/users
	*/

	// Protected routes (require authentication and approved account)
	r.Route("/api", func(r chi.Router) {
		// Apply auth middleware - validates JWT and loads user from DB

		r.Use(middleware.AuthMiddleware(usersRepo))
		r.Get("/auth/me", authHandler.Me)

		// Forms endpoints (require authentication + approved account)
		r.Route("/forms", func(r chi.Router) {
			// Require user to be approved (not pending)
			r.Use(middleware.RequireApproved)
			r.Get("/", formsHandler.ListForms) // GET /api/forms
			r.Route("/shrub", func(r chi.Router) {
				r.Post("/", formsHandler.CreateShrubForm)    // POST /api/forms/shrub
				r.Put("/{id}", formsHandler.UpdateShrubForm) // PUT /api/forms/shrub/{id}
				r.Get("/{id}", formsHandler.GetFormView)     // GET /api/forms/shrub/{id}
			})
			r.Route("/pesticide", func(r chi.Router) {
				r.Post("/", formsHandler.CreatePesticideForm)    // POST /api/forms/pesticide
				r.Put("/{id}", formsHandler.UpdatePesticideForm) // PUT /api/forms/pesticide/{id}
				r.Get("/{id}", formsHandler.GetFormView)         // GET /api/forms/pesticide/{id}
			})

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", formsHandler.GetFormView)   // GET /api/forms/{id}
				r.Delete("/", formsHandler.DeleteForm) // DELETE /api/forms/{id}
			})
		})

		// User endpoints (require authentication + approved account)
		r.Route("/users", func(r chi.Router) {
			// User can view/update own profile
			r.Get("/{id}", usersHandler.GetUser)    // GET /api/users/{id}
			r.Put("/{id}", usersHandler.UpdateUser) // PUT /api/users/{id}

			// Admin-only routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.AdminOnly) // Require admin role

				r.Get("/", usersHandler.ListUsers)                // GET /api/users
				r.Delete("/{id}", usersHandler.DeleteUser)        // DELETE /api/users/{id}
				r.Post("/{id}/approve", usersHandler.ApproveUser) // POST /api/users/{id}/approve
			})
		})

		// Admin routes for forms
		r.Route("/admin/forms", func(r chi.Router) {
			r.Use(middleware.AdminOnly) // Require admin role

			r.Get("/", formsHandler.ListAllForms) // GET /api/admin/forms - list ALL forms from all users
		})
	})

	return r
}

func main() {
	// Load environment variables
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to database
	database, err := db.New()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Initialize repositories and handlers
	formsRepo := forms.NewFormsRepository(database)
	formsHandler := handlers.NewFormsHandler(formsRepo)

	usersRepo := users.NewUsersRepository(database)
	usersHandler := handlers.NewUsersHandler(usersRepo)
	authHandler := handlers.NewAuthHandler(usersRepo)

	// Setup router
	router := setupRouter(formsHandler, usersHandler, authHandler, usersRepo)

	// Start server
	log.Printf("Server starting on localhost:%s", port)
	log.Printf("Database connected successfully")
	log.Printf("Available endpoints:")
	log.Printf("  GET    /health")
	log.Printf("")
	log.Printf("  Authentication:")
	log.Printf("  POST   /api/auth/login               (public - returns JWT token)")
	log.Printf("")
	log.Printf("  User endpoints:")
	log.Printf("  POST   /api/users                    (public - user registration)")
	log.Printf("  GET    /api/users/{id}               (auth required)")
	log.Printf("  PUT    /api/users/{id}               (auth required)")
	log.Printf("  GET    /api/users                    (admin only)")
	log.Printf("  DELETE /api/users/{id}               (admin only)")
	log.Printf("  POST   /api/users/{id}/approve       (admin only)")
	log.Printf("")
	log.Printf("  Form endpoints:")
	log.Printf("  GET    /api/forms                    (auth required - supports pagination & filtering)")
	log.Printf("  POST   /api/forms/shrub              (auth required)")
	log.Printf("  POST   /api/forms/pesticide          (auth required)")
	log.Printf("  GET    /api/forms/{id}               (auth required)")
	log.Printf("  PUT    /api/forms/{id}               (auth required)")
	log.Printf("  DELETE /api/forms/{id}               (auth required)")
	log.Printf("")
	log.Printf("  Admin-only Form endpoints:")
	log.Printf("  GET    /api/admin/forms              (admin only - list ALL forms from all users)")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
