package forms

import (
	"database/sql"
	"errors"
	"time"
)

type Form struct {
	ID        string
	CreatedBy string
	CreatedAt time.Time
	FormType  string
	UpdatedAt time.Time

	FirstName string
	LastName  string
	HomePhone string
}

type ShrubForm struct {
	Form
	ShrubDetails
}

type PesticideForm struct {
	Form
	PesticideDetails
}

type ShrubDetails struct {
	NumShrubs int
}

type PesticideDetails struct {
	PesticideName string
}

type shrubRow struct {
	NumShrubs sql.NullInt32
}

type pesticideRow struct {
	PesticideName sql.NullString
}

// Used to store forms in a common slice
type FormView struct {
	FormType  string
	Shrub     *ShrubForm
	Pesticide *PesticideForm
}

func (r shrubRow) ToDomain() (ShrubDetails, error) {
	if !r.NumShrubs.Valid {
		return ShrubDetails{}, errors.New("missing shrub details")
	}
	return ShrubDetails{
		NumShrubs: int(r.NumShrubs.Int32),
	}, nil
}

func (r pesticideRow) ToDomain() (PesticideDetails, error) {
	if !r.PesticideName.Valid {
		return PesticideDetails{}, errors.New("missing pesticide details")
	}
	return PesticideDetails{
		PesticideName: r.PesticideName.String,
	}, nil
}

func NewShrubFormView(form ShrubForm) *FormView {
	return &FormView{
		FormType: "shrub",
		Shrub:    &form,
	}
}

func NewPesticideFormView(form PesticideForm) *FormView {
	return &FormView{
		FormType:  "pesticide",
		Pesticide: &form,
	}
}
