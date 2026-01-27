// Package forms provides data access and domain models for landscape forms.
// It encapsulates persistence logic, enforces ownership rules, and ensures
// type-safe access to shrub and lawn forms.
package forms

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// FormsRepository provides database access for form records.
// All methods enforce ownership at the SQL layer and return sql.ErrNoRows
// when a form does not exist or is not owned by the given user.
type FormsRepository struct {
	db *sql.DB
}

// NewFormsRepository returns a repository backed by the given database connection.
func NewFormsRepository(database *sql.DB) *FormsRepository {
	return &FormsRepository{db: database}
}

// CreateFormInput contains the common fields required to create a new form.
type CreateShrubFormInput struct {
	CreatedBy    string
	FirstName    string
	LastName     string
	StreetNumber string
	StreetName   string
	Town         string
	ZipCode      string
	HomePhone    string
	OtherPhone   string
	CallBefore   bool
	IsHoliday    bool
	FleaOnly     bool
	Applications []PestApp
}
type CreateLawnFormInput struct {
	CreatedBy    string
	FirstName    string
	LastName     string
	StreetNumber string
	StreetName   string
	Town         string
	ZipCode      string
	HomePhone    string
	OtherPhone   string
	CallBefore   bool
	IsHoliday    bool
	LawnAreaSqFt int
	FertOnly     bool
	Applications []PestApp
}

// UpdateFormInput contains the fields that may be updated on an existing form.
type UpdateShrubFormInput struct {
	FirstName    string
	LastName     string
	StreetNumber string
	StreetName   string
	Town         string
	ZipCode      string
	HomePhone    string
	OtherPhone   string
	CallBefore   bool
	IsHoliday    bool
	FleaOnly     bool
}
type UpdateLawnFormInput struct {
	FirstName    string
	LastName     string
	StreetNumber string
	StreetName   string
	Town         string
	ZipCode      string
	HomePhone    string
	OtherPhone   string
	CallBefore   bool
	IsHoliday    bool
	LawnAreaSqFt int
	FertOnly     bool
}

// CreateShrubForm creates a new shrub form and its associated shrub details.
// Returns the created form's ID upon success
// The operation is atomic and will fail if shrub details are not provided.
func (r *FormsRepository) CreateShrubForm(
	ctx context.Context,
	shrubFormInput CreateShrubFormInput,
) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var formID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO forms (
			created_by,
			form_type,
			first_name,
			last_name,
			street_number,
			street_name,
			town,
			zip_code,
			home_phone,
			other_phone,
			call_before,
			is_holiday
		)
		VALUES ($1, 'shrub', $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`,
		shrubFormInput.CreatedBy,
		shrubFormInput.FirstName,
		shrubFormInput.LastName,
		shrubFormInput.StreetNumber,
		shrubFormInput.StreetName,
		shrubFormInput.Town,
		shrubFormInput.ZipCode,
		shrubFormInput.HomePhone,
		shrubFormInput.OtherPhone,
		shrubFormInput.CallBefore,
		shrubFormInput.IsHoliday,
	).Scan(
		&formID,
	)
	if err != nil {
		return "", fmt.Errorf("Failed to insert form: %s %s, %w", shrubFormInput.FirstName, shrubFormInput.LastName, err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO shrub_forms (
			form_id,
			flea_only
		)
		VALUES ($1, $2)
	`,
		formID,
		shrubFormInput.FleaOnly,
	)
	if err != nil {
		return "", fmt.Errorf("Failed to insert shrub form: %s %s, %w", shrubFormInput.FirstName, shrubFormInput.LastName, err)
	}

	// Insert pesticide applications if any
	for _, app := range shrubFormInput.Applications {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO pesticide_applications (
				form_id,
				chem_used,
				app_timestamp,
				rate,
				amount_applied,
				location_code
			)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			formID,
			app.ChemUsed,
			app.AppTimestamp,
			app.Rate,
			app.AmountApplied,
			app.LocationCode,
		)
		if err != nil {
			return "", fmt.Errorf("Failed to insert pesticide application for form %s %s: %w", shrubFormInput.FirstName, shrubFormInput.LastName, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("Failed to commit transaction for inserting shrub form: %s %s, %w", shrubFormInput.FirstName, shrubFormInput.LastName, err)
	}

	return formID, nil
}

// CreateLawnForm creates a new lawn form and its associated lawn details.
// Returns the created form's ID upon success
// The operation is atomic and will fail if lawn details are not provided.
func (r *FormsRepository) CreateLawnForm(
	ctx context.Context,
	lawnFormInput CreateLawnFormInput,
) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var formID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO forms (
			created_by,
			form_type,
			first_name,
			last_name,
			street_number,
			street_name,
			town,
			zip_code,
			home_phone,
			other_phone,
			call_before,
			is_holiday
		)
		VALUES ($1, 'lawn', $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`,
		lawnFormInput.CreatedBy,
		lawnFormInput.FirstName,
		lawnFormInput.LastName,
		lawnFormInput.StreetNumber,
		lawnFormInput.StreetName,
		lawnFormInput.Town,
		lawnFormInput.ZipCode,
		lawnFormInput.HomePhone,
		lawnFormInput.OtherPhone,
		lawnFormInput.CallBefore,
		lawnFormInput.IsHoliday,
	).Scan(
		&formID,
	)

	if err != nil {
		return "", fmt.Errorf("Failed to insert form: %s %s, %w", lawnFormInput.FirstName, lawnFormInput.LastName, err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO lawn_forms (
			form_id,
			lawn_area_sq_ft,
			fert_only
		)
		VALUES ($1, $2, $3)
	`,
		formID,
		lawnFormInput.LawnAreaSqFt,
		lawnFormInput.FertOnly,
	)
	if err != nil {
		return "", fmt.Errorf("Failed to insert lawn form: %s %s, %w", lawnFormInput.FirstName, lawnFormInput.LastName, err)
	}

	// Insert pesticide applications if any
	for _, app := range lawnFormInput.Applications {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO pesticide_applications (
				form_id,
				chem_used,
				app_timestamp,
				rate,
				amount_applied,
				location_code
			)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			formID,
			app.ChemUsed,
			app.AppTimestamp,
			app.Rate,
			app.AmountApplied,
			app.LocationCode,
		)
		if err != nil {
			return "", fmt.Errorf("Failed to insert pesticide application for form %s %s: %w", lawnFormInput.FirstName, lawnFormInput.LastName, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("Failed to commit transaction for inserting lawn form: %s %s, %w", lawnFormInput.FirstName, lawnFormInput.LastName, err)
	}

	return formID, nil
}

// ListFormsOptions contains optional filtering and pagination parameters
type ListFormsOptions struct {
	// Pagination
	Limit  int // Max number of results (0 = no limit)
	Offset int // Number of results to skip

	// Filtering
	FormType   string // Filter by form type: "shrub" or "lawn" (empty = all)
	SearchName string // Search in first_name or last_name (partial match)

	// Sorting
	SortBy string // "first_name", "last_name", or "created_at" (defaults to "created_at")
	Order  string // "ASC" or "DESC" (defaults to "DESC")
}

// ListFormsByUserId returns all forms owned by the given user with pagination and filtering.
// Results may be sorted by first name, last name, or creation time.
// Each returned FormView is fully hydrated with its subtype details.
func (r *FormsRepository) ListFormsByUserId(
	ctx context.Context,
	userID string,
	opts ListFormsOptions,
) ([]*FormView, error) {

	allowedSorts := map[string]string{
		"first_name": "f.first_name",
		"last_name":  "f.last_name",
		"created_at": "f.created_at",
	}

	sortColumn, ok := allowedSorts[opts.SortBy]
	if !ok {
		sortColumn = "f.created_at"
	}

	order := strings.ToUpper(opts.Order)
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	// Build WHERE clause
	whereConditions := []string{"f.created_by = $1"}
	args := []any{userID}
	argIndex := 2

	// Add form type filter
	if opts.FormType != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("f.form_type = $%d", argIndex))
		args = append(args, opts.FormType)
		argIndex++
	}

	// Add name search filter
	if opts.SearchName != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(f.first_name ILIKE $%d OR f.last_name ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+opts.SearchName+"%")
		argIndex++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Build query with pagination
	query := fmt.Sprintf(`
		SELECT
			f.id,
			f.created_by,
			f.created_at,
			f.form_type,
			f.updated_at,
			f.first_name,
			f.last_name,
			f.street_number,
			f.street_name,
			f.town,
			f.zip_code,
			f.home_phone,
			f.other_phone,
			f.call_before,
			f.is_holiday,
			sf.flea_only,
			lf.lawn_area_sq_ft,
			lf.fert_only
		FROM forms f
		LEFT JOIN shrub_forms sf ON f.id = sf.form_id
		LEFT JOIN lawn_forms lf ON f.id = lf.form_id
		WHERE %s
		ORDER BY %s %s
	`, whereClause, sortColumn, order)

	// Add pagination
	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, opts.Limit)
		argIndex++
	}
	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, opts.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows for forms list: %w", err)
	}
	defer rows.Close()

	var forms []*FormView
	for rows.Next() {
		var (
			form  Form
			shrub shrubRow
			lawn  lawnRow
		)

		err := rows.Scan(
			&form.ID,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.FormType,
			&form.UpdatedAt,
			&form.FirstName,
			&form.LastName,
			&form.StreetNumber,
			&form.StreetName,
			&form.Town,
			&form.ZipCode,
			&form.HomePhone,
			&form.OtherPhone,
			&form.CallBefore,
			&form.IsHoliday,
			&shrub.FleaOnly,
			&lawn.LawnAreaSqFt,
			&lawn.FertOnly,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		var view *FormView
		switch form.FormType {
		case "shrub":
			shrubDetails, err := shrub.ToDomain()
			if err != nil {
				return nil, fmt.Errorf("error casting row to shrub form %w", err)
			}
			view = NewShrubFormView(
				ShrubForm{
					Form:         form,
					ShrubDetails: shrubDetails,
				},
			)

		case "lawn":
			lawnDetails, err := lawn.ToDomain()
			if err != nil {
				return nil, fmt.Errorf("error casting row to lawn form: %w", err)
			}
			view = NewLawnFormView(
				LawnForm{
					Form:        form,
					LawnDetails: lawnDetails,
				},
			)
		default:
			return nil, fmt.Errorf("unknown form_type: %s", form.FormType)
		}
		forms = append(forms, view)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after list forms queries: %w", err)
	}

	return forms, nil
}

// ListAllForms returns all forms (admin only) with pagination and filtering.
// Does NOT filter by user - returns forms from all users.
// Each returned FormView is fully hydrated with its subtype details.
func (r *FormsRepository) ListAllForms(
	ctx context.Context,
	opts ListFormsOptions,
) ([]*FormView, error) {

	allowedSorts := map[string]string{
		"first_name": "f.first_name",
		"last_name":  "f.last_name",
		"created_at": "f.created_at",
	}

	sortColumn, ok := allowedSorts[opts.SortBy]
	if !ok {
		sortColumn = "f.created_at"
	}

	order := strings.ToUpper(opts.Order)
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	// Build WHERE clause
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	// Add form type filter
	if opts.FormType != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("f.form_type = $%d", argIndex))
		args = append(args, opts.FormType)
		argIndex++
	}

	// Add name search filter
	if opts.SearchName != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(f.first_name ILIKE $%d OR f.last_name ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+opts.SearchName+"%")
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Build query with pagination
	query := fmt.Sprintf(`
		SELECT
			f.id,
			f.created_by,
			f.created_at,
			f.form_type,
			f.updated_at,
			f.first_name,
			f.last_name,
			f.street_number,
			f.street_name,
			f.town,
			f.zip_code,
			f.home_phone,
			f.other_phone,
			f.call_before,
			f.is_holiday,
			sf.flea_only,
			lf.lawn_area_sq_ft,
			lf.fert_only
		FROM forms f
		LEFT JOIN shrub_forms sf ON f.id = sf.form_id
		LEFT JOIN lawn_forms lf ON f.id = lf.form_id
		%s
		ORDER BY %s %s
	`, whereClause, sortColumn, order)

	// Add pagination
	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, opts.Limit)
		argIndex++
	}
	if opts.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, opts.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying rows for forms list: %w", err)
	}
	defer rows.Close()

	var forms []*FormView
	for rows.Next() {
		var (
			form  Form
			shrub shrubRow
			lawn  lawnRow
		)

		err := rows.Scan(
			&form.ID,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.FormType,
			&form.UpdatedAt,
			&form.FirstName,
			&form.LastName,
			&form.StreetNumber,
			&form.StreetName,
			&form.Town,
			&form.ZipCode,
			&form.HomePhone,
			&form.OtherPhone,
			&form.CallBefore,
			&form.IsHoliday,
			&shrub.FleaOnly,
			&lawn.LawnAreaSqFt,
			&lawn.FertOnly,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		var view *FormView
		switch form.FormType {
		case "shrub":
			shrubDetails, err := shrub.ToDomain()
			if err != nil {
				return nil, fmt.Errorf("error casting row to shrub form: %w", err)
			}
			view = NewShrubFormView(
				ShrubForm{
					Form:         form,
					ShrubDetails: shrubDetails,
				},
			)

		case "lawn":
			lawnDetails, err := lawn.ToDomain()
			if err != nil {
				return nil, fmt.Errorf("error casting row to lawn form: %w", err)
			}
			view = NewLawnFormView(
				LawnForm{
					Form:        form,
					LawnDetails: lawnDetails,
				},
			)
		default:
			return nil, fmt.Errorf("unknown form_type: %s", form.FormType)
		}
		forms = append(forms, view)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after queries for forms list: %w", err)
	}

	return forms, nil
}

// GetFormViewById returns a single form owned by the given user.
// It returns sql.ErrNoRows if the form does not exist or is not owned by the user.
func (r *FormsRepository) GetFormViewById(
	ctx context.Context,
	formID string,
	userID string,
) (*FormView, error) {

	query := `
		SELECT
			f.id,
			f.created_by,
			f.created_at,
			f.form_type,
			f.updated_at,
			f.first_name,
			f.last_name,
			f.street_number,
			f.street_name,
			f.town,
			f.zip_code,
			f.home_phone,
			f.other_phone,
			f.call_before,
			f.is_holiday,
			sf.flea_only,
			lf.lawn_area_sq_ft,
			lf.fert_only
		FROM forms f
		LEFT JOIN shrub_forms sf ON f.id = sf.form_id
		LEFT JOIN lawn_forms lf ON f.id = lf.form_id
		WHERE f.id = $1
		  AND f.created_by = $2
	`

	var (
		form  Form
		shrub shrubRow
		lawn  lawnRow
	)

	err := r.db.QueryRowContext(ctx, query, formID, userID).Scan(
		&form.ID,
		&form.CreatedBy,
		&form.CreatedAt,
		&form.FormType,
		&form.UpdatedAt,
		&form.FirstName,
		&form.LastName,
		&form.StreetNumber,
		&form.StreetName,
		&form.Town,
		&form.ZipCode,
		&form.HomePhone,
		&form.OtherPhone,
		&form.CallBefore,
		&form.IsHoliday,
		&shrub.FleaOnly,
		&lawn.LawnAreaSqFt,
		&lawn.FertOnly,
	)
	if err != nil {
		// Important: let sql.ErrNoRows propagate
		return nil, err
	}

	var view *FormView
	switch form.FormType {
	case "shrub":
		shrubDetails, err := shrub.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("error casting row to shrub form: %w", err)
		}
		view = NewShrubFormView(
			ShrubForm{
				Form:         form,
				ShrubDetails: shrubDetails,
			},
		)

	case "lawn":
		lawnDetails, err := lawn.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("error casting row to lawn form: %w", err)
		}
		view = NewLawnFormView(
			LawnForm{
				Form:        form,
				LawnDetails: lawnDetails,
			},
		)
	default:
		return nil, fmt.Errorf("unknown form_type: %s", form.FormType)
	}

	return view, nil
}

// GetShrubFormById returns a single shrub form owned by the given user.
// It returns sql.ErrNoRows if the form does not exist or is not owned by the user.
func (r *FormsRepository) GetShrubFormById(
	ctx context.Context,
	formID string,
	userID string,
) (ShrubForm, error) {

	query := `
		SELECT
			f.id,
			f.created_by,
			f.created_at,
			f.form_type,
			f.updated_at,
			f.first_name,
			f.last_name,
			f.street_number,
			f.street_name,
			f.town,
			f.zip_code,
			f.home_phone,
			f.other_phone,
			f.call_before,
			f.is_holiday,
			sf.flea_only
		FROM forms f
		LEFT JOIN shrub_forms sf ON f.id = sf.form_id
		WHERE f.id = $1
		  AND f.created_by = $2
	`

	var shrubForm ShrubForm

	err := r.db.QueryRowContext(ctx, query, formID, userID).Scan(
		&shrubForm.ID,
		&shrubForm.CreatedBy,
		&shrubForm.CreatedAt,
		&shrubForm.FormType,
		&shrubForm.UpdatedAt,
		&shrubForm.FirstName,
		&shrubForm.LastName,
		&shrubForm.StreetNumber,
		&shrubForm.StreetName,
		&shrubForm.Town,
		&shrubForm.ZipCode,
		&shrubForm.HomePhone,
		&shrubForm.OtherPhone,
		&shrubForm.CallBefore,
		&shrubForm.IsHoliday,
		&shrubForm.FleaOnly,
	)
	if err != nil {
		// Important: let sql.ErrNoRows propagate
		return ShrubForm{}, err
	}

	return shrubForm, nil
}

// GetLawnFormById returns a single lawn form owned by the given user.
// It returns sql.ErrNoRows if the form does not exist or is not owned by the user.
func (r *FormsRepository) GetLawnFormById(
	ctx context.Context,
	formID string,
	userID string,
) (LawnForm, error) {

	query := `
		SELECT
			f.id,
			f.created_by,
			f.created_at,
			f.form_type,
			f.updated_at,
			f.first_name,
			f.last_name,
			f.street_number,
			f.street_name,
			f.town,
			f.zip_code,
			f.home_phone,
			f.other_phone,
			f.call_before,
			f.is_holiday,
			lf.lawn_area_sq_ft,
			lf.fert_only
		FROM forms f
		LEFT JOIN lawn_forms lf ON f.id = lf.form_id
		WHERE f.id = $1
		  AND f.created_by = $2
	`

	var lawnForm LawnForm

	err := r.db.QueryRowContext(ctx, query, formID, userID).Scan(
		&lawnForm.ID,
		&lawnForm.CreatedBy,
		&lawnForm.CreatedAt,
		&lawnForm.FormType,
		&lawnForm.UpdatedAt,
		&lawnForm.FirstName,
		&lawnForm.LastName,
		&lawnForm.StreetNumber,
		&lawnForm.StreetName,
		&lawnForm.Town,
		&lawnForm.ZipCode,
		&lawnForm.HomePhone,
		&lawnForm.OtherPhone,
		&lawnForm.CallBefore,
		&lawnForm.IsHoliday,
		&lawnForm.LawnAreaSqFt,
		&lawnForm.FertOnly,
	)
	if err != nil {
		// Important: let sql.ErrNoRows propagate
		return LawnForm{}, err
	}

	return lawnForm, nil
}

// UpdateShrubFormById updates a shrub form
// It returns sql.ErrNoRows if the form does not exist or is not owned by the user.
func (r *FormsRepository) UpdateShrubFormById(
	ctx context.Context,
	formID string,
	userID string,
	shrubFormInput UpdateShrubFormInput,
) (ShrubForm, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return ShrubForm{}, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	var shrubForm ShrubForm

	err = tx.QueryRowContext(ctx, `
		UPDATE forms
		SET first_name = $1,
			last_name = $2,
			street_number = $3,
			street_name = $4,
			town = $5,
			zip_code = $6,
			home_phone = $7,
			other_phone = $8,
			call_before = $9,
			is_holiday = $10
		WHERE id = $11 AND created_by = $12
		RETURNING
			id,
			created_by,
			created_at,
			form_type,
			updated_at,
			first_name,
			last_name,
			street_number,
			street_name,
			town,
			zip_code,
			home_phone,
			other_phone,
			call_before,
			is_holiday
	`,
		shrubFormInput.FirstName,
		shrubFormInput.LastName,
		shrubFormInput.StreetNumber,
		shrubFormInput.StreetName,
		shrubFormInput.Town,
		shrubFormInput.ZipCode,
		shrubFormInput.HomePhone,
		shrubFormInput.OtherPhone,
		shrubFormInput.CallBefore,
		shrubFormInput.IsHoliday,
		formID,
		userID,
	).Scan(
		&shrubForm.ID,
		&shrubForm.CreatedBy,
		&shrubForm.CreatedAt,
		&shrubForm.FormType,
		&shrubForm.UpdatedAt,
		&shrubForm.FirstName,
		&shrubForm.LastName,
		&shrubForm.StreetNumber,
		&shrubForm.StreetName,
		&shrubForm.Town,
		&shrubForm.ZipCode,
		&shrubForm.HomePhone,
		&shrubForm.OtherPhone,
		&shrubForm.CallBefore,
		&shrubForm.IsHoliday,
	)
	if err != nil {
		//sql.ErrNoRows
		return shrubForm, err
	}

	var query string
	query = `
		UPDATE shrub_forms
		SET flea_only = $1
		WHERE form_id = $2
		RETURNING flea_only
	`
	err = tx.QueryRowContext(ctx, query, shrubFormInput.FleaOnly, formID).Scan(
		&shrubForm.FleaOnly,
	)
	if err != nil {
		//sql.ErrNoRows
		return ShrubForm{}, err
	}

	if err := tx.Commit(); err != nil {
		return ShrubForm{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return shrubForm, nil
}

// UpdateLawnFormById updates a lawn form
// It returns sql.ErrNoRows if the form does not exist or is not owned by the user.
func (r *FormsRepository) UpdateLawnFormById(
	ctx context.Context,
	formID string,
	userID string,
	lawnFormInput UpdateLawnFormInput,
) (LawnForm, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return LawnForm{}, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	var lawnForm LawnForm

	err = tx.QueryRowContext(ctx, `
		UPDATE forms
		SET first_name = $1,
			last_name = $2,
			street_number = $3,
			street_name = $4,
			town = $5,
			zip_code = $6,
			home_phone = $7,
			other_phone = $8,
			call_before = $9,
			is_holiday = $10
		WHERE id = $11 AND created_by = $12
		RETURNING
			id,
			created_by,
			created_at,
			form_type,
			updated_at,
			first_name,
			last_name,
			street_number,
			street_name,
			town,
			zip_code,
			home_phone,
			other_phone,
			call_before,
			is_holiday
	`,
		lawnFormInput.FirstName,
		lawnFormInput.LastName,
		lawnFormInput.StreetNumber,
		lawnFormInput.StreetName,
		lawnFormInput.Town,
		lawnFormInput.ZipCode,
		lawnFormInput.HomePhone,
		lawnFormInput.OtherPhone,
		lawnFormInput.CallBefore,
		lawnFormInput.IsHoliday,
		formID,
		userID,
	).Scan(
		&lawnForm.ID,
		&lawnForm.CreatedBy,
		&lawnForm.CreatedAt,
		&lawnForm.FormType,
		&lawnForm.UpdatedAt,
		&lawnForm.FirstName,
		&lawnForm.LastName,
		&lawnForm.StreetNumber,
		&lawnForm.StreetName,
		&lawnForm.Town,
		&lawnForm.ZipCode,
		&lawnForm.HomePhone,
		&lawnForm.OtherPhone,
		&lawnForm.CallBefore,
		&lawnForm.IsHoliday,
	)
	if err != nil {
		//sql.ErrNoRows
		return LawnForm{}, err
	}

	var query string
	query = `
		UPDATE lawn_forms
		SET lawn_area_sq_ft = $1,
			fert_only = $2
		WHERE form_id = $3
		RETURNING lawn_area_sq_ft, fert_only
	`
	err = tx.QueryRowContext(ctx, query, lawnFormInput.LawnAreaSqFt, lawnFormInput.FertOnly, formID).Scan(
		&lawnForm.LawnAreaSqFt,
		&lawnForm.FertOnly,
	)
	if err != nil {
		//sql.ErrNoRows
		return LawnForm{}, err
	}

	if err := tx.Commit(); err != nil {
		return LawnForm{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return lawnForm, nil
}

// DeleteFormById deletes a form owned by the given user.
// Associated subtype records are removed via ON DELETE CASCADE.
// It returns sql.ErrNoRows if the form does not exist or is not owned by the user.
func (r *FormsRepository) DeleteFormById(
	ctx context.Context,
	formID string,
	userID string,
) error {

	err := r.db.QueryRowContext(ctx, `
		DELETE FROM forms
		WHERE id = $1 AND created_by = $2
		RETURNING id
	`, formID, userID).Scan(&formID)

	if err != nil {
		// sql.ErrNoRows → not found or not owned
		return err
	}

	return nil
}
