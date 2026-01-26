// Package chemicals provides data access and domain models for landscape chemicals.
package chemicals

import (
	"context"
	"database/sql"
	"fmt"
)

// ChemicalsRepository provides database access for chemical records.
type ChemicalsRepository struct {
	db *sql.DB
}

// NewChemicalsRepository returns a repository backed by the given database connection.
func NewChemicalsRepository(database *sql.DB) *ChemicalsRepository {
	return &ChemicalsRepository{db: database}
}

type Chemical struct {
	ID           int
	Category     string
	BrandName    string
	ChemicalName string
	EpaRegNo     string
	Recipe       string
	Unit         string
}

type ChemicalInput struct {
	Category     string
	BrandName    string
	ChemicalName string
	EpaRegNo     string
	Recipe       string
	Unit         string
}

// CreateChemical creates a new chemical record.
// Returns the created chemical's ID upon success.
// The operation is atomic.
func (r *ChemicalsRepository) CreateChemical(
	ctx context.Context,
	chemicalInput ChemicalInput,
) (string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var formID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO chemicals (
			category,
			brand_name,
			chemical_name,
			epa_reg_no,
			recipe,
			unit
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`,
		chemicalInput.Category,
		chemicalInput.BrandName,
		chemicalInput.ChemicalName,
		chemicalInput.EpaRegNo,
		chemicalInput.Recipe,
		chemicalInput.Unit,
	).Scan(
		&formID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to insert chemical %s: %w", chemicalInput.ChemicalName, err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction for inserting chemical %s: %w", chemicalInput.ChemicalName, err)
	}

	return formID, nil
}

// ListChemicalsByCategory returns all chemicals in a given category.
func (r *ChemicalsRepository) ListChemicalsByCategory(
	ctx context.Context,
	category string,
) ([]Chemical, error) {
	query := `
		SELECT
			c.id,
			c.category,
			c.brand_name,
			c.chemical_name,
			c.epa_reg_no,
			c.recipe,
			c.unit
		FROM chemicals c
		WHERE c.category = $1
	`

	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("error querying rows for chemicals list: %w", err)
	}
	defer rows.Close()

	var chemicals []Chemical
	for rows.Next() {
		var chemical Chemical
		err := rows.Scan(
			&chemical.ID,
			&chemical.Category,
			&chemical.BrandName,
			&chemical.ChemicalName,
			&chemical.EpaRegNo,
			&chemical.Recipe,
			&chemical.Unit,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		chemicals = append(chemicals, chemical)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after queries for chemicals list: %w", err)
	}

	return chemicals, nil
}

// UpdateChemicalById updates a chemical by ID.
// Returns the updated chemical upon success.
// Returns sql.ErrNoRows if the chemical does not exist.
func (r *ChemicalsRepository) UpdateChemicalById(
	ctx context.Context,
	ID int,
	chemicalInput ChemicalInput,
) (Chemical, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Chemical{}, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	var chemical Chemical

	err = tx.QueryRowContext(ctx, `
		UPDATE chemicals
		SET category = $1,
			brand_name = $2,
			chemical_name = $3,
			epa_reg_no = $4,
			recipe = $5,
			unit = $6
		WHERE id = $7
		RETURNING
			id,
			category,
			brand_name,
			chemical_name,
			epa_reg_no,
			recipe,
			unit
	`,
		chemicalInput.Category,
		chemicalInput.BrandName,
		chemicalInput.ChemicalName,
		chemicalInput.EpaRegNo,
		chemicalInput.Recipe,
		chemicalInput.Unit,
		ID,
	).Scan(
		&chemical.ID,
		&chemical.Category,
		&chemical.BrandName,
		&chemical.ChemicalName,
		&chemical.EpaRegNo,
		&chemical.Recipe,
		&chemical.Unit,
	)
	if err != nil {
		return chemical, err
	}

	if err := tx.Commit(); err != nil {
		return Chemical{}, fmt.Errorf("error committing transaction: %w", err)
	}

	return chemical, nil
}

// DeleteChemicalById deletes a chemical by ID.
// Returns sql.ErrNoRows if the chemical does not exist.
func (r *ChemicalsRepository) DeleteChemicalById(
	ctx context.Context,
	ID int,
) error {

	err := r.db.QueryRowContext(ctx, `
		DELETE FROM chemicals
		WHERE id = $1
		RETURNING id
	`, ID).Scan(&ID)

	if err != nil {
		return err
	}

	return nil
}
