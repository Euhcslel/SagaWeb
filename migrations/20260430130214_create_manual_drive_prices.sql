-- +goose Up
CREATE TABLE
    manual_drive_prices (
        id BIGSERIAL PRIMARY KEY,
        chain_meter_retail_price NUMERIC(10, 2) NOT NULL CHECK (chain_meter_retail_price >= 0),
        chain_meter_wholesale_price NUMERIC(10, 2) NOT NULL CHECK (chain_meter_wholesale_price >= 0),
        rcp_retail_price NUMERIC(10, 2) NOT NULL CHECK (rcp_retail_price >= 0),
        rcp_wholesale_price NUMERIC(10, 2) NOT NULL CHECK (rcp_wholesale_price >= 0),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- +goose Down
DROP TABLE IF EXISTS manual_drive_prices;
