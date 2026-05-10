-- +goose Up

-- ============================================================
-- Пользователи, компании, дилеры
-- ============================================================

INSERT INTO users (id, fullname, email, phone_number, created_at, password_hash, role) VALUES
(41, 'Анна Иванова', 'annaivanova@example.com', '85345678904', '2026-02-25 19:36:25.687363+05', decode('243261243130243170546f30396276746b342e782f796f6761356a7365687a4b57624a5154746c5a532e764e6a4145676d4b4c7a6351347469306a69', 'hex'), 'dealer'),
(38, 'Алексей Алексеев', 'alexeev@example.com', '89345678901', '2026-02-25 19:36:25.687363+05', decode('243261243130245569346245706c33396e4f4767626c305274584438755453706b334f6c364b654a4578706f536d3163772f5831306b4f5972446c79', 'hex'), 'manager'),
(39, 'Анастасия Артемова', 'artemova@example.com', '82345678902', '2026-02-25 19:36:25.687363+05', decode('24326124313024644c634939514e585858726536774b32524d364a5a7566525436346639704a374749714f595a676557444e6942424d356a3359734f', 'hex'), 'admin'),
(57, 'Елизавета Аникина', '9655460724@mail.ru', '84395834593', '2026-05-02 21:19:31.199579+05', decode('2432612431302430627133394c6e6572352e756f6b67494136415a6f4f5648653139514951776b4b644573727362332e666250706b61556f6d786643', 'hex'), 'dealer');

INSERT INTO companies (id, name) VALUES
(10, 'GateMarket LLC'),
(9, 'GateTrаde Inc.'),
(12, 'ИП Новикова'),
(13, 'ИП Купцова');

INSERT INTO dealers (user_id, company_id, address) VALUES
(41, 9, 'Тула, ул. Советская, д. 74'),
(57, 12, 'Екатеринбург, ул. Космонавтов, д. 15');

INSERT INTO dealer_manager_assignments (manager_id, dealer_id) VALUES
(38, 41),
(38, 57);

-- ============================================================
-- Справочники
-- ============================================================

INSERT INTO colors (id, code, hex) VALUES
(16, 'RAL8017', '45322f'),
(17, 'RAL7016', '4f5356'),
(19, 'RAL5005', '2A5F90'),
(20, 'RAL9016', 'F0efea'),
(21, 'RKKGrRR', 'F6f6f9');

INSERT INTO cycle_amounts (id, amount, wholesale_markup, retail_markup) VALUES
(13, '25000', 0.00, 0.00),
(14, '50000', 3.00, 3.00),
(15, '100000', 7.00, 7.00);

INSERT INTO lift_types (id, name, min_headroom, max_headroom, wholesale_markup, retail_markup) VALUES
(17, 'Вертикальный подъем (вал сверху)', 0, 0, 10.00, 10.00),
(18, 'Высокий подъем (вал сверху)', 600, 0, 7.00, 7.00),
(19, 'Вертикальный подъем (вал снизу)', 0, 0, 15.00, 15.00),
(20, 'Высокий подъем (вал снизу)', 1500, 0, 15.00, 15.00);

INSERT INTO products (id, name, wholesale_price, retail_price) VALUES
(5, 'Устройство защиты от обрыва троса', 3500.00, 5400.00),
(6, 'Устройство защиты от обрыва пружины', 3500.00, 5400.00);

INSERT INTO options (id, name, wholesale_price, retail_price) VALUES
(7, 'Калитка для секционных ворот, порог 100-150мм (в компл. с мех. замком и доводчиком)', 57750.00, 69300.00),
(8, 'Окно 660х330', 5400.00, 6700.00);

INSERT INTO rails (id, name, wholesale_price, retail_price, specifications) VALUES
(4, 'rail1', 1000.00, 2000.00, 'none');

INSERT INTO industrial_gate_drives (id, name, wholesale_price, retail_price, specifications) VALUES
(3, 'Комплект промышленного частотного сервопривода SGC40, в составе осевой сервопривод 40Н, цепь 8 м.п., пост четырехпозиционный DW4, энкодер, система перевода в ручной режим управления', 44800.23, 56000.65,'S=20м.кв., 40N.m, 220В, 450W Диаметр вала 25,4мм');

INSERT INTO residential_gate_drives (id, name, wholesale_price, retail_price, specifications) VALUES
(3, 'Привод Saga D-600 (Motor: 24V, 220В, 600Н)', 10320.00, 12900.00, 'S=9м.кв.');

INSERT INTO manual_drive_prices (
    id,
    chain_meter_retail_price,
    chain_meter_wholesale_price,
    rcp_retail_price,
    rcp_wholesale_price,
    created_at
) VALUES
(1, 540.00, 450.00, 14500.00, 10200.00, '2026-04-30 18:10:05.041553+05');

SELECT setval(pg_get_serial_sequence('users', 'id'), COALESCE((SELECT MAX(id) FROM users), 1), true);
SELECT setval(pg_get_serial_sequence('companies', 'id'), COALESCE((SELECT MAX(id) FROM companies), 1), true);
SELECT setval(pg_get_serial_sequence('colors', 'id'), COALESCE((SELECT MAX(id) FROM colors), 1), true);
SELECT setval(pg_get_serial_sequence('cycle_amounts', 'id'), COALESCE((SELECT MAX(id) FROM cycle_amounts), 1), true);
SELECT setval(pg_get_serial_sequence('lift_types', 'id'), COALESCE((SELECT MAX(id) FROM lift_types), 1), true);
SELECT setval(pg_get_serial_sequence('products', 'id'), COALESCE((SELECT MAX(id) FROM products), 1), true);
SELECT setval(pg_get_serial_sequence('options', 'id'), COALESCE((SELECT MAX(id) FROM options), 1), true);
SELECT setval(pg_get_serial_sequence('rails', 'id'), COALESCE((SELECT MAX(id) FROM rails), 1), true);
SELECT setval(pg_get_serial_sequence('industrial_gate_drives', 'id'), COALESCE((SELECT MAX(id) FROM industrial_gate_drives), 1), true);
SELECT setval(pg_get_serial_sequence('residential_gate_drives', 'id'), COALESCE((SELECT MAX(id) FROM residential_gate_drives), 1), true);
SELECT setval(pg_get_serial_sequence('manual_drive_prices', 'id'), COALESCE((SELECT MAX(id) FROM manual_drive_prices), 1), true);

-- +goose Down

-- ============================================================
-- Удаление пользовательских данных
-- ============================================================

DELETE FROM dealer_manager_assignments
WHERE manager_id = 38
  AND dealer_id IN (41, 57);

DELETE FROM dealers
WHERE user_id IN (41, 57);

DELETE FROM users
WHERE id IN (38, 39, 41, 57);

DELETE FROM companies
WHERE id IN (9, 10, 12, 13);

-- ============================================================
-- Удаление справочников
-- ============================================================

DELETE FROM manual_drive_prices
WHERE id = 1;

DELETE FROM residential_gate_drives
WHERE id IN (3, 4);

DELETE FROM industrial_gate_drives
WHERE id = 3;

DELETE FROM rails
WHERE id = 4;

DELETE FROM options
WHERE id IN (7, 8);

DELETE FROM products
WHERE id IN (5, 6);

DELETE FROM lift_types
WHERE id IN (17, 18, 19, 20);

DELETE FROM cycle_amounts
WHERE id IN (13, 14, 15);

DELETE FROM colors
WHERE id IN (16, 17, 18, 19, 20, 21);
