-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE product(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title TEXT NOT NULL,
    is_enabled BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE product_price(
    product_id INTEGER NOT NULL,
    currency_id INTEGER NOT NULL,
    price INTEGER NOT NULL,

    CONSTRAINT product_price_pk PRIMARY KEY(product_id, currency_id),
    CONSTRAINT product_price_product_id FOREIGN KEY(product_id) REFERENCES product(id),
);

CREATE TABLE category(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title TEXT NOT NULL,
    is_enabled BOOLEAN NOT NULL
);

CREATE TABLE product_category(
    product_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    
    CONSTRAINT category_pk PRIMARY KEY(product_id, category_id),
    CONSTRAINT product_id_fk FOREIGN KEY(product_id) REFERENCES product(id),
    CONSTRAINT category_id_fk FOREIGN KEY(category_id) REFERENCES category(id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE product_category;
DROP TABLE category;
DROP TABLE product;
