package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/users"
	"github.com/go-chi/chi/v5"
)

// UsersHandler handles all user-related HTTP requests
type UsersHandler struct {
	repo *users.UsersRepository
}

// NewUsersHandler creates a new users handler with the given repository
func NewUsersHandler(repo *users.UsersRepository) *UsersHandler {
	return &UsersHandler{repo: repo}
}

// CreateUser handles POST /api/users
// This is a public endpoint for user registration
func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateOrUpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Add validation
	// - Check required fields are not empty
	// - Validate username format
	// - Validate password strength (min length, etc.)
	// - Validate date_of_birth is in the past

	userInput := users.CreateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DoB:       req.DoB,
		Username:  req.Username,
		Password:  req.Password,
	}

	shortUserResponse, err := h.repo.CreateUser(r.Context(), userInput)
	if err != nil {
		// TODO: Check for duplicate username error and return appropriate message
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondJSON(w, http.StatusCreated, shortUserResponse)
}

// GetUser handles GET /api/users/{id}
// MIDDLEWARE REQUIRED: Authentication - Users can only view their own profile
// MIDDLEWARE REQUIRED: Admin can view any user profile
func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	if userID == "" {
		respondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// TODO: Add authorization check
	// - Extract authenticated user ID from context (set by auth middleware)
	// - Check if authenticated user ID matches requested user ID OR user is admin
	// - Return 403 Forbidden if not authorized

	getUserResponse, err := h.repo.GetUserById(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	resp := UserRepoToFullResponse(getUserResponse)
	respondJSON(w, http.StatusOK, resp)
}

// ListUsers handles GET /api/users?sort_by=last_name&order=DESC
// MIDDLEWARE REQUIRED: Admin only - Only admins can list all users
func (h *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Add admin authorization check
	// - Extract authenticated user from context (set by auth middleware)
	// - Check if user role is 'admin'
	// - Return 403 Forbidden if not admin

	// Parse query parameters
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "last_name"
	}

	order := r.URL.Query().Get("order")
	if order == "" {
		order = "DESC"
	}

	// TODO: Add pagination support (limit, offset)

	getUserResponses, err := h.repo.ListUsers(r.Context(), sortBy, order)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	// Convert to response format
	fullUserResponses := make([]FullUserResponse, 0, len(getUserResponses))
	for _, getUserResponse := range getUserResponses {
		fullUserResponses = append(fullUserResponses, UserRepoToFullResponse(getUserResponse))
	}

	respondJSON(w, http.StatusOK, ListUsersResponse{
		Users: fullUserResponses,
		Count: len(fullUserResponses),
	})
}

// UpdateUser handles PUT /api/users/{id}
// MIDDLEWARE REQUIRED: Authentication - Users can only update their own profile
// MIDDLEWARE REQUIRED: Admin can update any user profile
func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// TODO: Add authorization check
	// - Extract authenticated user ID from context (set by auth middleware)
	// - Check if authenticated user ID matches requested user ID OR user is admin
	// - Return 403 Forbidden if not authorized

	// Check if user exists
	_, err := h.repo.GetUserById(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	// Parse request
	var req CreateOrUpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Add validation
	// - Validate fields similar to CreateUser
	// - Password is optional on update (only update if provided)

	userInput := users.UpdateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DoB:       req.DoB,
		Username:  req.Username,
		Password:  req.Password, // Empty string if not provided
	}

	updatedUser, err := h.repo.UpdateUserById(r.Context(), userID, userInput)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}
		// TODO: Check for duplicate username error
		respondError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondJSON(w, http.StatusOK, updatedUser)
}

// DeleteUser handles DELETE /api/users/{id}
// MIDDLEWARE REQUIRED: Admin only - Only admins can delete users
func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Add admin authorization check
	// - Extract authenticated user from context (set by auth middleware)
	// - Check if user role is 'admin'
	// - Return 403 Forbidden if not admin

	userID := chi.URLParam(r, "id")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	deletedUserID, err := h.repo.DeleteUserById(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "User deleted successfully",
		"id":      deletedUserID,
	})
}

// ApproveUser handles POST /api/users/{id}/approve
// MIDDLEWARE REQUIRED: Admin only - Only admins can approve pending users
func (h *UsersHandler) ApproveUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Add admin authorization check
	// - Extract authenticated user from context (set by auth middleware)
	// - Check if user role is 'admin'
	// - Return 403 Forbidden if not admin

	userID := chi.URLParam(r, "id")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	approvedUser, err := h.repo.ApproveUserRegistration(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to approve user")
		return
	}

	respondJSON(w, http.StatusOK, approvedUser)
}
