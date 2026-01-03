// Package forms provides data access and domain models for landscape forms.
// It encapsulates persistence logic, enforces ownership rules, and ensures
// type-safe access to shrub and pesticide forms.
package forms

import (
	"context"
	"database/sql"
	"errors"
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
type CreateFormInput struct {
	CreatedBy string
	FirstName string
	LastName  string
	HomePhone string
}

// UpdateFormInput contains the fields that may be updated on an existing form.
type UpdateFormInput struct {
	FirstName string
	LastName  string
	HomePhone string
}

// CreateShrubForm creates a new shrub form and its associated shrub details.
// The operation is atomic and will fail if shrub details are not provided.
func (r *FormsRepository) CreateShrubForm(
	ctx context.Context,
	formInput CreateFormInput,
	details *ShrubDetails,
) (ShrubForm, error) {
	if details == nil {
		return ShrubForm{}, errors.New("shrub details required")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return ShrubForm{}, err
	}
	defer tx.Rollback()

	var res ShrubForm
	err = tx.QueryRowContext(ctx, `
		INSERT INTO forms (
			created_by,
			form_type,
			first_name,
			last_name,
			home_phone
		)
		VALUES ($1, 'shrub', $2, $3, $4)
		RETURNING id, created_by, created_at, form_type, updated_at, first_name, last_name, home_phone
	`,
		formInput.CreatedBy,
		formInput.FirstName,
		formInput.LastName,
		formInput.HomePhone,
	).Scan(
		&res.Form.ID,
		&res.Form.CreatedBy,
		&res.Form.CreatedAt,
		&res.Form.FormType,
		&res.Form.UpdatedAt,
		&res.Form.FirstName,
		&res.Form.LastName,
		&res.Form.HomePhone,
	)
	if err != nil {
		return ShrubForm{}, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO shrubs (
			form_id,
			num_shrubs
		)
		VALUES ($1, $2)
	`,
		res.Form.ID,
		details.NumShrubs,
	)
	if err != nil {
		return ShrubForm{}, err
	}
	res.NumShrubs = details.NumShrubs

	if err := tx.Commit(); err != nil {
		return ShrubForm{}, err
	}

	return res, nil
}

// CreatePesticideForm creates a new pesticide form and its associated pesticide details.
// The operation is atomic and will fail if pesticide details are not provided.
func (r *FormsRepository) CreatePesticideForm(
	ctx context.Context,
	formInput CreateFormInput,
	details *PesticideDetails,
) (PesticideForm, error) {
	if details == nil {
		return PesticideForm{}, errors.New("pesticide details required")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return PesticideForm{}, err
	}
	defer tx.Rollback()

	var res PesticideForm
	err = tx.QueryRowContext(ctx, `
		INSERT INTO forms (
			created_by,
			form_type,
			first_name,
			last_name,
			home_phone
		)
		VALUES ($1, 'pesticide', $2, $3, $4)
		RETURNING id, created_by, created_at, form_type, updated_at, first_name, last_name, home_phone
	`,
		formInput.CreatedBy,
		formInput.FirstName,
		formInput.LastName,
		formInput.HomePhone,
	).Scan(
		&res.Form.ID,
		&res.Form.CreatedBy,
		&res.Form.CreatedAt,
		&res.Form.FormType,
		&res.Form.UpdatedAt,
		&res.Form.FirstName,
		&res.Form.LastName,
		&res.Form.HomePhone,
	)

	if err != nil {
		return PesticideForm{}, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO pesticides (
			form_id,
			pesticide_name	
		)
		VALUES ($1, $2)
	`,
		res.Form.ID,
		details.PesticideName,
	)
	if err != nil {
		return PesticideForm{}, err
	}
	res.PesticideName = details.PesticideName

	if err := tx.Commit(); err != nil {
		return PesticideForm{}, err
	}

	return res, nil
}

// ListFormsByUserId returns all forms owned by the given user.
// Results may be sorted by first name, last name, or creation time.
// Each returned FormView is fully hydrated with its subtype details.
func (r *FormsRepository) ListFormsByUserId(
	ctx context.Context,
	userID string,
	sortBy string,
	order string,
) ([]*FormView, error) {

	allowedSorts := map[string]string{
		"first_name": "f.first_name",
		"last_name":  "f.last_name",
		"created_at": "f.created_at",
	}

	sortColumn, ok := allowedSorts[sortBy]
	if !ok {
		sortColumn = "f.created_at"
	}

	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

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
		WHERE f.created_by = $1
		ORDER BY %s %s
	`, sortColumn, order)

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		var view *FormView
		switch form.FormType {
		case "shrub":
			shrubDetails, err := shrub.ToDomain()
			if err != nil {
				return nil, err
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
				return nil, err
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
		return nil, err
	}

	return forms, nil
}

// GetFormById returns a single form owned by the given user.
// It returns sql.ErrNoRows if the form does not exist or is not owned by the user.
func (r *FormsRepository) GetFormById(
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
			return nil, err
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
			return nil, err
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

// UpdateFormById updates a form and its associated subtype fields.
// The existing form type determines which subtype is updated and cannot be changed.
// It returns sql.ErrNoRows if the form does not exist or is not owned by the user.
func (r *FormsRepository) UpdateFormById(
	ctx context.Context,
	formID string,
	userID string,
	formInput UpdateFormInput,
	shrub *ShrubDetails,
	pesticide *PesticideDetails,
) (*FormView, error) {
	if shrub != nil && pesticide != nil {
		return nil, errors.New("only one of shrub or pesticide details allowed")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var (
		view         *FormView
		form         Form
		shrubRow     shrubRow
		pesticideRow pesticideRow
	)

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
		formInput.FirstName,
		formInput.LastName,
		formInput.HomePhone,
		formID,
		userID,
	).Scan(
		&form.ID,
		&form.CreatedBy,
		&form.CreatedAt,
		&form.FormType,
		&form.UpdatedAt,
		&form.FirstName,
		&form.LastName,
		&form.HomePhone,
	)
	if err != nil {
		return nil, err
	}

	var query string
	switch form.FormType {
	case "shrub":
		if shrub == nil {
			return nil, errors.New("missing shrub details")
		}
		query = `
			UPDATE shrubs
			SET num_shrubs = $1
			WHERE form_id = $2
			RETURNING num_shrubs
		`
		err = tx.QueryRowContext(ctx, query, shrub.NumShrubs, formID).Scan(
			&shrubRow.NumShrubs,
		)
		if err != nil {
			return nil, err
		}

		shrubDetails, err := shrubRow.ToDomain()
		if err != nil {
			return nil, err
		}

		view = NewShrubFormView(
			ShrubForm{
				Form:         form,
				ShrubDetails: shrubDetails,
			},
		)

	case "pesticide":
		if pesticide == nil {
			return nil, errors.New("missing pesticide details")
		}
		query = `
			UPDATE pesticides 
			SET pesticide_name = $1
			WHERE form_id = $2
			RETURNING pesticide_name
		`

		err = tx.QueryRowContext(ctx, query, pesticide.PesticideName, formID).Scan(
			&pesticideRow.PesticideName,
		)
		if err != nil {
			return nil, err
		}

		pesticideDetails, err := pesticideRow.ToDomain()
		if err != nil {
			return nil, err
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

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return view, nil
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
