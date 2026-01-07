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

	formInput := forms.CreateFormInput{
		CreatedBy: userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		HomePhone: req.HomePhone,
	}

	shrubDetails := &forms.ShrubDetails{
		NumShrubs: req.NumShrubs,
	}

	shrubForm, err := h.repo.CreateShrubForm(r.Context(), formInput, shrubDetails)
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

	formInput := forms.CreateFormInput{
		CreatedBy: userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		HomePhone: req.HomePhone,
	}

	pesticideDetails := &forms.PesticideDetails{
		PesticideName: req.PesticideName,
	}

	pesticideForm, err := h.repo.CreatePesticideForm(r.Context(), formInput, pesticideDetails)
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
	formResponses := make([]FormResponse, 0, len(views))
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
	formResponses := make([]FormResponse, 0, len(views))
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

// GetForm handles GET /api/forms/{id}
func (h *FormsHandler) GetForm(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	formID := chi.URLParam(r, "id")
	if formID == "" {
		respondError(w, http.StatusBadRequest, "Form ID is required")
		return
	}

	view, err := h.repo.GetFormById(r.Context(), formID, userID)
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

// UpdateForm handles PUT /api/forms/{id}
// Determines form type and updates accordingly
func (h *FormsHandler) UpdateForm(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	formID := chi.URLParam(r, "id")
	if formID == "" {
		respondError(w, http.StatusBadRequest, "Form ID is required")
		return
	}

	// First, get the existing form to determine its type
	existingView, err := h.repo.GetFormById(r.Context(), formID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Form not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch form")
		return
	}

	// Parse request based on form type
	var (
		formInput        forms.UpdateFormInput
		shrubDetails     *forms.ShrubDetails
		pesticideDetails *forms.PesticideDetails
	)

	if existingView.Shrub != nil {
		var req UpdateShrubFormRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// TODO: Add validation

		formInput = forms.UpdateFormInput{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			HomePhone: req.HomePhone,
		}
		shrubDetails = &forms.ShrubDetails{
			NumShrubs: req.NumShrubs,
		}
	} else if existingView.Pesticide != nil {
		var req UpdatePesticideFormRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// TODO: Add validation

		formInput = forms.UpdateFormInput{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			HomePhone: req.HomePhone,
		}
		pesticideDetails = &forms.PesticideDetails{
			PesticideName: req.PesticideName,
		}
	}

	updatedView, err := h.repo.UpdateFormById(r.Context(), formID, userID, formInput, shrubDetails, pesticideDetails)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Form not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update form")
		return
	}

	resp := formViewToResponse(updatedView)
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
