CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    pending BOOLEAN NOT NULL DEFAULT TRUE,
    role TEXT NOT NULL DEFAULT 'employee',
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    date_of_birth DATE NOT NULL DEFAULT '2000-01-01',
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

-- Forms table
CREATE TABLE forms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    form_type TEXT NOT NULL CHECK (form_type IN ('shrub', 'lawn')),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Client info
    first_name TEXT NOT NULL,
    last_name  TEXT NOT NULL,
    street_number TEXT NOT NULL,
    street_name TEXT NOT NULL,
    town TEXT NOT NULL,
    zip_code TEXT NOT NULL CHECK (zip_code ~ '^\d{5}(-\d{4})?$'),
    home_phone TEXT NOT NULL,
    other_phone TEXT NOT NULL,

    -- General form info
    call_before BOOLEAN NOT NULL DEFAULT FALSE,
    is_holiday BOOLEAN NOT NULL DEFAULT FALSE,
    num_pest_applications INT NOT NULL DEFAULT 0
);

-- chemical list for forms
CREATE TABLE chemicals (
    id SMALLSERIAL PRIMARY KEY,
    category TEXT NOT NULL CHECK (category IN ('lawn', 'shrub')),
    brand_name TEXT NOT NULL,
    chemical_name TEXT NOT NULL,
    epa_reg_no TEXT NOT NULL,
    recipe TEXT NOT NULL,
    unit TEXT NOT NULL
);

CREATE TABLE pesticide_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    chem_used SMALLINT NOT NULL REFERENCES chemicals(id),
    app_timestamp TIMESTAMPTZ NOT NULL,
    rate TEXT NOT NULL,
    amount_applied NUMERIC(10, 2) NOT NULL,
    location_code VARCHAR(2) NOT NULL
);

-- Shrub forms table
CREATE TABLE shrub_forms (
    form_id UUID PRIMARY KEY REFERENCES forms(id) ON DELETE CASCADE,

    -- Shrubs info
    flea_only BOOLEAN NOT NULL DEFAULT FALSE
);

-- Lawn forms table
CREATE TABLE lawn_forms (
    form_id UUID PRIMARY KEY REFERENCES forms(id) ON DELETE CASCADE,

    -- Lawns info
    lawn_area_sq_ft INT NOT NULL,
    fert_only BOOLEAN NOT NULL DEFAULT FALSE
);



-- Notes for each form (optional)
CREATE TABLE notes (
    id SMALLSERIAL PRIMARY KEY,
    form_id UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    note VARCHAR(25) NOT NULL
);


-- Indices
CREATE INDEX idx_forms_user_created_at ON forms(created_by, created_at DESC);
CREATE INDEX idx_forms_name_lower ON forms (LOWER(first_name), LOWER(last_name));
CREATE INDEX idx_forms_name ON forms(first_name, last_name);
CREATE INDEX idx_forms_street_name ON forms(street_name);
CREATE INDEX idx_forms_street_name_lower ON forms(LOWER(street_name));
CREATE INDEX idx_forms_town ON forms(town);
CREATE INDEX idx_forms_town_lower ON forms(LOWER(town));
CREATE INDEX idx_forms_street_number ON forms(street_number);
CREATE INDEX idx_forms_zip_code ON forms(zip_code);
CREATE INDEX idx_forms_home_phone ON forms(home_phone);

-- Triggers
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_forms_updated
BEFORE UPDATE ON forms
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_users_updated
BEFORE UPDATE ON users 
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE FUNCTION prevent_form_type_change()
RETURNS TRIGGER AS $$
BEGIN
  IF OLD.form_type IS DISTINCT FROM NEW.form_type THEN
    RAISE EXCEPTION 'form_type cannot be changed once set';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_form_type_change
BEFORE UPDATE ON forms
FOR EACH ROW
EXECUTE FUNCTION prevent_form_type_change();

CREATE OR REPLACE FUNCTION enforce_shrub_form()
RETURNS TRIGGER AS $$
BEGIN
  IF (SELECT form_type FROM forms WHERE id = NEW.form_id) <> 'shrub' THEN
    RAISE EXCEPTION 'Form % is not a shrub form', NEW.form_id;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_shrub_type_check
BEFORE INSERT OR UPDATE ON shrub_forms
FOR EACH ROW
EXECUTE FUNCTION enforce_shrub_form();

CREATE OR REPLACE FUNCTION enforce_lawn_form()
RETURNS TRIGGER AS $$
BEGIN
  IF (SELECT form_type FROM forms WHERE id = NEW.form_id) <> 'lawn' THEN
    RAISE EXCEPTION 'Form % is not a lawn form', NEW.form_id;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_lawn_type_check
BEFORE INSERT OR UPDATE ON lawn_forms 
FOR EACH ROW
EXECUTE FUNCTION enforce_lawn_form();
