-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ADD CONSTRAINT orders_customerId_fkey FOREIGN KEY (customerId) REFERENCES customers(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    ADD CONSTRAINT orders_productId_fkey FOREIGN KEY (productId) REFERENCES products(id)
        ON DELETE CASCADE ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP CONSTRAINT orders_customerId_fkey,
    DROP CONSTRAINT orders_productId_fkey;
-- +goose StatementEnd
