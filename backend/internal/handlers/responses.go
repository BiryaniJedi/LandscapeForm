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
		ID:        shrubForm.ID,
		CreatedBy: shrubForm.CreatedBy,
		CreatedAt: shrubForm.CreatedAt,
		UpdatedAt: shrubForm.UpdatedAt,
		FirstName: shrubForm.FirstName,
		LastName:  shrubForm.LastName,
		HomePhone: shrubForm.HomePhone,
		NumShrubs: shrubForm.NumShrubs,
	}
}

func pesticideFormToResponse(pesticideForm forms.LawnForm) PesticideFormResponse {
	return PesticideFormResponse{
		ID:            pesticideForm.ID,
		CreatedBy:     pesticideForm.CreatedBy,
		CreatedAt:     pesticideForm.CreatedAt,
		UpdatedAt:     pesticideForm.UpdatedAt,
		FirstName:     pesticideForm.FirstName,
		LastName:      pesticideForm.LastName,
		HomePhone:     pesticideForm.HomePhone,
		PesticideName: pesticideForm.PesticideName,
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

	if view.Lawn != nil {
		resp.ID = view.Lawn.Form.ID
		resp.CreatedBy = view.Lawn.Form.CreatedBy
		resp.CreatedAt = view.Lawn.Form.CreatedAt
		resp.UpdatedAt = view.Lawn.Form.UpdatedAt
		resp.FirstName = view.Lawn.Form.FirstName
		resp.LastName = view.Lawn.Form.LastName
		resp.HomePhone = view.Lawn.Form.HomePhone
		resp.PesticideName = &view.Lawn.PesticideName
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
