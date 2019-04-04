package procesor_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/honestbee/habitat-catalog-cache/procesor"
)

const (
	causeFetchStoreFailedExternalKey   = "cause-fetch-store-failed-external-key"
	causeFetchBrandFailedExternalKey   = "cause-fetch-brand-failed-external-key"
	causeSyncBrandFailedExternalKey    = "cause-sync-brand-failed-external-key"
	causeSyncStoreFailedExternalKey    = "cause-sync-store-failed-external-key"
	causeFetchProductFailedExternalKey = "cause-fetch-product-failed-external-key"
	causeSyncProductFailedExternalKey  = "cause-sync-product-failed-external-key"
	causeSyncCatalogFailedExternalKey  = "cause-sync-catalog-failed-external-key"
	causeFetchBrandFailedBrandID       = 449449
	causeSyncBrandFailedBrandID        = 450450
	casueSyncStoreFailedStoreID        = 3345678
	causeFetchProductFailedStoreID     = 3345679
	causeSyncProductFailedStoreID      = 3345677
	causeSyncProductFailedProductID    = 5566
	causeSyncCatalogFailedCatalogID    = 5577
)

type mockSeeker struct{}

func (*mockSeeker) FetchStores(exKey string) ([]*model.Store, error) {
	switch exKey {
	case causeSyncCatalogFailedExternalKey:
		return []*model.Store{
			&model.Store{
				ID:        1,
				CatalogID: causeSyncCatalogFailedCatalogID,
			},
		}, nil
	case causeFetchStoreFailedExternalKey:
		return nil, errors.New("error occured")
	case causeFetchBrandFailedExternalKey:
		return []*model.Store{
			&model.Store{
				ID:      1,
				BrandID: causeFetchBrandFailedBrandID,
			},
		}, nil
	case causeSyncBrandFailedExternalKey:
		return []*model.Store{
			&model.Store{
				ID:      1,
				BrandID: causeSyncBrandFailedBrandID,
			},
		}, nil
	case causeSyncStoreFailedExternalKey:
		return []*model.Store{
			&model.Store{
				ID:      casueSyncStoreFailedStoreID,
				BrandID: 1,
			},
		}, nil
	case causeFetchProductFailedExternalKey:
		return []*model.Store{
			&model.Store{
				ID:        causeFetchProductFailedStoreID,
				BrandID:   1,
				CatalogID: 1,
			},
		}, nil
	case causeSyncProductFailedExternalKey:
		return []*model.Store{
			&model.Store{
				ID:        causeSyncProductFailedStoreID,
				BrandID:   1,
				CatalogID: 1,
			},
		}, nil
	}
	return []*model.Store{
		&model.Store{
			ID:        1,
			BrandID:   1,
			CatalogID: 1,
		},
		&model.Store{
			ID:        2,
			BrandID:   2,
			CatalogID: 2,
		},
	}, nil
}
func (*mockSeeker) FetchBrand(id int) (*model.Brand, error) {
	switch id {
	case causeFetchBrandFailedBrandID:
		return nil, errors.New("fetch error occured")
	case causeSyncBrandFailedBrandID:
		return &model.Brand{
			ID: causeSyncBrandFailedBrandID,
		}, nil
	case 1:
		return &model.Brand{ID: 1}, nil
	case 2:
		return &model.Brand{ID: 2}, nil
	}
	return &model.Brand{}, nil
}
func (*mockSeeker) FetchProducts(storeID int, page int) ([]*model.Product, int, error) {
	switch storeID {
	case causeFetchProductFailedStoreID:
		return nil, 0, errors.New("error occured")
	case causeSyncProductFailedStoreID:
		return []*model.Product{
			&model.Product{
				ID:       causeSyncProductFailedProductID,
				Barcodes: []string{"3345678"},
			},
		}, 0, nil
	case 1:
		switch page {
		case 1:
			return []*model.Product{
				&model.Product{ID: 1, Barcodes: []string{"3345678"}},
				&model.Product{ID: 2, Barcodes: []string{"3345678"}},
				&model.Product{ID: 3, Barcodes: []string{"3345678"}},
			}, 3, nil
		case 2:
			return []*model.Product{
				&model.Product{ID: 4, Barcodes: []string{"3345678"}},
				&model.Product{ID: 5, Barcodes: []string{"3345678"}},
				&model.Product{ID: 6, Barcodes: []string{"3345678"}},
			}, 3, nil
		case 3:
			return []*model.Product{
				&model.Product{ID: 7, Barcodes: []string{"3345678"}},
				&model.Product{ID: 8, Barcodes: []string{"3345678"}},
				&model.Product{ID: 9, Barcodes: []string{"3345678"}},
			}, 3, nil
		}
	case 2:
		switch page {
		case 1:
			return []*model.Product{
				&model.Product{ID: 10, Barcodes: []string{"3345678"}},
				&model.Product{ID: 11, Barcodes: []string{"3345678"}},
				&model.Product{ID: 12, Barcodes: []string{"3345678"}},
			}, 3, nil
		case 2:
			return []*model.Product{
				&model.Product{ID: 13, Barcodes: []string{"3345678"}},
				&model.Product{ID: 14, Barcodes: []string{"3345678"}},
				&model.Product{ID: 15, Barcodes: []string{"3345678"}},
			}, 3, nil
		case 3:
			return []*model.Product{
				&model.Product{ID: 16, Barcodes: []string{"3345678"}},
				&model.Product{ID: 17, Barcodes: []string{"3345678"}},
				&model.Product{ID: 18},
			}, 3, nil
		}
	}
	return []*model.Product{}, 0, nil
}

type mockService struct {
	recMux *sync.RWMutex
	rec    map[string][]interface{}
}

func newMockService() *mockService {
	return &mockService{
		recMux: new(sync.RWMutex),
		rec:    make(map[string][]interface{}),
	}
}

func (m *mockService) putRecord(name string, value interface{}) {
	m.recMux.Lock()
	defer m.recMux.Unlock()
	m.rec[name] = append(m.rec[name], value)
}

func (m *mockService) getRecordLen(name string) int {
	m.recMux.RLock()
	defer m.recMux.RUnlock()
	return len(m.rec[name])
}

func (m *mockService) clean() {
	m.recMux.Lock()
	m.rec = make(map[string][]interface{})
	m.recMux.Unlock()
}

func (*mockService) GetExternalKeys() ([]string, error) { return []string{}, nil }
func (m *mockService) SyncCatalog(c *model.Catalog) error {
	switch c.ID {
	case causeSyncCatalogFailedCatalogID:
		return errors.New("error occured")
	}
	m.putRecord("catalogs", c)
	return nil
}
func (*mockService) GetBrand(int) (*model.Brand, error) { return nil, nil }
func (m *mockService) SyncBrand(b *model.Brand) error {
	switch b.ID {
	case causeSyncBrandFailedBrandID:
		return errors.New("error occured")
	}
	m.putRecord("brands", b)
	return nil
}
func (*mockService) PushJob(job *model.Job) error { return nil }
func (*mockService) PopJob() (*model.Job, error)  { return nil, nil }
func (m *mockService) SyncProduct(p *model.Product) error {
	switch p.ID {
	case causeSyncProductFailedProductID:
		return errors.New("error occured")
	}
	m.putRecord("products", p)
	return nil
}
func (*mockService) SyncBarcode([]*model.Barcode) error                  { return nil }
func (*mockService) GetProductIDs(int, []string) ([]int, error)          { return nil, nil }
func (*mockService) GetProducts([]int) ([]*model.Product, error)         { return nil, nil }
func (*mockService) GetStoreByBrandID(int) (*model.Store, error)         { return nil, nil }
func (*mockService) GetStoreByExternalKey(string) (*model.Store, error)  { return nil, nil }
func (*mockService) SelectStores(string, string) ([]*model.Store, error) { return nil, nil }
func (*mockService) GetStoreByID(id int) (*model.Store, error)           { return nil, nil }
func (m *mockService) SyncStore(s *model.Store) error {
	switch s.ID {
	case casueSyncStoreFailedStoreID:
		return errors.New("error occured")
	}
	m.putRecord("stores", s)
	return nil
}
func (*mockService) Close() error                              { return nil }
func (*mockService) UpdateSyncStatus(*model.SyncStatus) error  { return nil }
func (*mockService) GetSyncStatus() (*model.SyncStatus, error) { return nil, nil }

func TestProcess(t *testing.T) {
	service := newMockService()
	proc, err := procesor.New(&config.Procesor{
		PoolSize:  3,
		WorkerNum: 3,
	}, service, new(mockSeeker))
	if err != nil {
		t.Fatalf("new procesor failed:%v", err)
	}

	testCases := []struct {
		description string
		job         *model.Job
		want        map[string]int
		wantErr     bool
	}{
		{
			description: "ExternalKeyJob fetch store failed",
			job: &model.Job{
				Type:  model.ExternalKeyJob,
				Value: causeFetchStoreFailedExternalKey,
			},
			wantErr: true,
			want: map[string]int{
				"stores":   0,
				"brands":   0,
				"products": 0,
				"catalogs": 0,
			},
		},
		{
			description: "ExternalKeyJob fetch brand failed",
			job: &model.Job{
				Type:  model.ExternalKeyJob,
				Value: causeFetchBrandFailedExternalKey,
			},
			wantErr: true,
			want: map[string]int{
				"stores":   0,
				"brands":   0,
				"products": 0,
				"catalogs": 0,
			},
		},
		{
			description: "ExternalKeyJob sync brand failed",
			job: &model.Job{
				Type:  model.ExternalKeyJob,
				Value: causeSyncBrandFailedExternalKey,
			},
			wantErr: true,
			want: map[string]int{
				"stores":   0,
				"brands":   0,
				"products": 0,
				"catalogs": 0,
			},
		},
		{
			description: "ExternalKeyJob sync store failed",
			job: &model.Job{
				Type:  model.ExternalKeyJob,
				Value: causeSyncStoreFailedExternalKey,
			},
			wantErr: true,
			want: map[string]int{
				"stores":   0,
				"brands":   1,
				"products": 0,
				"catalogs": 1,
			},
		},
		{
			description: "ExternalKeyJob fetch product failed",
			job: &model.Job{
				Type:  model.ExternalKeyJob,
				Value: causeFetchProductFailedExternalKey,
			},
			wantErr: true,
			want: map[string]int{
				"stores":   1,
				"brands":   1,
				"products": 0,
				"catalogs": 1,
			},
		},
		{
			description: "ExternalKeyJob sync product failed",
			job: &model.Job{
				Type:  model.ExternalKeyJob,
				Value: causeSyncProductFailedExternalKey,
			},
			wantErr: true,
			want: map[string]int{
				"stores":   1,
				"brands":   1,
				"products": 0,
				"catalogs": 1,
			},
		},
		{
			description: "ExternalKeyJob sync catalog failed",
			job: &model.Job{
				Type:  model.ExternalKeyJob,
				Value: causeSyncCatalogFailedExternalKey,
			},
			wantErr: true,
			want: map[string]int{
				"stores":   0,
				"brands":   1,
				"products": 0,
				"catalogs": 0,
			},
		},
		{
			description: "ExternalKeyJob get job value failed",
			job: &model.Job{
				Type:  model.ExternalKeyJob,
				Value: 3345678,
			},
			wantErr: true,
			want: map[string]int{
				"stores":   0,
				"brands":   0,
				"products": 0,
				"catalogs": 0,
			},
		},
		{
			description: "ExternalKeyJob process success with one external key",
			job: &model.Job{
				Type:  model.ExternalKeyJob,
				Value: "3345678",
			},
			wantErr: false,
			want: map[string]int{
				"stores":   2,
				"brands":   2,
				"catalogs": 2,
				"products": 17,
			},
		},
		{
			description: "StoreIDJob get job value failed",
			job: &model.Job{
				Type:  model.StoreIDJob,
				Value: 3345678,
			},
			wantErr: true,
			want: map[string]int{
				"stores":   0,
				"brands":   0,
				"catalogs": 0,
				"products": 0,
			},
		},
		{
			description: "StoreIDJob fetch stores failed",
			job: &model.Job{
				Type: model.StoreIDJob,
				Value: &model.StoreIDJobValue{
					ExternalKey: causeFetchStoreFailedExternalKey,
					ID:          1,
				},
			},
			wantErr: true,
			want: map[string]int{
				"stores":   0,
				"brands":   0,
				"catalogs": 0,
				"products": 0,
			},
		},
		{
			description: "StoreIDJob success case",
			job: &model.Job{
				Type: model.StoreIDJob,
				Value: &model.StoreIDJobValue{
					ExternalKey: "3345678",
					ID:          1,
				},
			},
			wantErr: false,
			want: map[string]int{
				"stores":   1,
				"brands":   1,
				"catalogs": 1,
				"products": 9,
			},
		},
		{
			description: "BrandIDJob get job value failed",
			job: &model.Job{
				Type:  model.BrandIDJob,
				Value: "3345678",
			},
			wantErr: true,
			want: map[string]int{
				"stores":   0,
				"brands":   0,
				"catalogs": 0,
				"products": 0,
			},
		},
		{
			description: "BrandIDJob fetch brand failed",
			job: &model.Job{
				Type:  model.BrandIDJob,
				Value: causeFetchBrandFailedBrandID,
			},
			wantErr: true,
			want: map[string]int{
				"stores":   0,
				"brands":   0,
				"catalogs": 0,
				"products": 0,
			},
		},
		{
			description: "BrandIDJob success case",
			job: &model.Job{
				Type:  model.BrandIDJob,
				Value: 1,
			},
			wantErr: false,
			want: map[string]int{
				"stores":   0,
				"brands":   1,
				"catalogs": 0,
				"products": 0,
			},
		},
		{
			description: "unsupported job type",
			job: &model.Job{
				Type:  3345678,
				Value: 1,
			},
			wantErr: true,
			want: map[string]int{
				"stores":   0,
				"brands":   0,
				"catalogs": 0,
				"products": 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.description), func(t *testing.T) {
			gotErr := proc.Process(tc.job)
			if tc.wantErr {
				if gotErr == nil {
					t.Errorf("want an error, got nil")
				}
			} else {
				if gotErr != nil {
					t.Errorf("want no error, got:%v", gotErr)
				}
			}
			for name, length := range tc.want {
				if length != service.getRecordLen(name) {
					t.Errorf("name:%s want length:%d, got:%d", name, length, service.getRecordLen(name))
				}
			}
			service.clean()
		})
	}
}
