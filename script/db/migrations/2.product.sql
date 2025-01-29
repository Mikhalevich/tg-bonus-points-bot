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

    CONSTRAINT product_price_pk PRIMARY KEY(product_id, currency_id)
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

INSERT INTO currency(code, exp, decimal_sep, min_amount, max_amount, is_enabled) VALUES('BYN', 2, ',', 0, 0, TRUE);
INSERT INTO currency(code, exp, decimal_sep, min_amount, max_amount, is_enabled) VALUES('USD', 2, ',', 0, 0, TRUE);

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

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE product_category;
DROP TABLE category;
DROP TABLE product;
