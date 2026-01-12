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
	err := db.QueryRow(`
		INSERT INTO users (first_name, last_name, username, password_hash)
		VALUES ('Test', 'User', 'TestUserName1', 'TestPass')
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
			CreatedBy: userID,
			FirstName: "Alice",
			LastName:  "Gardener",
			HomePhone: "555-1234",
			NumShrubs: 6,
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
	require.Equal(t, "555-1234", got.HomePhone)
	require.Equal(t, 6, got.NumShrubs)
}

func TestCreateAndGetPesticideForm(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	createdPesticideFormId, err := repo.CreatePesticideForm(
		ctx,
		CreatePesticideFormInput{
			CreatedBy:     userID,
			FirstName:     "Bob",
			LastName:      "Johnson",
			HomePhone:     "555-5678",
			PesticideName: "Roundup",
		},
	)
	require.NoError(t, err)

	// Validate returned formID
	require.NotEmpty(t, createdPesticideFormId)

	// Fetch from DB
	got, err := repo.GetFormViewById(ctx, createdPesticideFormId, userID)
	require.NoError(t, err)

	require.NotNil(t, got.Pesticide)
	require.Equal(t, "Bob", got.Pesticide.Form.FirstName)
	require.Equal(t, "Johnson", got.Pesticide.Form.LastName)
	require.Equal(t, "555-5678", got.Pesticide.Form.HomePhone)
	require.Equal(t, "Roundup", got.Pesticide.PesticideDetails.PesticideName)
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
		CreatedBy: userID,
		FirstName: "Charlie",
		LastName:  "Brown",
		HomePhone: "555-1111",
		NumShrubs: 3,
	})
	require.NoError(t, err)

	// Create pesticide form
	_, err = repo.CreatePesticideForm(ctx, CreatePesticideFormInput{
		CreatedBy:     userID,
		FirstName:     "Alice",
		LastName:      "Anderson",
		HomePhone:     "555-2222",
		PesticideName: "Weed-B-Gone",
	})
	require.NoError(t, err)

	// Create another shrub form
	_, err = repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy: userID,
		FirstName: "Bob",
		LastName:  "White",
		HomePhone: "555-3333",
		NumShrubs: 10,
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
	require.Equal(t, "pesticide", forms[1].FormType)
	require.Equal(t, "shrub", forms[2].FormType)
}

func TestListFormsByUserId_SortByFirstName(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	// Create forms with different first names
	_, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy: userID,
		FirstName: "Zoe",
		LastName:  "Smith",
		HomePhone: "555-0001",
		NumShrubs: 1,
	})
	require.NoError(t, err)

	_, err = repo.CreatePesticideForm(ctx, CreatePesticideFormInput{
		CreatedBy:     userID,
		FirstName:     "Alice",
		LastName:      "Jones",
		HomePhone:     "555-0002",
		PesticideName: "Spray",
	})
	require.NoError(t, err)

	_, err = repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy: userID,
		FirstName: "Michael",
		LastName:  "Brown",
		HomePhone: "555-0003",
		NumShrubs: 5,
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
		return fv.Pesticide.Form.FirstName
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
		CreatedBy: userID,
		FirstName: "John",
		LastName:  "Zimmerman",
		HomePhone: "555-0001",
		NumShrubs: 1,
	})
	require.NoError(t, err)

	_, err = repo.CreatePesticideForm(ctx, CreatePesticideFormInput{
		CreatedBy:     userID,
		FirstName:     "Jane",
		LastName:      "Adams",
		HomePhone:     "555-0002",
		PesticideName: "Bug Killer",
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
		return fv.Pesticide.Form.LastName
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
		CreatedBy: user1ID,
		FirstName: "User1",
		LastName:  "Form",
		HomePhone: "555-1111",
		NumShrubs: 1,
	})
	require.NoError(t, err)

	// User 2 creates a form
	_, err = repo.CreatePesticideForm(ctx, CreatePesticideFormInput{
		CreatedBy:     user2ID,
		FirstName:     "User2",
		LastName:      "Form",
		HomePhone:     "555-2222",
		PesticideName: "Spray",
	})
	require.NoError(t, err)

	// Helper to get first name from FormView
	getFirstName := func(fv *FormView) string {
		if fv.Shrub != nil {
			return fv.Shrub.Form.FirstName
		}
		return fv.Pesticide.Form.FirstName
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
	form, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy: user1ID,
		FirstName: "User1",
		LastName:  "Form",
		HomePhone: "555-1111",
		NumShrubs: 5,
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
	_, err = repo.GetFormViewById(ctx, form.ID, user2ID)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err) // Should return ErrNoRows for authorization failure
}

func TestUpdateShrubFormById(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	// Create shrub form
	created, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy: userID,
		FirstName: "Original",
		LastName:  "Name",
		HomePhone: "555-0000",
		NumShrubs: 5,
	})
	require.NoError(t, err)

	// Update the form
	updated, err := repo.UpdateShrubFormById(
		ctx,
		created.ID,
		userID,
		UpdateShrubFormInput{
			FirstName: "Updated",
			LastName:  "NewName",
			HomePhone: "555-9999",
			NumShrubs: 10,
		},
	)
	require.NoError(t, err)
	require.NotNil(t, updated.Shrub)
	require.Equal(t, "Updated", updated.Shrub.Form.FirstName)
	require.Equal(t, "NewName", updated.Shrub.Form.LastName)
	require.Equal(t, "555-9999", updated.Shrub.Form.HomePhone)
	require.Equal(t, 10, updated.Shrub.ShrubDetails.NumShrubs)

	// Verify updated_at changed
	require.True(t, updated.Shrub.Form.UpdatedAt.After(created.Form.CreatedAt))
}

func TestUpdatePesticideFormById(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	// Create pesticide form
	created, err := repo.CreatePesticideForm(ctx, CreatePesticideFormInput{
		CreatedBy:     userID,
		FirstName:     "Jane",
		LastName:      "Doe",
		HomePhone:     "555-1111",
		PesticideName: "OldSpray",
	})
	require.NoError(t, err)

	// Update the form
	updated, err := repo.UpdatePesticideFormById(
		ctx,
		created.ID,
		userID,
		UpdatePesticideFormInput{
			FirstName:     "Janet",
			LastName:      "Smith",
			HomePhone:     "555-2222",
			PesticideName: "NewSpray",
		},
	)
	require.NoError(t, err)
	require.NotNil(t, updated.Pesticide)
	require.Equal(t, "Janet", updated.Pesticide.Form.FirstName)
	require.Equal(t, "Smith", updated.Pesticide.Form.LastName)
	require.Equal(t, "555-2222", updated.Pesticide.Form.HomePhone)
	require.Equal(t, "NewSpray", updated.Pesticide.PesticideDetails.PesticideName)
}

func TestUpdateShrubFormById_WrongUser(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	// Create user 1 and their form
	user1ID := createTestUser(t, db)
	form, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy: user1ID,
		FirstName: "User1",
		LastName:  "Form",
		HomePhone: "555-1111",
		NumShrubs: 5,
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
		form.ID,
		user2ID,
		UpdateShrubFormInput{
			FirstName: "Hacked",
			LastName:  "Name",
			HomePhone: "555-9999",
			NumShrubs: 100,
		},
	)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestUpdatePesticideFormById_WrongUser(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	// Create user 1 and their form
	user1ID := createTestUser(t, db)
	form, err := repo.CreatePesticideForm(ctx, CreatePesticideFormInput{
		CreatedBy:     user1ID,
		FirstName:     "User1",
		LastName:      "Form",
		HomePhone:     "555-1111",
		PesticideName: "genPest",
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
	_, err = repo.UpdatePesticideFormById(
		ctx,
		form.ID,
		user2ID,
		UpdatePesticideFormInput{
			FirstName:     "Hacked",
			LastName:      "Name",
			HomePhone:     "555-9999",
			PesticideName: "genPest2",
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
	form, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy: userID,
		FirstName: "Delete",
		LastName:  "Me",
		HomePhone: "555-0000",
		NumShrubs: 1,
	})
	require.NoError(t, err)

	// Delete form
	err = repo.DeleteFormById(ctx, form.ID, userID)
	require.NoError(t, err)

	// Verify it's gone
	_, err = repo.GetFormViewById(ctx, form.ID, userID)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)

	// Verify shrub details also deleted (cascade)
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM shrubs WHERE form_id = $1`, form.ID).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestDeleteFormById_WrongUser(t *testing.T) {
	ctx := context.Background()
	db := db.TestDB(t)
	repo := NewFormsRepository(db)

	// Create user 1 and their form
	user1ID := createTestUser(t, db)
	form, err := repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy: user1ID,
		FirstName: "User1",
		LastName:  "Form",
		HomePhone: "555-1111",
		NumShrubs: 5,
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
	err = repo.DeleteFormById(ctx, form.ID, user2ID)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)

	// Verify form still exists for user 1
	_, err = repo.GetFormViewById(ctx, form.ID, user1ID)
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
