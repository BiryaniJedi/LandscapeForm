CREATE TABLE form_chemicals (
    form_id UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    chemical_id UUID NOT NULL REFERENCES chemicals(id),
    amount_applied TEXT NOT NULL,
    PRIMARY KEY (form_id, chemical_id)
);
