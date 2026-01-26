package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/chemicals"
	"github.com/go-chi/chi/v5"
)

// ChemicalsHandler handles all chemical-related HTTP requests
type ChemicalsHandler struct {
	repo *chemicals.ChemicalsRepository
}

// NewChemicalsHandler creates a new chemicals handler with the given repository
func NewChemicalsHandler(repo *chemicals.ChemicalsRepository) *ChemicalsHandler {
	return &ChemicalsHandler{repo: repo}
}

// CreateChemicalRequest represents the request body for creating a chemical
type CreateChemicalRequest struct {
	Category     string `json:"category"`
	BrandName    string `json:"brand_name"`
	ChemicalName string `json:"chemical_name"`
	EpaRegNo     string `json:"epa_reg_no"`
	Recipe       string `json:"recipe"`
	Unit         string `json:"unit"`
}

// ChemicalResponse represents the response for a chemical
type ChemicalResponse struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	BrandName    string `json:"brand_name"`
	ChemicalName string `json:"chemical_name"`
	EpaRegNo     string `json:"epa_reg_no"`
	Recipe       string `json:"recipe"`
	Unit         string `json:"unit"`
}

// ListChemicalsResponse represents the response for listing chemicals
type ListChemicalsResponse struct {
	Chemicals []ChemicalResponse `json:"chemicals"`
	Count     int                `json:"count"`
}

// CreateChemical handles POST /api/admin/chemicals
func (h *ChemicalsHandler) CreateChemical(w http.ResponseWriter, r *http.Request) {
	var req CreateChemicalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate category
	if req.Category != "lawn" && req.Category != "shrub" {
		respondError(w, http.StatusBadRequest, "category must be 'lawn' or 'shrub'")
		return
	}

	// Validate required fields
	if req.BrandName == "" || req.ChemicalName == "" {
		respondError(w, http.StatusBadRequest, "brand_name and chemical_name are required")
		return
	}

	chemicalInput := chemicals.ChemicalInput{
		Category:     req.Category,
		BrandName:    req.BrandName,
		ChemicalName: req.ChemicalName,
		EpaRegNo:     req.EpaRegNo,
		Recipe:       req.Recipe,
		Unit:         req.Unit,
	}

	chemicalId, err := h.repo.CreateChemical(r.Context(), chemicalInput)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, CreateFormResponse{ID: chemicalId})
}

// ListChemicals handles GET /api/admin/chemicals
func (h *ChemicalsHandler) ListChemicals(w http.ResponseWriter, r *http.Request) {
	// For now, list all chemicals (both lawn and shrub)
	// We could add filtering by category if needed
	category := r.URL.Query().Get("category")

	var allChemicals []chemicals.Chemical
	var err error

	if category != "" {
		// Validate category
		if category != "lawn" && category != "shrub" {
			respondError(w, http.StatusBadRequest, "category must be 'lawn' or 'shrub'")
			return
		}
		allChemicals, err = h.repo.ListChemicalsByCategory(r.Context(), category)
	} else {
		// List both lawn and shrub
		lawnChems, err1 := h.repo.ListChemicalsByCategory(r.Context(), "lawn")
		if err1 != nil {
			respondError(w, http.StatusInternalServerError, err1.Error())
			return
		}
		shrubChems, err2 := h.repo.ListChemicalsByCategory(r.Context(), "shrub")
		if err2 != nil {
			respondError(w, http.StatusInternalServerError, err2.Error())
			return
		}
		allChemicals = append(lawnChems, shrubChems...)
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chemicalResponses := make([]ChemicalResponse, 0, len(allChemicals))
	for _, chem := range allChemicals {
		chemicalResponses = append(chemicalResponses, ChemicalResponse{
			ID:           chem.ID,
			Category:     chem.Category,
			BrandName:    chem.BrandName,
			ChemicalName: chem.ChemicalName,
			EpaRegNo:     chem.EpaRegNo,
			Recipe:       chem.Recipe,
			Unit:         chem.Unit,
		})
	}

	respondJSON(w, http.StatusOK, ListChemicalsResponse{
		Chemicals: chemicalResponses,
		Count:     len(chemicalResponses),
	})
}

// ListChemicalsByCategory handles GET /api/chemicals/category/{category}
func (h *ChemicalsHandler) ListChemicalsByCategory(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")
	if category == "" {
		respondError(w, http.StatusBadRequest, "Category is required")
		return
	}

	// Validate category
	if category != "lawn" && category != "shrub" {
		respondError(w, http.StatusBadRequest, "category must be 'lawn' or 'shrub'")
		return
	}

	chemicalsList, err := h.repo.ListChemicalsByCategory(r.Context(), category)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chemicalResponses := make([]ChemicalResponse, 0, len(chemicalsList))
	for _, chem := range chemicalsList {
		chemicalResponses = append(chemicalResponses, ChemicalResponse{
			ID:           chem.ID,
			Category:     chem.Category,
			BrandName:    chem.BrandName,
			ChemicalName: chem.ChemicalName,
			EpaRegNo:     chem.EpaRegNo,
			Recipe:       chem.Recipe,
			Unit:         chem.Unit,
		})
	}

	respondJSON(w, http.StatusOK, ListChemicalsResponse{
		Chemicals: chemicalResponses,
		Count:     len(chemicalResponses),
	})
}

// UpdateChemical handles PUT /api/admin/chemicals/{id}
func (h *ChemicalsHandler) UpdateChemical(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		respondError(w, http.StatusBadRequest, "Chemical ID is required")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid chemical ID")
		return
	}

	var req CreateChemicalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate category
	if req.Category != "lawn" && req.Category != "shrub" {
		respondError(w, http.StatusBadRequest, "category must be 'lawn' or 'shrub'")
		return
	}

	// Validate required fields
	if req.BrandName == "" || req.ChemicalName == "" {
		respondError(w, http.StatusBadRequest, "brand_name and chemical_name are required")
		return
	}

	chemicalInput := chemicals.ChemicalInput{
		Category:     req.Category,
		BrandName:    req.BrandName,
		ChemicalName: req.ChemicalName,
		EpaRegNo:     req.EpaRegNo,
		Recipe:       req.Recipe,
		Unit:         req.Unit,
	}

	chemical, err := h.repo.UpdateChemicalById(r.Context(), id, chemicalInput)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Chemical not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, ChemicalResponse{
		ID:           chemical.ID,
		Category:     chemical.Category,
		BrandName:    chemical.BrandName,
		ChemicalName: chemical.ChemicalName,
		EpaRegNo:     chemical.EpaRegNo,
		Recipe:       chemical.Recipe,
		Unit:         chemical.Unit,
	})
}

// DeleteChemical handles DELETE /api/admin/chemicals/{id}
func (h *ChemicalsHandler) DeleteChemical(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		respondError(w, http.StatusBadRequest, "Chemical ID is required")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid chemical ID")
		return
	}

	err = h.repo.DeleteChemicalById(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Chemical not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, "Chemical deleted successfully")
}
