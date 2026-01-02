package forms

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type FormsRepository struct {
	db *sql.DB
}

func NewFormsRepository(database *sql.DB) *FormsRepository {
	return &FormsRepository{db: database}
}

type CreateFormInput struct {
	CreatedBy string
	FirstName string
	LastName  string
	HomePhone string
}

func (r *FormsRepository) CreateShrubForm(
	ctx context.Context,
	formInput CreateFormInput,
	details *ShrubDetails,
) (ShrubForm, error) {
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

func (r *FormsRepository) CreatePesticideForm(
	ctx context.Context,
	formInput CreateFormInput,
	details *PesticideDetails,
) (PesticideForm, error) {
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

func (r *FormsRepository) ListFormsByUserId(
	ctx context.Context,
	userID string,
	sortBy string,
	order string,
) ([]FormView, error) {

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

	var forms []FormView
	for rows.Next() {
		var (
			form          Form
			numShrubs     sql.NullInt32
			pesticideName sql.NullString
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
			&numShrubs,
			&pesticideName,
		)
		if err != nil {
			return nil, err
		}

		item := FormView{
			FormType: form.FormType,
			Form:     &form,
		}

		switch form.FormType {
		case "shrub":
			if numShrubs.Valid {
				item.Shrub = &ShrubForm{
					Form:      form,
					NumShrubs: int(numShrubs.Int32),
				}
			}

		case "pesticide":
			if pesticideName.Valid {
				item.Pesticide = &PesticideForm{
					Form:          form,
					PesticideName: pesticideName.String,
				}
			}
		default:
			return nil, fmt.Errorf("unknown form_type: %s", form.FormType)
		}
		forms = append(forms, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return forms, nil
}

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
		form          Form
		numShrubs     sql.NullInt32
		pesticideName sql.NullString
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
		&numShrubs,
		&pesticideName,
	)
	if err != nil {
		// Important: let sql.ErrNoRows propagate
		return nil, err
	}

	view := &FormView{
		FormType: form.FormType,
		Form:     &form,
	}

	switch form.FormType {
	case "shrub":
		if !numShrubs.Valid {
			// Defensive: DB invariant violated
			return nil, fmt.Errorf("shrub form %s missing shrub details", form.ID)
		}
		view.Shrub = &ShrubForm{
			Form:      form,
			NumShrubs: int(numShrubs.Int32),
		}

	case "pesticide":
		if !pesticideName.Valid {
			// Defensive: DB invariant violated
			return nil, fmt.Errorf("pesticide form %s missing pesticide details", form.ID)
		}
		view.Pesticide = &PesticideForm{
			Form:          form,
			PesticideName: pesticideName.String,
		}

	default:
		// Defensive: unknown form type
		return nil, fmt.Errorf("unknown form_type: %s", form.FormType)
	}

	return view, nil
}

// TODO implement this
func (r *FormsRepository) UpdateFormById(
	ctx context.Context,
	formID string,
	userID string,
	formInput CreateFormInput,
	shrub *ShrubDetails,
	pesticide *PesticideDetails,
) (*FormView, error) {
	return nil, nil
}
