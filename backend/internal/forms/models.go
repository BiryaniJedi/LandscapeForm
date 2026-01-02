package forms

import "time"

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
	NumShrubs int
}

type PesticideForm struct {
	Form
	PesticideName string
}

type ShrubDetails struct {
	FormID    string
	NumShrubs int
}

type PesticideDetails struct {
	FormID        string
	PesticideName string
}

// Used to store forms in a common slice
type FormView struct {
	FormType string

	Form *Form

	Shrub     *ShrubForm
	Pesticide *PesticideForm
}
