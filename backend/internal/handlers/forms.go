package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/forms"
	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
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
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Add validation
	// - Check required fields are not empty
	// - Validate phone number format

	// Convert applications from request to domain model
	var applications []forms.PestApp
	for _, appReq := range req.Applications {
		appTime, err := time.Parse(time.RFC3339, appReq.AppTimestamp)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid application timestamp format: "+err.Error())
			return
		}

		applications = append(applications, forms.PestApp{
			ChemUsed:      appReq.ChemUsed,
			AppTimestamp:  appTime,
			Rate:          appReq.Rate,
			AmountApplied: decimal.NewFromFloat(appReq.AmountApplied),
			LocationCode:  appReq.LocationCode,
		})
	}

	shrubFormInput := forms.CreateShrubFormInput{
		CreatedBy:    userID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		StreetNumber: req.StreetNumber,
		StreetName:   req.StreetName,
		Town:         req.Town,
		ZipCode:      req.ZipCode,
		HomePhone:    req.HomePhone,
		OtherPhone:   req.OtherPhone,
		CallBefore:   req.CallBefore,
		IsHoliday:    req.IsHoliday,
		FleaOnly:     req.FleaOnly,
		Applications: applications,
	}

	shrubFormId, err := h.repo.CreateShrubForm(r.Context(), shrubFormInput)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, CreateFormResponse{shrubFormId})
}

// CreateLawnForm handles POST /api/forms/lawn
func (h *FormsHandler) CreateLawnForm(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	var req CreateLawnFormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Debug logging
	fmt.Printf("Received lawn form request with %d applications\n", len(req.Applications))
	for i, app := range req.Applications {
		fmt.Printf("  App %d: ChemUsed=%d, Rate=%s, Amount=%.2f, Location=%s, Timestamp=%s\n",
			i, app.ChemUsed, app.Rate, app.AmountApplied, app.LocationCode, app.AppTimestamp)
	}

	// TODO: Add validation
	// - Check required fields are not empty
	// - Validate phone number format
	// - Validate lawn_area_sq_ft > 0

	// Convert applications from request to domain model
	var applications []forms.PestApp
	for _, appReq := range req.Applications {
		appTime, err := time.Parse(time.RFC3339, appReq.AppTimestamp)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid application timestamp format: "+err.Error())
			return
		}

		applications = append(applications, forms.PestApp{
			ChemUsed:      appReq.ChemUsed,
			AppTimestamp:  appTime,
			Rate:          appReq.Rate,
			AmountApplied: decimal.NewFromFloat(appReq.AmountApplied),
			LocationCode:  appReq.LocationCode,
		})
	}

	lawnFormInput := forms.CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		StreetNumber: req.StreetNumber,
		StreetName:   req.StreetName,
		Town:         req.Town,
		ZipCode:      req.ZipCode,
		HomePhone:    req.HomePhone,
		OtherPhone:   req.OtherPhone,
		CallBefore:   req.CallBefore,
		IsHoliday:    req.IsHoliday,
		LawnAreaSqFt: req.LawnAreaSqFt,
		FertOnly:     req.FertOnly,
		Applications: applications,
	}

	lawnFormId, err := h.repo.CreateLawnForm(r.Context(), lawnFormInput)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, CreateFormResponse{lawnFormId})
}

// ListForms handles GET /api/forms?sort_by=created_at&order=DESC&limit=10&offset=0&type=shrub&search=john
func (h *FormsHandler) ListForms(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	// Parse query parameters
	opts := parseListFormsOptions(r)

	views, err := h.repo.ListFormsByUserId(r.Context(), userID, opts)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
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
		respondError(w, http.StatusInternalServerError, err.Error())
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
	opts.FormType = r.URL.Query().Get("type")     // "shrub" or "lawn"
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

// GetShrubForm handles GET /api/forms/shrub/{id}
func (h *FormsHandler) GetShrubForm(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	formID := chi.URLParam(r, "id")
	if formID == "" {
		respondError(w, http.StatusBadRequest, "Form ID is required")
		return
	}

	shrubForm, err := h.repo.GetShrubFormById(r.Context(), formID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, shrubFormToResponse(shrubForm))
}

// GetLawnForm handles GET /api/forms/lawn/{id}
func (h *FormsHandler) GetLawnForm(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	formID := chi.URLParam(r, "id")
	if formID == "" {
		respondError(w, http.StatusBadRequest, "Form ID is required")
		return
	}

	lawnForm, err := h.repo.GetLawnFormById(r.Context(), formID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, lawnFormToResponse(lawnForm))
}

// GetFormView handles GET /api/forms/{id}
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
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
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
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	shrubFormInput := forms.UpdateShrubFormInput{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		StreetNumber: req.StreetNumber,
		StreetName:   req.StreetName,
		Town:         req.Town,
		ZipCode:      req.ZipCode,
		HomePhone:    req.HomePhone,
		OtherPhone:   req.OtherPhone,
		CallBefore:   req.CallBefore,
		IsHoliday:    req.IsHoliday,
		FleaOnly:     req.FleaOnly,
	}

	shrubForm, err := h.repo.UpdateShrubFormById(r.Context(), formID, userID, shrubFormInput)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to update form")
		return
	}

	respondJSON(w, http.StatusOK, shrubFormToResponse(shrubForm))
}

// UpdateLawnForm handles PUT /api/forms/lawn/{id}
func (h *FormsHandler) UpdateLawnForm(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	formID := chi.URLParam(r, "id")
	if formID == "" {
		respondError(w, http.StatusBadRequest, "Form ID is required")
		return
	}

	// Parse request
	var req UpdateLawnFormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	lawnFormInput := forms.UpdateLawnFormInput{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		StreetNumber: req.StreetNumber,
		StreetName:   req.StreetName,
		Town:         req.Town,
		ZipCode:      req.ZipCode,
		HomePhone:    req.HomePhone,
		OtherPhone:   req.OtherPhone,
		CallBefore:   req.CallBefore,
		IsHoliday:    req.IsHoliday,
		LawnAreaSqFt: req.LawnAreaSqFt,
		FertOnly:     req.FertOnly,
	}

	lawnForm, err := h.repo.UpdateLawnFormById(r.Context(), formID, userID, lawnFormInput)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, lawnFormToResponse(lawnForm))
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
			respondError(w, http.StatusNotFound, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, "Form deleted successfully")
}
