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

func setupRouter(formsHandler *handlers.FormsHandler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Recovery)      // Recover from panics
	r.Use(middleware.Logger)        // Log all requests
	r.Use(middleware.CORS)          // Enable CORS
	r.Use(chimiddleware.RequestID)  // Add request ID to each request
	r.Use(chimiddleware.RealIP)     // Get real client IP

	// Public routes (no auth required)
	r.Get("/health", healthHandler)

	// TODO: Add authentication routes here
	// r.Post("/api/auth/register", authHandler.Register)
	// r.Post("/api/auth/login", authHandler.Login)

	// Protected routes (require authentication)
	r.Route("/api", func(r chi.Router) {
		// Apply auth middleware to all /api routes
		//r.Use(middleware.AuthMiddleware)

		// Forms endpoints
		r.Route("/forms", func(r chi.Router) {
			r.Get("/", formsHandler.ListForms)           // GET /api/forms?sort_by=created_at&order=DESC
			r.Post("/shrub", formsHandler.CreateShrubForm)       // POST /api/forms/shrub
			r.Post("/pesticide", formsHandler.CreatePesticideForm) // POST /api/forms/pesticide

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", formsHandler.GetForm)       // GET /api/forms/{id}
				r.Put("/", formsHandler.UpdateForm)    // PUT /api/forms/{id}
				r.Delete("/", formsHandler.DeleteForm) // DELETE /api/forms/{id}
			})
		})

		// TODO: Add admin routes here (for viewing all forms)
		// r.Route("/admin", func(r chi.Router) {
		//     r.Use(middleware.AdminOnly)
		//     r.Get("/forms", adminHandler.ListAllForms)
		// })
	})

	return r
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
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

	// Initialize repository and handlers
	formsRepo := forms.NewFormsRepository(database)
	formsHandler := handlers.NewFormsHandler(formsRepo)

	// Setup router
	router := setupRouter(formsHandler)

	// Start server
	log.Printf("Server starting on localhost:%s", port)
	log.Printf("Database connected successfully")
	log.Printf("Available endpoints:")
	log.Printf("  GET    /health")
	log.Printf("  GET    /api/forms")
	log.Printf("  POST   /api/forms/shrub")
	log.Printf("  POST   /api/forms/pesticide")
	log.Printf("  GET    /api/forms/{id}")
	log.Printf("  PUT    /api/forms/{id}")
	log.Printf("  DELETE /api/forms/{id}")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
