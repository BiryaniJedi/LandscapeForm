CREATE TABLE chemicals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand_name TEXT NOT NULL,
    chemical_name TEXT NOT NULL,
    epa_reg_num TEXT NOT NULL,
    recipe TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
