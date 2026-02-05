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

/*
// CreateUser creates a new user account. This is a public endpoint for user registration.
func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateOrUpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userInput := users.CreateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DoB:       req.DoB,
		Username:  req.Username,
		Password:  req.Password,
	}

	shortUserResponse, err := h.repo.CreateUser(r.Context(), userInput)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondJSON(w, http.StatusCreated, shortUserResponse)
}
*/

// GetUser retrieves a user by ID. Returns the user if found, otherwise returns 404.
func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	if userID == "" {
		respondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

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

// ListUsers retrieves all users with optional sorting. Accepts sort_by and order query parameters.
// Defaults to sorting by last_name DESC.
func (h *UsersHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "last_name"
	}

	order := r.URL.Query().Get("order")
	if order == "" {
		order = "DESC"
	}

	getUserResponses, err := h.repo.ListUsers(r.Context(), sortBy, order)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	fullUserResponses := make([]FullUserResponse, 0, len(getUserResponses))
	for _, getUserResponse := range getUserResponses {
		fullUserResponses = append(fullUserResponses, UserRepoToFullResponse(getUserResponse))
	}

	respondJSON(w, http.StatusOK, ListUsersResponse{
		Users: fullUserResponses,
		Count: len(fullUserResponses),
	})
}

// UpdateUser updates user information by ID. Returns 404 if user not found.
func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "User ID is required")
		return
	}

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

	var req CreateOrUpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

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
		respondError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondJSON(w, http.StatusOK, updatedUser)
}

// DeleteUser deletes a user by ID. Returns the deleted user ID upon success.
func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

// ApproveUser approves a pending user registration by ID. Returns the approved user upon success.
func (h *UsersHandler) ApproveUser(w http.ResponseWriter, r *http.Request) {
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
