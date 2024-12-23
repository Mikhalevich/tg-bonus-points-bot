-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE order_status AS ENUM (
    'created',
    'in_progress',
    'ready',
    'completed',
    'canceled'
);

CREATE TABLE orders(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id BIGINT NOT NULL,
    status order_status NOT NULL DEFAULT 'created',
    verification_code TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    in_progress_at TIMESTAMPTZ,
    ready_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    canceled_at TIMESTAMPTZ
);

CREATE INDEX chat_id_created_at_idx ON orders(chat_id, created_at);
CREATE INDEX chat_id_status_idx ON orders(chat_id, status);
CREATE UNIQUE INDEX only_one_active_order_unique_idx ON orders(chat_id) WHERE status IN ('created', 'in_progress', 'ready');

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE order;
DROP TYPE order_status;
