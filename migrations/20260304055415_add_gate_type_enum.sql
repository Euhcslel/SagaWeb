-- +goose Up
CREATE TYPE gate_type AS ENUM ('ind', 'res');

ALTER TABLE sales_and_gates
ADD COLUMN gate_type gate_type NOT NULL DEFAULT 'ind';

UPDATE sales_and_gates sg
SET
    gate_type = 'res'
FROM
    gate_types gt
WHERE
    sg.gate_type_id = gt.id
    AND gt.name = 'Бытовые ворота';

ALTER TABLE sales_and_gates
DROP COLUMN gate_type_id;

ALTER TABLE sizes
ADD COLUMN gate_type gate_type NOT NULL DEFAULT 'ind';

UPDATE sizes s
SET
    gate_type = 'res'
FROM
    gate_types gt
WHERE
    s.gate_type_id = gt.id
    AND gt.name = 'Бытовые ворота';

ALTER TABLE sizes
DROP COLUMN gate_type_id;

DROP TABLE gate_types;

-- +goose Down
CREATE TABLE
    gate_types (id SERIAL PRIMARY KEY, name TEXT NOT NULL);

INSERT INTO
    gate_types (name)
VALUES
    ('Промышленные ворота'),
    ('Бытовые ворота');

ALTER TABLE sales_and_gates
ADD COLUMN gate_type_id INT;

UPDATE sales_and_gates sg
SET
    gate_type_id = gt.id
FROM
    gate_types gt
WHERE
    (
        sg.gate_type = 'res'
        AND gt.name = 'Бытовые ворота'
    )
    OR (
        sg.gate_type = 'ind'
        AND gt.name = 'Промышленные ворота'
    );

ALTER TABLE sales_and_gates
DROP COLUMN gate_type;

ALTER TABLE sales_and_gates
ALTER COLUMN gate_type_id
SET
    NOT NULL;

ALTER TABLE sizes
ADD COLUMN gate_type_id INT;

UPDATE sizes s
SET
    gate_type_id = gt.id
FROM
    gate_types gt
WHERE
    (
        s.gate_type = 'res'
        AND gt.name = 'Бытовые ворота'
    )
    OR (
        s.gate_type = 'ind'
        AND gt.name = 'Промышленные ворота'
    );

ALTER TABLE sizes
DROP COLUMN gate_type;

ALTER TABLE sizes
ALTER COLUMN gate_type_id
SET
    NOT NULL;

DROP TYPE gate_type;
