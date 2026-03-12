-- +goose Up
CREATE TABLE
    products_price_history (
        id BIGSERIAL PRIMARY KEY,
        product_id INT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
        wholesale_price NUMERIC(10, 2) NOT NULL CHECK (wholesale_price >= 0),
        retail_price NUMERIC(10, 2) NOT NULL CHECK (retail_price >= 0),
        set_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_products_price_history_product_date ON products_price_history (product_id, set_at);

CREATE TABLE
    options_price_history (
        id BIGSERIAL PRIMARY KEY,
        option_id INT NOT NULL REFERENCES options (id) ON DELETE CASCADE,
        wholesale_price NUMERIC(10, 2) NOT NULL CHECK (wholesale_price >= 0),
        retail_price NUMERIC(10, 2) NOT NULL CHECK (retail_price >= 0),
        set_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_options_price_history_option_date ON options_price_history (option_id, set_at);

CREATE TABLE
    sizes_price_history (
        id BIGSERIAL PRIMARY KEY,
        width BIGINT NOT NULL,
        height BIGINT NOT NULL,
        gate_type gate_type NOT NULL,
        wholesale_price NUMERIC(10, 2) NOT NULL CHECK (wholesale_price >= 0),
        retail_price NUMERIC(10, 2) NOT NULL CHECK (retail_price >= 0),
        set_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (width, height, gate_type) REFERENCES sizes (width, height, gate_type)
    );

CREATE INDEX idx_sizes_price_history_size_date ON sizes_price_history (width, height, gate_type, set_at);

CREATE TABLE
    rails_price_history (
        id BIGSERIAL PRIMARY KEY,
        rail_id INT NOT NULL REFERENCES rails (id) ON DELETE CASCADE,
        wholesale_price NUMERIC(10, 2) NOT NULL CHECK (wholesale_price >= 0),
        retail_price NUMERIC(10, 2) NOT NULL CHECK (retail_price >= 0),
        set_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_rails_price_history_rail_date ON rails_price_history (rail_id, set_at);

CREATE TABLE
    res_drives_price_history (
        id BIGSERIAL PRIMARY KEY,
        drive_id INT NOT NULL REFERENCES residential_gate_drives (id) ON DELETE CASCADE,
        wholesale_price NUMERIC(10, 2) NOT NULL CHECK (wholesale_price >= 0),
        retail_price NUMERIC(10, 2) NOT NULL CHECK (retail_price >= 0),
        set_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_res_drives_price_history_drive_date ON res_drives_price_history (drive_id, set_at);

CREATE TABLE
    ind_drives_price_history (
        id BIGSERIAL PRIMARY KEY,
        drive_id INT NOT NULL REFERENCES industrial_gate_drives (id) ON DELETE CASCADE,
        wholesale_price NUMERIC(10, 2) NOT NULL CHECK (wholesale_price >= 0),
        retail_price NUMERIC(10, 2) NOT NULL CHECK (retail_price >= 0),
        set_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_ind_drives_price_history_drive_date ON ind_drives_price_history (drive_id, set_at);

CREATE TABLE
    lift_type_markup_history (
        id BIGSERIAL PRIMARY KEY,
        lift_type_id INT NOT NULL REFERENCES lift_types (id) ON DELETE CASCADE,
        wholesale_markup NUMERIC(10, 2) NOT NULL CHECK (wholesale_markup >= 0),
        retail_markup NUMERIC(10, 2) NOT NULL CHECK (retail_markup >= 0),
        set_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_lift_type_markup_history_lift_type_date ON lift_type_markup_history (lift_type_id, set_at);

CREATE TABLE
    cycle_amount_markup_history (
        id BIGSERIAL PRIMARY KEY,
        cycle_amount_id INT NOT NULL REFERENCES cycle_amount (id) ON DELETE CASCADE,
        wholesale_markup NUMERIC(10, 2) NOT NULL CHECK (wholesale_markup >= 0),
        retail_markup NUMERIC(10, 2) NOT NULL CHECK (retail_markup >= 0),
        set_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_cycle_amount_markup_history_cycle_amount_date ON cycle_amount_markup_history (cycle_amount_id, set_at);

-- +goose Down
DROP TABLE IF EXISTS cycle_amount_markup_history;

DROP TABLE IF EXISTS lift_type_markup_history;

DROP TABLE IF EXISTS ind_drives_price_history;

DROP TABLE IF EXISTS res_drives_price_history;

DROP TABLE IF EXISTS rails_price_history;

DROP TABLE IF EXISTS sizes_price_history;

DROP TABLE IF EXISTS options_price_history;

DROP TABLE IF EXISTS products_price_history;
