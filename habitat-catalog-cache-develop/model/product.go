package model

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type productService interface {
	SyncProduct(*Product) error
	GetProducts(ids []int) ([]*Product, error)
}

type productOps struct {
	readTimeout           time.Duration
	writeTimeout          time.Duration
	transactionMaxTimeout time.Duration
	db                    *sqlx.DB
}

func (p *productOps) GetProducts(ids []int) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.readTimeout)
	defer cancel()

	ret := make([]*Product, 0)
	query := `SELECT 
		id,
		title,
		description,
		image_url,
		preview_image_url,
		slug,
		barcode,
		barcodes,
		unit_type,
		sold_by,
		amount_per_unit,
		size,
		status,
		image_url_basename,
		currency,
		max_quantity,
		customer_notes_enabled,
		price,
		normal_price,
		external_id,
		catalog_id,
		brand_id,
		brand_slug,
		is_alcohol 
	FROM products`
	if ids != nil && len(ids) > 0 {
		query += " WHERE "
		for _, id := range ids {
			query += "id=" + strconv.Itoa(id) + " OR "
		}
		query = strings.TrimRight(query, " OR ")
	}
	err := p.db.SelectContext(ctx, &ret, query)
	if err != nil {
		return nil, errors.Wrapf(err, "model: [GetProducts] select product failed")
	}

	return ret, nil
}

func (p *productOps) SyncProduct(product *Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.writeTimeout)
	defer cancel()

	_, err := p.db.NamedExecContext(ctx, `
		INSERT INTO products (id,title,description,image_url,preview_image_url,slug,barcode,barcodes,unit_type,sold_by,amount_per_unit,size,status,image_url_basename,currency,max_quantity,customer_notes_enabled,price,normal_price,external_id,catalog_id,brand_id,brand_slug,is_alcohol)
		VALUES (:id,:title,:description,:image_url,:preview_image_url,:slug,:barcode,:barcodes,:unit_type,:sold_by,:amount_per_unit,:size,:status,:image_url_basename,:currency,:max_quantity,:customer_notes_enabled,:price,:normal_price,:external_id,:catalog_id,:brand_id,:brand_slug,:is_alcohol)
		ON CONFLICT (id)
		DO UPDATE SET
			title=:title,
			description=:description,
			image_url=:image_url,
			preview_image_url=:preview_image_url,
			slug=:slug,
			barcode=:barcode,
			barcodes=:barcodes,
			unit_type=:unit_type,
			sold_by=:sold_by,
			amount_per_unit=:amount_per_unit,
			size=:size,
			status=:status,
			image_url_basename=:image_url_basename,
			currency=:currency,
			max_quantity=:max_quantity,
			customer_notes_enabled=:customer_notes_enabled,
			price=:price,
			normal_price=:normal_price,
			external_id=:external_id,
			catalog_id=:catalog_id,
			brand_id=:brand_id,
			brand_slug=:brand_slug,
			is_alcohol=:is_alcohol
			;`,
		product,
	)
	return errors.Wrapf(err, "model: [SyncProduct] upsert failed")
}

// Product is the backend /api/products?catalog_id=xxx returns format and database format.
type Product struct {
	ID                   int            `json:"id,omitempty" db:"id"`
	Title                string         `json:"title,omitempty" db:"title"`
	Description          string         `json:"description,omitempty" db:"description"`
	ImageURL             string         `json:"imageUrl,omitempty" db:"image_url"`
	PreviewImageURL      string         `json:"previewImageUrl,omitempty" db:"preview_image_url"`
	Slug                 string         `json:"slug,omitempty" db:"slug"`
	Barcodes             pq.StringArray `json:"barcodes,omitempty" db:"barcodes"`
	Barcode              string         `json:"-" db:"barcode"`
	UnitType             string         `json:"unitType,omitempty" db:"unit_type"`
	SoldBy               string         `json:"soldBy,omitempty" db:"sold_by"`
	AmountPerUnit        string         `json:"amountPerUnit,omitempty" db:"amount_per_unit"`
	Size                 string         `json:"size,omitempty" db:"size"`
	Status               string         `json:"status,omitempty" db:"status"`
	ImageURLBasename     string         `json:"imageUrlBasename,omitempty" db:"image_url_basename"`
	Currency             string         `json:"currency,omitempty" db:"currency"`
	MaxQuantity          string         `json:"maxQuantity,omitempty" db:"max_quantity" default:"0.0"`
	CustomerNotesEnabled bool           `json:"customerNotesEnabled,omitempty" db:"customer_notes_enabled"`
	Price                string         `json:"price,omitempty" db:"price" default:"0.0"`
	NormalPrice          string         `json:"normalPrice,omitempty" db:"normal_price" default:"0.0"`
	BrandSlug            string         `json:"brandSlug,omitempty" db:"brand_slug"`
	ExternalID           string         `json:"externalId,omitempty" db:"external_id"`
	CatalogID            int            `json:"-" db:"catalog_id"`
	BrandID              int            `json:"brandId,omitempty" db:"brand_id"`
	Alcohol              bool           `json:"alcohol,omitempty" db:"is_alcohol"`
}
