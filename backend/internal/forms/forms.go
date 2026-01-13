// Package forms provides data access and domain models for landscape forms.
// It encapsulates persistence logic, enforces ownership rules, and ensures
// type-safe access to shrub and pesticide forms.
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
	CreatedBy string
	FirstName string
	LastName  string
	HomePhone string
	NumShrubs int
}
type CreatePesticideFormInput struct {
	CreatedBy     string
	FirstName     string
	LastName      string
	HomePhone     string
	PesticideName string
}

// UpdateFormInput contains the fields that may be updated on an existing form.
type UpdateShrubFormInput struct {
	FirstName string
	LastName  string
	HomePhone string
	NumShrubs int
}
type UpdatePesticideFormInput struct {
	FirstName     string
	LastName      string
	HomePhone     string
	PesticideName string
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
			home_phone
		)
		VALUES ($1, 'shrub', $2, $3, $4)
		RETURNING id
	`,
		shrubFormInput.CreatedBy,
		shrubFormInput.FirstName,
		shrubFormInput.LastName,
		shrubFormInput.HomePhone,
	).Scan(
		&formID,
	)
	if err != nil {
		return "", fmt.Errorf("Failed to insert form: %s %s, %w", shrubFormInput.FirstName, shrubFormInput.LastName, err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO shrubs (
			form_id,
			num_shrubs
		)
		VALUES ($1, $2)
	`,
		formID,
		shrubFormInput.NumShrubs,
	)
	if err != nil {
		return "", fmt.Errorf("Failed to insert shrub form: %s %s, %w", shrubFormInput.FirstName, shrubFormInput.LastName, err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("Failed to commit transaction for inserting shrub form: %s %s, %w", shrubFormInput.FirstName, shrubFormInput.LastName, err)
	}

	return formID, nil
}

// CreatePesticideForm creates a new pesticide form and its associated pesticide details.
// Returns the created form's ID upon success
// The operation is atomic and will fail if pesticide details are not provided.
func (r *FormsRepository) CreatePesticideForm(
	ctx context.Context,
	pesticideFormInput CreatePesticideFormInput,
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
			home_phone
		)
		VALUES ($1, 'pesticide', $2, $3, $4)
		RETURNING id
	`,
		pesticideFormInput.CreatedBy,
		pesticideFormInput.FirstName,
		pesticideFormInput.LastName,
		pesticideFormInput.HomePhone,
	).Scan(
		&formID,
	)

	if err != nil {
		return "", fmt.Errorf("Failed to insert form: %s %s, %w", pesticideFormInput.FirstName, pesticideFormInput.LastName, err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO pesticides (
			form_id,
			pesticide_name	
		)
		VALUES ($1, $2)
	`,
		formID,
		pesticideFormInput.PesticideName,
	)
	if err != nil {
		return "", fmt.Errorf("Failed to insert pesticide form: %s %s, %w", pesticideFormInput.FirstName, pesticideFormInput.LastName, err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("Failed to commit transaction for inserting pesticide form: %s %s, %w", pesticideFormInput.FirstName, pesticideFormInput.LastName, err)
	}

	return formID, nil
}

// ListFormsOptions contains optional filtering and pagination parameters
type ListFormsOptions struct {
	// Pagination
	Limit  int // Max number of results (0 = no limit)
	Offset int // Number of results to skip

	// Filtering
	FormType   string // Filter by form type: "shrub" or "pesticide" (empty = all)
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
			f.home_phone,
			s.num_shrubs,
			p.pesticide_name
		FROM forms f
		LEFT JOIN shrubs s ON f.id = s.form_id
		LEFT JOIN pesticides p ON f.id = p.form_id
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
			form      Form
			shrub     shrubRow
			pesticide pesticideRow
		)

		err := rows.Scan(
			&form.ID,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.FormType,
			&form.UpdatedAt,
			&form.FirstName,
			&form.LastName,
			&form.HomePhone,
			&shrub.NumShrubs,
			&pesticide.PesticideName,
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

		case "pesticide":
			pesticideDetails, err := pesticide.ToDomain()
			if err != nil {
				return nil, fmt.Errorf("error casting row to pesticide form: %w", err)
			}
			view = NewPesticideFormView(
				PesticideForm{
					Form:             form,
					PesticideDetails: pesticideDetails,
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
			f.home_phone,
			s.num_shrubs,
			p.pesticide_name
		FROM forms f
		LEFT JOIN shrubs s ON f.id = s.form_id
		LEFT JOIN pesticides p ON f.id = p.form_id
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
			form      Form
			shrub     shrubRow
			pesticide pesticideRow
		)

		err := rows.Scan(
			&form.ID,
			&form.CreatedBy,
			&form.CreatedAt,
			&form.FormType,
			&form.UpdatedAt,
			&form.FirstName,
			&form.LastName,
			&form.HomePhone,
			&shrub.NumShrubs,
			&pesticide.PesticideName,
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

		case "pesticide":
			pesticideDetails, err := pesticide.ToDomain()
			if err != nil {
				return nil, fmt.Errorf("error casting row to pesticide form: %w", err)
			}
			view = NewPesticideFormView(
				PesticideForm{
					Form:             form,
					PesticideDetails: pesticideDetails,
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
			f.home_phone,
			s.num_shrubs,
			p.pesticide_name
		FROM forms f
		LEFT JOIN shrubs s ON f.id = s.form_id
		LEFT JOIN pesticides p ON f.id = p.form_id
		WHERE f.id = $1
		  AND f.created_by = $2
	`

	var (
		form      Form
		shrub     shrubRow
		pesticide pesticideRow
	)

	err := r.db.QueryRowContext(ctx, query, formID, userID).Scan(
		&form.ID,
		&form.CreatedBy,
		&form.CreatedAt,
		&form.FormType,
		&form.UpdatedAt,
		&form.FirstName,
		&form.LastName,
		&form.HomePhone,
		&shrub.NumShrubs,
		&pesticide.PesticideName,
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

	case "pesticide":
		pesticideDetails, err := pesticide.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("error casting row to pesticide form: %w", err)
		}
		view = NewPesticideFormView(
			PesticideForm{
				Form:             form,
				PesticideDetails: pesticideDetails,
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
			f.home_phone,
			s.num_shrubs
		FROM forms f
		LEFT JOIN shrubs s ON f.id = s.form_id
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
		&shrubForm.HomePhone,
		&shrubForm.NumShrubs,
	)
	if err != nil {
		// Important: let sql.ErrNoRows propagate
		return ShrubForm{}, err
	}

	return shrubForm, nil
}

// GetPesticideFormById returns a single pesticide form owned by the given user.
// It returns sql.ErrNoRows if the form does not exist or is not owned by the user.
func (r *FormsRepository) GetPesticideFormById(
	ctx context.Context,
	formID string,
	userID string,
) (PesticideForm, error) {

	query := `
		SELECT
			f.id,
			f.created_by,
			f.created_at,
			f.form_type,
			f.updated_at,
			f.first_name,
			f.last_name,
			f.home_phone,
			p.pesticide_name
		FROM forms f
		LEFT JOIN pesticides p ON f.id = p.form_id
		WHERE f.id = $1
		  AND f.created_by = $2
	`

	var pesticideForm PesticideForm

	err := r.db.QueryRowContext(ctx, query, formID, userID).Scan(
		&pesticideForm.ID,
		&pesticideForm.CreatedBy,
		&pesticideForm.CreatedAt,
		&pesticideForm.FormType,
		&pesticideForm.UpdatedAt,
		&pesticideForm.FirstName,
		&pesticideForm.LastName,
		&pesticideForm.HomePhone,
		&pesticideForm.PesticideName,
	)
	if err != nil {
		// Important: let sql.ErrNoRows propagate
		return PesticideForm{}, err
	}

	return pesticideForm, nil
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
			home_phone = $3
		WHERE id = $4 AND created_by = $5
		RETURNING 
			id, 
			created_by, 
			created_at,
			form_type,
			updated_at, 
			first_name, 
			last_name, 
			home_phone
	`,
		shrubFormInput.FirstName,
		shrubFormInput.LastName,
		shrubFormInput.HomePhone,
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
		&shrubForm.HomePhone,
	)
	if err != nil {
		//sql.ErrNoRows
		return shrubForm, err
	}

	var query string
	query = `
		UPDATE shrubs
		SET num_shrubs = $1
		WHERE form_id = $2
		RETURNING num_shrubs
	`
	err = tx.QueryRowContext(ctx, query, shrubFormInput.NumShrubs, formID).Scan(
		&shrubForm.NumShrubs,
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

// UpdatePesticideFormById updates a pesticide form
// It returns sql.ErrNoRows if the form does not exist or is not owned by the user.
func (r *FormsRepository) UpdatePesticideFormById(
	ctx context.Context,
	formID string,
	userID string,
	pesticideFormInput UpdatePesticideFormInput,
) (PesticideForm, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return PesticideForm{}, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	var pesticideForm PesticideForm

	err = tx.QueryRowContext(ctx, `
		UPDATE forms
		SET first_name = $1,
			last_name = $2,
			home_phone = $3
		WHERE id = $4 AND created_by = $5
		RETURNING 
			id, 
			created_by, 
			created_at,
			form_type,
			updated_at, 
			first_name, 
			last_name, 
			home_phone
	`,
		pesticideFormInput.FirstName,
		pesticideFormInput.LastName,
		pesticideFormInput.HomePhone,
		formID,
		userID,
	).Scan(
		&pesticideForm.ID,
		&pesticideForm.CreatedBy,
		&pesticideForm.CreatedAt,
		&pesticideForm.FormType,
		&pesticideForm.UpdatedAt,
		&pesticideForm.FirstName,
		&pesticideForm.LastName,
		&pesticideForm.HomePhone,
	)
	if err != nil {
		//sql.ErrNoRows
		return PesticideForm{}, err
	}

	var query string
	query = `
		UPDATE pesticides
		SET pesticide_name = $1
		WHERE form_id = $2
		RETURNING pesticide_name
	`
	err = tx.QueryRowContext(ctx, query, pesticideFormInput.PesticideName, formID).Scan(
		&pesticideForm.PesticideName,
	)
	if err != nil {
		//sql.ErrNoRows
		return PesticideForm{}, err
	}

	if err := tx.Commit(); err != nil {
		return PesticideForm{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return pesticideForm, nil
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
