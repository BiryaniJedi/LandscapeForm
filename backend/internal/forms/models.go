package forms

import (
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

type Form struct {
	ID        string
	CreatedBy string
	CreatedAt time.Time
	FormType  string
	UpdatedAt time.Time

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

	FirstAppDate time.Time
	LastAppDate  time.Time

	AppTimes []PestApp
	Notes    []Note
}

type PestApp struct {
	ID            int
	ChemUsed      int
	AppTimestamp  time.Time
	Rate          string
	AmountApplied decimal.Decimal
	LocationCode  string
}

type Note struct {
	ID      int
	Message string
}

type ShrubForm struct {
	Form
	ShrubDetails
}

type LawnForm struct {
	Form
	LawnDetails
}

type ShrubDetails struct {
	FleaOnly bool
}

type LawnDetails struct {
	LawnAreaSqFt int
	FertOnly     bool
}

type shrubRow struct {
	FleaOnly sql.NullBool
}

type lawnRow struct {
	LawnAreaSqFt sql.NullInt32
	FertOnly     sql.NullBool
}

// Used to store forms in a common slice
type FormView struct {
	FormType string
	Shrub    *ShrubForm
	Lawn     *LawnForm
}

func (r shrubRow) ToDomain() (ShrubDetails, error) {
	if !r.FleaOnly.Valid {
		return ShrubDetails{}, errors.New("missing shrub details")
	}
	return ShrubDetails{
		FleaOnly: r.FleaOnly.Bool,
	}, nil
}

func (r lawnRow) ToDomain() (LawnDetails, error) {
	if !(r.FertOnly.Valid && r.LawnAreaSqFt.Valid) {
		return LawnDetails{}, errors.New("missing lawn details")
	}
	return LawnDetails{
		FertOnly:     r.FertOnly.Bool,
		LawnAreaSqFt: int(r.LawnAreaSqFt.Int32),
	}, nil
}

func NewShrubFormView(form ShrubForm) *FormView {
	return &FormView{
		FormType: "shrub",
		Shrub:    &form,
	}
}

func NewLawnFormView(form LawnForm) *FormView {
	return &FormView{
		FormType: "lawn",
		Lawn:     &form,
	}
}
