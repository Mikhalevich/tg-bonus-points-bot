-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE order_status AS ENUM (
    'waiting_payment',
    'payment_in_progress',
    'confirmed',
    'in_progress',
    'ready',
    'completed',
    'canceled',
    'rejected'
);

CREATE TABLE currency(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    code TEXT NOT NULL,
    exp INTEGER NOT NULL,
    decimal_sep TEXT NOT NULL,
    min_amount INTEGER NOT NULL,
    max_amount INTEGER NOT NULL,
    is_enabled BOOLEAN NOT NULL
);

CREATE UNIQUE INDEX currency_code_unique_idx ON currency(code);

CREATE TABLE orders(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id BIGINT NOT NULL,
    status order_status NOT NULL,
    verification_code TEXT,
    currency_id INTEGER NOT NULL,
    daily_position INTEGER,

    CONSTRAINT orders_currency_id_fk FOREIGN KEY(currency_id) REFERENCES currency(id)
);

CREATE INDEX orders_chat_id_status_idx ON orders(chat_id, status);
CREATE UNIQUE INDEX orders_only_one_active_order_unique_idx ON orders(chat_id) WHERE status IN ('waiting_payment', 'payment_in_progress', 'confirmed', 'in_progress', 'ready');

CREATE TABLE order_status_timeline(
    order_id INTEGER NOT NULL,
    status order_status NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT order_status_timeline_pk PRIMARY KEY(order_id, status),
    CONSTRAINT order_status_timeline_order_id_fk FOREIGN KEY(order_id) REFERENCES orders(id)
);

CREATE INDEX order_status_timeline_order_id_idx ON order_status_timeline(order_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE order_status_timeline;
DROP TABLE orders;
DROP TABLE currency;
DROP TYPE order_status;
