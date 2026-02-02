package forms

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/db"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Load test-specific environment variables
	_ = godotenv.Load("../../.env.testing")

	os.Exit(m.Run())
}

func createTestUser(t *testing.T, db *sql.DB) string {
	t.Helper()

	var id string
	// Generate unique username using random UUID suffix
	err := db.QueryRow(`
		INSERT INTO users (first_name, last_name, username, password_hash)
		VALUES ('Test', 'User', 'TestUser_' || gen_random_uuid()::text, 'TestPass')
		RETURNING id
	`).Scan(&id)

	require.NoError(t, err)
	return id
}

func TestCreateAndGetShrubForm(t *testing.T) {
	ctx := context.Background()

	db := db.TestDB(t) // assumes your existing db.TestDB helper
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	createdShrubFormId, err := repo.CreateShrubForm(
		ctx,
		CreateShrubFormInput{
			CreatedBy:    userID,
			FirstName:    "Alice",
			LastName:     "Gardener",
			StreetNumber: "123",
			StreetName:   "Main St",
			Town:         "Springfield",
			ZipCode:      "12345",
			HomePhone:    "555-1234",
			OtherPhone:   "555-5678",
			CallBefore:   true,
			IsHoliday:    false,
			FleaOnly:     true,
		},
	)
	require.NoError(t, err)

	// Validate returned formID
	require.NotEmpty(t, createdShrubFormId)

	// Fetch from DB
	got, err := repo.GetShrubFormById(ctx, createdShrubFormId, userID)
	require.NoError(t, err)

	require.Equal(t, "Alice", got.FirstName)
	require.Equal(t, "Gardener", got.LastName)
	require.Equal(t, "123", got.StreetNumber)
	require.Equal(t, "Main St", got.StreetName)
	require.Equal(t, "Springfield", got.Town)
	require.Equal(t, "12345", got.ZipCode)
	require.Equal(t, "555-1234", got.HomePhone)
	require.Equal(t, "555-5678", got.OtherPhone)
	require.Equal(t, true, got.CallBefore)
	require.Equal(t, false, got.IsHoliday)
	require.Equal(t, true, got.FleaOnly)
}

func TestCreateAndGetLawnForm(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	createdLawnFormId, err := repo.CreateLawnForm(
		ctx,
		CreateLawnFormInput{
			CreatedBy:    userID,
			FirstName:    "Bob",
			LastName:     "Johnson",
			StreetNumber: "456",
			StreetName:   "Oak Ave",
			Town:         "Shelbyville",
			ZipCode:      "54321",
			HomePhone:    "555-5678",
			OtherPhone:   "555-9999",
			CallBefore:   false,
			IsHoliday:    true,
			LawnAreaSqFt: 5000,
			FertOnly:     false,
		},
	)
	require.NoError(t, err)

	// Validate returned formID
	require.NotEmpty(t, createdLawnFormId)

	// Fetch from DB
	got, err := repo.GetFormViewById(ctx, createdLawnFormId, userID)
	require.NoError(t, err)

	require.NotNil(t, got.Lawn)
	require.Equal(t, "Bob", got.Lawn.Form.FirstName)
	require.Equal(t, "Johnson", got.Lawn.Form.LastName)
	require.Equal(t, "456", got.Lawn.Form.StreetNumber)
	require.Equal(t, "Oak Ave", got.Lawn.Form.StreetName)
	require.Equal(t, "Shelbyville", got.Lawn.Form.Town)
	require.Equal(t, "54321", got.Lawn.Form.ZipCode)
	require.Equal(t, "555-5678", got.Lawn.Form.HomePhone)
	require.Equal(t, "555-9999", got.Lawn.Form.OtherPhone)
	require.Equal(t, false, got.Lawn.Form.CallBefore)
	require.Equal(t, true, got.Lawn.Form.IsHoliday)
	require.Equal(t, 5000, got.Lawn.LawnDetails.LawnAreaSqFt)
	require.Equal(t, false, got.Lawn.LawnDetails.FertOnly)
}

func TestListFormsByUserId_Empty(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)
	listOptions := ListFormsOptions{
		SortBy: "created_at",
		Order:  "DESC",
	}

	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)
	require.Empty(t, forms)
}

func TestListFormsByUserId_MultipleForms(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	// Create shrub form
	_, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    userID,
		FirstName:    "Charlie",
		LastName:     "Brown",
		StreetNumber: "111",
		StreetName:   "Pine St",
		Town:         "Town1",
		ZipCode:      "11111",
		HomePhone:    "555-1111",
		OtherPhone:   "555-1112",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     false,
	})
	require.NoError(t, err)

	// Create lawn form
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Alice",
		LastName:     "Anderson",
		StreetNumber: "222",
		StreetName:   "Elm St",
		Town:         "Town2",
		ZipCode:      "22222",
		HomePhone:    "555-2222",
		OtherPhone:   "555-2223",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 3000,
		FertOnly:     true,
	})
	require.NoError(t, err)

	// Create another shrub form
	_, err = repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    userID,
		FirstName:    "Bob",
		LastName:     "White",
		StreetNumber: "333",
		StreetName:   "Maple St",
		Town:         "Town3",
		ZipCode:      "33333",
		HomePhone:    "555-3333",
		OtherPhone:   "555-3334",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     true,
	})
	require.NoError(t, err)

	listOptions := ListFormsOptions{
		SortBy: "created_at",
		Order:  "DESC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)
	require.Len(t, forms, 3)

	// Check types are correct
	require.Equal(t, "shrub", forms[0].FormType)
	require.Equal(t, "lawn", forms[1].FormType)
	require.Equal(t, "shrub", forms[2].FormType)
}

func TestListFormsByUserId_SortByFirstName(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	// Create forms with different first names
	_, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    userID,
		FirstName:    "Zoe",
		LastName:     "Smith",
		StreetNumber: "100",
		StreetName:   "A St",
		Town:         "Town",
		ZipCode:      "10001",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     false,
	})
	require.NoError(t, err)

	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Alice",
		LastName:     "Jones",
		StreetNumber: "200",
		StreetName:   "B St",
		Town:         "Town",
		ZipCode:      "10002",
		HomePhone:    "555-0002",
		OtherPhone:   "555-0022",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 2000,
		FertOnly:     false,
	})
	require.NoError(t, err)

	_, err = repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    userID,
		FirstName:    "Michael",
		LastName:     "Brown",
		StreetNumber: "300",
		StreetName:   "C St",
		Town:         "Town",
		ZipCode:      "10003",
		HomePhone:    "555-0003",
		OtherPhone:   "555-0033",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     true,
	})
	require.NoError(t, err)

	// Sort by first_name ASC

	listOptions := ListFormsOptions{
		SortBy: "first_name",
		Order:  "ASC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)
	require.Len(t, forms, 3)

	// Helper to get first name from FormView
	getFirstName := func(fv *FormView) string {
		if fv.Shrub != nil {
			return fv.Shrub.Form.FirstName
		}
		return fv.Lawn.Form.FirstName
	}

	require.Equal(t, "Alice", getFirstName(forms[0]))
	require.Equal(t, "Michael", getFirstName(forms[1]))
	require.Equal(t, "Zoe", getFirstName(forms[2]))

	// Sort by first_name DESC
	listOptions.Order = "DESC"
	forms, err = repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)
	require.Len(t, forms, 3)
	require.Equal(t, "Zoe", getFirstName(forms[0]))
	require.Equal(t, "Michael", getFirstName(forms[1]))
	require.Equal(t, "Alice", getFirstName(forms[2]))
}

func TestListFormsByUserId_SortByLastName(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	_, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    userID,
		FirstName:    "John",
		LastName:     "Zimmerman",
		StreetNumber: "111",
		StreetName:   "Z St",
		Town:         "Town",
		ZipCode:      "20001",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     false,
	})
	require.NoError(t, err)

	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Jane",
		LastName:     "Adams",
		StreetNumber: "222",
		StreetName:   "A St",
		Town:         "Town",
		ZipCode:      "20002",
		HomePhone:    "555-0002",
		OtherPhone:   "555-0022",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 1500,
		FertOnly:     false,
	})
	require.NoError(t, err)

	listOptions := ListFormsOptions{
		SortBy: "last_name",
		Order:  "ASC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)
	require.Len(t, forms, 2)

	// Helper to get last name from FormView
	getLastName := func(fv *FormView) string {
		if fv.Shrub != nil {
			return fv.Shrub.Form.LastName
		}
		return fv.Lawn.Form.LastName
	}

	require.Equal(t, "Adams", getLastName(forms[0]))
	require.Equal(t, "Zimmerman", getLastName(forms[1]))
}

func TestListFormsByUserId_OnlyOwnForms(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	// Create two users
	user1ID := createTestUser(t, db)

	var user2ID string
	err := db.QueryRow(`
		INSERT INTO users (first_name, last_name, username, password_hash)
		VALUES ('Test', 'User', 'TestUserName2', 'TestPass')
		RETURNING id
	`).Scan(&user2ID)
	require.NoError(t, err)

	// User 1 creates a form
	_, err = repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    user1ID,
		FirstName:    "User1",
		LastName:     "Form",
		StreetNumber: "11",
		StreetName:   "User St",
		Town:         "Town",
		ZipCode:      "30001",
		HomePhone:    "555-1111",
		OtherPhone:   "555-1112",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     false,
	})
	require.NoError(t, err)

	// User 2 creates a form
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    user2ID,
		FirstName:    "User2",
		LastName:     "Form",
		StreetNumber: "22",
		StreetName:   "User Ave",
		Town:         "Town",
		ZipCode:      "30002",
		HomePhone:    "555-2222",
		OtherPhone:   "555-2223",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 4000,
		FertOnly:     false,
	})
	require.NoError(t, err)

	// Helper to get first name from FormView
	getFirstName := func(fv *FormView) string {
		if fv.Shrub != nil {
			return fv.Shrub.Form.FirstName
		}
		return fv.Lawn.Form.FirstName
	}

	// User 1 should only see their own form
	user1Forms, err := repo.ListFormsByUserId(ctx, user1ID, ListFormsOptions{SortBy: "created_at", Order: "DESC"})
	require.NoError(t, err)
	require.Len(t, user1Forms, 1)
	require.Equal(t, "User1", getFirstName(user1Forms[0]))

	// User 2 should only see their own form
	user2Forms, err := repo.ListFormsByUserId(ctx, user2ID, ListFormsOptions{SortBy: "created_at", Order: "DESC"})
	require.NoError(t, err)
	require.Len(t, user2Forms, 1)
	require.Equal(t, "User2", getFirstName(user2Forms[0]))
}

func TestGetFormById_NotFound(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	// Try to get non-existent form
	_, err := repo.GetFormViewById(ctx, "00000000-0000-0000-0000-000000000000", userID)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestGetFormById_WrongUser(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	// Create user 1 and their form
	user1ID := createTestUser(t, db)
	shrubFormId, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    user1ID,
		FirstName:    "User1",
		LastName:     "Form",
		StreetNumber: "100",
		StreetName:   "Test St",
		Town:         "TestTown",
		ZipCode:      "40001",
		HomePhone:    "555-1111",
		OtherPhone:   "555-1112",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     true,
	})
	require.NoError(t, err)

	// Create user 2
	var user2ID string
	err = db.QueryRow(`
		INSERT INTO users (first_name, last_name, username, password_hash)
		VALUES ('Test', 'User', 'TestUserName3', 'TestPass')
		RETURNING id
	`).Scan(&user2ID)
	require.NoError(t, err)

	// User 2 tries to access User 1's form
	_, err = repo.GetFormViewById(ctx, shrubFormId, user2ID)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err) // Should return ErrNoRows for authorization failure
}

func TestUpdateShrubFormById(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	// Create shrub form
	shrubFormId, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    userID,
		FirstName:    "Original",
		LastName:     "Name",
		StreetNumber: "999",
		StreetName:   "Old St",
		Town:         "OldTown",
		ZipCode:      "90001",
		HomePhone:    "555-0000",
		OtherPhone:   "555-0001",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     false,
	})
	require.NoError(t, err)

	// Update the form
	updated, err := repo.UpdateShrubFormById(
		ctx,
		shrubFormId,
		userID,
		UpdateShrubFormInput{
			FirstName:    "Updated",
			LastName:     "NewName",
			StreetNumber: "888",
			StreetName:   "New Ave",
			Town:         "NewTown",
			ZipCode:      "90002",
			HomePhone:    "555-9999",
			OtherPhone:   "555-9998",
			CallBefore:   true,
			IsHoliday:    true,
			FleaOnly:     true,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, "Updated", updated.FirstName)
	require.Equal(t, "NewName", updated.LastName)
	require.Equal(t, "888", updated.StreetNumber)
	require.Equal(t, "New Ave", updated.StreetName)
	require.Equal(t, "NewTown", updated.Town)
	require.Equal(t, "90002", updated.ZipCode)
	require.Equal(t, "555-9999", updated.HomePhone)
	require.Equal(t, "555-9998", updated.OtherPhone)
	require.Equal(t, true, updated.CallBefore)
	require.Equal(t, true, updated.IsHoliday)
	require.Equal(t, true, updated.FleaOnly)

	// Verify updated_at changed
	require.True(t, updated.UpdatedAt.After(updated.CreatedAt))
}

func TestUpdateLawnFormById(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	// Create lawn form
	lawnFormId, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Jane",
		LastName:     "Doe",
		StreetNumber: "500",
		StreetName:   "Old Lawn St",
		Town:         "GrassTown",
		ZipCode:      "70001",
		HomePhone:    "555-1111",
		OtherPhone:   "555-1112",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 3000,
		FertOnly:     false,
	})
	require.NoError(t, err)

	// Update the form
	updated, err := repo.UpdateLawnFormById(
		ctx,
		lawnFormId,
		userID,
		UpdateLawnFormInput{
			FirstName:    "Janet",
			LastName:     "Smith",
			StreetNumber: "600",
			StreetName:   "New Lawn Ave",
			Town:         "MeadowTown",
			ZipCode:      "70002",
			HomePhone:    "555-2222",
			OtherPhone:   "555-2223",
			CallBefore:   true,
			IsHoliday:    false,
			LawnAreaSqFt: 4500,
			FertOnly:     true,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, "Janet", updated.FirstName)
	require.Equal(t, "Smith", updated.LastName)
	require.Equal(t, "600", updated.StreetNumber)
	require.Equal(t, "New Lawn Ave", updated.StreetName)
	require.Equal(t, "MeadowTown", updated.Town)
	require.Equal(t, "70002", updated.ZipCode)
	require.Equal(t, "555-2222", updated.HomePhone)
	require.Equal(t, "555-2223", updated.OtherPhone)
	require.Equal(t, true, updated.CallBefore)
	require.Equal(t, false, updated.IsHoliday)
	require.Equal(t, 4500, updated.LawnAreaSqFt)
	require.Equal(t, true, updated.FertOnly)
}

func TestUpdateShrubFormById_WrongUser(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	// Create user 1 and their form
	user1ID := createTestUser(t, db)
	shrubFormId, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    user1ID,
		FirstName:    "User1",
		LastName:     "Form",
		StreetNumber: "111",
		StreetName:   "User St",
		Town:         "UserTown",
		ZipCode:      "80001",
		HomePhone:    "555-1111",
		OtherPhone:   "555-1112",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     false,
	})
	require.NoError(t, err)

	// Create user 2
	var user2ID string
	err = db.QueryRow(`
		INSERT INTO users (first_name, last_name, username, password_hash)
		VALUES ('Test', 'User', 'TestUserName4', 'TestPass')
		RETURNING id
	`).Scan(&user2ID)
	require.NoError(t, err)

	// User 2 tries to update User 1's form
	_, err = repo.UpdateShrubFormById(
		ctx,
		shrubFormId,
		user2ID,
		UpdateShrubFormInput{
			FirstName:    "Hacked",
			LastName:     "Name",
			StreetNumber: "999",
			StreetName:   "Hack St",
			Town:         "HackTown",
			ZipCode:      "99999",
			HomePhone:    "555-9999",
			OtherPhone:   "555-9998",
			CallBefore:   true,
			IsHoliday:    true,
			FleaOnly:     true,
		},
	)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestUpdateLawnFormById_WrongUser(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	// Create user 1 and their form
	user1ID := createTestUser(t, db)
	lawnFormId, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    user1ID,
		FirstName:    "User1",
		LastName:     "Form",
		StreetNumber: "222",
		StreetName:   "Lawn St",
		Town:         "LawnTown",
		ZipCode:      "80002",
		HomePhone:    "555-1111",
		OtherPhone:   "555-1112",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 2500,
		FertOnly:     false,
	})
	require.NoError(t, err)

	// Create user 2
	var user2ID string
	err = db.QueryRow(`
		INSERT INTO users (first_name, last_name, username, password_hash)
		VALUES ('Test', 'User', 'TestUserName4', 'TestPass')
		RETURNING id
	`).Scan(&user2ID)
	require.NoError(t, err)

	// User 2 tries to update User 1's form
	_, err = repo.UpdateLawnFormById(
		ctx,
		lawnFormId,
		user2ID,
		UpdateLawnFormInput{
			FirstName:    "Hacked",
			LastName:     "Name",
			StreetNumber: "888",
			StreetName:   "Hack Ave",
			Town:         "HackCity",
			ZipCode:      "88888",
			HomePhone:    "555-9999",
			OtherPhone:   "555-9998",
			CallBefore:   true,
			IsHoliday:    true,
			LawnAreaSqFt: 9999,
			FertOnly:     true,
		},
	)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestDeleteFormById_Success(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	// Create form
	shrubFormId, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    userID,
		FirstName:    "Delete",
		LastName:     "Me",
		StreetNumber: "999",
		StreetName:   "Delete St",
		Town:         "DeleteTown",
		ZipCode:      "99999",
		HomePhone:    "555-0000",
		OtherPhone:   "555-0001",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     false,
	})
	require.NoError(t, err)

	// Delete form
	err = repo.DeleteFormById(ctx, shrubFormId, userID)
	require.NoError(t, err)

	// Verify it's gone
	_, err = repo.GetFormViewById(ctx, shrubFormId, userID)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)

	// Verify shrub details also deleted (cascade)
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM shrub_forms WHERE form_id = $1`, shrubFormId).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestDeleteFormById_WrongUser(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	// Create user 1 and their form
	user1ID := createTestUser(t, db)
	shrubFormId, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    user1ID,
		FirstName:    "User1",
		LastName:     "Form",
		StreetNumber: "777",
		StreetName:   "Persist St",
		Town:         "SafeTown",
		ZipCode:      "77777",
		HomePhone:    "555-1111",
		OtherPhone:   "555-1112",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     true,
	})
	require.NoError(t, err)

	// Create user 2
	var user2ID string
	err = db.QueryRow(`
		INSERT INTO users (first_name, last_name, username, password_hash)
		VALUES ('Test', 'User', 'TestUserName5', 'TestPass')
		RETURNING id
	`).Scan(&user2ID)
	require.NoError(t, err)

	// User 2 tries to delete User 1's form
	err = repo.DeleteFormById(ctx, shrubFormId, user2ID)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)

	// Verify form still exists for user 1
	_, err = repo.GetFormViewById(ctx, shrubFormId, user1ID)
	require.NoError(t, err)
}

func TestDeleteFormById_NotFound(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	// Try to delete non-existent form
	err := repo.DeleteFormById(ctx, "00000000-0000-0000-0000-000000000000", userID)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}
