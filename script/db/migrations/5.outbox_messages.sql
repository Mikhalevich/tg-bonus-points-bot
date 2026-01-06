-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE message_type AS ENUM (
    'plain',
    'markdown',
    'png'
);

CREATE TABLE outbox_messages(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id BIGINT NOT NULL,
    reply_msg_id BIGINT,
    msg_text TEXT NOT NULL,
    msg_type message_type NOT NULL,
    payload BYTEA,
    buttons JSONB NOT NULL
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE outbox_messages;
DROP TYPE message_type;
