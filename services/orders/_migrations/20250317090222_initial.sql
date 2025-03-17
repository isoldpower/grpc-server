-- +goose Up
-- +goose StatementBegin
CREATE TABLE customers (
   id VARCHAR(255) NOT NULL UNIQUE,
   name VARCHAR(255) DEFAULT '',
   PRIMARY KEY (id)
);

CREATE TABLE orders (
    id VARCHAR(255) NOT NULL UNIQUE,
    customerId VARCHAR(255) NOT NULL,
    quantity int,
    FOREIGN KEY (customerId) REFERENCES customers(id),
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE customers CASCADE;
DROP TABLE orders CASCADE;
-- +goose StatementEnd
