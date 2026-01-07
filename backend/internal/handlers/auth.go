package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/auth"
	"github.com/BiryaniJedi/LandscapeForm-backend/internal/users"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	repo *users.UsersRepository
}

// NewAuthHandler creates a new auth handler with the given repository
func NewAuthHandler(repo *users.UsersRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest represents the register request body
type RegisterRequest struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
}

// LoginOrRegisterResponse represents the login response body
type LoginOrRegisterResponse struct {
	Token string           `json:"token"`
	User  FullUserResponse `json:"user"`
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	// Get user by username
	user, err := h.repo.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to authenticate")
		return
	}

	// Check if user is pending approval
	if user.Pending {
		respondError(w, http.StatusForbidden, "Account pending admin approval")
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Prepare response (don't include password hash)
	userResponse := FullUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Pending:   user.Pending,
		Role:      user.Role,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		DoB:       user.DateOfBirth,
		Username:  user.Username,
	}

	respondJSON(w, http.StatusOK, LoginOrRegisterResponse{
		Token: token,
		User:  userResponse,
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" || req.LastName == "" || req.FirstName == "" {
		respondError(w, http.StatusBadRequest, "First name, Last name, Username, and Password are required")
		return
	}

	// Create user
	user, err := h.repo.CreateUser(r.Context(), users.CreateUserInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DoB:       req.DateOfBirth,
		Username:  req.Username,
		Password:  req.Password,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to authenticate")
		return
	}

	userFull, err := h.repo.GetUserById(r.Context(), user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch user")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(userFull.ID, userFull.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Prepare response (don't include password hash)
	userResponse := FullUserResponse{
		ID:        userFull.ID,
		CreatedAt: userFull.CreatedAt,
		UpdatedAt: userFull.UpdatedAt,
		Pending:   userFull.Pending,
		Role:      userFull.Role,
		FirstName: userFull.FirstName,
		LastName:  userFull.LastName,
		DoB:       userFull.DateOfBirth,
		Username:  userFull.Username,
	}

	respondJSON(w, http.StatusCreated, LoginOrRegisterResponse{
		Token: token,
		User:  userResponse,
	})
}
