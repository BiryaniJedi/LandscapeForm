package forms

import (
	"context"
	"database/sql"
	"os"
	"testing"

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
		INSERT INTO users (email)
		VALUES ('test@example.com')
		RETURNING id
	`).Scan(&id)

	require.NoError(t, err)
	return id
}

func TestCreateAndGetShrubForm(t *testing.T) {
	ctx := context.Background()

	db := testDB(t) // assumes your existing testDB helper
	repo := NewFormsRepository(db)

	userID := createTestUser(t, db)

	created, err := repo.CreateShrubForm(
		ctx,
		CreateFormInput{
			CreatedBy: userID,
			FirstName: "Alice",
			LastName:  "Gardener",
			HomePhone: "555-1234",
		},
		&ShrubDetails{
			NumShrubs: 6,
		},
	)
	require.NoError(t, err)

	// Validate returned data
	require.NotEmpty(t, created.ID)
	require.Equal(t, "shrub", created.FormType)
	require.Equal(t, 6, created.ShrubDetails.NumShrubs)

	// Fetch from DB
	got, err := repo.GetFormById(ctx, created.ID, userID)
	require.NoError(t, err)

	require.NotNil(t, got.Shrub)
	require.Equal(t, "Alice", got.Shrub.Form.FirstName)
	require.Equal(t, "Gardener", got.Shrub.Form.LastName)
	require.Equal(t, "555-1234", got.Shrub.Form.HomePhone)
	require.Equal(t, 6, got.Shrub.ShrubDetails.NumShrubs)
}
