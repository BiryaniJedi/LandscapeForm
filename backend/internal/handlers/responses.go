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

func pestAppToResponse(pestApp forms.PestApp) PesticideApplicationResponse {
	return PesticideApplicationResponse{
		ID:            pestApp.ID,
		ChemUsed:      pestApp.ChemUsed,
		AppTimestamp:  pestApp.AppTimestamp,
		Rate:          pestApp.Rate,
		AmountApplied: pestApp.AmountApplied,
		LocationCode:  pestApp.LocationCode,
	}
}

func pestAppsToResponse(pestApps []forms.PestApp) []PesticideApplicationResponse {
	var responses []PesticideApplicationResponse
	for _, pestApp := range pestApps {
		responses = append(responses, pestAppToResponse(pestApp))
	}
	return responses
}

func shrubFormToResponse(shrubForm forms.ShrubForm) ShrubFormResponse {
	return ShrubFormResponse{
		ID:           shrubForm.ID,
		CreatedBy:    shrubForm.CreatedBy,
		CreatedAt:    shrubForm.CreatedAt,
		UpdatedAt:    shrubForm.UpdatedAt,
		FormType:     shrubForm.FormType,
		FirstName:    shrubForm.FirstName,
		LastName:     shrubForm.LastName,
		StreetNumber: shrubForm.StreetNumber,
		StreetName:   shrubForm.StreetName,
		Town:         shrubForm.Town,
		ZipCode:      shrubForm.ZipCode,
		HomePhone:    shrubForm.HomePhone,
		OtherPhone:   shrubForm.OtherPhone,
		CallBefore:   shrubForm.CallBefore,
		IsHoliday:    shrubForm.IsHoliday,
		FirstAppDate: shrubForm.FirstAppDate,
		LastAppDate:  shrubForm.LastAppDate,
		FleaOnly:     shrubForm.FleaOnly,
		PestApps:     pestAppsToResponse(shrubForm.AppTimes),
	}
}

func lawnFormToResponse(lawnForm forms.LawnForm) LawnFormResponse {
	return LawnFormResponse{
		ID:           lawnForm.ID,
		CreatedBy:    lawnForm.CreatedBy,
		CreatedAt:    lawnForm.CreatedAt,
		UpdatedAt:    lawnForm.UpdatedAt,
		FormType:     lawnForm.FormType,
		FirstName:    lawnForm.FirstName,
		LastName:     lawnForm.LastName,
		StreetNumber: lawnForm.StreetNumber,
		StreetName:   lawnForm.StreetName,
		Town:         lawnForm.Town,
		ZipCode:      lawnForm.ZipCode,
		HomePhone:    lawnForm.HomePhone,
		OtherPhone:   lawnForm.OtherPhone,
		CallBefore:   lawnForm.CallBefore,
		IsHoliday:    lawnForm.IsHoliday,
		FirstAppDate: lawnForm.FirstAppDate,
		LastAppDate:  lawnForm.LastAppDate,
		LawnAreaSqFt: lawnForm.LawnAreaSqFt,
		FertOnly:     lawnForm.FertOnly,
		PestApps:     pestAppsToResponse(lawnForm.AppTimes),
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
		resp.StreetNumber = view.Shrub.Form.StreetNumber
		resp.StreetName = view.Shrub.Form.StreetName
		resp.Town = view.Shrub.Form.Town
		resp.ZipCode = view.Shrub.Form.ZipCode
		resp.HomePhone = view.Shrub.Form.HomePhone
		resp.OtherPhone = view.Shrub.Form.OtherPhone
		resp.CallBefore = view.Shrub.Form.CallBefore
		resp.IsHoliday = view.Shrub.Form.IsHoliday
		resp.FirstAppDate = view.Shrub.Form.FirstAppDate
		resp.LastAppDate = view.Shrub.Form.LastAppDate
		resp.PestApps = pestAppsToResponse(view.Shrub.Form.AppTimes)
		resp.FleaOnly = &view.Shrub.FleaOnly
	}

	if view.Lawn != nil {
		resp.ID = view.Lawn.Form.ID
		resp.CreatedBy = view.Lawn.Form.CreatedBy
		resp.CreatedAt = view.Lawn.Form.CreatedAt
		resp.UpdatedAt = view.Lawn.Form.UpdatedAt
		resp.FirstName = view.Lawn.Form.FirstName
		resp.LastName = view.Lawn.Form.LastName
		resp.StreetNumber = view.Lawn.Form.StreetNumber
		resp.StreetName = view.Lawn.Form.StreetName
		resp.Town = view.Lawn.Form.Town
		resp.ZipCode = view.Lawn.Form.ZipCode
		resp.HomePhone = view.Lawn.Form.HomePhone
		resp.OtherPhone = view.Lawn.Form.OtherPhone
		resp.CallBefore = view.Lawn.Form.CallBefore
		resp.IsHoliday = view.Lawn.Form.IsHoliday
		resp.FirstAppDate = view.Lawn.Form.FirstAppDate
		resp.LastAppDate = view.Lawn.Form.LastAppDate
		resp.LawnAreaSqFt = &view.Lawn.LawnAreaSqFt
		resp.PestApps = pestAppsToResponse(view.Lawn.Form.AppTimes)
		resp.FertOnly = &view.Lawn.FertOnly
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
