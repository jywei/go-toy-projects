package model

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type barcodeService interface {
	SyncBarcode([]*Barcode) error
	GetProductIDs(catalogID int, barcodes []string) ([]int, error)
}

type barcodeOps struct {
	readTimeout           time.Duration
	writeTimeout          time.Duration
	transactionMaxTimeout time.Duration
	db                    *sqlx.DB
}

// GetProductIDs returns with two conditions, if barcodes is empty, it will returns product ids by catalog id.
func (b *barcodeOps) GetProductIDs(catalogID int, barcodes []string) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.readTimeout)
	defer cancel()

	query := `SELECT product_id FROM barcodes WHERE catalog_id=` + strconv.Itoa(catalogID) + ` AND is_active=true`
	if barcodes != nil && len(barcodes) > 0 {
		query += " AND ("
		for _, barcode := range barcodes {
			query += "barcode='" + barcode + "' OR "
		}
		query = strings.TrimRight(query, " OR ")
		query += ")"
	}
	ret := make([]int, 0, len(barcodes))
	err := b.db.SelectContext(ctx, &ret, query)
	if err != nil {
		return nil, errors.Wrapf(err, "model: [GetProductIDs] select product ids failed")
	}

	return ret, nil
}

// SyncBarcode must passing all the same product id and catalog id,
// and should passing all barcodes relate on a product id and catalog id at once.
func (b *barcodeOps) SyncBarcode(nbs []*Barcode) error {
	if len(nbs) == 0 {
		return nil
	}

	nbarcodes := make(map[string]*Barcode)
	sameProductID := nbs[0].ProductID
	sameCatalogID := nbs[0].CatalogID
	for _, nb := range nbs {
		if _, exist := nbarcodes[nb.Barcode]; !exist {
			nbarcodes[nb.Barcode] = nb
		}
		if nb.ProductID != sameProductID {
			return errors.Errorf("model: [SyncBarcode] facing not all the same product id:%d, %d", sameProductID, nb.ProductID)
		}
		if nb.CatalogID != sameCatalogID {
			return errors.Errorf("model: [SyncBarcode] facing not all the same catalog id:%d, %d", sameCatalogID, nb.CatalogID)
		}
		sameProductID = nb.ProductID
		sameCatalogID = nb.CatalogID
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.transactionMaxTimeout)
	defer cancel()
	tx, err := b.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "model: [SyncBarcode] transaction begin on productID:%d, catalogID:%d failed", sameProductID, sameCatalogID)
	}

	condition := " WHERE product_id=" + strconv.Itoa(sameProductID) + " AND catalog_id=" + strconv.Itoa(sameCatalogID)
	obs := make([]*Barcode, 0, len(nbs))
	query := `SELECT id,product_id,barcode,catalog_id,is_active FROM barcodes` + condition
	if err := tx.SelectContext(ctx, &obs, query); err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "model: [SyncBarcode] select barcode on productID:%d, catalogID:%d failed", sameProductID, sameCatalogID)
	}

	if _, err := tx.ExecContext(ctx, "UPDATE barcodes SET is_active=false"+condition); err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "model: [SyncBarcode] disable all is active on productID:%d, catalogID:%d failed", sameProductID, sameCatalogID)
	}

	for _, ob := range obs {
		if _, exist := nbarcodes[ob.Barcode]; exist {
			_, err := tx.ExecContext(ctx, "UPDATE barcodes SET is_active=true"+condition+" AND barcode='"+ob.Barcode+"'")
			if err != nil {
				tx.Rollback()
				return errors.Wrapf(err, "model: [SyncBarcode] enable on productID:%d, catalogID:%d, barcode:%s failed", sameProductID, sameCatalogID, ob.Barcode)
			}
			delete(nbarcodes, ob.Barcode)
		}
	}
	for _, nb := range nbarcodes {
		nb.IsActive = true
		_, err := tx.NamedExecContext(ctx, "INSERT INTO barcodes (product_id,barcode,catalog_id,is_active) VALUES (:product_id,:barcode,:catalog_id,:is_active)", nb)
		if err != nil {
			tx.Rollback()
			return errors.Wrapf(err, "model: [SyncBarcode] insert on productID:%d, catalogID:%d, barcode:%s failed", sameProductID, sameCatalogID, nb.Barcode)
		}
	}

	return errors.Wrapf(tx.Commit(), "model: [SyncBarcode] commit failed")
}

// Barcode is the backend middle mapping table of products and barcodes.
type Barcode struct {
	ID        int    `json:"id,omitempty" db:"id"`
	ProductID int    `json:"productId,omitempty" db:"product_id"`
	Barcode   string `json:"barcode,omitempty" db:"barcode"`
	CatalogID int    `json:"catalogId,omitempty" db:"catalog_id"`
	IsActive  bool   `json:"isActive,omitempty" db:"is_active"`
}
