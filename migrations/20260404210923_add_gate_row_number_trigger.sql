-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_gate_row_number_before_insert()
RETURNS TRIGGER AS $$
BEGIN
  SELECT COALESCE(MAX(row_number), 0) + 1
  INTO NEW.row_number
  FROM sales_and_gates
  WHERE sale_id = NEW.sale_id;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_set_gate_row_number_before_insert
BEFORE INSERT ON sales_and_gates
FOR EACH ROW
EXECUTE FUNCTION set_gate_row_number_before_insert();

-- +goose Down
DROP TRIGGER IF EXISTS trg_set_gate_row_number_before_insert ON sales_and_gates;

-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_gate_row_number_before_insert();
-- +goose StatementEnd
