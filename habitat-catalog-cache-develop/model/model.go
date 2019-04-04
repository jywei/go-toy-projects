package model

import (
	"context"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // blank import for sqlx.Connect usage
	"github.com/pkg/errors"
)

var (
	// ErrNoRows represent the query got sql.ErrNoRows error
	ErrNoRows = errors.New("no such rows")
)

// Service is the model interface for defining all database operations.
type Service interface {
	catalogService
	brandService
	storeService
	productService
	jobService
	barcodeService
	syncStatusService
	Close() error
}

type service struct {
	*catalogOps
	*brandOps
	*storeOps
	*productOps
	*jobOps
	*barcodeOps
	*syncStatusOps
	close func() error
}

// New returns a Service instance.
func New(dconf *config.Database, cconf *config.Cache) (Service, error) {
	db, err := newDB(dconf)
	if err != nil {
		return nil, errors.Wrapf(err, "model: [New] new database failed")
	}

	readTimeout := time.Duration(dconf.ReadTimeoutSec) * time.Second
	writeTimeout := time.Duration(dconf.WriteTimeoutSec) * time.Second
	txTimeout := time.Duration(dconf.TransactionMaxTimeoutSec) * time.Second

	cache, err := newCache(cconf)
	if err != nil {
		return nil, errors.Wrapf(err, "model: [New] new cache failed")
	}

	getConnTimeout := time.Duration(cconf.GetConnectionTimeoutSec) * time.Second

	return &service{
		catalogOps:    &catalogOps{readTimeout: readTimeout, writeTimeout: writeTimeout, transactionMaxTimeout: txTimeout, db: db},
		brandOps:      &brandOps{readTimeout: readTimeout, writeTimeout: writeTimeout, transactionMaxTimeout: txTimeout, db: db},
		storeOps:      &storeOps{readTimeout: readTimeout, writeTimeout: writeTimeout, transactionMaxTimeout: txTimeout, db: db},
		productOps:    &productOps{readTimeout: readTimeout, writeTimeout: writeTimeout, transactionMaxTimeout: txTimeout, db: db},
		jobOps:        &jobOps{getConnTimeout: getConnTimeout, pool: cache},
		barcodeOps:    &barcodeOps{readTimeout: readTimeout, writeTimeout: writeTimeout, transactionMaxTimeout: txTimeout, db: db},
		syncStatusOps: &syncStatusOps{getConnTimeout: getConnTimeout, pool: cache},
		close: func() error {
			derr := errors.Wrapf(db.Close(), "model: [Close] db close failed")
			cerr := errors.Wrapf(cache.Close(), "model: [Close] cache close failed")
			retErr := errors.Errorf("%v, %v", derr, cerr)
			if derr == nil {
				retErr = cerr
			} else if cerr == nil {
				retErr = derr
			}
			return retErr
		},
	}, nil
}

func newDB(conf *config.Database) (*sqlx.DB, error) {
	connSchema := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		conf.User,
		conf.Name,
		conf.Pwd,
		conf.Host,
		conf.Port,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(conf.ConnectTimeoutSec)*time.Second,
	)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "postgres", connSchema)
	if err != nil {
		return nil, errors.Wrapf(err, "model: [newDB] connect db failed")
	}

	db.SetMaxIdleConns(conf.MaxIdle)
	db.SetMaxOpenConns(conf.MaxActive)

	return db, nil
}

func newCache(conf *config.Cache) (*redis.Pool, error) {
	connectTimeout := time.Duration(conf.ConnectTimeoutSec) * time.Second
	readTimeout := time.Duration(conf.ReadTimeoutSec) * time.Second
	writeTimeout := time.Duration(conf.WriteTimeoutSec) * time.Second
	idleTimeout := time.Duration(conf.IdleTimeoutSec) * time.Second

	return &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: idleTimeout,
		Wait:        conf.Wait,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				conf.Host+":"+conf.Port,
				redis.DialConnectTimeout(connectTimeout),
				redis.DialReadTimeout(readTimeout),
				redis.DialWriteTimeout(writeTimeout),
				redis.DialPassword(conf.Pwd),
				redis.DialDatabase(conf.Index),
			)
			if err != nil {
				return nil, errors.Wrapf(
					err,
					"model: [redis Dial] dial failed",
				)
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return errors.Wrapf(
				err,
				"model: [redis TestOnBorrow] PING failed",
			)
		},
	}, nil
}

func (s *service) Close() error {
	return s.close()
}
