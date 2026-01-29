package forms

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/BiryaniJedi/LandscapeForm-backend/internal/db"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// Helper to create a test chemical
func createTestChemical(t *testing.T, db *sql.DB, category string) int {
	t.Helper()

	var id int
	err := db.QueryRow(`
		INSERT INTO chemicals (category, brand_name, chemical_name, epa_reg_no, recipe, unit)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, category, "Test Brand", "Test Chemical "+category, "12345-67", "Test Recipe", "oz").Scan(&id)

	require.NoError(t, err)
	return id
}

// Helper to get first name from FormView
func getFirstName(fv *FormView) string {
	if fv.Shrub != nil {
		return fv.Shrub.Form.FirstName
	}
	return fv.Lawn.Form.FirstName
}

// Helper to get last name from FormView
func getLastName(fv *FormView) string {
	if fv.Shrub != nil {
		return fv.Shrub.Form.LastName
	}
	return fv.Lawn.Form.LastName
}

func TestListFormsByUserId_SortByFirstAppDate(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	userID := createTestUser(t, testDB)
	chemID := createTestChemical(t, testDB, "lawn")

	now := time.Now()

	// Create form with earliest application (3 days ago)
	_, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Early",
		LastName:     "Bird",
		StreetNumber: "100",
		StreetName:   "First St",
		Town:         "Town",
		ZipCode:      "10001",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-72 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create form with latest application (1 day ago)
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Late",
		LastName:     "Riser",
		StreetNumber: "200",
		StreetName:   "Last St",
		Town:         "Town",
		ZipCode:      "10002",
		HomePhone:    "555-0002",
		OtherPhone:   "555-0022",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 2000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-24 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create form with middle application (2 days ago)
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Middle",
		LastName:     "Ground",
		StreetNumber: "300",
		StreetName:   "Middle St",
		Town:         "Town",
		ZipCode:      "10003",
		HomePhone:    "555-0003",
		OtherPhone:   "555-0033",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 3000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-48 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create form WITHOUT application (should appear at end with NULLS LAST)
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "No",
		LastName:     "Application",
		StreetNumber: "400",
		StreetName:   "Empty St",
		Town:         "Town",
		ZipCode:      "10004",
		HomePhone:    "555-0004",
		OtherPhone:   "555-0044",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 4000,
		FertOnly:     false,
		Applications: []PestApp{},
	})
	require.NoError(t, err)

	// Test ASC order (oldest first, nulls last)
	listOptions := ListFormsOptions{
		SortBy: "first_app_date",
		Order:  "ASC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)
	require.Len(t, forms, 4)

	require.Equal(t, "Early", getFirstName(forms[0]))   // 3 days ago
	require.Equal(t, "Middle", getFirstName(forms[1]))  // 2 days ago
	require.Equal(t, "Late", getFirstName(forms[2]))    // 1 day ago
	require.Equal(t, "No", getFirstName(forms[3]))      // NULL (no application)

	// Test DESC order (newest first, nulls last)
	listOptions.Order = "DESC"
	forms, err = repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)
	require.Len(t, forms, 4)

	require.Equal(t, "Late", getFirstName(forms[0]))    // 1 day ago
	require.Equal(t, "Middle", getFirstName(forms[1]))  // 2 days ago
	require.Equal(t, "Early", getFirstName(forms[2]))   // 3 days ago
	require.Equal(t, "No", getFirstName(forms[3]))      // NULL (no application)
}

func TestListFormsByUserId_FilterByDateLow(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	userID := createTestUser(t, testDB)
	chemID := createTestChemical(t, testDB, "lawn")

	now := time.Now()
	twoDaysAgo := now.Add(-48 * time.Hour)

	// Create form with application 3 days ago (before cutoff)
	_, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Old",
		LastName:     "Form",
		StreetNumber: "100",
		StreetName:   "Old St",
		Town:         "Town",
		ZipCode:      "10001",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-72 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create form with application 1 day ago (after cutoff)
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Recent",
		LastName:     "Form",
		StreetNumber: "200",
		StreetName:   "Recent St",
		Town:         "Town",
		ZipCode:      "10002",
		HomePhone:    "555-0002",
		OtherPhone:   "555-0022",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 2000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-24 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Filter for forms with first application >= 2 days ago
	listOptions := ListFormsOptions{
		DateLow: twoDaysAgo,
		SortBy:  "created_at",
		Order:   "DESC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	// Should only get the recent form
	require.Len(t, forms, 1)
	require.Equal(t, "Recent", getFirstName(forms[0]))
}

func TestListFormsByUserId_FilterByDateHigh(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	userID := createTestUser(t, testDB)
	chemID := createTestChemical(t, testDB, "lawn")

	now := time.Now()
	twoDaysAgo := now.Add(-48 * time.Hour)

	// Create form with last application 3 days ago (before cutoff)
	_, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Old",
		LastName:     "Form",
		StreetNumber: "100",
		StreetName:   "Old St",
		Town:         "Town",
		ZipCode:      "10001",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-72 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create form with last application 1 day ago (after cutoff)
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Recent",
		LastName:     "Form",
		StreetNumber: "200",
		StreetName:   "Recent St",
		Town:         "Town",
		ZipCode:      "10002",
		HomePhone:    "555-0002",
		OtherPhone:   "555-0022",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 2000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-24 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Filter for forms with last application <= 2 days ago
	listOptions := ListFormsOptions{
		DateHigh: twoDaysAgo,
		SortBy:   "created_at",
		Order:    "DESC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	// Should only get the old form
	require.Len(t, forms, 1)
	require.Equal(t, "Old", getFirstName(forms[0]))
}

func TestListFormsByUserId_FilterByDateRange(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	userID := createTestUser(t, testDB)
	chemID := createTestChemical(t, testDB, "lawn")

	now := time.Now()

	// Create form with applications 5 days ago (outside range)
	_, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "TooOld",
		LastName:     "Form",
		StreetNumber: "100",
		StreetName:   "Ancient St",
		Town:         "Town",
		ZipCode:      "10001",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-120 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create form in range (first app 3 days ago, last app 2 days ago)
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "InRange",
		LastName:     "Form",
		StreetNumber: "200",
		StreetName:   "Valid St",
		Town:         "Town",
		ZipCode:      "10002",
		HomePhone:    "555-0002",
		OtherPhone:   "555-0022",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 2000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-72 * time.Hour), // 3 days ago
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-48 * time.Hour), // 2 days ago
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create form too recent (first app 12 hours ago)
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "TooNew",
		LastName:     "Form",
		StreetNumber: "300",
		StreetName:   "Recent St",
		Town:         "Town",
		ZipCode:      "10003",
		HomePhone:    "555-0003",
		OtherPhone:   "555-0033",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 3000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-12 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Filter for forms with first app >= 4 days ago AND last app <= 1 day ago
	listOptions := ListFormsOptions{
		DateLow:  now.Add(-96 * time.Hour),  // 4 days ago
		DateHigh: now.Add(-24 * time.Hour),  // 1 day ago
		SortBy:   "created_at",
		Order:    "DESC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	// Should only get the InRange form
	require.Len(t, forms, 1)
	require.Equal(t, "InRange", getFirstName(forms[0]))
}

func TestListFormsByUserId_FilterByZipCode(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	userID := createTestUser(t, testDB)

	// Create forms with different zip codes
	_, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Alice",
		LastName:     "Boston",
		StreetNumber: "100",
		StreetName:   "Comm Ave",
		Town:         "Boston",
		ZipCode:      "02134",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
	})
	require.NoError(t, err)

	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Bob",
		LastName:     "Cambridge",
		StreetNumber: "200",
		StreetName:   "Mass Ave",
		Town:         "Cambridge",
		ZipCode:      "02139",
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
		FirstName:    "Carol",
		LastName:     "Boston2",
		StreetNumber: "300",
		StreetName:   "Beacon St",
		Town:         "Boston",
		ZipCode:      "02134",
		HomePhone:    "555-0003",
		OtherPhone:   "555-0033",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     true,
	})
	require.NoError(t, err)

	// Filter by zip code 02134
	listOptions := ListFormsOptions{
		ZipCode: "02134",
		SortBy:  "first_name",
		Order:   "ASC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	// Should get only Boston forms
	require.Len(t, forms, 2)
	require.Equal(t, "Alice", getFirstName(forms[0]))
	require.Equal(t, "Carol", getFirstName(forms[1]))
}

func TestListFormsByUserId_FilterByJewishHolidayYes(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	userID := createTestUser(t, testDB)

	// Create form with is_holiday = true
	_, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Holiday",
		LastName:     "Person",
		StreetNumber: "100",
		StreetName:   "Synagogue St",
		Town:         "Town",
		ZipCode:      "10001",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    true,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
	})
	require.NoError(t, err)

	// Create form with is_holiday = false
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Regular",
		LastName:     "Person",
		StreetNumber: "200",
		StreetName:   "Normal St",
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

	// Filter for jewish_holiday = yes
	listOptions := ListFormsOptions{
		JewishHoliday: "yes",
		SortBy:        "created_at",
		Order:         "DESC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	// Should only get holiday form
	require.Len(t, forms, 1)
	require.Equal(t, "Holiday", getFirstName(forms[0]))
}

func TestListFormsByUserId_FilterByJewishHolidayNo(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	userID := createTestUser(t, testDB)

	// Create form with is_holiday = true
	_, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Holiday",
		LastName:     "Person",
		StreetNumber: "100",
		StreetName:   "Synagogue St",
		Town:         "Town",
		ZipCode:      "10001",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    true,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
	})
	require.NoError(t, err)

	// Create form with is_holiday = false
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Regular",
		LastName:     "Person",
		StreetNumber: "200",
		StreetName:   "Normal St",
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

	// Filter for jewish_holiday = no
	listOptions := ListFormsOptions{
		JewishHoliday: "no",
		SortBy:        "created_at",
		Order:         "DESC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	// Should only get regular form
	require.Len(t, forms, 1)
	require.Equal(t, "Regular", getFirstName(forms[0]))
}

func TestListFormsByUserId_FilterByFormType(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	userID := createTestUser(t, testDB)

	// Create lawn form
	_, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Lawn",
		LastName:     "Owner",
		StreetNumber: "100",
		StreetName:   "Grass St",
		Town:         "Town",
		ZipCode:      "10001",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
	})
	require.NoError(t, err)

	// Create shrub form
	_, err = repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    userID,
		FirstName:    "Shrub",
		LastName:     "Owner",
		StreetNumber: "200",
		StreetName:   "Bush St",
		Town:         "Town",
		ZipCode:      "10002",
		HomePhone:    "555-0002",
		OtherPhone:   "555-0022",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     true,
	})
	require.NoError(t, err)

	// Filter for lawn forms only
	listOptions := ListFormsOptions{
		FormType: "lawn",
		SortBy:   "created_at",
		Order:    "DESC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	require.Len(t, forms, 1)
	require.Equal(t, "Lawn", getFirstName(forms[0]))
	require.Equal(t, "lawn", forms[0].FormType)

	// Filter for shrub forms only
	listOptions.FormType = "shrub"
	forms, err = repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	require.Len(t, forms, 1)
	require.Equal(t, "Shrub", getFirstName(forms[0]))
	require.Equal(t, "shrub", forms[0].FormType)
}

func TestListFormsByUserId_FilterBySearchName(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	userID := createTestUser(t, testDB)

	// Create forms with different names
	_, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Alice",
		LastName:     "Smith",
		StreetNumber: "100",
		StreetName:   "A St",
		Town:         "Town",
		ZipCode:      "10001",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
	})
	require.NoError(t, err)

	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Bob",
		LastName:     "Johnson",
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
		FirstName:    "Charlie",
		LastName:     "Smith",
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

	// Search by last name "smith"
	listOptions := ListFormsOptions{
		SearchName: "smith",
		SortBy:     "first_name",
		Order:      "ASC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	require.Len(t, forms, 2)
	require.Equal(t, "Alice", getFirstName(forms[0]))
	require.Equal(t, "Charlie", getFirstName(forms[1]))

	// Search by first name "bob" (case insensitive)
	listOptions.SearchName = "BoB"
	forms, err = repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	require.Len(t, forms, 1)
	require.Equal(t, "Bob", getFirstName(forms[0]))
}

func TestListFormsByUserId_FilterByChemicals(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	userID := createTestUser(t, testDB)
	chem1 := createTestChemical(t, testDB, "lawn")
	chem2 := createTestChemical(t, testDB, "lawn")
	chem3 := createTestChemical(t, testDB, "lawn")

	now := time.Now()

	// Create form with chem1
	_, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "User1",
		LastName:     "Form1",
		StreetNumber: "100",
		StreetName:   "St1",
		Town:         "Town",
		ZipCode:      "10001",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chem1,
				AppTimestamp:  now,
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create form with chem2
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "User2",
		LastName:     "Form2",
		StreetNumber: "200",
		StreetName:   "St2",
		Town:         "Town",
		ZipCode:      "10002",
		HomePhone:    "555-0002",
		OtherPhone:   "555-0022",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 2000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chem2,
				AppTimestamp:  now,
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create form with chem1 and chem3
	_, err = repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    userID,
		FirstName:    "User3",
		LastName:     "Form3",
		StreetNumber: "300",
		StreetName:   "St3",
		Town:         "Town",
		ZipCode:      "10003",
		HomePhone:    "555-0003",
		OtherPhone:   "555-0033",
		CallBefore:   false,
		IsHoliday:    false,
		FleaOnly:     true,
		Applications: []PestApp{
			{
				ChemUsed:      chem1,
				AppTimestamp:  now,
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
			{
				ChemUsed:      chem3,
				AppTimestamp:  now,
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Filter by chem1 only
	listOptions := ListFormsOptions{
		ChemicalIDs: []int{chem1},
		SortBy:      "first_name",
		Order:       "ASC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	require.Len(t, forms, 2)
	require.Equal(t, "User1", getFirstName(forms[0]))
	require.Equal(t, "User3", getFirstName(forms[1]))

	// Filter by chem1 OR chem2
	listOptions.ChemicalIDs = []int{chem1, chem2}
	forms, err = repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	require.Len(t, forms, 3)
	require.Equal(t, "User1", getFirstName(forms[0]))
	require.Equal(t, "User2", getFirstName(forms[1]))
	require.Equal(t, "User3", getFirstName(forms[2]))

	// Filter by chem3 only
	listOptions.ChemicalIDs = []int{chem3}
	forms, err = repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	require.Len(t, forms, 1)
	require.Equal(t, "User3", getFirstName(forms[0]))
}

func TestListFormsByUserId_CombinedFilters(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	userID := createTestUser(t, testDB)
	chemID := createTestChemical(t, testDB, "lawn")

	now := time.Now()

	// Create lawn form in Boston with holiday, recent application
	_, err := repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Match",
		LastName:     "AllFilters",
		StreetNumber: "100",
		StreetName:   "Perfect St",
		Town:         "Boston",
		ZipCode:      "02134",
		HomePhone:    "555-0001",
		OtherPhone:   "555-0011",
		CallBefore:   false,
		IsHoliday:    true,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-24 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create lawn form in Boston but NOT holiday
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Wrong",
		LastName:     "Holiday",
		StreetNumber: "200",
		StreetName:   "Wrong St",
		Town:         "Boston",
		ZipCode:      "02134",
		HomePhone:    "555-0002",
		OtherPhone:   "555-0022",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 2000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-24 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create shrub form (wrong type)
	_, err = repo.CreateShrubForm(ctx, CreateShrubFormInput{
		CreatedBy:    userID,
		FirstName:    "Wrong",
		LastName:     "Type",
		StreetNumber: "300",
		StreetName:   "Wrong St",
		Town:         "Boston",
		ZipCode:      "02134",
		HomePhone:    "555-0003",
		OtherPhone:   "555-0033",
		CallBefore:   false,
		IsHoliday:    true,
		FleaOnly:     true,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-24 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Create lawn form in Cambridge (wrong zip)
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    userID,
		FirstName:    "Wrong",
		LastName:     "Zip",
		StreetNumber: "400",
		StreetName:   "Wrong St",
		Town:         "Cambridge",
		ZipCode:      "02139",
		HomePhone:    "555-0004",
		OtherPhone:   "555-0044",
		CallBefore:   false,
		IsHoliday:    true,
		LawnAreaSqFt: 4000,
		FertOnly:     false,
		Applications: []PestApp{
			{
				ChemUsed:      chemID,
				AppTimestamp:  now.Add(-24 * time.Hour),
				Rate:          "2 oz/1000 sq ft",
				AmountApplied: decimal.NewFromFloat(2.0),
				LocationCode:  "FL",
			},
		},
	})
	require.NoError(t, err)

	// Apply all filters: lawn forms, zip 02134, holiday=yes, with chemID, sorted by first_app_date
	listOptions := ListFormsOptions{
		FormType:      "lawn",
		ZipCode:       "02134",
		JewishHoliday: "yes",
		ChemicalIDs:   []int{chemID},
		DateLow:       now.Add(-48 * time.Hour),
		SortBy:        "first_app_date",
		Order:         "DESC",
	}
	forms, err := repo.ListFormsByUserId(ctx, userID, listOptions)
	require.NoError(t, err)

	// Should only match the first form
	require.Len(t, forms, 1)
	require.Equal(t, "Match", getFirstName(forms[0]))
	require.Equal(t, "AllFilters", getLastName(forms[0]))
}

func TestListAllForms_WithFilters(t *testing.T) {
	ctx := context.Background()
	testDB := db.TestDB(t)
	repo := NewFormsRepository(testDB)

	// Create two users
	user1ID := createTestUser(t, testDB)

	var user2ID string
	err := testDB.QueryRow(`
		INSERT INTO users (first_name, last_name, username, password_hash)
		VALUES ('Test', 'User', 'TestUserName_ListAll', 'TestPass')
		RETURNING id
	`).Scan(&user2ID)
	require.NoError(t, err)

	// User 1 creates a lawn form in Boston
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    user1ID,
		FirstName:    "User1",
		LastName:     "Boston",
		StreetNumber: "100",
		StreetName:   "User1 St",
		Town:         "Boston",
		ZipCode:      "02134",
		HomePhone:    "555-1111",
		OtherPhone:   "555-1112",
		CallBefore:   false,
		IsHoliday:    true,
		LawnAreaSqFt: 1000,
		FertOnly:     false,
	})
	require.NoError(t, err)

	// User 2 creates a lawn form in Cambridge
	_, err = repo.CreateLawnForm(ctx, CreateLawnFormInput{
		CreatedBy:    user2ID,
		FirstName:    "User2",
		LastName:     "Cambridge",
		StreetNumber: "200",
		StreetName:   "User2 St",
		Town:         "Cambridge",
		ZipCode:      "02139",
		HomePhone:    "555-2222",
		OtherPhone:   "555-2223",
		CallBefore:   false,
		IsHoliday:    false,
		LawnAreaSqFt: 2000,
		FertOnly:     false,
	})
	require.NoError(t, err)

	// ListAllForms should return both (no user filter)
	listOptions := ListFormsOptions{
		SortBy: "first_name",
		Order:  "ASC",
	}
	forms, err := repo.ListAllForms(ctx, listOptions)
	require.NoError(t, err)
	require.Len(t, forms, 2)

	// Filter by zip code in ListAllForms
	listOptions.ZipCode = "02134"
	forms, err = repo.ListAllForms(ctx, listOptions)
	require.NoError(t, err)
	require.Len(t, forms, 1)
	require.Equal(t, "User1", getFirstName(forms[0]))

	// Filter by jewish_holiday = no in ListAllForms
	listOptions.ZipCode = ""
	listOptions.JewishHoliday = "no"
	forms, err = repo.ListAllForms(ctx, listOptions)
	require.NoError(t, err)
	require.Len(t, forms, 1)
	require.Equal(t, "User2", getFirstName(forms[0]))
}
