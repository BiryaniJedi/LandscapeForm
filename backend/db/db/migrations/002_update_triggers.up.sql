DROP TRIGGER IF EXISTS trg_shrub_type_check ON shrub_forms;
DROP TRIGGER IF EXISTS trg_pesticide_type_check ON lawn_forms;

-- Recreate triggers with same logic
CREATE TRIGGER trg_shrub_form_type_check
BEFORE INSERT OR UPDATE ON shrub_forms
FOR EACH ROW
EXECUTE FUNCTION enforce_shrub_form();

CREATE TRIGGER trg_lawn_form_type_check
BEFORE INSERT OR UPDATE ON lawn_forms
FOR EACH ROW
EXECUTE FUNCTION enforce_pesticide_form();
