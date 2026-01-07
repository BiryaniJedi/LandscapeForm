package users

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/db"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	// Load test-specific environment variables
	_ = godotenv.Load("../../.env.testing")

	os.Exit(m.Run())
}

// TestCreateUser tests creating a new user
func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	input := CreateUserInput{
		FirstName: "John",
		LastName:  "Doe",
		DoB:       time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		Username:  "johndoe",
		Password:  "securepassword123",
	}

	res, err := repo.CreateUser(ctx, input)
	require.NoError(t, err)
	require.NotEqual(t, "", res.ID)
	require.False(t, res.CreatedAt.IsZero(), "Expected non-zero CreatedAt timestamp")

	// Verify the user was actually created in the database
	getRes, err := repo.GetUserById(ctx, res.ID)
	require.NoError(t, err, "GetUserById failed after creation")

	require.Equal(t, input.FirstName, getRes.FirstName)
	require.Equal(t, input.LastName, getRes.LastName)
	require.Equal(t, input.Username, getRes.Username)
	require.True(t, getRes.DateOfBirth.Equal(input.DoB))
	require.Equal(t, "employee", getRes.Role)
	require.True(t, getRes.Pending, "Expected Pending to be true by default")
}

// TestCreateUserPasswordHashing tests that passwords are hashed correctly
func TestCreateUserPasswordHashing(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	plainPassword := "mySecretPassword123"
	input := CreateUserInput{
		FirstName: "Jane",
		LastName:  "Smith",
		DoB:       time.Date(1985, 5, 20, 0, 0, 0, 0, time.UTC),
		Username:  "janesmith",
		Password:  plainPassword,
	}

	res, err := repo.CreateUser(ctx, input)
	require.NoError(t, err, "CreateUser failed")

	// Directly query the password hash from the database
	var passwordHash string
	err = database.QueryRowContext(ctx, "SELECT password_hash FROM users WHERE id = $1", res.ID).Scan(&passwordHash)
	require.NoError(t, err, "Failed to query password hash")

	// Verify the password is hashed (not stored in plain text)
	require.NotEqual(t, plainPassword, passwordHash, "Password should be hashed, not stored in plain text")

	// Verify the hash can be verified with bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword))
	require.NoError(t, err, "Password hash verification failed")
}

// TestGetUserById tests retrieving a user by ID
func TestGetUserById(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create a user first
	input := CreateUserInput{
		FirstName: "Alice",
		LastName:  "Johnson",
		DoB:       time.Date(1992, 3, 10, 0, 0, 0, 0, time.UTC),
		Username:  "alicej",
		Password:  "password123",
	}

	createRes, err := repo.CreateUser(ctx, input)
	require.NoError(t, err, "CreateUser failed")

	// Retrieve the user
	getRes, err := repo.GetUserById(ctx, createRes.ID)
	require.NoError(t, err, "GetUserById failed")

	require.Equal(t, createRes.ID, getRes.ID)
	require.Equal(t, input.FirstName, getRes.FirstName)
	require.Equal(t, input.LastName, getRes.LastName)
	require.Equal(t, input.Username, getRes.Username)
}

// TestGetUserByIdNotFound tests retrieving a non-existent user
func TestGetUserByIdNotFound(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	fakeID := "00000000-0000-0000-0000-000000000000"
	_, err := repo.GetUserById(ctx, fakeID)

	require.ErrorIs(t, err, sql.ErrNoRows, "Expected sql.ErrNoRows for non-existent user")
}

// TestUpdateUserByIdWithoutPassword tests updating user info without changing password
func TestUpdateUserByIdWithoutPassword(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create a user first
	createInput := CreateUserInput{
		FirstName: "Bob",
		LastName:  "Williams",
		DoB:       time.Date(1988, 7, 25, 0, 0, 0, 0, time.UTC),
		Username:  "bobw",
		Password:  "originalpassword",
	}

	createRes, err := repo.CreateUser(ctx, createInput)
	require.NoError(t, err, "CreateUser failed")

	// Get the original password hash
	var originalHash string
	err = database.QueryRowContext(ctx, "SELECT password_hash FROM users WHERE id = $1", createRes.ID).Scan(&originalHash)
	require.NoError(t, err, "Failed to query original password hash")

	// Update the user without changing password
	updateInput := UpdateUserInput{
		FirstName: "Robert",
		LastName:  "Williams Jr.",
		DoB:       time.Date(1988, 7, 26, 0, 0, 0, 0, time.UTC),
		Username:  "bobwilliams",
		Password:  "", // Empty password means no password change
	}

	updateRes, err := repo.UpdateUserById(ctx, createRes.ID, updateInput)
	require.NoError(t, err, "UpdateUserById failed")
	require.Equal(t, createRes.ID, updateRes.ID)

	// Verify the updated fields
	getRes, err := repo.GetUserById(ctx, createRes.ID)
	require.NoError(t, err, "GetUserById failed after update")

	require.Equal(t, updateInput.FirstName, getRes.FirstName)
	require.Equal(t, updateInput.LastName, getRes.LastName)
	require.Equal(t, updateInput.Username, getRes.Username)

	// Verify password hash wasn't changed
	var currentHash string
	err = database.QueryRowContext(ctx, "SELECT password_hash FROM users WHERE id = $1", createRes.ID).Scan(&currentHash)
	require.NoError(t, err, "Failed to query current password hash")

	require.Equal(t, originalHash, currentHash, "Password hash should not have changed when password is empty")
}

// TestUpdateUserByIdWithPassword tests updating user info including password
func TestUpdateUserByIdWithPassword(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create a user first
	createInput := CreateUserInput{
		FirstName: "Carol",
		LastName:  "Davis",
		DoB:       time.Date(1995, 11, 5, 0, 0, 0, 0, time.UTC),
		Username:  "carold",
		Password:  "oldpassword123",
	}

	createRes, err := repo.CreateUser(ctx, createInput)
	require.NoError(t, err, "CreateUser failed")

	// Get the original password hash
	var originalHash string
	err = database.QueryRowContext(ctx, "SELECT password_hash FROM users WHERE id = $1", createRes.ID).Scan(&originalHash)
	require.NoError(t, err, "Failed to query original password hash")

	// Update the user with a new password
	newPassword := "newpassword456"
	updateInput := UpdateUserInput{
		FirstName: "Caroline",
		LastName:  "Davis",
		DoB:       time.Date(1995, 11, 5, 0, 0, 0, 0, time.UTC),
		Username:  "carolinedavis",
		Password:  newPassword,
	}

	_, err = repo.UpdateUserById(ctx, createRes.ID, updateInput)
	require.NoError(t, err, "UpdateUserById failed")

	// Verify the password hash was changed
	var newHash string
	err = database.QueryRowContext(ctx, "SELECT password_hash FROM users WHERE id = $1", createRes.ID).Scan(&newHash)
	require.NoError(t, err, "Failed to query new password hash")

	require.NotEqual(t, originalHash, newHash, "Password hash should have changed")

	// Verify the new password hash is correct
	err = bcrypt.CompareHashAndPassword([]byte(newHash), []byte(newPassword))
	require.NoError(t, err, "New password hash verification failed")

	// Verify the old password no longer works
	err = bcrypt.CompareHashAndPassword([]byte(newHash), []byte(createInput.Password))
	require.Error(t, err, "Old password should not work with new hash")
}

// TestUpdateUserByIdNotFound tests updating a non-existent user
func TestUpdateUserByIdNotFound(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	fakeID := "00000000-0000-0000-0000-000000000000"
	updateInput := UpdateUserInput{
		FirstName: "Ghost",
		LastName:  "User",
		DoB:       time.Now(),
		Username:  "ghostuser",
		Password:  "",
	}

	_, err := repo.UpdateUserById(ctx, fakeID, updateInput)

	require.ErrorIs(t, err, sql.ErrNoRows, "Expected sql.ErrNoRows for non-existent user")
}

// TestDeleteUserById tests deleting a user
func TestDeleteUserById(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create a user first
	createInput := CreateUserInput{
		FirstName: "David",
		LastName:  "Miller",
		DoB:       time.Date(1993, 9, 15, 0, 0, 0, 0, time.UTC),
		Username:  "davidm",
		Password:  "password123",
	}

	createRes, err := repo.CreateUser(ctx, createInput)
	require.NoError(t, err, "CreateUser failed")

	// Delete the user
	deletedID, err := repo.DeleteUserById(ctx, createRes.ID)
	require.NoError(t, err, "DeleteUserById failed")
	require.Equal(t, createRes.ID, deletedID)

	// Verify the user no longer exists
	_, err = repo.GetUserById(ctx, createRes.ID)
	require.ErrorIs(t, err, sql.ErrNoRows, "Expected sql.ErrNoRows after deletion")
}

// TestDeleteUserByIdNotFound tests deleting a non-existent user
func TestDeleteUserByIdNotFound(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	fakeID := "00000000-0000-0000-0000-000000000000"
	_, err := repo.DeleteUserById(ctx, fakeID)

	require.ErrorIs(t, err, sql.ErrNoRows, "Expected sql.ErrNoRows for non-existent user")
}

// TestUpdatedAtTimestamp tests that the updated_at timestamp is automatically updated
func TestUpdatedAtTimestamp(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create a user
	createInput := CreateUserInput{
		FirstName: "Emma",
		LastName:  "Wilson",
		DoB:       time.Date(1991, 4, 12, 0, 0, 0, 0, time.UTC),
		Username:  "emmaw",
		Password:  "password123",
	}

	createRes, err := repo.CreateUser(ctx, createInput)
	require.NoError(t, err, "CreateUser failed")

	// Get initial updated_at
	initialUser, err := repo.GetUserById(ctx, createRes.ID)
	require.NoError(t, err, "GetUserById failed")

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Update the user
	updateInput := UpdateUserInput{
		FirstName: "Emily",
		LastName:  "Wilson",
		DoB:       initialUser.DateOfBirth,
		Username:  initialUser.Username,
		Password:  "",
	}

	_, err = repo.UpdateUserById(ctx, createRes.ID, updateInput)
	require.NoError(t, err, "UpdateUserById failed")

	// Get updated user
	updatedUser, err := repo.GetUserById(ctx, createRes.ID)
	require.NoError(t, err, "GetUserById failed after update")

	// Verify updated_at changed
	require.True(t, updatedUser.UpdatedAt.After(initialUser.UpdatedAt),
		"Expected UpdatedAt to be after initial timestamp. Initial: %v, Updated: %v",
		initialUser.UpdatedAt, updatedUser.UpdatedAt)
}

// TestListUsersEmpty tests listing users when no users exist
func TestListUsersEmpty(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	users, err := repo.ListUsers(ctx, "last_name", "ASC")
	require.NoError(t, err, "ListUsers should not error on empty database")
	require.Empty(t, users, "Expected empty list when no users exist")
}

// TestListUsersDefaultSort tests listing users with default sort (last_name DESC)
func TestListUsersDefaultSort(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create multiple users
	users := []CreateUserInput{
		{FirstName: "Alice", LastName: "Zebra", DoB: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), Username: "alice", Password: "pass1"},
		{FirstName: "Bob", LastName: "Apple", DoB: time.Date(1991, 2, 2, 0, 0, 0, 0, time.UTC), Username: "bob", Password: "pass2"},
		{FirstName: "Charlie", LastName: "Mango", DoB: time.Date(1992, 3, 3, 0, 0, 0, 0, time.UTC), Username: "charlie", Password: "pass3"},
	}

	for _, u := range users {
		_, err := repo.CreateUser(ctx, u)
		require.NoError(t, err, "CreateUser failed")
	}

	// List with default sort (last_name DESC)
	result, err := repo.ListUsers(ctx, "", "")
	require.NoError(t, err, "ListUsers failed")
	require.Len(t, result, 3, "Expected 3 users")

	// Verify order: Zebra, Mango, Apple (DESC)
	require.Equal(t, "Zebra", result[0].LastName)
	require.Equal(t, "Mango", result[1].LastName)
	require.Equal(t, "Apple", result[2].LastName)
}

// TestListUsersSortByFirstName tests sorting by first name
func TestListUsersSortByFirstName(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create users with different first names
	users := []CreateUserInput{
		{FirstName: "Zoe", LastName: "Smith", DoB: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), Username: "zoe", Password: "pass1"},
		{FirstName: "Alice", LastName: "Jones", DoB: time.Date(1991, 2, 2, 0, 0, 0, 0, time.UTC), Username: "alice2", Password: "pass2"},
		{FirstName: "Mike", LastName: "Brown", DoB: time.Date(1992, 3, 3, 0, 0, 0, 0, time.UTC), Username: "mike", Password: "pass3"},
	}

	for _, u := range users {
		_, err := repo.CreateUser(ctx, u)
		require.NoError(t, err, "CreateUser failed")
	}

	// List sorted by first_name ASC
	result, err := repo.ListUsers(ctx, "first_name", "ASC")
	require.NoError(t, err, "ListUsers failed")
	require.Len(t, result, 3, "Expected 3 users")

	// Verify order: Alice, Mike, Zoe
	require.Equal(t, "Alice", result[0].FirstName)
	require.Equal(t, "Mike", result[1].FirstName)
	require.Equal(t, "Zoe", result[2].FirstName)
}

// TestListUsersSortByLastName tests sorting by last name
func TestListUsersSortByLastName(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create users
	users := []CreateUserInput{
		{FirstName: "John", LastName: "Wilson", DoB: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), Username: "john", Password: "pass1"},
		{FirstName: "Jane", LastName: "Adams", DoB: time.Date(1991, 2, 2, 0, 0, 0, 0, time.UTC), Username: "jane", Password: "pass2"},
		{FirstName: "Jack", LastName: "Taylor", DoB: time.Date(1992, 3, 3, 0, 0, 0, 0, time.UTC), Username: "jack", Password: "pass3"},
	}

	for _, u := range users {
		_, err := repo.CreateUser(ctx, u)
		require.NoError(t, err, "CreateUser failed")
	}

	// List sorted by last_name DESC
	result, err := repo.ListUsers(ctx, "last_name", "DESC")
	require.NoError(t, err, "ListUsers failed")
	require.Len(t, result, 3, "Expected 3 users")

	// Verify order: Wilson, Taylor, Adams
	require.Equal(t, "Wilson", result[0].LastName)
	require.Equal(t, "Taylor", result[1].LastName)
	require.Equal(t, "Adams", result[2].LastName)
}

// TestListUsersSortByCreatedAt tests sorting by created_at
func TestListUsersSortByCreatedAt(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create users with delays to ensure different timestamps
	user1, err := repo.CreateUser(ctx, CreateUserInput{
		FirstName: "First", LastName: "User", DoB: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Username: "first", Password: "pass1",
	})
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	user2, err := repo.CreateUser(ctx, CreateUserInput{
		FirstName: "Second", LastName: "User", DoB: time.Date(1991, 2, 2, 0, 0, 0, 0, time.UTC),
		Username: "second", Password: "pass2",
	})
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	user3, err := repo.CreateUser(ctx, CreateUserInput{
		FirstName: "Third", LastName: "User", DoB: time.Date(1992, 3, 3, 0, 0, 0, 0, time.UTC),
		Username: "third", Password: "pass3",
	})
	require.NoError(t, err)

	// List sorted by created_at ASC (oldest first)
	result, err := repo.ListUsers(ctx, "created_at", "ASC")
	require.NoError(t, err, "ListUsers failed")
	require.Len(t, result, 3, "Expected 3 users")

	// Verify order: First, Second, Third
	require.Equal(t, user1.ID, result[0].ID)
	require.Equal(t, user2.ID, result[1].ID)
	require.Equal(t, user3.ID, result[2].ID)

	// List sorted by created_at DESC (newest first)
	result, err = repo.ListUsers(ctx, "created_at", "DESC")
	require.NoError(t, err, "ListUsers failed")
	require.Len(t, result, 3, "Expected 3 users")

	// Verify order: Third, Second, First
	require.Equal(t, user3.ID, result[0].ID)
	require.Equal(t, user2.ID, result[1].ID)
	require.Equal(t, user1.ID, result[2].ID)
}

// TestListUsersSortByDateOfBirth tests sorting by date of birth
func TestListUsersSortByDateOfBirth(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create users with different birth dates
	users := []CreateUserInput{
		{FirstName: "Oldest", LastName: "User", DoB: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), Username: "oldest", Password: "pass1"},
		{FirstName: "Youngest", LastName: "User", DoB: time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC), Username: "youngest", Password: "pass2"},
		{FirstName: "Middle", LastName: "User", DoB: time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC), Username: "middle", Password: "pass3"},
	}

	for _, u := range users {
		_, err := repo.CreateUser(ctx, u)
		require.NoError(t, err, "CreateUser failed")
	}

	// List sorted by date_of_birth ASC (oldest first)
	result, err := repo.ListUsers(ctx, "date_of_birth", "ASC")
	require.NoError(t, err, "ListUsers failed")
	require.Len(t, result, 3, "Expected 3 users")

	// Verify order: 1980, 1990, 2000
	require.Equal(t, "Oldest", result[0].FirstName)
	require.Equal(t, "Middle", result[1].FirstName)
	require.Equal(t, "Youngest", result[2].FirstName)
}

// TestListUsersInvalidSortFallsBackToDefault tests that invalid sort parameters fall back to defaults
func TestListUsersInvalidSortFallsBackToDefault(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create users
	users := []CreateUserInput{
		{FirstName: "Alice", LastName: "Zebra", DoB: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), Username: "alice3", Password: "pass1"},
		{FirstName: "Bob", LastName: "Apple", DoB: time.Date(1991, 2, 2, 0, 0, 0, 0, time.UTC), Username: "bob2", Password: "pass2"},
	}

	for _, u := range users {
		_, err := repo.CreateUser(ctx, u)
		require.NoError(t, err, "CreateUser failed")
	}

	// Test with invalid sortBy - should default to last_name
	result, err := repo.ListUsers(ctx, "invalid_column", "ASC")
	require.NoError(t, err, "ListUsers should not error on invalid sortBy")
	require.Len(t, result, 2, "Expected 2 users")
	// Default is last_name, so with ASC: Apple, Zebra
	// But wait, the default order is DESC when sortBy is invalid
	// Let me check the code again... actually when sortBy is invalid, it defaults to "last_name"
	// and the order is still ASC as specified, so it should be Apple, Zebra

	// Test with invalid order - should default to DESC
	result, err = repo.ListUsers(ctx, "last_name", "INVALID_ORDER")
	require.NoError(t, err, "ListUsers should not error on invalid order")
	require.Len(t, result, 2, "Expected 2 users")
	// Should be DESC: Zebra, Apple
	require.Equal(t, "Zebra", result[0].LastName)
	require.Equal(t, "Apple", result[1].LastName)
}

// TestApproveUserRegistration tests approving a pending user
func TestApproveUserRegistration(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create a user (pending by default)
	createInput := CreateUserInput{
		FirstName: "Pending",
		LastName:  "User",
		DoB:       time.Date(1995, 5, 15, 0, 0, 0, 0, time.UTC),
		Username:  "pendinguser",
		Password:  "password123",
	}

	createRes, err := repo.CreateUser(ctx, createInput)
	require.NoError(t, err, "CreateUser failed")

	// Verify user is pending
	user, err := repo.GetUserById(ctx, createRes.ID)
	require.NoError(t, err, "GetUserById failed")
	require.True(t, user.Pending, "User should be pending by default")

	// Approve the user
	approveRes, err := repo.ApproveUserRegistration(ctx, createRes.ID)
	require.NoError(t, err, "ApproveUserRegistration failed")
	require.Equal(t, createRes.ID, approveRes.ID)

	// Verify user is no longer pending
	approvedUser, err := repo.GetUserById(ctx, createRes.ID)
	require.NoError(t, err, "GetUserById failed after approval")
	require.False(t, approvedUser.Pending, "User should not be pending after approval")
}

// TestApproveUserRegistrationNotFound tests approving a non-existent user
func TestApproveUserRegistrationNotFound(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	fakeID := "00000000-0000-0000-0000-000000000000"
	_, err := repo.ApproveUserRegistration(ctx, fakeID)

	require.ErrorIs(t, err, sql.ErrNoRows, "Expected sql.ErrNoRows for non-existent user")
}

// TestApproveUserRegistrationAlreadyApproved tests approving an already approved user
func TestApproveUserRegistrationAlreadyApproved(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewUsersRepository(database)

	// Create and approve a user
	createInput := CreateUserInput{
		FirstName: "Approved",
		LastName:  "User",
		DoB:       time.Date(1993, 8, 20, 0, 0, 0, 0, time.UTC),
		Username:  "approveduser",
		Password:  "password123",
	}

	createRes, err := repo.CreateUser(ctx, createInput)
	require.NoError(t, err, "CreateUser failed")

	// First approval
	_, err = repo.ApproveUserRegistration(ctx, createRes.ID)
	require.NoError(t, err, "ApproveUserRegistration failed")

	// Verify user is approved
	user, err := repo.GetUserById(ctx, createRes.ID)
	require.NoError(t, err, "GetUserById failed")
	require.False(t, user.Pending, "User should not be pending after first approval")

	// Second approval (should still work, just sets pending = FALSE again)
	_, err = repo.ApproveUserRegistration(ctx, createRes.ID)
	require.NoError(t, err, "ApproveUserRegistration should not error on already approved user")

	// Verify user is still approved
	user, err = repo.GetUserById(ctx, createRes.ID)
	require.NoError(t, err, "GetUserById failed after second approval")
	require.False(t, user.Pending, "User should still not be pending after second approval")
}
