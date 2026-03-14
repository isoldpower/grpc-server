-- +goose Up
-- +goose StatementBegin
-- create products table
CREATE TABLE products (
    id VARCHAR(255) NOT NULL,
    title VARCHAR(255) UNIQUE NOT NULL DEFAULT '',
    description VARCHAR(255) DEFAULT '',
    PRIMARY KEY (id)
);
-- create column in orders that will reference the product
ALTER TABLE orders
    ADD COLUMN productId VARCHAR(255);
-- protect already-created elements by specifying fallback product id
UPDATE orders SET productId = 'some_existing_product_id' WHERE productId IS NULL;
-- constraint the newly created column to the products table
ALTER TABLE orders
    ALTER COLUMN productId SET NOT NULL,
    ADD CONSTRAINT fk_orders_productId FOREIGN KEY (productId) REFERENCES products(id)
    ON DELETE CASCADE ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP CONSTRAINT fk_orders_productId;
ALTER TABLE orders
    DROP COLUMN productId;
DROP TABLE products;
-- +goose StatementEnd
