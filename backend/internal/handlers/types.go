package handlers

import "time"

// Request types for creating and updating forms

type CreateShrubFormRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	HomePhone string `json:"home_phone"`
	NumShrubs int    `json:"num_shrubs"`
}

type CreatePesticideFormRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	HomePhone     string `json:"home_phone"`
	PesticideName string `json:"pesticide_name"`
}

type UpdateShrubFormRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	HomePhone string `json:"home_phone"`
	NumShrubs int    `json:"num_shrubs"`
}

type UpdatePesticideFormRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	HomePhone     string `json:"home_phone"`
	PesticideName string `json:"pesticide_name"`
}

// Response types

type FormResponse struct {
	ID        string    `json:"id"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FormType  string    `json:"form_type"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	HomePhone string    `json:"home_phone"`
	// Shrub-specific fields (null if pesticide form)
	NumShrubs *int `json:"num_shrubs,omitempty"`
	// Pesticide-specific fields (null if shrub form)
	PesticideName *string `json:"pesticide_name,omitempty"`
}

type ListFormsResponse struct {
	Forms []FormResponse `json:"forms"`
	Count int            `json:"count"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
