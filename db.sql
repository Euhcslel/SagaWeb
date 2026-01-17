--Пользователи

CREATE TABLE statuses (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE INDEX idx_statuses_name ON statuses (name);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    fullname TEXT NOT NULL,
    email TEXT,
    phone_number TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_phone ON users(phone_number);

CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE INDEX idx_companies_name ON companies(name);

CREATE TABLE dealers (
    id BIGINT PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    company_id INT REFERENCES companies (id) ON DELETE SET NULL
);

CREATE INDEX idx_dealers_company ON dealers(company_id);

CREATE TABLE sessions (
    user_id BIGINT NOT NULL PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + INTERVAL '7 days'
);

CREATE INDEX idx_sessions_token ON sessions(token_hash);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);

CREATE TABLE managers_and_dealers (
    manager_id BIGINT NOT NULL REFERENCES users (id),
    dealer_id BIGINT NOT NULL REFERENCES dealers (id),
    PRIMARY KEY (manager_id, dealer_id)
);

CREATE INDEX idx_managers_and_dealers_manager ON managers_and_dealers (manager_id);
CREATE INDEX idx_managers_and_dealers_dealer ON managers_and_dealers (dealer_id);

CREATE TABLE dealers_reg_requests (
    id BIGSERIAL PRIMARY KEY,
    company TEXT,
    fullname TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    status_id INT NOT NULL REFERENCES statuses (id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Продажи

CREATE TABLE sales (
    id BIGSERIAL PRIMARY KEY,
    client BIGINT NOT NULL REFERENCES users(id),
    manager BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sales_client ON sales (client);
CREATE INDEX idx_sales_manager ON sales (manager);

CREATE TABLE units (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE INDEX idx_units_name ON units (name);

CREATE TABLE colors (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    hex TEXT NOT NULL UNIQUE
);

CREATE INDEX idx_colors_code ON colors (code);

CREATE TABLE cycle_amount (
    id BIGSERIAL PRIMARY KEY,
    amount INT NOT NULL UNIQUE,
    wholesale_markup INT,
    retail_markup INT
);

CREATE INDEX idx_cycle_amount_amount ON cycle_amount (amount);

CREATE TABLE lift_types (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    wholesale_markup INT,
    retail_markup INT
);

CREATE INDEX idx_lift_types_name ON lift_types (name);

CREATE TABLE gate_type (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE sales_and_gates (
    sale_id BIGINT NOT NULL REFERENCES sales (id) ON DELETE CASCADE,
    row_number BIGINT NOT NULL,
    gate_type INT NOT NULL REFERENCES gate_type (id),
    width INT NOT NULL,
    height INT NOT NULL,
    lintel_height INT NOT NULL,
    lift_type_id BIGINT NOT NULL REFERENCES lift_types (id),
    color_in_id BIGINT NOT NULL REFERENCES colors (id),
    color_out_id BIGINT NOT NULL REFERENCES colors (id),
    cycle_amount_id BIGINT NOT NULL REFERENCES cycle_amount (id),
    total_price INT NOT NULL,
    status_id INT NOT NULL REFERENCES statuses (id),
    PRIMARY KEY(sale_id, row_number),
    CONSTRAINT positive_dimensions CHECK (width > 0 AND height > 0),
    CONSTRAINT non_negative_lintel CHECK (lintel_height >= 0),
    CONSTRAINT positive_price CHECK (total_price > 0)
);

CREATE OR REPLACE FUNCTION change_rows_numbers()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE sales_and_gates
    SET row_number = row_number - 1
    WHERE sale_id = OLD.sale_id AND row_number > OLD.row_number;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER sales_and_gates_after_delete
AFTER DELETE ON sales_and_gates
FOR EACH ROW
EXECUTE FUNCTION change_rows_numbers();

CREATE INDEX idx_sales_gates_lift_type ON sales_and_gates (lift_type_id);
CREATE INDEX idx_sales_gates_color_in ON sales_and_gates (color_in_id);
CREATE INDEX idx_sales_gates_color_out ON sales_and_gates (color_out_id);
CREATE INDEX idx_sales_gates_cycle_amount ON sales_and_gates (cycle_amount_id);
CREATE INDEX idx_sales_gates_status ON sales_and_gates (status_id);

-- Ворота и привода

CREATE TABLE industrial_gate_drives (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    unit_id INT REFERENCES units (id),
    wholesale_price INT CHECK (wholesale_price >= 0),
    retail_price INT CHECK (retail_price >= 0),
    specifications TEXT
);

CREATE INDEX idx_industrial_drives_name ON industrial_gate_drives (name);

CREATE TABLE rails (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    unit_id INT REFERENCES units (id),
    wholesale_price INT CHECK (wholesale_price >= 0),
    retail_price INT CHECK (retail_price >= 0),
    specifications TEXT
);

CREATE INDEX idx_rails_name ON rails (name);

CREATE TABLE residential_gate_drives (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    unit_id INT REFERENCES units (id),
    wholesale_price INT CHECK (wholesale_price >= 0),
    retail_price INT CHECK (retail_price >= 0),
    specifications TEXT
);

CREATE INDEX idx_residential_drives_name ON residential_gate_drives (name);

CREATE TABLE industrial_gates_and_sales_drive (
    sale_id BIGINT NOT NULL,
    row_number BIGINT NOT NULL,
    drive_id INT NOT NULL REFERENCES industrial_gate_drives (id),
    FOREIGN KEY (sale_id, row_number)
    REFERENCES sales_and_gates (sale_id, row_number)
    ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (sale_id, row_number)
);

CREATE TABLE residential_gates_and_sales_drive_rail (
    sale_id BIGINT NOT NULL,
    row_number BIGINT NOT NULL,
    drive_id INT NOT NULL REFERENCES residential_gate_drives (id),
    rail_id INT NOT NULL REFERENCES rails (id),
    FOREIGN KEY (sale_id, row_number)
    REFERENCES sales_and_gates (sale_id, row_number)
    ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (sale_id, row_number)
);

CREATE TABLE gates_and_sales_manual_drive (
    sale_id BIGINT NOT NULL,
    row_number BIGINT NOT NULL,
    chain_length INT CHECK (chain_length >= 0),
    FOREIGN KEY (sale_id, row_number)
    REFERENCES sales_and_gates (sale_id, row_number)
    ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (sale_id, row_number)
);

CREATE OR REPLACE FUNCTION delete_old_drives()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_TABLE_NAME = 'industrial_gates_and_sales_drive' THEN
        DELETE FROM residential_gates_and_sales_drive_rail
        WHERE sale_id = NEW.sale_id AND row_number = NEW.row_number;

        DELETE FROM gates_and_sales_manual_drive
        WHERE sale_id = NEW.sale_id AND row_number = NEW.row_number;

    ELSIF TG_TABLE_NAME = 'residential_gates_and_sales_drive_rail' THEN
        DELETE FROM industrial_gates_and_sales_drive
        WHERE sale_id = NEW.sale_id AND row_number = NEW.row_number;

        DELETE FROM gates_and_sales_manual_drive
        WHERE sale_id = NEW.sale_id AND row_number = NEW.row_number;

    ELSIF TG_TABLE_NAME = 'gates_and_sales_manual_drive' THEN
        DELETE FROM industrial_gates_and_sales_drive
        WHERE sale_id = NEW.sale_id AND row_number = NEW.row_number;

        DELETE FROM residential_gates_and_sales_drive_rail
        WHERE sale_id = NEW.sale_id AND row_number = NEW.row_number;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_industrial_drives
BEFORE INSERT ON industrial_gates_and_sales_drive
FOR EACH ROW
EXECUTE FUNCTION delete_old_drives();

CREATE TRIGGER before_insert_residential_drives
BEFORE INSERT ON residential_gates_and_sales_drive_rail
FOR EACH ROW
EXECUTE FUNCTION delete_old_drives();

CREATE TRIGGER before_insert_manual_drives
BEFORE INSERT ON gates_and_sales_manual_drive
FOR EACH ROW
EXECUTE FUNCTION delete_old_drives();

-- Монтажи

CREATE TABLE montage_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE INDEX idx_montage_types_name ON montage_types (name);

CREATE TABLE montages_and_sales (
    sale_id BIGINT NOT NULL PRIMARY KEY REFERENCES sales (id) ON DELETE CASCADE,
    type_id INT REFERENCES montage_types (id),
    montage_date DATE NOT NULL,
    price INT CHECK (price > 0)
);

CREATE INDEX idx_montages_sales_type ON montages_and_sales (type_id);
CREATE INDEX idx_montages_sales_date ON montages_and_sales (montage_date);

-- Дополнительные опции

CREATE TABLE options (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    unit_id INT REFERENCES units (id),
    wholesale_price INT CHECK (wholesale_price >= 0),
    retail_price INT CHECK (retail_price >= 0),
    for_sale BOOLEAN NOT NULL,
    condition TEXT
);

CREATE INDEX idx_options_for_sale ON options (for_sale);
CREATE INDEX idx_options_name ON options (name);

CREATE TABLE gates_and_sales_options (
    sale_id BIGINT NOT NULL,
    row_number BIGINT NOT NULL,
    option_id INT NOT NULL REFERENCES options (id),
    amount INT NOT NULL CHECK (amount > 0),
    FOREIGN KEY (sale_id, row_number)
    REFERENCES sales_and_gates (sale_id, row_number)
    ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (sale_id, row_number)
);

CREATE TABLE standart_equipment (
    gate_type_id INT NOT NULL REFERENCES gate_type (id) ON DELETE CASCADE,
    option_id INT NOT NULL REFERENCES options (id),
    amount INT CHECK (amount > 0),
    PRIMARY KEY(gate_type_id, option_id)
);

-- Товары

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    unit_id INT REFERENCES units (id),
    wholesale_price INT CHECK (wholesale_price >= 0),
    retail_price INT CHECK (retail_price >= 0)
);

CREATE INDEX idx_products_name ON products (name);

CREATE TABLE sales_and_products (
    sale_id BIGINT NOT NULL REFERENCES sales (id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products (id),
    amount INT NOT NULL CHECK (amount > 0),
    PRIMARY KEY(sale_id, product_id)
);

-- Документы

CREATE TABLE sales_and_offers (
    sale_id BIGINT REFERENCES sales (id) ON DELETE CASCADE,
    offer_number TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(sale_id, offer_number)
);

CREATE TABLE sales_and_contracts (
    sale_id BIGINT REFERENCES sales (id) ON DELETE CASCADE,
    contract_number TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(sale_id, contract_number)
);

CREATE TABLE sales_and_bills (
    sale_id BIGINT REFERENCES sales (id) ON DELETE CASCADE,
    bill_number TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(sale_id, bill_number)
);

CREATE TABLE sales_and_documents (
    sale_id BIGINT REFERENCES sales (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(sale_id, name)
);
