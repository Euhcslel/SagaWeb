INSERT INTO cycle_amount (amount, wholesale_markup, retail_markup) VALUES
(25000, 0, 0), (50000, 3, 3), (100000, 7, 7);

INSERT INTO units (name) VALUES ('шт.'), ('комплект'), ('м.п.');

INSERT INTO colors (code, hex) VALUES
('RAL8017', '45322f'),
('RAL7016', '4f5356'),
('RAL7004', '969991'),
('RAL5005', '2A5F90'),
('RAL9016', 'F0efea');

INSERT INTO lift_types (name, wholesale_markup, retail_markup, min_headroom, max_headroom) VALUES
('Вертикальный подъем (вал сверху)', 10, 10, 0, 0),
('Высокий подъем (вал сверху)', 7, 7, 600, 0),
('Вертикальный подъем (вал снизу)', 15, 15, 0, 0),
('Высокий подъем (вал снизу)', 15, 15, 1500, 0);

INSERT INTO gate_type (name) VALUES ('Промышленные ворота'), ('Бытовые ворота');

INSERT INTO statuses (name) VALUES
('В производстве'),
('Подтвержден'),
('Ожидает подтверждения'),
('Оплачен'),
('Ожидает оплаты'),
('Завершен');

INSERT INTO companies (name) VALUES
('GateTrade Inc.'),
('GateMarket LLC');

INSERT INTO users (username, fullname, email, phone_number, password, role_id) VALUES
('john_doe', 'John Doe', 'john@example.com', '+12345678901', ENCODE(DIGEST('password123', 'sha256'), 'hex'), (SELECT id FROM roles WHERE name = 'manager')),
('jane_smith', 'Jane Smith', 'jane@example.com', '+12345678902', ENCODE(DIGEST('password456', 'sha256'), 'hex'), (SELECT id FROM roles WHERE name = 'admin')),
('mike_client', 'Mike Wilson', 'mike@example.com', '+12345678903', ENCODE(DIGEST('password789', 'sha256'), 'hex'), (SELECT id FROM roles WHERE name = 'client')),
('sarah_dealer', 'Sarah Johnson', 'sarah@example.com', '+12345678904', ENCODE(DIGEST('dealer_pass1', 'sha256'), 'hex'), (SELECT id FROM roles WHERE name = 'dealer')),
('alex_dealer', 'Alex Brown', 'alex@example.com', '+12345678905', ENCODE(DIGEST('dealer_pass2', 'sha256'), 'hex'), (SELECT id FROM roles WHERE name = 'dealer'));

INSERT INTO dealers (user_id, company_id) VALUES
((SELECT id FROM users WHERE username = 'sarah_dealer'), (SELECT id FROM companies WHERE name = 'GateTrade Inc.')),
((SELECT id FROM users WHERE username = 'alex_dealer'), (SELECT id FROM companies WHERE name = 'GateMarket LLC'));

 INSERT INTO sales (client_id, manager_id) VALUES
((SELECT id FROM users WHERE username = 'sarah_dealer'),
 (SELECT id FROM users WHERE username = 'john_doe')),
 ((SELECT id FROM users WHERE username = 'sarah_dealer'),
 (SELECT id FROM users WHERE username = 'john_doe')),
((SELECT id FROM users WHERE username = 'alex_dealer'),
 (SELECT id FROM users WHERE username = 'jane_smith'));

INSERT INTO dealers_reg_requests (company, fullname, email, phone_number, status_id) VALUES
('Gate Empire LLC', 'Robert Wilson', 'robert@example.com', '+19876543210',
 (SELECT id FROM statuses WHERE name = 'Ожидает подтверждения')),
('Gate World Co', 'Emily Davis', 'emily@example.com', '+19876543211',
 (SELECT id FROM statuses WHERE name = 'Ожидает подтверждения')),
('Gate Masters', 'David Lee', 'david@example.com', '+19876543212',
 (SELECT id FROM statuses WHERE name = 'Ожидает подтверждения'));

 INSERT INTO managers_and_dealers (manager_id, dealer_id) VALUES
((SELECT id FROM users WHERE username = 'john_doe'),
 (SELECT id FROM users WHERE username = 'sarah_dealer')),
((SELECT id FROM users WHERE username = 'jane_smith'),
 (SELECT id FROM users WHERE username = 'alex_dealer'));

INSERT INTO residential_gate_drives (name, unit_id, wholesale_price, retail_price, specifications) VALUES
('Привод Saga D-600 (Motor: 24V, 220В, 600Н)', (SELECT id FROM units WHERE name = 'шт.'), 10320, 12900, 'S=9м.кв.');

INSERT INTO industrial_gate_drives (name, unit_id, wholesale_price, retail_price, specifications) VALUES
('Комплект промышленного частотного сервопривода SGC40, в составе осевой сервопривод 40Н, цепь 8 м.п.,
пост четырехпозиционный DW4, энкодер, система перевода в ручной режим управления',
(SELECT id FROM units WHERE name = 'шт.'), 44800, 56000, 'S=20м.кв., 40N.m, 220В, 450W
Диаметр вала 25,4мм');

insert into roles (name) values
('admin'), ('client'), ('dealer');

INSERT INTO options (name, unit_id, wholesale_price, retail_price, for_sale, condition) VALUES
('Калитка для секционных ворот, порог 100-150мм (в компл. с мех. замком и доводчиком)', (SELECT id FROM units WHERE name = 'шт.'), 57750, 69300, true, ''),
('Окно 660х330', (SELECT id FROM units WHERE name = 'шт.'), 5400, 6700, true, '');

INSERT INTO products (name, unit_id, wholesale_price, retail_price) VALUES
('Устройство защиты от обрыва троса', (SELECT id FROM units WHERE name = 'комплект'), 3500, 5400),
('Устройство защиты от обрыва пружины', (SELECT id FROM units WHERE name = 'комплект'), 3500, 5400);

