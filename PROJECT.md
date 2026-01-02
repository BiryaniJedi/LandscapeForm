# Tables
Basic for dev, will update with all fields later

## Database Schema
```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Forms table
CREATE TABLE forms (
    -- Form info
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    form_type TEXT NOT NULL CHECK (form_type IN ('shrub', 'pesticide')),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Client info
    first_name TEXT NOT NULL,
    last_name  TEXT NOT NULL,
    home_phone TEXT NOT NULL
);

-- Shrubs table
CREATE TABLE shrubs (
    form_id UUID PRIMARY KEY REFERENCES forms(id) ON DELETE CASCADE,

    -- Shrubs info
    num_shrubs INT NOT NULL CHECK (num_shrubs >= 0)
);

-- Pesticides table
CREATE TABLE pesticides (
    form_id UUID PRIMARY KEY REFERENCES forms(id) ON DELETE CASCADE,

    -- Pesticides info
    pesticide_name TEXT NOT NULL
);

-- Indices
CREATE INDEX idx_forms_user_created_at ON forms(created_by, created_at DESC);
CREATE INDEX idx_forms_name_lower ON forms (LOWER(first_name), LOWER(last_name));
CREATE INDEX idx_forms_name ON forms(first_name, last_name);

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
BEFORE INSERT OR UPDATE ON shrubs
FOR EACH ROW
EXECUTE FUNCTION enforce_shrub_form();

CREATE OR REPLACE FUNCTION enforce_pesticide_form()
RETURNS TRIGGER AS $$
BEGIN
  IF (SELECT form_type FROM forms WHERE id = NEW.form_id) <> 'pesticide' THEN
    RAISE EXCEPTION 'Form % is not a pesticide form', NEW.form_id;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_pesticide_type_check
BEFORE INSERT OR UPDATE ON pesticides
FOR EACH ROW
EXECUTE FUNCTION enforce_pesticide_form();
```

# Routes
```ts
POST   "/forms/pesticide"
POST   "/forms/shrub"  
GET    "/forms"  
GET    "/forms/{id}"  
PUT    "/forms/{id}"  
DELETE "/forms/{id}"  
GET    "/forms/{id}/pdf"
POST   "/forms/import/pdf"  
```

# User
Forms will sort by:
* `first_name`
* `last_name`
* `created_at`<br>
Forms will filter by:
* `first_name`
* `last_name`
* `created_at`<br>