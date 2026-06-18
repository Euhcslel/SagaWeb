-- +goose Up

-- ============================================================
-- Пользователи, компании, дилеры
-- ============================================================

INSERT INTO users (id, fullname, email, phone_number, created_at, password_hash, role) VALUES
(41, 'Анна Иванова', 'dealer1@test.com', '85345678904', '2026-02-25 19:36:25.687363+05', decode('243261243130246a56587266525451753955307073762f4b6f38363065436f46355a3852357237505949737a472e4c34336d66307550632f69597a36', 'hex'), 'dealer'),
(38, 'Алексей Алексеев', 'manager@test.com', '89345678901', '2026-02-25 19:36:25.687363+05', decode('24326124313024434c4b585a4436306a756169372f4a68784f483843653143376d67355843486a7741395858476d583057596354533138642e705465', 'hex'), 'manager'),
(39, 'Анастасия Артемова', 'admin@test.com', '82345678902', '2026-02-25 19:36:25.687363+05', decode('24326124313024764e35707a38794a425579696f515457626d72666c4f47636a516b455a626e515148736e6852393563433458694476325869574a53', 'hex'), 'admin'),
(57, 'Елизавета Аникина', 'dealer2@test.com', '84395834593', '2026-05-02 21:19:31.199579+05', decode('243261243130246a56587266525451753955307073762f4b6f38363065436f46355a3852357237505949737a472e4c34336d66307550632f69597a36', 'hex'), 'dealer'),
(60, 'Дмитрий Логистов', 'logist@test.com', '89001234567', '2026-02-25 19:36:25.687363+05', decode('2432612431302430592e6d2e6435725146343952665366512e4c77322e42684a5970497038656e58566c3150515954712e71475849633543695a584b', 'hex'), 'logistician');

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

INSERT INTO colors (code, hex) VALUES
('RAL8017', '45322f'),
('RAL7016', '4f5356'),
('RAL5005', '2A5F90'),
('RAL9016', 'F0efea');

INSERT INTO cycle_amounts (amount, wholesale_markup, retail_markup) VALUES
('25000', 0.00, 0.00),
('50000', 3.00, 3.00),
('100000', 7.00, 7.00);

INSERT INTO lift_types (name, min_headroom, max_headroom, wholesale_markup, retail_markup) VALUES
('Вертикальный подъем (вал сверху)', 0, 0, 10.00, 10.00),
('Высокий подъем (вал сверху)', 600, 0, 7.00, 7.00),
('Вертикальный подъем (вал снизу)', 0, 0, 15.00, 15.00),
('Высокий подъем (вал снизу)', 1500, 0, 15.00, 15.00);

INSERT INTO products (name, wholesale_price, retail_price) VALUES
('Устройство защиты от обрыва троса', 3500.00, 5400.00),
('Устройство защиты от обрыва пружины', 3500.00, 5400.00),
('Пульт управления SG 433 МГц 4 кн.', 960.00, 1200.00),
('Приемник SGRE-2 внешний 1-о канальный', 1520.00, 1900.00),
('Кромка безопасности (оптосенсоры)', 4700.00, 7500.00);

INSERT INTO options (name, wholesale_price, retail_price) VALUES
('Калитка для секционных ворот', 57750.00, 69300.00),
('Окно 660х330', 5400.00, 6700.00),
('Подвесы для притолоки 500-1000мм', 2000.00, 2400.00),
('Подвесы для притолоки 1000-1500мм', 4000.00, 4800.00);

INSERT INTO rails (name, wholesale_price, retail_price, specifications) VALUES
('Направляющая BR-3300 с ремнем', 4280.00, 5100.00, 'L=3300мм, H=2500мм'),
('Направляющая BR-3600 с ремнем', 4770.00, 5650.00, 'L=3600мм, H=2800мм'),
('Направляющая BR-4200 с ремнем', 5740.00, 6800.00, 'L=4200мм, H=3400мм');

INSERT INTO industrial_gate_drives (name, wholesale_price, retail_price, specifications) VALUES
('Комплект промышленного привода SGC40', 44800.23, 56000.65,'S=20м.кв., 40N.m, 220В, 450W Диаметр вала 25,4мм'),
('Комплект промышленного привода SGC60', 50000.00, 62500.00, 'S=30м.кв., 60N.m, 220В, 650W Диаметр вала 25,4мм'),
('Комплект промышленного привода SGC90', 60400.00, 75500.00, 'S=40м.кв., 90N.m, 220В, 1100W Диаметр вала 25,4мм'),
('Комплект промышленного привода SGC120', 64000.00, 80000.00, 'S=50м.кв., 120N.m, 220В, 1500W Диаметр вала 32мм');

INSERT INTO residential_gate_drives (name, wholesale_price, retail_price, specifications) VALUES
('Привод Saga D-600 (Motor: 24V, 220В, 600Н)', 10320.00, 12900.00, 'S=9м.кв.'),
('Привод Saga D-800 (Motor: 24V, 220В, 800Н)', 11680.00, 14600.00, 'S=12м.кв.'),
('Привод Saga D-1000 (Motor: 24V, 220В, 800Н)', 12560.00, 15700.00, 'S=15м.кв.'),
('Привод Saga D-1200 (Motor: 24V, 220В, 1200Н)', 13760.00, 17200.00, 'S=18м.кв.');

INSERT INTO manual_drive_prices (
    chain_meter_retail_price,
    chain_meter_wholesale_price,
    rcp_retail_price,
    rcp_wholesale_price,
    created_at
) VALUES
(540.00, 450.00, 14500.00, 10200.00, '2026-04-30 18:10:05.041553+05');

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
TRUNCATE TABLE
    order_gate_options,
    order_gate_industrial_drives,
    order_gate_residential_drives,
    order_gate_manual_drives,
    order_gates,
    order_products,
    order_bills,
    order_appendices,
    order_documents,
    orders,
    dealer_manager_assignments,
    dealers,
    users,
    companies,
    colors,
    cycle_amounts,
    lift_types,
    products,
    options,
    rails,
    industrial_gate_drives,
    residential_gate_drives,
    manual_drive_prices
RESTART IDENTITY CASCADE;
