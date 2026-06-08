-- +goose Up
-- ============================================================
-- Перечисляемые типы
-- ============================================================

CREATE TYPE drive_type AS ENUM (
    'industrial',
    'residential',
    'manual'
);

CREATE TYPE gate_type AS ENUM (
    'industrial',
    'residential'
);

CREATE TYPE order_status AS ENUM (
    'pending',
    'paid',
    'cancelled',
    'confirmed',
    'done'
);

CREATE TYPE registration_request_status AS ENUM (
    'pending',
    'approved',
    'rejected'
);

CREATE TYPE role AS ENUM (
    'dealer',
    'manager',
    'admin',
    'logistician'
);

-- ============================================================
-- Пользователи, роли и сессии
-- ============================================================

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    fullname TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    password_hash BYTEA NOT NULL,
    role role NOT NULL,

    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE TABLE sessions (
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT sessions_pkey PRIMARY KEY (user_id),

    CONSTRAINT fk_sessions_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);

-- ============================================================
-- Компании, дилеры и менеджеры
-- ============================================================

CREATE TABLE companies (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE dealers (
    user_id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    address TEXT,

    CONSTRAINT dealers_pkey PRIMARY KEY (user_id),

    CONSTRAINT fk_dealers_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT fk_dealers_company
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
);

CREATE TABLE dealer_contract (
    dealer_id       BIGINT NOT NULL PRIMARY KEY,
    contract_number TEXT NOT NULL,
    signed_at       DATE NOT NULL,
    path            TEXT,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_dealer_contracts_dealer
        FOREIGN KEY (dealer_id)
        REFERENCES dealers(user_id)
);

CREATE TABLE dealer_manager_assignments (
    manager_id BIGINT NOT NULL,
    dealer_id BIGINT NOT NULL,

    CONSTRAINT dealer_manager_assignments_pkey
        PRIMARY KEY (manager_id, dealer_id),

    CONSTRAINT fk_dealer_manager_assignments_manager
        FOREIGN KEY (manager_id)
        REFERENCES users(id),

    CONSTRAINT fk_dealer_manager_assignments_dealer
        FOREIGN KEY (dealer_id)
        REFERENCES dealers(user_id)
);

CREATE TABLE dealer_registration_requests (
    id BIGSERIAL PRIMARY KEY,
    company TEXT NOT NULL,
    fullname TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    status registration_request_status DEFAULT 'pending'::registration_request_status NOT NULL
);

-- ============================================================
-- Справочники для расчёта стоимости ворот
-- ============================================================

CREATE TABLE colors (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL,
    hex TEXT NOT NULL
);

CREATE TABLE cycle_amounts (
    id BIGSERIAL PRIMARY KEY,
    amount TEXT NOT NULL,
    wholesale_markup NUMERIC(10,2),
    retail_markup NUMERIC(10,2)
);

CREATE TABLE lift_types (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    min_headroom INTEGER NOT NULL,
    max_headroom INTEGER NOT NULL,
    wholesale_markup NUMERIC(10,2),
    retail_markup NUMERIC(10,2)
);

CREATE TABLE gate_sizes (
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    wholesale_price NUMERIC(10,2),
    retail_price NUMERIC(10,2),
    gate_type gate_type DEFAULT 'industrial'::gate_type NOT NULL,

    CONSTRAINT gate_sizes_pkey PRIMARY KEY (width, height, gate_type)
);

CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    wholesale_price NUMERIC(10,2),
    retail_price NUMERIC(10,2),
    image_path TEXT NOT NULL DEFAULT 'no_image'
);

CREATE TABLE options (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    wholesale_price NUMERIC(10,2),
    retail_price NUMERIC(10,2),
    image_path TEXT NOT NULL DEFAULT 'no_image'
);

CREATE TABLE rails (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    wholesale_price NUMERIC(10,2),
    retail_price NUMERIC(10,2),
    specifications TEXT,
    image_path TEXT NOT NULL DEFAULT 'no_image'
);

CREATE TABLE industrial_gate_drives (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    wholesale_price NUMERIC(10,2),
    retail_price NUMERIC(10,2),
    specifications TEXT,
    image_path TEXT NOT NULL DEFAULT 'no_image'
);

CREATE TABLE residential_gate_drives (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    wholesale_price NUMERIC(10,2),
    retail_price NUMERIC(10,2),
    specifications TEXT,
    image_path TEXT NOT NULL DEFAULT 'no_image'
);

CREATE TABLE manual_drive_prices (
    id BIGSERIAL PRIMARY KEY,
    chain_meter_retail_price NUMERIC(10,2) NOT NULL,
    chain_meter_wholesale_price NUMERIC(10,2) NOT NULL,
    rcp_retail_price NUMERIC(10,2) NOT NULL,
    rcp_wholesale_price NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    CONSTRAINT manual_drive_prices_chain_meter_retail_price_check
        CHECK (chain_meter_retail_price >= 0),

    CONSTRAINT manual_drive_prices_chain_meter_wholesale_price_check
        CHECK (chain_meter_wholesale_price >= 0),

    CONSTRAINT manual_drive_prices_rcp_retail_price_check
        CHECK (rcp_retail_price >= 0),

    CONSTRAINT manual_drive_prices_rcp_wholesale_price_check
        CHECK (rcp_wholesale_price >= 0)
);

-- ============================================================
-- Заказы
-- ============================================================

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    dealer_id BIGINT,
    manager_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    status order_status DEFAULT 'pending'::order_status NOT NULL,
    manufacture_date DATE NULL,
    finalized_at TIMESTAMPTZ NULL,

    CONSTRAINT fk_orders_dealer
        FOREIGN KEY (dealer_id)
        REFERENCES dealers(user_id),

    CONSTRAINT fk_orders_manager
        FOREIGN KEY (manager_id)
        REFERENCES users(id)
);

CREATE TABLE order_gates (
    order_id BIGINT NOT NULL,
    row_number BIGINT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    lift_type_id BIGINT NOT NULL,
    color_out_id BIGINT NOT NULL,
    cycle_amount_id BIGINT NOT NULL,
    headroom INTEGER NOT NULL,
    gate_type gate_type DEFAULT 'industrial'::gate_type NOT NULL,
    drive_type drive_type DEFAULT 'industrial'::drive_type NOT NULL,
    amount INTEGER NOT NULL,
    gate_retail_price NUMERIC(10,2) NOT NULL,
    gate_wholesale_price NUMERIC(10,2) NOT NULL,

    CONSTRAINT order_gates_pkey
        PRIMARY KEY (order_id, row_number),

    CONSTRAINT fk_order_gates_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_order_gates_lift_type
        FOREIGN KEY (lift_type_id)
        REFERENCES lift_types(id),

    CONSTRAINT fk_order_gates_color_out
        FOREIGN KEY (color_out_id)
        REFERENCES colors(id),

    CONSTRAINT fk_order_gates_cycle_amount
        FOREIGN KEY (cycle_amount_id)
        REFERENCES cycle_amounts(id)
);

CREATE INDEX idx_order_gates_order_id
    ON order_gates USING btree (order_id);

CREATE TABLE order_products (
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    amount INTEGER NOT NULL,

    CONSTRAINT order_products_pkey
        PRIMARY KEY (order_id, product_id),

    CONSTRAINT fk_order_products_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_order_products_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
);

CREATE TABLE order_gate_options (
    order_id BIGINT NOT NULL,
    row_number BIGINT NOT NULL,
    option_id BIGINT NOT NULL,
    amount INTEGER NOT NULL,

    CONSTRAINT order_gate_options_pkey
        PRIMARY KEY (order_id, row_number, option_id),

    CONSTRAINT fk_order_gate_options_order_gate
        FOREIGN KEY (order_id, row_number)
        REFERENCES order_gates(order_id, row_number)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_order_gate_options_option
        FOREIGN KEY (option_id)
        REFERENCES options(id)
);

-- ============================================================
-- Приводы ворот в заказе
-- ============================================================

CREATE TABLE order_gate_manual_drives (
    order_id BIGINT NOT NULL,
    row_number BIGINT NOT NULL,
    chain_length INTEGER,

    CONSTRAINT order_gate_manual_drives_pkey
        PRIMARY KEY (order_id, row_number),

    CONSTRAINT fk_order_gate_manual_drives_order_gate
        FOREIGN KEY (order_id, row_number)
        REFERENCES order_gates(order_id, row_number)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE order_gate_industrial_drives (
    order_id BIGINT NOT NULL,
    row_number BIGINT NOT NULL,
    drive_id BIGINT NOT NULL,

    CONSTRAINT order_gate_industrial_drives_pkey
        PRIMARY KEY (order_id, row_number),

    CONSTRAINT fk_order_gate_industrial_drives_order_gate
        FOREIGN KEY (order_id, row_number)
        REFERENCES order_gates(order_id, row_number)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_order_gate_industrial_drives_drive
        FOREIGN KEY (drive_id)
        REFERENCES industrial_gate_drives(id)
);

CREATE TABLE order_gate_residential_drives (
    order_id BIGINT NOT NULL,
    row_number BIGINT NOT NULL,
    drive_id BIGINT NOT NULL,
    rail_id BIGINT NOT NULL,

    CONSTRAINT order_gate_residential_drives_pkey
        PRIMARY KEY (order_id, row_number),

    CONSTRAINT fk_order_gate_residential_drives_order_gate
        FOREIGN KEY (order_id, row_number)
        REFERENCES order_gates(order_id, row_number)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_order_gate_residential_drives_drive
        FOREIGN KEY (drive_id)
        REFERENCES residential_gate_drives(id),

    CONSTRAINT fk_order_gate_residential_drives_rail
        FOREIGN KEY (rail_id)
        REFERENCES rails(id)
);

-- ============================================================
-- Документы по заказам
-- ============================================================

CREATE TABLE order_bills (
    order_id BIGINT NOT NULL,
    bill_number TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT order_bills_pkey
        PRIMARY KEY (order_id, bill_number),

    CONSTRAINT fk_order_bills_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE order_appendices (
    order_id         BIGINT NOT NULL,
    appendix_number  TEXT NOT NULL,
    path             TEXT NOT NULL,
    created_at       TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT order_appendices_pkey
        PRIMARY KEY (order_id, appendix_number),

    CONSTRAINT fk_order_appendices_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE order_documents (
    order_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT order_documents_pkey
        PRIMARY KEY (order_id, name),

    CONSTRAINT fk_order_documents_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE order_offers (
    order_id BIGINT NOT NULL,
    offer_number TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT order_offers_pkey
        PRIMARY KEY (order_id, offer_number),

    CONSTRAINT fk_order_offers_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

-- ============================================================
-- История изменения цен и наценок
-- ============================================================

CREATE TABLE product_price_history (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    wholesale_price NUMERIC(10,2) NOT NULL,
    retail_price NUMERIC(10,2) NOT NULL,
    set_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT product_price_history_wholesale_price_check
        CHECK (wholesale_price >= 0),

    CONSTRAINT product_price_history_retail_price_check
        CHECK (retail_price >= 0),

    CONSTRAINT product_price_history_product_fkey
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_product_price_history_product_date
    ON product_price_history USING btree (product_id, set_at);

CREATE TABLE option_price_history (
    id BIGSERIAL PRIMARY KEY,
    option_id BIGINT NOT NULL,
    wholesale_price NUMERIC(10,2) NOT NULL,
    retail_price NUMERIC(10,2) NOT NULL,
    set_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT option_price_history_wholesale_price_check
        CHECK (wholesale_price >= 0),

    CONSTRAINT option_price_history_retail_price_check
        CHECK (retail_price >= 0),

    CONSTRAINT option_price_history_option_fkey
        FOREIGN KEY (option_id)
        REFERENCES options(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_option_price_history_option_date
    ON option_price_history USING btree (option_id, set_at);

CREATE TABLE gate_size_price_history (
    id BIGSERIAL PRIMARY KEY,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    gate_type gate_type NOT NULL,
    wholesale_price NUMERIC(10,2) NOT NULL,
    retail_price NUMERIC(10,2) NOT NULL,
    set_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT gate_size_price_history_wholesale_price_check
        CHECK (wholesale_price >= 0),

    CONSTRAINT gate_size_price_history_retail_price_check
        CHECK (retail_price >= 0),

    CONSTRAINT gate_size_price_history_gate_size_fkey
        FOREIGN KEY (width, height, gate_type)
        REFERENCES gate_sizes(width, height, gate_type)
);

CREATE INDEX idx_gate_size_price_history_size_date
    ON gate_size_price_history USING btree (width, height, gate_type, set_at);

CREATE TABLE rail_price_history (
    id BIGSERIAL PRIMARY KEY,
    rail_id BIGINT NOT NULL,
    wholesale_price NUMERIC(10,2) NOT NULL,
    retail_price NUMERIC(10,2) NOT NULL,
    set_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT rail_price_history_wholesale_price_check
        CHECK (wholesale_price >= 0),

    CONSTRAINT rail_price_history_retail_price_check
        CHECK (retail_price >= 0),

    CONSTRAINT rail_price_history_rail_fkey
        FOREIGN KEY (rail_id)
        REFERENCES rails(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_rail_price_history_rail_date
    ON rail_price_history USING btree (rail_id, set_at);

CREATE TABLE residential_gate_drive_price_history (
    id BIGSERIAL PRIMARY KEY,
    drive_id BIGINT NOT NULL,
    wholesale_price NUMERIC(10,2) NOT NULL,
    retail_price NUMERIC(10,2) NOT NULL,
    set_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT residential_gate_drive_price_history_wholesale_price_check
        CHECK (wholesale_price >= 0),

    CONSTRAINT residential_gate_drive_price_history_retail_price_check
        CHECK (retail_price >= 0),

    CONSTRAINT residential_gate_drive_price_history_drive_fkey
        FOREIGN KEY (drive_id)
        REFERENCES residential_gate_drives(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_residential_gate_drive_price_history_drive_date
    ON residential_gate_drive_price_history USING btree (drive_id, set_at);

CREATE TABLE industrial_gate_drive_price_history (
    id BIGSERIAL PRIMARY KEY,
    drive_id BIGINT NOT NULL,
    wholesale_price NUMERIC(10,2) NOT NULL,
    retail_price NUMERIC(10,2) NOT NULL,
    set_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT industrial_gate_drive_price_history_wholesale_price_check
        CHECK (wholesale_price >= 0),

    CONSTRAINT industrial_gate_drive_price_history_retail_price_check
        CHECK (retail_price >= 0),

    CONSTRAINT industrial_gate_drive_price_history_drive_fkey
        FOREIGN KEY (drive_id)
        REFERENCES industrial_gate_drives(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_industrial_gate_drive_price_history_drive_date
    ON industrial_gate_drive_price_history USING btree (drive_id, set_at);

CREATE TABLE lift_type_markup_history (
    id BIGSERIAL PRIMARY KEY,
    lift_type_id BIGINT NOT NULL,
    wholesale_markup NUMERIC(10,2) NOT NULL,
    retail_markup NUMERIC(10,2) NOT NULL,
    set_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT lift_type_markup_history_wholesale_markup_check
        CHECK (wholesale_markup >= 0),

    CONSTRAINT lift_type_markup_history_retail_markup_check
        CHECK (retail_markup >= 0),

    CONSTRAINT lift_type_markup_history_lift_type_fkey
        FOREIGN KEY (lift_type_id)
        REFERENCES lift_types(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_lift_type_markup_history_lift_type_date
    ON lift_type_markup_history USING btree (lift_type_id, set_at);

CREATE TABLE cycle_amount_markup_history (
    id BIGSERIAL PRIMARY KEY,
    cycle_amount_id BIGINT NOT NULL,
    wholesale_markup NUMERIC(10,2) NOT NULL,
    retail_markup NUMERIC(10,2) NOT NULL,
    set_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT cycle_amount_markup_history_wholesale_markup_check
        CHECK (wholesale_markup >= 0),

    CONSTRAINT cycle_amount_markup_history_retail_markup_check
        CHECK (retail_markup >= 0),

    CONSTRAINT cycle_amount_markup_history_cycle_amount_fkey
        FOREIGN KEY (cycle_amount_id)
        REFERENCES cycle_amounts(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_cycle_amount_markup_history_cycle_amount_date
    ON cycle_amount_markup_history USING btree (cycle_amount_id, set_at);

-- ============================================================
-- Функции записи истории цен
-- ============================================================

-- +goose StatementBegin
CREATE FUNCTION log_cycle_amount_markup_history()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
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
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION log_industrial_gate_drive_price_history()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO industrial_gate_drive_price_history (
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
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION log_lift_type_markup_history()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
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
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION log_option_price_history()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO option_price_history (
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
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION log_product_price_history()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO product_price_history (
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
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION log_rail_price_history()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO rail_price_history (
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
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION log_residential_gate_drive_price_history()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO residential_gate_drive_price_history (
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
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION log_gate_size_price_history()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO gate_size_price_history (
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
$$;
-- +goose StatementEnd

-- ============================================================
-- Функция автоматического номера ворот внутри заказа
-- ============================================================

-- +goose StatementBegin
CREATE FUNCTION set_gate_row_number_before_insert()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    SELECT COALESCE(MAX(row_number), 0) + 1
    INTO NEW.row_number
    FROM order_gates
    WHERE order_id = NEW.order_id;

    RETURN NEW;
END;
$$;
-- +goose StatementEnd

-- ============================================================
-- Триггеры истории цен
-- ============================================================

CREATE TRIGGER trg_cycle_amount_markup_history_insert
AFTER INSERT ON cycle_amounts
FOR EACH ROW
EXECUTE FUNCTION log_cycle_amount_markup_history();

CREATE TRIGGER trg_cycle_amount_markup_history_update
AFTER UPDATE OF wholesale_markup, retail_markup ON cycle_amounts
FOR EACH ROW
WHEN (
    OLD.wholesale_markup IS DISTINCT FROM NEW.wholesale_markup
    OR OLD.retail_markup IS DISTINCT FROM NEW.retail_markup
)
EXECUTE FUNCTION log_cycle_amount_markup_history();

CREATE TRIGGER trg_industrial_gate_drive_price_history_insert
AFTER INSERT ON industrial_gate_drives
FOR EACH ROW
EXECUTE FUNCTION log_industrial_gate_drive_price_history();

CREATE TRIGGER trg_industrial_gate_drive_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON industrial_gate_drives
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_industrial_gate_drive_price_history();

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

CREATE TRIGGER trg_option_price_history_insert
AFTER INSERT ON options
FOR EACH ROW
EXECUTE FUNCTION log_option_price_history();

CREATE TRIGGER trg_option_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON options
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_option_price_history();

CREATE TRIGGER trg_product_price_history_insert
AFTER INSERT ON products
FOR EACH ROW
EXECUTE FUNCTION log_product_price_history();

CREATE TRIGGER trg_product_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON products
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_product_price_history();

CREATE TRIGGER trg_rail_price_history_insert
AFTER INSERT ON rails
FOR EACH ROW
EXECUTE FUNCTION log_rail_price_history();

CREATE TRIGGER trg_rail_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON rails
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_rail_price_history();

CREATE TRIGGER trg_residential_gate_drive_price_history_insert
AFTER INSERT ON residential_gate_drives
FOR EACH ROW
EXECUTE FUNCTION log_residential_gate_drive_price_history();

CREATE TRIGGER trg_residential_gate_drive_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON residential_gate_drives
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_residential_gate_drive_price_history();

CREATE TRIGGER trg_gate_size_price_history_insert
AFTER INSERT ON gate_sizes
FOR EACH ROW
EXECUTE FUNCTION log_gate_size_price_history();

CREATE TRIGGER trg_gate_size_price_history_update
AFTER UPDATE OF wholesale_price, retail_price ON gate_sizes
FOR EACH ROW
WHEN (
    OLD.wholesale_price IS DISTINCT FROM NEW.wholesale_price
    OR OLD.retail_price IS DISTINCT FROM NEW.retail_price
)
EXECUTE FUNCTION log_gate_size_price_history();

-- ============================================================
-- Триггер автоматического номера ворот внутри заказа
-- ============================================================

CREATE TRIGGER trg_set_gate_row_number_before_insert
BEFORE INSERT ON order_gates
FOR EACH ROW
EXECUTE FUNCTION set_gate_row_number_before_insert();


-- +goose Down
DROP TRIGGER IF EXISTS trg_set_gate_row_number_before_insert ON order_gates;

DROP TRIGGER IF EXISTS trg_gate_size_price_history_update ON gate_sizes;
DROP TRIGGER IF EXISTS trg_gate_size_price_history_insert ON gate_sizes;

DROP TRIGGER IF EXISTS trg_residential_gate_drive_price_history_update ON residential_gate_drives;
DROP TRIGGER IF EXISTS trg_residential_gate_drive_price_history_insert ON residential_gate_drives;

DROP TRIGGER IF EXISTS trg_rail_price_history_update ON rails;
DROP TRIGGER IF EXISTS trg_rail_price_history_insert ON rails;

DROP TRIGGER IF EXISTS trg_product_price_history_update ON products;
DROP TRIGGER IF EXISTS trg_product_price_history_insert ON products;

DROP TRIGGER IF EXISTS trg_option_price_history_update ON options;
DROP TRIGGER IF EXISTS trg_option_price_history_insert ON options;

DROP TRIGGER IF EXISTS trg_lift_type_markup_history_update ON lift_types;
DROP TRIGGER IF EXISTS trg_lift_type_markup_history_insert ON lift_types;

DROP TRIGGER IF EXISTS trg_industrial_gate_drive_price_history_update ON industrial_gate_drives;
DROP TRIGGER IF EXISTS trg_industrial_gate_drive_price_history_insert ON industrial_gate_drives;

DROP TRIGGER IF EXISTS trg_cycle_amount_markup_history_update ON cycle_amounts;
DROP TRIGGER IF EXISTS trg_cycle_amount_markup_history_insert ON cycle_amounts;

DROP FUNCTION IF EXISTS set_gate_row_number_before_insert();
DROP FUNCTION IF EXISTS log_gate_size_price_history();
DROP FUNCTION IF EXISTS log_residential_gate_drive_price_history();
DROP FUNCTION IF EXISTS log_rail_price_history();
DROP FUNCTION IF EXISTS log_product_price_history();
DROP FUNCTION IF EXISTS log_option_price_history();
DROP FUNCTION IF EXISTS log_lift_type_markup_history();
DROP FUNCTION IF EXISTS log_industrial_gate_drive_price_history();
DROP FUNCTION IF EXISTS log_cycle_amount_markup_history();

DROP TABLE IF EXISTS order_gate_residential_drives;
DROP TABLE IF EXISTS order_gate_industrial_drives;
DROP TABLE IF EXISTS order_gate_manual_drives;
DROP TABLE IF EXISTS order_gate_options;
DROP TABLE IF EXISTS order_products;
DROP TABLE IF EXISTS order_gates;

DROP TABLE IF EXISTS order_offers;
DROP TABLE IF EXISTS order_documents;
DROP TABLE IF EXISTS order_appendices;
DROP TABLE IF EXISTS order_bills;
DROP TABLE IF EXISTS order_contracts;
DROP TABLE IF EXISTS orders;

DROP TABLE IF EXISTS cycle_amount_markup_history;
DROP TABLE IF EXISTS lift_type_markup_history;
DROP TABLE IF EXISTS industrial_gate_drive_price_history;
DROP TABLE IF EXISTS residential_gate_drive_price_history;
DROP TABLE IF EXISTS rail_price_history;
DROP TABLE IF EXISTS gate_size_price_history;
DROP TABLE IF EXISTS option_price_history;
DROP TABLE IF EXISTS product_price_history;

DROP TABLE IF EXISTS manual_drive_prices;
DROP TABLE IF EXISTS residential_gate_drives;
DROP TABLE IF EXISTS industrial_gate_drives;
DROP TABLE IF EXISTS rails;
DROP TABLE IF EXISTS options;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS gate_sizes;
DROP TABLE IF EXISTS lift_types;
DROP TABLE IF EXISTS cycle_amounts;
DROP TABLE IF EXISTS colors;

DROP TABLE IF EXISTS dealer_registration_requests;
DROP TABLE IF EXISTS dealer_manager_assignments;
DROP TABLE IF EXISTS dealer_contract;
DROP TABLE IF EXISTS dealers;
DROP TABLE IF EXISTS companies;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS role;
DROP TYPE IF EXISTS registration_request_status;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS gate_type;
DROP TYPE IF EXISTS drive_type;
