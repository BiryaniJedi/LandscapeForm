package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/auth"
	"github.com/BiryaniJedi/LandscapeForm-backend/internal/middleware"
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

// AuthUserResponse represents the login response body
type AuthUserResponse struct {
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

	// Set token as HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 hours
	})

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

	respondJSON(w, http.StatusOK, AuthUserResponse{
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
		fmt.Printf("User id test error: %s\n", user.ID)
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

	// Set token as HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 hours
	})

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

	respondJSON(w, http.StatusCreated, AuthUserResponse{
		Token: token,
		User:  userResponse,
	})
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Immediately expire
	})

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// Me handles GET /api/auth/me
// Route protected already, two layers redundant
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("r.Context(): %#v\n", r.Context())

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		// fmt.Printf("userID not ok: %v\n", userID)
		http.Error(w, "unauthorized: userID not ok", http.StatusUnauthorized)
		return
	}
	// fmt.Printf("userID: %v\n", userID)

	/* Degbug Info
	role, ok := middleware.GetUserRole(r.Context())
	if !ok {
		fmt.Printf("userRole not ok: %v\n", role)
		http.Error(w, "unauthorized: role not OK", http.StatusUnauthorized)
		return
	}
	fmt.Printf("userRole: %v\n", role)

	pending, ok := middleware.GetUserPending(r.Context())
	if !ok {
		fmt.Printf("userPending not ok: %v\n", pending)
		http.Error(w, "unauthorized: pending not OK", http.StatusUnauthorized)
		return
	}

	fmt.Printf("userPending: %v\n", pending)

	fmt.Printf("User id: %v\nRole: %v\nPending: %v\n", userID, role, pending)
	*/

	user, err := h.repo.GetUserById(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	/*  Not used anymore, using middleware
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	tokenStr := cookie.Value
	fmt.Printf("tokenStr: %s\f", tokenStr)

	claims, err := auth.ValidateToken(tokenStr)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.repo.GetUserById(r.Context(), claims.UserID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	fmt.Printf("User: %+v\n", user)
	*/

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

	respondJSON(w, http.StatusOK, userResponse)
}
