package model

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

const (
	syncStatusKey = "catalog_cache_service_sync_status"
	// SyncStatusSuccess means the sync up job success
	SyncStatusSuccess = "success"
	// SyncStatusFailed means the sync up job failed
	SyncStatusFailed = "failed"
)

type syncStatusService interface {
	UpdateSyncStatus(*SyncStatus) error
	GetSyncStatus() (*SyncStatus, error)
}

type syncStatusOps struct {
	getConnTimeout time.Duration
	pool           *redis.Pool
}

// SyncStatus is the sync-up status.
type SyncStatus struct {
	LastSyncTime   string `json:"lastSyncTime"`
	LastSyncStatus string `json:"lastSyncStatus"`
}

// Duplication is the report of duplication.
type Duplication struct {
	Barcode    string `json:"barcode"`
	StoreID    int    `json:"storeId"`
	ProductIDs []int  `json:"productIds"`
}

func (s *syncStatusOps) do(cmd string, args ...interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.getConnTimeout)
	defer cancel()

	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "model: [do] get redis connection failed")
	}
	defer conn.Close()

	reply, err := conn.Do(cmd, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "model: [do] failed on cmd:%s, args:%v", cmd, args)
	}

	return reply, nil
}

func (s *syncStatusOps) UpdateSyncStatus(status *SyncStatus) error {
	ds, err := json.Marshal(status)
	if err != nil {
		return errors.Wrapf(err, "model: [PushUpdateSyncStatusJob] json marshal failed")
	}

	_, err = s.do("SET", syncStatusKey, base64.StdEncoding.EncodeToString(ds))
	return errors.Wrapf(err, "model: [PushUpdateSyncStatusJob] SET failed")
}

func (s *syncStatusOps) GetSyncStatus() (*SyncStatus, error) {
	status, err := redis.String(s.do("GET", syncStatusKey))
	if err != nil {
		return nil, errors.Wrapf(err, "model: [GetSyncStatus] GET failed")
	}
	ds, err := base64.StdEncoding.DecodeString(status)
	if err != nil {
		return nil, errors.Wrapf(err, "model: [GetSyncStatus] base64 decode failed")
	}

	ret := new(SyncStatus)
	if err = json.Unmarshal(ds, ret); err != nil {
		return nil, errors.Wrapf(err, "model: [GetSyncStatus] json unmarshal failed")
	}

	return ret, nil
}
