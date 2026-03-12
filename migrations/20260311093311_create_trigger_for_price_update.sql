-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_products_price_history()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO products_price_history (
        product_id,
        wholesale_price,
        retail_price
    )
    VALUES (
        NEW.id,
        NEW.wholesale_price,
        NEW.retail_price
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_products_price_history_insert
AFTER INSERT ON products
FOR EACH ROW
EXECUTE FUNCTION log_products_price_history();

CREATE TRIGGER trg_products_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON products
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_products_price_history();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_options_price_history()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO options_price_history (
        option_id,
        wholesale_price,
        retail_price
    )
    VALUES (
        NEW.id,
        NEW.wholesale_price,
        NEW.retail_price
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_options_price_history_insert
AFTER INSERT ON options
FOR EACH ROW
EXECUTE FUNCTION log_options_price_history();

CREATE TRIGGER trg_options_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON options
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_options_price_history();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_rails_price_history()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO rails_price_history (
        rail_id,
        wholesale_price,
        retail_price
    )
    VALUES (
        NEW.id,
        NEW.wholesale_price,
        NEW.retail_price
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_rails_price_history_insert
AFTER INSERT ON rails
FOR EACH ROW
EXECUTE FUNCTION log_rails_price_history();

CREATE TRIGGER trg_rails_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON rails
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_rails_price_history();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_res_drives_price_history()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO res_drives_price_history (
        drive_id,
        wholesale_price,
        retail_price
    )
    VALUES (
        NEW.id,
        NEW.wholesale_price,
        NEW.retail_price
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_res_drives_price_history_insert
AFTER INSERT ON residential_gate_drives
FOR EACH ROW
EXECUTE FUNCTION log_res_drives_price_history();

CREATE TRIGGER trg_res_drives_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON residential_gate_drives
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_res_drives_price_history();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_ind_drives_price_history()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO ind_drives_price_history (
        drive_id,
        wholesale_price,
        retail_price
    )
    VALUES (
        NEW.id,
        NEW.wholesale_price,
        NEW.retail_price
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_ind_drives_price_history_insert
AFTER INSERT ON industrial_gate_drives
FOR EACH ROW
EXECUTE FUNCTION log_ind_drives_price_history();

CREATE TRIGGER trg_ind_drives_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON industrial_gate_drives
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_ind_drives_price_history();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_sizes_price_history()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO sizes_price_history (
        width,
        height,
        gate_type,
        wholesale_price,
        retail_price
    )
    VALUES (
        NEW.width,
        NEW.height,
        NEW.gate_type,
        NEW.wholesale_price,
        NEW.retail_price
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_sizes_price_history_insert
AFTER INSERT ON sizes
FOR EACH ROW
EXECUTE FUNCTION log_sizes_price_history();

CREATE TRIGGER trg_sizes_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON sizes
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_sizes_price_history();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_lift_type_markup_history()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO lift_type_markup_history (
        lift_type_id,
        wholesale_markup,
        retail_markup
    )
    VALUES (
        NEW.id,
        NEW.wholesale_markup,
        NEW.retail_markup
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_lift_type_markup_history_insert
AFTER INSERT ON lift_types
FOR EACH ROW
EXECUTE FUNCTION log_lift_type_markup_history();

CREATE TRIGGER trg_lift_type_markup_history_update
AFTER UPDATE OF wholesale_markup, retail_markup ON lift_types
FOR EACH ROW
WHEN (
    OLD.wholesale_markup IS DISTINCT FROM NEW.wholesale_markup
    OR OLD.retail_markup IS DISTINCT FROM NEW.retail_markup
)
EXECUTE FUNCTION log_lift_type_markup_history();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_cycle_amount_markup_history()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO cycle_amount_markup_history (
        cycle_amount_id,
        wholesale_markup,
        retail_markup
    )
    VALUES (
        NEW.id,
        NEW.wholesale_markup,
        NEW.retail_markup
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_cycle_amount_markup_history_insert
AFTER INSERT ON cycle_amount
FOR EACH ROW
EXECUTE FUNCTION log_cycle_amount_markup_history();

CREATE TRIGGER trg_cycle_amount_markup_history_update
AFTER UPDATE OF wholesale_markup, retail_markup ON cycle_amount
FOR EACH ROW
WHEN (
    OLD.wholesale_markup IS DISTINCT FROM NEW.wholesale_markup
    OR OLD.retail_markup IS DISTINCT FROM NEW.retail_markup
)
EXECUTE FUNCTION log_cycle_amount_markup_history()

-- +goose Down

DROP TRIGGER IF EXISTS trg_cycle_amount_markup_history_update ON cycle_amount;
DROP TRIGGER IF EXISTS trg_cycle_amount_markup_history_insert ON cycle_amount;

DROP TRIGGER IF EXISTS trg_lift_type_markup_history_update ON lift_types;
DROP TRIGGER IF EXISTS trg_lift_type_markup_history_insert ON lift_types;

DROP TRIGGER IF EXISTS trg_sizes_price_history_update ON sizes;
DROP TRIGGER IF EXISTS trg_sizes_price_history_insert ON sizes;

DROP TRIGGER IF EXISTS trg_ind_drives_price_history_update ON industrial_gate_drives;
DROP TRIGGER IF EXISTS trg_ind_drives_price_history_insert ON industrial_gate_drives;

DROP TRIGGER IF EXISTS trg_res_drives_price_history_update ON residential_gate_drives;
DROP TRIGGER IF EXISTS trg_res_drives_price_history_insert ON residential_gate_drives;

DROP TRIGGER IF EXISTS trg_rails_price_history_update ON rails;
DROP TRIGGER IF EXISTS trg_rails_price_history_insert ON rails;

DROP TRIGGER IF EXISTS trg_options_price_history_update ON options;
DROP TRIGGER IF EXISTS trg_options_price_history_insert ON options;

DROP TRIGGER IF EXISTS trg_products_price_history_update ON products;
DROP TRIGGER IF EXISTS trg_products_price_history_insert ON products;

DROP FUNCTION IF EXISTS log_cycle_amount_markup_history();
DROP FUNCTION IF EXISTS log_lift_type_markup_history();
DROP FUNCTION IF EXISTS log_sizes_price_history();
DROP FUNCTION IF EXISTS log_ind_drives_price_history();
DROP FUNCTION IF EXISTS log_res_drives_price_history();
DROP FUNCTION IF EXISTS log_rails_price_history();
DROP FUNCTION IF EXISTS log_options_price_history();
DROP FUNCTION IF EXISTS log_products_price_history();
