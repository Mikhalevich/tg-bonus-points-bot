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

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE schedule;
DROP TYPE day_of_week;
DROP TABLE store;
