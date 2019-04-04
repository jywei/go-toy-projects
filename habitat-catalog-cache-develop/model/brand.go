package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type brandService interface {
	GetBrand(id int) (*Brand, error)
	SyncBrand(*Brand) error
}

type brandOps struct {
	readTimeout           time.Duration
	writeTimeout          time.Duration
	transactionMaxTimeout time.Duration
	db                    *sqlx.DB
}

func (b *brandOps) GetBrand(id int) (*Brand, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.readTimeout)
	defer cancel()

	ret := new(Brand)
	err := b.db.GetContext(
		ctx, ret, `SELECT 
			id,
			name,
			slug,
			description,
			brand_color,
			currency,
			country_id,
			minimum_order_free_delivery,
			default_delivery_fee,
			price_markup_percentage,
			free_delivery_eligible,
			minimum_spend,
			minimum_spend_extra_fee,
			default_concierge_fee,
			delivery_types,
			shipping_modes,
			estimated_delivery_time 
		FROM brands WHERE id=$1`, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, errors.Wrapf(err, "model: [GetBrand] get brand failed")
		}
	}

	return ret, nil
}

func (b *brandOps) SyncBrand(brand *Brand) error {
	ctx, cancel := context.WithTimeout(context.Background(), b.writeTimeout)
	defer cancel()

	_, err := b.db.NamedExecContext(ctx, `
		INSERT INTO brands (id,name,slug,description,brand_color,currency,country_id,minimum_order_free_delivery,default_delivery_fee,price_markup_percentage,free_delivery_eligible,minimum_spend,minimum_spend_extra_fee,default_concierge_fee,delivery_types,shipping_modes,estimated_delivery_time)
		VALUES (:id,:name,:slug,:description,:brand_color,:currency,:country_id,:minimum_order_free_delivery,:default_delivery_fee,:price_markup_percentage,:free_delivery_eligible,:minimum_spend,:minimum_spend_extra_fee,:default_concierge_fee,:delivery_types,:shipping_modes,:estimated_delivery_time)
		ON CONFLICT (id)
		DO UPDATE SET
			name=:name,
			slug=:slug,
			description=:description,
			brand_color=:brand_color,
			currency=:currency,
			country_id=:country_id,
			minimum_order_free_delivery=:minimum_order_free_delivery,
			default_delivery_fee=:default_delivery_fee,
			price_markup_percentage=:price_markup_percentage,
			free_delivery_eligible=:free_delivery_eligible,
			minimum_spend=:minimum_spend,
			minimum_spend_extra_fee=:minimum_spend_extra_fee,
			default_concierge_fee=:default_concierge_fee,
			delivery_types=:delivery_types,
			estimated_delivery_time=:estimated_delivery_time,
			shipping_modes=:shipping_modes
			;`,
		brand,
	)
	return errors.Wrapf(err, "model: [SyncBrand] upsert failed")
}

// Brand is the backend /api/brands/xxx returns format and database format.
type Brand struct {
	ID                       int            `json:"id,omitempty" db:"id"`
	Name                     string         `json:"name,omitempty" db:"name"`
	Slug                     string         `json:"slug,omitempty" db:"slug"`
	Description              string         `json:"description,omitempty" db:"description"`
	About                    string         `json:"about,omitempty" db:"-"`
	ImageURL                 string         `json:"imageUrl,omitempty" db:"-"`
	ProductsImageURL         string         `json:"productsImageUrl,omitempty" db:"-"`
	BrandColor               string         `json:"brandColor,omitempty" db:"brand_color"`
	Currency                 string         `json:"currency,omitempty" db:"currency"`
	CountryID                int            `json:"countryId,omitempty" db:"country_id"`
	MinimumOrderFreeDelivery string         `json:"minimumOrderFreeDelivery,omitempty" db:"minimum_order_free_delivery"`
	DefaultDeliveryFee       string         `json:"defaultDeliveryFee,omitempty" db:"default_delivery_fee"`
	SameStorePrice           bool           `json:"sameStorePrice,omitempty" db:"-"`
	BrandType                string         `json:"brandType,omitempty" db:"-"`
	PromotionText            string         `json:"promotionText,omitempty" db:"-"`
	ParentBrandID            int            `json:"parentBrandId,omitempty" db:"-"`
	ProductsCount            int            `json:"productsCount,omitempty" db:"-"`
	StoreID                  int            `json:"storeId,omitempty" db:"-"`
	PriceMarkupPercentage    string         `json:"priceMarkupPercentage,omitempty" db:"price_markup_percentage"`
	FreeDeliveryEligible     bool           `json:"freeDeliveryEligible,omitempty" db:"free_delivery_eligible"`
	EstimatedDeliveryTime    int            `json:"estimatedDeliveryTime,omitempty" db:"estimated_delivery_time"`
	Closed                   bool           `json:"closed,omitempty" db:"-"`
	CatalogID                int            `json:"catalogId,omitempty" db:"-"`
	OpensAt                  time.Time      `json:"opensAt,omitempty" db:"-"`
	MinimumSpend             string         `json:"minimumSpend,omitempty" db:"minimum_spend"`
	MinimumSpendExtraFee     string         `json:"minimumSpendExtraFee,omitempty" db:"minimum_spend_extra_fee"`
	ShippingModes            pq.StringArray `json:"shippingModes,omitempty" db:"shipping_modes"`
	DeliveryTypes            pq.StringArray `json:"deliveryTypes,omitempty" db:"delivery_types"`
	DefaultConciergeFee      string         `json:"defaultConciergeFee,omitempty" db:"default_concierge_fee"`
}

// Validate currects the brand data.
func (b *Brand) Validate() {
	if b.MinimumOrderFreeDelivery == "" {
		b.MinimumOrderFreeDelivery = "0.0"
	}
	if b.DefaultDeliveryFee == "" {
		b.DefaultDeliveryFee = "0.0"
	}
	if b.MinimumSpend == "" {
		b.MinimumSpend = "0.0"
	}
	if b.MinimumSpendExtraFee == "" {
		b.MinimumSpendExtraFee = "0.0"
	}
	if b.DefaultConciergeFee == "" {
		b.DefaultConciergeFee = "0.0"
	}
}

// ValidateWithStore replaces brand data if store has corresponding data.
func (b *Brand) ValidateWithStore(store *Store) {
	b.StoreID = store.ID
	b.CatalogID = store.CatalogID

	if store.MinimumSpendExtraFee != "" && store.MinimumSpendExtraFee != "0.0" {
		b.MinimumSpendExtraFee = store.MinimumSpendExtraFee
	}
	if store.MinimumSpend != "" && store.MinimumSpend != "0.0" {
		b.MinimumSpend = store.MinimumSpend
	}
	if store.DefaultDeliveryFee != "" && store.DefaultDeliveryFee != "0.0" {
		b.DefaultDeliveryFee = store.DefaultDeliveryFee
	}
	if store.MinimumOrderFreeDelivery != "" && store.MinimumOrderFreeDelivery != "0.0" {
		b.MinimumOrderFreeDelivery = store.MinimumOrderFreeDelivery
	}
	if len(store.ShippingModes) != 0 {
		b.ShippingModes = store.ShippingModes
	}
	if len(store.DeliveryTypes) != 0 {
		b.DeliveryTypes = store.DeliveryTypes
	}
}
