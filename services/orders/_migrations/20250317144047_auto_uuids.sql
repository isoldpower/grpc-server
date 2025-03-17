-- +goose Up
-- +goose StatementBegin

-- migrate customers table to UUIDs and convert currently-set to UUIDs
ALTER TABLE customers
    ALTER COLUMN id TYPE UUID USING (
        CASE WHEN id ~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$'
            THEN id::UUID
            ELSE gen_random_uuid()
    END),
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- migrate products table to UUIDs and convert currently-set to UUIDs
ALTER TABLE products
    ALTER COLUMN id TYPE UUID USING (
        CASE WHEN id ~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$'
            THEN id::UUID
            ELSE gen_random_uuid()
    END),
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- migrate orders table to UUIDs and convert currently-set to UUIDs
ALTER TABLE orders
    ALTER COLUMN id TYPE UUID USING (
        CASE WHEN id ~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$'
            THEN id::UUID
            ELSE gen_random_uuid()
    END),
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- update customerId column
UPDATE orders SET customerId = (SELECT id FROM customers LIMIT 1)
    WHERE customerId !~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$';

ALTER TABLE orders
    ALTER COLUMN customerId TYPE UUID USING (
        CASE WHEN customerId ~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$'
            THEN customerId::UUID
    END);

-- update productId column
UPDATE orders SET productId = (SELECT id FROM products LIMIT 1)
    WHERE productId !~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$';

ALTER TABLE orders
    ALTER COLUMN productId TYPE UUID USING (
        CASE WHEN productId ~ '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$'
            THEN productId::UUID
    END);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- drop foreign key constraints
-- ALTER TABLE orders
--     DROP CONSTRAINT orders_customerId_fkey,
--     DROP CONSTRAINT orders_productId_fkey;

-- convert orders back to VARCHAR(255)
ALTER TABLE orders
    ALTER COLUMN id TYPE VARCHAR(255) USING (id::TEXT),
    ALTER COLUMN id DROP DEFAULT,
    ALTER COLUMN customerId TYPE VARCHAR(255) USING (customerId::TEXT),
    ALTER COLUMN productId TYPE VARCHAR(255) USING (productId::TEXT);

-- convert products back to VARCHAR(255)
ALTER TABLE products
    ALTER COLUMN id TYPE VARCHAR(255) USING (id::TEXT),
    ALTER COLUMN id DROP DEFAULT;

-- convert customers back to VARCHAR(255)
ALTER TABLE customers
    ALTER COLUMN id TYPE VARCHAR(255) USING (id::TEXT);
ALTER TABLE customers
    ALTER COLUMN id DROP DEFAULT;

-- re-add foreign key constraints
-- ALTER TABLE orders
--     ADD CONSTRAINT orders_customerId_fkey FOREIGN KEY (customerId) REFERENCES customers(id)
--         ON DELETE CASCADE ON UPDATE CASCADE,
--     ADD CONSTRAINT orders_productId_fkey FOREIGN KEY (productId) REFERENCES products(id)
--         ON DELETE CASCADE ON UPDATE CASCADE;

-- +goose StatementEnd
