package handlers

import (
	"github.com/shopspring/decimal"
	"time"
)

// Pesticide Applications

type PesticideApplicationRequest struct {
	ChemUsed      int     `json:"chem_used"`
	AppTimestamp  string  `json:"app_timestamp"`
	Rate          string  `json:"rate"`
	AmountApplied float64 `json:"amount_applied"`
	LocationCode  string  `json:"location_code"`
}

// Forms

type CreateShrubFormRequest struct {
	FirstName    string                        `json:"first_name"`
	LastName     string                        `json:"last_name"`
	StreetNumber string                        `json:"street_number"`
	StreetName   string                        `json:"street_name"`
	Town         string                        `json:"town"`
	ZipCode      string                        `json:"zip_code"`
	HomePhone    string                        `json:"home_phone"`
	OtherPhone   string                        `json:"other_phone"`
	CallBefore   bool                          `json:"call_before"`
	IsHoliday    bool                          `json:"is_holiday"`
	FleaOnly     bool                          `json:"flea_only"`
	Applications []PesticideApplicationRequest `json:"applications,omitempty"`
}

type CreateLawnFormRequest struct {
	FirstName    string                        `json:"first_name"`
	LastName     string                        `json:"last_name"`
	StreetNumber string                        `json:"street_number"`
	StreetName   string                        `json:"street_name"`
	Town         string                        `json:"town"`
	ZipCode      string                        `json:"zip_code"`
	HomePhone    string                        `json:"home_phone"`
	OtherPhone   string                        `json:"other_phone"`
	CallBefore   bool                          `json:"call_before"`
	IsHoliday    bool                          `json:"is_holiday"`
	LawnAreaSqFt int                           `json:"lawn_area_sq_ft"`
	FertOnly     bool                          `json:"fert_only"`
	Applications []PesticideApplicationRequest `json:"applications,omitempty"`
}

type UpdateShrubFormRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	StreetNumber string `json:"street_number"`
	StreetName   string `json:"street_name"`
	Town         string `json:"town"`
	ZipCode      string `json:"zip_code"`
	HomePhone    string `json:"home_phone"`
	OtherPhone   string `json:"other_phone"`
	CallBefore   bool   `json:"call_before"`
	IsHoliday    bool   `json:"is_holiday"`
	FleaOnly     bool   `json:"flea_only"`
}

type UpdateLawnFormRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	StreetNumber string `json:"street_number"`
	StreetName   string `json:"street_name"`
	Town         string `json:"town"`
	ZipCode      string `json:"zip_code"`
	HomePhone    string `json:"home_phone"`
	OtherPhone   string `json:"other_phone"`
	CallBefore   bool   `json:"call_before"`
	IsHoliday    bool   `json:"is_holiday"`
	LawnAreaSqFt int    `json:"lawn_area_sq_ft"`
	FertOnly     bool   `json:"fert_only"`
}

// Forms response types

type FormViewResponse struct {
	ID           string     `json:"id"`
	CreatedBy    string     `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	FormType     string     `json:"form_type"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	StreetNumber string     `json:"street_number"`
	StreetName   string     `json:"street_name"`
	Town         string     `json:"town"`
	ZipCode      string     `json:"zip_code"`
	HomePhone    string     `json:"home_phone"`
	OtherPhone   string     `json:"other_phone"`
	CallBefore   bool      `json:"call_before"`
	IsHoliday    bool      `json:"is_holiday"`
	FirstAppDate time.Time `json:"first_app_date"`
	LastAppDate  time.Time `json:"last_app_date"`
	// Shrub-specific fields (null if lawn form)
	FleaOnly *bool `json:"flea_only,omitempty"`
	// Lawn-specific fields (null if shrub form)
	LawnAreaSqFt *int                           `json:"lawn_area_sq_ft,omitempty"`
	FertOnly     *bool                          `json:"fert_only,omitempty"`
	PestApps     []PesticideApplicationResponse `json:"pest_apps"`
}

type ShrubFormResponse struct {
	ID           string                         `json:"id"`
	CreatedBy    string                         `json:"created_by"`
	CreatedAt    time.Time                      `json:"created_at"`
	UpdatedAt    time.Time                      `json:"updated_at"`
	FormType     string                         `json:"form_type"`
	FirstName    string                         `json:"first_name"`
	LastName     string                         `json:"last_name"`
	StreetNumber string                         `json:"street_number"`
	StreetName   string                         `json:"street_name"`
	Town         string                         `json:"town"`
	ZipCode      string                         `json:"zip_code"`
	HomePhone    string                         `json:"home_phone"`
	OtherPhone   string                         `json:"other_phone"`
	CallBefore   bool                           `json:"call_before"`
	IsHoliday    bool                           `json:"is_holiday"`
	FirstAppDate time.Time                      `json:"first_app_date"`
	LastAppDate  time.Time                      `json:"last_app_date"`
	FleaOnly     bool                           `json:"flea_only"`
	PestApps     []PesticideApplicationResponse `json:"pest_apps"`
}

type LawnFormResponse struct {
	ID           string                         `json:"id"`
	CreatedBy    string                         `json:"created_by"`
	CreatedAt    time.Time                      `json:"created_at"`
	UpdatedAt    time.Time                      `json:"updated_at"`
	FormType     string                         `json:"form_type"`
	FirstName    string                         `json:"first_name"`
	LastName     string                         `json:"last_name"`
	StreetNumber string                         `json:"street_number"`
	StreetName   string                         `json:"street_name"`
	Town         string                         `json:"town"`
	ZipCode      string                         `json:"zip_code"`
	HomePhone    string                         `json:"home_phone"`
	OtherPhone   string                         `json:"other_phone"`
	CallBefore   bool                           `json:"call_before"`
	IsHoliday    bool                           `json:"is_holiday"`
	FirstAppDate time.Time                      `json:"first_app_date"`
	LastAppDate  time.Time                      `json:"last_app_date"`
	LawnAreaSqFt int                            `json:"lawn_area_sq_ft"`
	FertOnly     bool                           `json:"fert_only"`
	PestApps     []PesticideApplicationResponse `json:"pest_apps"`
}

type PesticideApplicationResponse struct {
	ID            int             `json:"id"`
	ChemUsed      int             `json:"chem_used"`
	AppTimestamp  time.Time       `json:"app_timestamp"`
	Rate          string          `json:"rate"`
	AmountApplied decimal.Decimal `json:"amount_applied"`
	LocationCode  string          `json:"location_code"`
}
type ListFormsResponse struct {
	Forms []FormViewResponse `json:"forms"`
	Count int                `json:"count"`
}

type CreateFormResponse struct {
	ID string `json:"id"`
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
