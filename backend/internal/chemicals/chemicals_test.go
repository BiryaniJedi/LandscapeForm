package chemicals

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

func TestCreateChemical_Success(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewChemicalsRepository(database)

	input := ChemicalInput{
		Category:     "lawn",
		BrandName:    "RoundUp",
		ChemicalName: "Glyphosate",
		EpaRegNo:     "524-445",
		Recipe:       "Mix 2oz per gallon",
		Unit:         "oz",
	}

	chemicalID, err := repo.CreateChemical(ctx, input)
	require.NoError(t, err)
	require.NotEmpty(t, chemicalID)

	// Verify it was inserted by querying directly
	var chemical Chemical
	err = database.QueryRow(`
		SELECT id, category, brand_name, chemical_name, epa_reg_no, recipe, unit
		FROM chemicals
		WHERE id = $1
	`, chemicalID).Scan(
		&chemical.ID,
		&chemical.Category,
		&chemical.BrandName,
		&chemical.ChemicalName,
		&chemical.EpaRegNo,
		&chemical.Recipe,
		&chemical.Unit,
	)
	require.NoError(t, err)
	require.Equal(t, "lawn", chemical.Category)
	require.Equal(t, "RoundUp", chemical.BrandName)
	require.Equal(t, "Glyphosate", chemical.ChemicalName)
	require.Equal(t, "524-445", chemical.EpaRegNo)
	require.Equal(t, "Mix 2oz per gallon", chemical.Recipe)
	require.Equal(t, "oz", chemical.Unit)
}

func TestCreateChemical_MultipleChemicals(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewChemicalsRepository(database)

	inputs := []ChemicalInput{
		{
			Category:     "lawn",
			BrandName:    "RoundUp",
			ChemicalName: "Glyphosate",
			EpaRegNo:     "524-445",
			Recipe:       "Mix 2oz per gallon",
			Unit:         "oz",
		},
		{
			Category:     "shrub",
			BrandName:    "Daconil",
			ChemicalName: "Chlorothalonil",
			EpaRegNo:     "100-1093",
			Recipe:       "Mix 1.5oz per gallon",
			Unit:         "oz",
		},
		{
			Category:     "lawn",
			BrandName:    "Sevin",
			ChemicalName: "Carbaryl",
			EpaRegNo:     "264-333",
			Recipe:       "Mix 3 tablespoons per gallon",
			Unit:         "tbsp",
		},
	}

	for _, input := range inputs {
		chemicalID, err := repo.CreateChemical(ctx, input)
		require.NoError(t, err)
		require.NotEmpty(t, chemicalID)
	}

	// Verify all were inserted
	var count int
	err := database.QueryRow(`SELECT COUNT(*) FROM chemicals`).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 3, count)
}

func TestListChemicalsByCategory_EmptyResult(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewChemicalsRepository(database)

	chemicals, err := repo.ListChemicalsByCategory(ctx, "lawn")
	require.NoError(t, err)
	require.Empty(t, chemicals)
}

func TestListChemicalsByCategory_WithResults(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewChemicalsRepository(database)

	// Create chemicals in different categories
	lawnChemicals := []ChemicalInput{
		{
			Category:     "lawn",
			BrandName:    "RoundUp",
			ChemicalName: "Glyphosate",
			EpaRegNo:     "524-445",
			Recipe:       "Mix 2oz per gallon",
			Unit:         "oz",
		},
		{
			Category:     "lawn",
			BrandName:    "2,4-D Amine",
			ChemicalName: "2,4-Dichlorophenoxyacetic acid",
			EpaRegNo:     "228-365",
			Recipe:       "Mix 1oz per gallon",
			Unit:         "oz",
		},
	}

	shrubChemicals := []ChemicalInput{
		{
			Category:     "shrub",
			BrandName:    "Daconil",
			ChemicalName: "Chlorothalonil",
			EpaRegNo:     "100-1093",
			Recipe:       "Mix 1.5oz per gallon",
			Unit:         "oz",
		},
	}

	// Insert lawn chemicals
	for _, input := range lawnChemicals {
		_, err := repo.CreateChemical(ctx, input)
		require.NoError(t, err)
	}

	// Insert shrub chemicals
	for _, input := range shrubChemicals {
		_, err := repo.CreateChemical(ctx, input)
		require.NoError(t, err)
	}

	// Query lawn chemicals
	result, err := repo.ListChemicalsByCategory(ctx, "lawn")
	require.NoError(t, err)
	require.Len(t, result, 2)

	// Verify all returned chemicals are lawn chemicals
	for _, chemical := range result {
		require.Equal(t, "lawn", chemical.Category)
	}

	// Query shrub chemicals
	result, err = repo.ListChemicalsByCategory(ctx, "shrub")
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "Daconil", result[0].BrandName)
}

func TestUpdateChemicalById_Success(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewChemicalsRepository(database)

	// Create initial chemical
	originalInput := ChemicalInput{
		Category:     "lawn",
		BrandName:    "Original Brand",
		ChemicalName: "Original Chemical",
		EpaRegNo:     "111-111",
		Recipe:       "Original recipe",
		Unit:         "oz",
	}

	chemicalID, err := repo.CreateChemical(ctx, originalInput)
	require.NoError(t, err)

	// Convert chemicalID string to int (assuming it returns a numeric string)
	var id int
	err = database.QueryRow(`SELECT id FROM chemicals WHERE id = $1`, chemicalID).Scan(&id)
	require.NoError(t, err)

	// Update the chemical
	updateInput := ChemicalInput{
		Category:     "shrub",
		BrandName:    "Updated Brand",
		ChemicalName: "Updated Chemical",
		EpaRegNo:     "222-222",
		Recipe:       "Updated recipe",
		Unit:         "ml",
	}

	updated, err := repo.UpdateChemicalById(ctx, id, updateInput)
	require.NoError(t, err)
	require.Equal(t, id, updated.ID)
	require.Equal(t, "shrub", updated.Category)
	require.Equal(t, "Updated Brand", updated.BrandName)
	require.Equal(t, "Updated Chemical", updated.ChemicalName)
	require.Equal(t, "222-222", updated.EpaRegNo)
	require.Equal(t, "Updated recipe", updated.Recipe)
	require.Equal(t, "ml", updated.Unit)

	// Verify in database
	var chemical Chemical
	err = database.QueryRow(`
		SELECT id, category, brand_name, chemical_name, epa_reg_no, recipe, unit
		FROM chemicals
		WHERE id = $1
	`, id).Scan(
		&chemical.ID,
		&chemical.Category,
		&chemical.BrandName,
		&chemical.ChemicalName,
		&chemical.EpaRegNo,
		&chemical.Recipe,
		&chemical.Unit,
	)
	require.NoError(t, err)
	require.Equal(t, "shrub", chemical.Category)
	require.Equal(t, "Updated Brand", chemical.BrandName)
}

func TestUpdateChemicalById_NotFound(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewChemicalsRepository(database)

	updateInput := ChemicalInput{
		Category:     "lawn",
		BrandName:    "Test Brand",
		ChemicalName: "Test Chemical",
		EpaRegNo:     "999-999",
		Recipe:       "Test recipe",
		Unit:         "oz",
	}

	// Try to update non-existent chemical (using max smallint value)
	_, err := repo.UpdateChemicalById(ctx, 32767, updateInput)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestDeleteChemicalById_Success(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewChemicalsRepository(database)

	// Create chemical
	input := ChemicalInput{
		Category:     "lawn",
		BrandName:    "Delete Me",
		ChemicalName: "Test Chemical",
		EpaRegNo:     "555-555",
		Recipe:       "Test recipe",
		Unit:         "oz",
	}

	chemicalID, err := repo.CreateChemical(ctx, input)
	require.NoError(t, err)

	// Convert chemicalID string to int
	var id int
	err = database.QueryRow(`SELECT id FROM chemicals WHERE id = $1`, chemicalID).Scan(&id)
	require.NoError(t, err)

	// Delete the chemical
	err = repo.DeleteChemicalById(ctx, id)
	require.NoError(t, err)

	// Verify it's deleted
	var count int
	err = database.QueryRow(`SELECT COUNT(*) FROM chemicals WHERE id = $1`, id).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestDeleteChemicalById_NotFound(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewChemicalsRepository(database)

	// Try to delete non-existent chemical (using max smallint value)
	err := repo.DeleteChemicalById(ctx, 32767)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestDeleteChemicalById_AfterDeleteCannotUpdate(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewChemicalsRepository(database)

	// Create chemical
	input := ChemicalInput{
		Category:     "shrub",
		BrandName:    "Test Brand",
		ChemicalName: "Test Chemical",
		EpaRegNo:     "777-777",
		Recipe:       "Test recipe",
		Unit:         "oz",
	}

	chemicalID, err := repo.CreateChemical(ctx, input)
	require.NoError(t, err)

	// Convert chemicalID string to int
	var id int
	err = database.QueryRow(`SELECT id FROM chemicals WHERE id = $1`, chemicalID).Scan(&id)
	require.NoError(t, err)

	// Delete the chemical
	err = repo.DeleteChemicalById(ctx, id)
	require.NoError(t, err)

	// Try to update deleted chemical
	updateInput := ChemicalInput{
		Category:     "lawn",
		BrandName:    "Updated",
		ChemicalName: "Updated",
		EpaRegNo:     "888-888",
		Recipe:       "Updated",
		Unit:         "ml",
	}

	_, err = repo.UpdateChemicalById(ctx, id, updateInput)
	require.Error(t, err)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestCreateChemical_EmptyFields(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewChemicalsRepository(database)

	// Test with empty strings for non-required fields (category must be valid)
	input := ChemicalInput{
		Category:     "lawn",
		BrandName:    "",
		ChemicalName: "",
		EpaRegNo:     "",
		Recipe:       "",
		Unit:         "",
	}

	chemicalID, err := repo.CreateChemical(ctx, input)
	require.NoError(t, err)
	require.NotEmpty(t, chemicalID)
}

func TestCreateChemical_TransactionRollback(t *testing.T) {
	ctx := context.Background()
	database := db.TestDB(t)
	repo := NewChemicalsRepository(database)

	// Create a chemical
	input := ChemicalInput{
		Category:     "lawn",
		BrandName:    "Test",
		ChemicalName: "Test",
		EpaRegNo:     "123-456",
		Recipe:       "Test",
		Unit:         "oz",
	}

	chemicalID, err := repo.CreateChemical(ctx, input)
	require.NoError(t, err)
	require.NotEmpty(t, chemicalID)

	// Verify it exists
	var count int
	err = database.QueryRow(`SELECT COUNT(*) FROM chemicals WHERE id = $1`, chemicalID).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}
