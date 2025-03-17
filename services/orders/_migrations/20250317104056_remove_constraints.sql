-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    DROP CONSTRAINT fk_orders_productId;
ALTER TABLE orders
    DROP CONSTRAINT orders_customerid_fkey;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    ADD CONSTRAINT fk_orders_productId FOREIGN KEY (productid) REFERENCES products(id)
    ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE orders
    ADD CONSTRAINT orders_customerid_fkey FOREIGN KEY (customerid) REFERENCES customers(id)
    ON DELETE CASCADE ON UPDATE CASCADE;
-- +goose StatementEnd
