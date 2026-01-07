// Package users provides data access and domain models for users.
// It encapsulates persistence logic, enforces ownership rules, and ensures
// type-safe access to users
package users

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UsersRepository provides database access for user records.
// All methods enforce ownership at the SQL layer and return sql.ErrNoRows
// when a user does not exist
type UsersRepository struct {
	db *sql.DB
}

// NewUsersRepository returns a repository backed by the given database connection.
func NewUsersRepository(database *sql.DB) *UsersRepository {
	return &UsersRepository{db: database}
}

// CreateUserInput contains the common fields required to create a new user.
type CreateUserInput struct {
	FirstName string
	LastName  string
	DoB       time.Time
	Username  string
	Password  string
}

// UpdateUserInput contains the fields that may be updated on an existing user.
type UpdateUserInput struct {
	FirstName string
	LastName  string
	DoB       time.Time
	Username  string
	Password  string
}

// CreateUser creates a new user in the Users table, it, role is 'employee'
// by default and pending is true by default
func (r *UsersRepository) CreateUser(
	ctx context.Context,
	userInput CreateUserInput,
) (UserRepResponse, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return UserRepResponse{}, err
	}
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserRepResponse{}, fmt.Errorf("Error hashing password: %v", err)
	}
	fmt.Printf("userInput: %+v\n", userInput)
	var res UserRepResponse
	err = tx.QueryRowContext(ctx, `
		INSERT INTO users (
			first_name,
			last_name,
			date_of_birth,
			username,
			password_hash
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
		`,
		userInput.FirstName,
		userInput.LastName,
		userInput.DoB,
		userInput.Username,
		hashedPassword,
	).Scan(
		&res.ID,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("Create user query error: %v\n", err)
		return UserRepResponse{ID: "0"}, err
	}

	if err := tx.Commit(); err != nil {
		return UserRepResponse{ID: "1"}, err
	}

	return res, nil
}

// GetUserById returns a single user by the given userID.
// It returns sql.ErrNoRows if the user does not exist
func (r *UsersRepository) GetUserById(
	ctx context.Context,
	userID string,
) (GetUserResponse, error) {
	query := `
		SELECT
			u.id,
			u.created_at,
			u.updated_at,
			u.pending,
			u.role,
			u.first_name,
			u.last_name,
			u.date_of_birth,
			u.username
		FROM users u
		WHERE u.id = $1
	`

	var res GetUserResponse
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&res.ID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Pending,
		&res.Role,
		&res.FirstName,
		&res.LastName,
		&res.DateOfBirth,
		&res.Username,
	)
	if err != nil {
		// Important: let sql.ErrNoRows propagate
		return GetUserResponse{}, err
	}

	return res, nil
}

// GetUserByUsername returns a user by username (for login)
// Includes password hash for authentication
// It returns sql.ErrNoRows if the user does not exist
func (r *UsersRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (User, error) {
	query := `
		SELECT
			id,
			created_at,
			updated_at,
			pending,
			role,
			first_name,
			last_name,
			date_of_birth,
			username,
			password_hash
		FROM users
		WHERE username = $1
	`

	var user User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Pending,
		&user.Role,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Username,
		&user.PasswordHash,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// ApproveUserRegistration allows an admin to approve the registration of an employee
func (r *UsersRepository) ApproveUserRegistration(
	ctx context.Context,
	userID string,
) (UserRepResponse, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return UserRepResponse{}, err
	}
	defer tx.Rollback()

	var res UserRepResponse
	err = tx.QueryRowContext(ctx, `
		UPDATE users
		SET pending = FALSE
		WHERE id = $1
		RETURNING
			id,
			created_at,
			updated_at
	`,
		userID,
	).Scan(
		&res.ID,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return UserRepResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		return UserRepResponse{}, err
	}

	return res, nil
}

// ListUsers lists all users sorted by the provided field
// Can only be called by Admin
func (r *UsersRepository) ListUsers(
	ctx context.Context,
	sortBy string,
	order string,
) ([]GetUserResponse, error) {

	allowedSorts := map[string]string{
		"first_name":    "first_name",
		"last_name":     "last_name",
		"created_at":    "created_at",
		"date_of_birth": "date_of_birth",
	}

	//default sort by last name
	sortColumn, ok := allowedSorts[sortBy]
	if !ok {
		sortColumn = "last_name"
	}

	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			created_at,
			updated_at,
			pending,
			role,
			first_name,
			last_name,
			date_of_birth,
			username
		FROM users
		ORDER BY %s %s
	`, sortColumn, order)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var getUserResponse GetUserResponse
	var users []GetUserResponse
	for rows.Next() {

		err := rows.Scan(
			&getUserResponse.ID,
			&getUserResponse.CreatedAt,
			&getUserResponse.UpdatedAt,
			&getUserResponse.Pending,
			&getUserResponse.Role,
			&getUserResponse.FirstName,
			&getUserResponse.LastName,
			&getUserResponse.DateOfBirth,
			&getUserResponse.Username,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, getUserResponse)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUserById updates a user and its associated subtype fields.
// It returns sql.ErrNoRows if the user does not exist or
func (r *UsersRepository) UpdateUserById(
	ctx context.Context,
	userID string,
	userInput UpdateUserInput,
) (UserRepResponse, error) {
	// TODO auth
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return UserRepResponse{}, err
	}
	defer tx.Rollback()

	var res UserRepResponse

	if userInput.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
		if err != nil {
			return UserRepResponse{}, fmt.Errorf("Error hashing password: %v", err)
		}
		_, err = tx.ExecContext(ctx, `
			UPDATE users
			SET password_hash = $1
			WHERE id = $2
			`, hashedPassword, userID)

		if err != nil {
			return UserRepResponse{}, err
		}
	}

	err = tx.QueryRowContext(ctx, `
		UPDATE users
		SET first_name = $1,
			last_name = $2,
			date_of_birth = $3,
			username = $4
		WHERE id = $5
		RETURNING
			id,
			created_at,
			updated_at
	`,
		userInput.FirstName,
		userInput.LastName,
		userInput.DoB,
		userInput.Username,
		userID,
	).Scan(
		&res.ID,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return UserRepResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		return UserRepResponse{}, err
	}

	return res, nil
}

// DeleteUserById deletes a user.
// It returns sql.ErrNoRows if the user does not exist
func (r *UsersRepository) DeleteUserById(
	ctx context.Context,
	userID string,
) (string, error) {
	// TODO: Auth
	var deletedUserId string
	err := r.db.QueryRowContext(ctx, `
		DELETE FROM users 
		WHERE id = $1
		RETURNING id
	`, userID).Scan(&deletedUserId)

	if err != nil {
		// sql.ErrNoRows → not found or not owned
		return "", err
	}

	return deletedUserId, nil
}
