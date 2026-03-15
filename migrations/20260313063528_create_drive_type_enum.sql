-- +goose Up
CREATE TYPE drive_type AS ENUM ('industrial', 'residential', 'manual');

ALTER TABLE sales_and_gates
ADD COLUMN drive_type drive_type NOT NULL DEFAULT 'industrial';

UPDATE sales_and_gates sg
SET
    drive_type = 'residential'
WHERE
    sg.gate_type = 'res';

-- +goose Down
ALTER TABLE sales_and_gates
DROP COLUMN drive_type;

DROP TYPE drive_type;
