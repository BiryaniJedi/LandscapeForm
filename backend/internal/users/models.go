package users

import (
	"time"
)

type User struct {
	ID           string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Pending      bool
	Role         string
	FirstName    string
	LastName     string
	DateOfBirth  time.Time
	Username     string
	PasswordHash string
}

type UserRepResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetUserResponse struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Pending     bool      `json:"pending"`
	Role        string    `json:"role"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Username    string    `json:"username"`
}
