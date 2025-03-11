-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE store(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    description TEXT NOT NULL,
    default_currency_id INTEGER NOT NULL,

    CONSTRAINT store_default_currency_id_fk FOREIGN KEY(default_currency_id) REFERENCES currency(id)
);

CREATE INDEX store_default_currency_id_idx ON store(default_currency_id);

CREATE TYPE day_of_week AS ENUM (
    'Monday',
    'Tuesday',
    'Wednesday',
    'Thursday',
    'Friday',
    'Saturday',
    'Sunday'
);

CREATE TABLE store_schedule(
    store_id INTEGER NOT NULL,
    day_of_week day_of_week NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,

    CONSTRAINT store_schedule_pk PRIMARY KEY(store_id, day_of_week),
    CONSTRAINT store_schedule_store_id_fk FOREIGN KEY(store_id) REFERENCES store(id)
);

-- test data
INSERT INTO currency(code, exp, decimal_sep, min_amount, max_amount, is_enabled) VALUES('BYN', 2, ',', 0, 0, TRUE);
INSERT INTO currency(code, exp, decimal_sep, min_amount, max_amount, is_enabled) VALUES('USD', 2, ',', 0, 0, TRUE);

INSERT INTO store(description, default_currency_id) VALUES('test description', 1);

INSERT INTO store_schedule(store_id, day_of_week, start_time, end_time) VALUES(1, 'Monday', '2025-02-07 08:00:00+03'::timestamptz, '2025-02-07 23:00:00+03'::timestamptz);
INSERT INTO store_schedule(store_id, day_of_week, start_time, end_time) VALUES(1, 'Tuesday', '2025-02-07 08:00:00+03'::timestamptz, '2025-02-07 23:00:00+03'::timestamptz);
INSERT INTO store_schedule(store_id, day_of_week, start_time, end_time) VALUES(1, 'Wednesday', '2025-02-07 08:00:00+03'::timestamptz, '2025-02-07 23:00:00+03'::timestamptz);
INSERT INTO store_schedule(store_id, day_of_week, start_time, end_time) VALUES(1, 'Thursday', '2025-02-07 08:00:00+03'::timestamptz, '2025-02-07 23:00:00+03'::timestamptz);
INSERT INTO store_schedule(store_id, day_of_week, start_time, end_time) VALUES(1, 'Friday', '2025-02-07 08:00:00+03'::timestamptz, '2025-02-07 23:00:00+03'::timestamptz);
INSERT INTO store_schedule(store_id, day_of_week, start_time, end_time) VALUES(1, 'Saturday', '2025-02-07 08:00:00+03'::timestamptz, '2025-02-07 23:00:00+03'::timestamptz);
INSERT INTO store_schedule(store_id, day_of_week, start_time, end_time) VALUES(1, 'Sunday', '2025-02-07 08:00:00+03'::timestamptz, '2025-02-07 23:00:00+03'::timestamptz);

INSERT INTO product(title, is_enabled, created_at, updated_at) VALUES('latte', TRUE, NOW(), NOW());
INSERT INTO product(title, is_enabled, created_at, updated_at) VALUES('americano', TRUE, NOW(), NOW());
INSERT INTO product(title, is_enabled, created_at, updated_at) VALUES('cappuccino', FALSE, NOW(), NOW());
INSERT INTO product(title, is_enabled, created_at, updated_at) VALUES('chips', TRUE, NOW(), NOW());

INSERT INTO product_price(product_id, currency_id, price) VALUES(1, 1, 200);
INSERT INTO product_price(product_id, currency_id, price) VALUES(2, 1, 100);
INSERT INTO product_price(product_id, currency_id, price) VALUES(3, 1, 140);
INSERT INTO product_price(product_id, currency_id, price) VALUES(4, 1, 400);

INSERT INTO product_price(product_id, currency_id, price) VALUES(1, 2, 100);
INSERT INTO product_price(product_id, currency_id, price) VALUES(2, 2, 50);
INSERT INTO product_price(product_id, currency_id, price) VALUES(3, 2, 70);
INSERT INTO product_price(product_id, currency_id, price) VALUES(4, 2, 200);

INSERT INTO category(title, is_enabled) VALUES('coffee', TRUE);
INSERT INTO category(title, is_enabled) VALUES('food', TRUE);
INSERT INTO category(title, is_enabled) VALUES('seafood', FALSE);

INSERT INTO product_category(product_id, category_id) VALUES(1, 1);
INSERT INTO product_category(product_id, category_id) VALUES(2, 1);
INSERT INTO product_category(product_id, category_id) VALUES(3, 3);
INSERT INTO product_category(product_id, category_id) VALUES(4, 2);

-- orders

INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 100, '2025-02-07 23:00:01+03'::timestamptz, '2025-02-07 23:00:01+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 200, '2025-02-07 23:00:02+03'::timestamptz, '2025-02-07 23:00:02+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 300, '2025-02-07 23:00:03+03'::timestamptz, '2025-02-07 23:00:03+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 400, '2025-02-07 23:00:04+03'::timestamptz, '2025-02-07 23:00:04+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 500, '2025-02-07 23:00:05+03'::timestamptz, '2025-02-07 23:00:05+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 600, '2025-02-07 23:00:06+03'::timestamptz, '2025-02-07 23:00:06+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 700, '2025-02-07 23:00:07+03'::timestamptz, '2025-02-07 23:00:07+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 800, '2025-02-07 23:00:08+03'::timestamptz, '2025-02-07 23:00:08+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 900, '2025-02-07 23:00:09+03'::timestamptz, '2025-02-07 23:00:09+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 1000, '2025-02-07 23:00:10+03'::timestamptz, '2025-02-07 23:00:10+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 1100, '2025-02-07 23:00:11+03'::timestamptz, '2025-02-07 23:00:11+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 1200, '2025-02-07 23:00:12+03'::timestamptz, '2025-02-07 23:00:12+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 1300, '2025-02-07 23:00:13+03'::timestamptz, '2025-02-07 23:00:13+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 1400, '2025-02-07 23:00:14+03'::timestamptz, '2025-02-07 23:00:14+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 1500, '2025-02-07 23:00:15+03'::timestamptz, '2025-02-07 23:00:15+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 1600, '2025-02-07 23:00:16+03'::timestamptz, '2025-02-07 23:00:16+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 1700, '2025-02-07 23:00:17+03'::timestamptz, '2025-02-07 23:00:17+03'::timestamptz);
INSERT INTO orders(chat_id, status, verification_code, currency_id, daily_position, total_price, created_at, updated_at) VALUES(707363692, 'completed', '111', 1, 1, 1800, '2025-02-07 23:00:18+03'::timestamptz, '2025-02-07 23:00:18+03'::timestamptz);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE schedule;
DROP TYPE day_of_week;
DROP TABLE store;
