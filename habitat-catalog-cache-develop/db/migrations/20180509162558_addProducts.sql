-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE TABLE products (
    id integer primary key,
    title text not null,
    description text not null,
    image_url text not null,
    preview_image_url text not null,
    slug varchar(128) not null,
    barcode varchar(32) not null,
    barcodes varchar(32)[] not null,
    unit_type varchar(64) not null,
    sold_by varchar(64) not null,
    amount_per_unit varchar(32) not null,
    size varchar(64) not null,
    status varchar(64) not null,
    image_url_basename text not null,
    currency varchar(8) not null,
    max_quantity varchar(32) not null,
    customer_notes_enabled boolean not null,
    price varchar(32) not null,
    normal_price varchar(32) not null,
    catalog_id integer not null references catalogs(id),
    external_id varchar(32) not null,
    brand_id integer not null references brands(id),
    brand_slug varchar(64) not null,
    is_alcohol boolean not null
);
CREATE INDEX products_catalog_id_index ON products (catalog_id);
CREATE INDEX products_brand_id_index ON products (brand_id);
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
