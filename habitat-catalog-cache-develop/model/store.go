package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type storeService interface {
	GetStoreByBrandID(int) (*Store, error)
	SyncStore(*Store) error
	GetStoreByID(id int) (*Store, error)
	GetExternalKeys() ([]string, error)
	SelectStores(externalKey, storeType string) ([]*Store, error)
}

type storeOps struct {
	readTimeout           time.Duration
	writeTimeout          time.Duration
	transactionMaxTimeout time.Duration
	db                    *sqlx.DB
}

func (s *storeOps) GetExternalKeys() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.readTimeout)
	defer cancel()

	ret := make([]string, 0)
	err := s.db.SelectContext(ctx, &ret, "SELECT external_key FROM stores")
	if err != nil {
		return nil, errors.Wrapf(err, "model: [GetExternalKeys] select context failed")
	}

	return ret, nil
}

const (
	selectAllFromStoresSQL = `SELECT 
		id,
		name,
		pick_up_point,
		slug,
		brand_id,
		address_id,
		catalog_id,
		priority,
		notes,
		description,
		image_url,
		closed,
		temporarily_closed,
		opens_at,
		estimated_delivery_time,
		buffer_time,
		shipping_modes,
		delivery_types,
		store_type,
		minimum_order_free_delivery,
		default_delivery_fee,
		free_delivery_eligible,
		minimum_spend,
		minimum_spend_extra_fee,
		external_key
	FROM stores`
)

func (s *storeOps) SelectStores(externalKey, storeType string) ([]*Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.readTimeout)
	defer cancel()

	query := selectAllFromStoresSQL
	if externalKey != "" && storeType != "" {
		query += " WHERE external_key='" + externalKey + "'" + " AND store_type='" + storeType + "'"
	} else if storeType != "" {
		query += " WHERE store_type='" + storeType + "'"
	} else if externalKey != "" {
		query += " WHERE external_key='" + externalKey + "'"
	}

	ret := make([]*Store, 0)
	if err := s.db.SelectContext(ctx, &ret, query); err != nil {
		return nil, errors.Wrapf(err, "model: [SelectStores] select stores failed")
	}

	return ret, nil
}

func (s *storeOps) GetStoreByID(id int) (*Store, error) {
	return s.getStore(selectAllFromStoresSQL+" WHERE id=$1 LIMIT 1", id)
}

// GetStoreByBrandID returns the first store of brand id.
// followed current backend logic for habitat.
func (s *storeOps) GetStoreByBrandID(id int) (*Store, error) {
	return s.getStore(selectAllFromStoresSQL+" WHERE brand_id=$1 LIMIT 1", id)
}

func (s *storeOps) getStore(query string, args ...interface{}) (*Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.readTimeout)
	defer cancel()

	ret := new(Store)
	err := s.db.GetContext(ctx, ret, query, args...)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNoRows
		default:
			return nil, errors.Wrapf(err, "model: [getStore] get store failed")
		}
	}

	return ret, nil
}

func (s *storeOps) SyncStore(store *Store) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.writeTimeout)
	defer cancel()

	_, err := s.db.NamedExecContext(ctx, `
		INSERT INTO stores (id,name,pick_up_point,slug,brand_id,address_id,catalog_id,priority,notes,description,image_url,closed,temporarily_closed,opens_at,estimated_delivery_time,buffer_time,shipping_modes,delivery_types,store_type,minimum_order_free_delivery,default_delivery_fee,free_delivery_eligible,minimum_spend,minimum_spend_extra_fee,external_key)
		VALUES (:id,:name,:pick_up_point,:slug,:brand_id,:address_id,:catalog_id,:priority,:notes,:description,:image_url,:closed,:temporarily_closed,:opens_at,:estimated_delivery_time,:buffer_time,:shipping_modes,:delivery_types,:store_type,:minimum_order_free_delivery,:default_delivery_fee,:free_delivery_eligible,:minimum_spend,:minimum_spend_extra_fee,:external_key)
		ON CONFLICT (id)
		DO UPDATE SET
			name=:name,
			pick_up_point=:pick_up_point,
			slug=:slug,
			brand_id=:brand_id,
			address_id=:address_id,
			catalog_id=:catalog_id,
			priority=:priority,
			notes=:notes,
			description=:description,
			image_url=:image_url,
			closed=:closed,
			temporarily_closed=:temporarily_closed,
			opens_at=:opens_at,
			estimated_delivery_time=:estimated_delivery_time,
			buffer_time=:buffer_time,
			shipping_modes=:shipping_modes,
			delivery_types=:delivery_types,
			store_type=:store_type,
			minimum_order_free_delivery=:minimum_order_free_delivery,
			default_delivery_fee=:default_delivery_fee,
			free_delivery_eligible=:free_delivery_eligible,
			minimum_spend=:minimum_spend,
			minimum_spend_extra_fee=:minimum_spend_extra_fee,
			external_key=:external_key
			;`,
		store,
	)
	return errors.Wrapf(err, "model: [SyncStore] upsert failed")
}

// Store is the backend /api/stores?externalKey=xxx returns format and database format.
type Store struct {
	ID                       int            `json:"id,omitempty" db:"id"`
	Name                     string         `json:"name,omitempty" db:"name"`
	PickupPoint              string         `json:"pickupPoint,omitempty" db:"pick_up_point"`
	Slug                     string         `json:"slug,omitempty" db:"slug"`
	BrandID                  int            `json:"brandId,omitempty" db:"brand_id"`
	AddressID                int            `json:"addressId,omitempty" db:"address_id"`
	CatalogID                int            `json:"catalogId,omitempty" db:"catalog_id"`
	Priority                 string         `json:"priority,omitempty" db:"priority"`
	Notes                    string         `json:"notes,omitempty" db:"notes"`
	Description              string         `json:"description,omitempty" db:"description"`
	ImageURL                 string         `json:"imageUrl,omitempty" db:"image_url"`
	Closed                   bool           `json:"closed,omitempty" db:"closed"`
	TemporarilyClosed        bool           `json:"temporarilyClosed,omitempty" db:"temporarily_closed"`
	OpensAt                  time.Time      `json:"opensAt,omitempty" db:"opens_at"`
	EstimatedDeliveryTime    int            `json:"estimatedDeliveryTime,omitempty" db:"estimated_delivery_time"`
	BufferTime               int            `json:"bufferTime,omitempty" db:"buffer_time"`
	StoreType                string         `json:"storeType,omitempty" db:"store_type"`
	DeliveryTypes            pq.StringArray `json:"deliveryTypes,omitempty" db:"delivery_types"`
	ShippingModes            pq.StringArray `json:"shippingModes,omitempty" db:"shipping_modes"`
	MinimumOrderFreeDelivery string         `json:"minimumOrderFreeDelivery,omitempty" db:"minimum_order_free_delivery"`
	DefaultDeliveryFee       string         `json:"defaultDeliveryFee,omitempty" db:"default_delivery_fee"`
	FreeDeliveryEligible     bool           `json:"freeDeliveryEligible,omitempty" db:"free_delivery_eligible"`
	MinimumSpend             string         `json:"minimumSpend,omitempty" db:"minimum_spend"`
	MinimumSpendExtraFee     string         `json:"minimumSpendExtraFee,omitempty" db:"minimum_spend_extra_fee"`
	ExternalKey              string         `json:"-" db:"external_key"`
	Brand                    *Brand         `json:"brand,omitempty"`
}
