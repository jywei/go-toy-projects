package model

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type catalogService interface {
	SyncCatalog(*Catalog) error
}

type catalogOps struct {
	readTimeout           time.Duration
	writeTimeout          time.Duration
	transactionMaxTimeout time.Duration
	db                    *sqlx.DB
}

func (c *catalogOps) SyncCatalog(catalog *Catalog) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.writeTimeout)
	defer cancel()

	_, err := c.db.NamedExecContext(ctx, `
		INSERT INTO catalogs (id)
		VALUES (:id)
		ON CONFLICT (id)
		DO NOTHING;`,
		catalog,
	)
	return errors.Wrapf(err, "model: [SyncCatalog] INSERT failed")
}

// Catalog is the backend middleware of store and product.
type Catalog struct {
	ID int `json:"-" db:"id"`
}
