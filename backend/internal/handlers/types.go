package handlers

import "time"

// Forms

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

// Forms response types

type FormViewResponse struct {
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
	Forms []FormViewResponse `json:"forms"`
	Count int                `json:"count"`
}

// Users
type CreateOrUpdateUserRequest struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	DoB       time.Time `json:"date_of_birth"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
}

type ShortUserResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}

type FullUserResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Pending   bool      `json:"pending"`
	Role      string    `json:"role"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	DoB       time.Time `json:"date_of_birth"`
	Username  string    `json:"username"`
}

type ListUsersResponse struct {
	Users []FullUserResponse `json:"users"`
	Count int                `json:"count"`
}

// Generic Responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
