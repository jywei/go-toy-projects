-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE TABLE stores (
    id integer primary key,
    name varchar(128) not null,
    pick_up_point varchar(64) not null,
    slug varchar(128) not null,
    brand_id integer not null references brands(id),
    address_id integer not null,
    catalog_id integer not null references catalogs(id),
    priority varchar(16) not null,
    notes text not null,
    description text not null,
    image_url text not null,
    closed boolean not null,
    temporarily_closed boolean not null,
    opens_at timestamp,
    estimated_delivery_time integer,
    buffer_time integer not null,
    delivery_types varchar(64)[] not null,
    shipping_modes varchar(64)[] not null,
    store_type varchar(128) not null,
    minimum_spend_extra_fee varchar(32),
    minimum_spend varchar(32),
    free_delivery_eligible boolean,
    default_delivery_fee varchar(32),
    minimum_order_free_delivery varchar(32),
    external_key varchar(128) not null
);
CREATE INDEX stores_catalog_id_index ON stores (catalog_id);
CREATE INDEX stores_brand_id_index ON stores (brand_id);
CREATE INDEX stores_external_key_index ON stores (external_key);
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
-- +goose StatementBegin
DROP TABLE stores;
-- +goose StatementEnd
