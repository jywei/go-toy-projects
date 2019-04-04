package autoexec_test

import (
	"fmt"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/honestbee/habitat-catalog-cache/autoexec"
	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/rs/zerolog"
)

type mockService struct {
	keys    []string
	counter int
	mux     *sync.Mutex

	wait chan struct{}
}

func newMockService() *mockService {
	return &mockService{
		wait: make(chan struct{}, 10),
		mux:  new(sync.Mutex),
	}
}

func (m *mockService) waiting() {
	<-m.wait
}

func (m *mockService) setExternalKeys(keys []string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.keys = keys
}

func (m *mockService) getCounter() int {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.counter
}

func (m *mockService) GetExternalKeys() ([]string, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.keys, nil
}
func (*mockService) SyncCatalog(*model.Catalog) error   { return nil }
func (*mockService) GetBrand(int) (*model.Brand, error) { return nil, nil }
func (*mockService) SyncBrand(*model.Brand) error       { return nil }
func (m *mockService) PushJob(*model.Job) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.counter++
	m.wait <- struct{}{}
	return nil
}
func (*mockService) PopJob() (*model.Job, error)                         { return nil, nil }
func (*mockService) SyncProduct(*model.Product) error                    { return nil }
func (*mockService) GetProducts([]int) ([]*model.Product, error)         { return nil, nil }
func (*mockService) SyncBarcode([]*model.Barcode) error                  { return nil }
func (*mockService) GetProductIDs(int, []string) ([]int, error)          { return nil, nil }
func (*mockService) GetStoreByBrandID(int) (*model.Store, error)         { return nil, nil }
func (*mockService) GetStoreByExternalKey(string) (*model.Store, error)  { return nil, nil }
func (*mockService) SelectStores(string, string) ([]*model.Store, error) { return nil, nil }
func (*mockService) GetStoreByID(id int) (*model.Store, error)           { return nil, nil }
func (*mockService) SyncStore(*model.Store) error                        { return nil }
func (*mockService) Close() error                                        { return nil }
func (*mockService) UpdateSyncStatus(*model.SyncStatus) error            { return nil }
func (*mockService) GetSyncStatus() (*model.SyncStatus, error)           { return nil, nil }

func TestAutoExec(t *testing.T) {
	logger := zerolog.New(ioutil.Discard)

	testCases := []struct {
		description string
		exKeys      []string
		want        int
		wait        func()
	}{
		{
			description: "sync up three different jobs",
			exKeys:      []string{"1", "2", "3"},
			want:        3,
		},
		{
			description: "sync up three same jobs",
			exKeys:      []string{"1", "1", "1"},
			want:        1,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.description), func(t *testing.T) {
			service := newMockService()
			service.setExternalKeys(tc.exKeys)
			auto, err := autoexec.New(&config.AutoExec{
				SyncupPeriodSec: 1,
			}, service, &logger)
			if err != nil {
				t.Fatalf("new autoexec failed:%v", err)
			}

			service.waiting()
			auto.Close()

			got := service.getCounter()
			if tc.want != got {
				t.Errorf("count want:%d, got:%d", tc.want, got)
			}
		})
	}
}
