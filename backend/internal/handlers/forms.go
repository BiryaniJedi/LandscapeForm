package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/forms"
	"github.com/go-chi/chi/v5"
)

// FormsHandler handles all form-related HTTP requests
type FormsHandler struct {
	repo *forms.FormsRepository
}

// NewFormsHandler creates a new forms handler with the given repository
func NewFormsHandler(repo *forms.FormsRepository) *FormsHandler {
	return &FormsHandler{repo: repo}
}

// getUserID safely extracts userID from context
// Returns a test user ID if not found (for testing without auth)
func getUserID(r *http.Request) string {
	if userID, ok := r.Context().Value("userID").(string); ok {
		return userID
	}
	// Fallback for testing without auth middleware
	// This UUID must exist in the users table
	return "00000000-0000-0000-0000-000000000001"
}

// CreateShrubForm handles POST /api/forms/shrub
func (h *FormsHandler) CreateShrubForm(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	var req CreateShrubFormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Add validation
	// - Check required fields are not empty
	// - Validate phone number format
	// - Validate num_shrubs > 0

	shrubFormInput := forms.CreateShrubFormInput{
		CreatedBy: userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		HomePhone: req.HomePhone,
		NumShrubs: req.NumShrubs,
	}

	shrubForm, err := h.repo.CreateShrubForm(r.Context(), shrubFormInput)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create shrub form")
		return
	}

	view := forms.NewShrubFormView(shrubForm)
	resp := formViewToResponse(view)
	respondJSON(w, http.StatusCreated, resp)
}

// CreatePesticideForm handles POST /api/forms/pesticide
func (h *FormsHandler) CreatePesticideForm(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	var req CreatePesticideFormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Add validation
	// - Check required fields are not empty
	// - Validate phone number format
	// - Validate pesticide_name is not empty

	pesticideFormInput := forms.CreatePesticideFormInput{
		CreatedBy:     userID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		HomePhone:     req.HomePhone,
		PesticideName: req.PesticideName,
	}

	pesticideForm, err := h.repo.CreatePesticideForm(r.Context(), pesticideFormInput)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create pesticide form")
		return
	}

	view := forms.NewPesticideFormView(pesticideForm)
	resp := formViewToResponse(view)
	respondJSON(w, http.StatusCreated, resp)
}

// ListForms handles GET /api/forms?sort_by=created_at&order=DESC&limit=10&offset=0&type=shrub&search=john
func (h *FormsHandler) ListForms(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	// Parse query parameters
	opts := parseListFormsOptions(r)

	views, err := h.repo.ListFormsByUserId(r.Context(), userID, opts)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch forms")
		return
	}

	// Convert to response format
	formResponses := make([]FormViewResponse, 0, len(views))
	for _, view := range views {
		formResponses = append(formResponses, formViewToResponse(view))
	}

	respondJSON(w, http.StatusOK, ListFormsResponse{
		Forms: formResponses,
		Count: len(formResponses),
	})
}

// ListAllForms handles GET /api/admin/forms - returns ALL forms from all users (admin only)
func (h *FormsHandler) ListAllForms(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	opts := parseListFormsOptions(r)

	views, err := h.repo.ListAllForms(r.Context(), opts)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch forms")
		return
	}

	// Convert to response format
	formResponses := make([]FormViewResponse, 0, len(views))
	for _, view := range views {
		formResponses = append(formResponses, formViewToResponse(view))
	}

	respondJSON(w, http.StatusOK, ListFormsResponse{
		Forms: formResponses,
		Count: len(formResponses),
	})
}

// parseListFormsOptions parses query parameters for list forms endpoints
func parseListFormsOptions(r *http.Request) forms.ListFormsOptions {
	opts := forms.ListFormsOptions{}

	// Pagination
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			opts.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			opts.Offset = offset
		}
	}

	// Pagination by page number (alternative to offset)
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 && opts.Limit > 0 {
			opts.Offset = (page - 1) * opts.Limit
		}
	}

	// Filtering
	opts.FormType = r.URL.Query().Get("type")     // "shrub" or "pesticide"
	opts.SearchName = r.URL.Query().Get("search") // search in first_name or last_name

	// Sorting
	opts.SortBy = r.URL.Query().Get("sort_by")
	if opts.SortBy == "" {
		opts.SortBy = "created_at"
	}

	opts.Order = r.URL.Query().Get("order")
	if opts.Order == "" {
		opts.Order = "DESC"
	}

	return opts
}

// GetFormView handles GET /api/forms/{id}, /api/forms/shrub/{id}, and /api/forms/pesticide/{id}
func (h *FormsHandler) GetFormView(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	formID := chi.URLParam(r, "id")
	if formID == "" {
		respondError(w, http.StatusBadRequest, "Form ID is required")
		return
	}

	view, err := h.repo.GetFormViewById(r.Context(), formID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Form not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch form")
		return
	}

	resp := formViewToResponse(view)
	respondJSON(w, http.StatusOK, resp)
}

// UpdateShrubForm handles PUT /api/forms/shrub/{id}
func (h *FormsHandler) UpdateShrubForm(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	formID := chi.URLParam(r, "id")
	if formID == "" {
		respondError(w, http.StatusBadRequest, "Form ID is required")
		return
	}

	// Parse request
	var req UpdateShrubFormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	shrubFormInput := forms.UpdateShrubFormInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		HomePhone: req.HomePhone,
		NumShrubs: req.NumShrubs,
	}

	formView, err := h.repo.UpdateShrubFormById(r.Context(), formID, userID, shrubFormInput)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Form not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update form")
		return
	}

	resp := formViewToResponse(formView)
	respondJSON(w, http.StatusOK, resp)
}

// UpdatePesticideForm handles PUT /api/forms/pesticide/{id}
func (h *FormsHandler) UpdatePesticideForm(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	formID := chi.URLParam(r, "id")
	if formID == "" {
		respondError(w, http.StatusBadRequest, "Form ID is required")
		return
	}

	// Parse request
	var req UpdatePesticideFormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	pesticideFormInput := forms.UpdatePesticideFormInput{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		HomePhone:     req.HomePhone,
		PesticideName: req.PesticideName,
	}

	formView, err := h.repo.UpdatePesticideFormById(r.Context(), formID, userID, pesticideFormInput)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Form not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update form")
		return
	}

	resp := formViewToResponse(formView)
	respondJSON(w, http.StatusOK, resp)
}

// DeleteForm handles DELETE /api/forms/{id}
func (h *FormsHandler) DeleteForm(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	formID := chi.URLParam(r, "id")
	if formID == "" {
		respondError(w, http.StatusBadRequest, "Form ID is required")
		return
	}

	err := h.repo.DeleteFormById(r.Context(), formID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Form not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete form")
		return
	}

	respondSuccess(w, "Form deleted successfully")
}
