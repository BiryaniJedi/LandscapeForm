package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/forms"
	"github.com/BiryaniJedi/LandscapeForm-backend/internal/users"
)

// respondJSON writes a JSON response with the given status code
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, we've already written headers, so just log
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// respondError writes a JSON error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

// respondSuccess writes a JSON success message
func respondSuccess(w http.ResponseWriter, message string) {
	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: message,
	})
}

func shrubFormToResponse(shrubForm forms.ShrubForm) ShrubFormResponse {
	return ShrubFormResponse{
		ID: shrubForm.ID,
	}
}

// formViewToResponse converts a FormView from the repository to a FormResponse for the API
func formViewToResponse(view *forms.FormView) FormViewResponse {
	resp := FormViewResponse{
		FormType: view.FormType,
	}

	if view.Shrub != nil {
		resp.ID = view.Shrub.Form.ID
		resp.CreatedBy = view.Shrub.Form.CreatedBy
		resp.CreatedAt = view.Shrub.Form.CreatedAt
		resp.UpdatedAt = view.Shrub.Form.UpdatedAt
		resp.FirstName = view.Shrub.Form.FirstName
		resp.LastName = view.Shrub.Form.LastName
		resp.HomePhone = view.Shrub.Form.HomePhone
		resp.NumShrubs = &view.Shrub.NumShrubs
	}

	if view.Pesticide != nil {
		resp.ID = view.Pesticide.Form.ID
		resp.CreatedBy = view.Pesticide.Form.CreatedBy
		resp.CreatedAt = view.Pesticide.Form.CreatedAt
		resp.UpdatedAt = view.Pesticide.Form.UpdatedAt
		resp.FirstName = view.Pesticide.Form.FirstName
		resp.LastName = view.Pesticide.Form.LastName
		resp.HomePhone = view.Pesticide.Form.HomePhone
		resp.PesticideName = &view.Pesticide.PesticideName
	}

	return resp
}

func UserRepoToFullResponse(user users.GetUserResponse) FullUserResponse {
	return FullUserResponse{
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
}
