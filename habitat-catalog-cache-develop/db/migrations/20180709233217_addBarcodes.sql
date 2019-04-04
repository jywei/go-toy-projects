-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE TABLE barcodes (
    id serial primary key,
    product_id integer not null,
    barcode varchar(32) not null,
    catalog_id integer not null,
    is_active boolean not null
);
CREATE INDEX barcodes_barcode_index ON barcodes (barcode);
CREATE INDEX barcodes_catalog_id_index ON barcodes (catalog_id);
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
-- +goose StatementBegin
DROP TABLE barcodes;
-- +goose StatementEnd
