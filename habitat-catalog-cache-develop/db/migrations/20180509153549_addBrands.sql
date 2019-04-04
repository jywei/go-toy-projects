-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- +goose StatementBegin
CREATE TABLE brands (
    id integer primary key,
    name varchar(128) not null,
    slug varchar(128) not null,
    description text not null,
    brand_color varchar(64) not null,
    currency varchar(8) not null,
    country_id integer not null,
    price_markup_percentage varchar(16) not null,
    minimum_order_free_delivery varchar(32),
    default_delivery_fee varchar(32),
    free_delivery_eligible boolean,
    estimated_delivery_time integer,
    minimum_spend varchar(32),
    minimum_spend_extra_fee varchar(32),
    default_concierge_fee varchar(32),
    delivery_types varchar(64)[] not null,
    shipping_modes varchar(64)[] not null
);
-- +goose StatementEnd

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
-- +goose StatementBegin
DROP TABLE brands;
-- +goose StatementEnd
