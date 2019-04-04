-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE catalogs (
    id integer,
    UNIQUE (id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE catalogs;
