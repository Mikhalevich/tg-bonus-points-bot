-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE order_products(
    order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    count INTEGER NOT NULL,
    price INTEGER NOT NULL,

    CONSTRAINT order_products_pk PRIMARY KEY(order_id, product_id),
    CONSTRAINT order_products_order_id_fk FOREIGN KEY(order_id) REFERENCES orders(id),
    CONSTRAINT order_products_product_id_fk FOREIGN KEY(product_id) REFERENCES product(id)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE order_products;
