package procesor

import (
	"sync"
	"time"

	"github.com/honestbee/habitat-catalog-cache/config"
	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/honestbee/habitat-catalog-cache/seeker"
	"github.com/pkg/errors"
)

type errorWorkerGroup struct {
	wg      *sync.WaitGroup
	errOnce *sync.Once
	err     error
	pool    chan func() error
}

func newErrorWorkerGroup(workerNum, poolSize int) *errorWorkerGroup {
	workerGroup := &errorWorkerGroup{
		wg:      new(sync.WaitGroup),
		errOnce: new(sync.Once),
		pool:    make(chan func() error, poolSize),
	}

	for i := 0; i < workerNum; i++ {
		go workerGroup.worker()
	}

	return workerGroup
}

func (ew *errorWorkerGroup) worker() {
	for job := range ew.pool {
		if err := job(); err != nil {
			ew.errOnce.Do(func() {
				ew.err = err
			})
		}
		ew.wg.Done()
	}
}

func (ew *errorWorkerGroup) waitAndClose() error {
	ew.wg.Wait()
	close(ew.pool)
	return ew.err
}

func (ew *errorWorkerGroup) fork(f func() error) {
	ew.wg.Add(1)
	ew.pool <- f
}

// Procesor is the processing unit who does all cache service actions.
type Procesor interface {
	Process(job *model.Job) error
}

type procesor struct {
	service   model.Service
	seeker    seeker.Seeker
	workerNum int
	poolSize  int
}

// New returns a procesor instance.
func New(conf *config.Procesor, service model.Service, seeker seeker.Seeker) (Procesor, error) {
	return &procesor{
		service:   service,
		seeker:    seeker,
		workerNum: conf.WorkerNum,
		poolSize:  conf.PoolSize,
	}, nil
}

// Process do the cache service caching job.
func (p *procesor) Process(job *model.Job) error {
	exe := newExecutor(p.service, p.seeker, p.workerNum, p.poolSize)
	switch job.Type {
	case model.ExternalKeyJob:
		var key string
		if err := job.GetValue(&key); err != nil {
			return errors.Wrapf(err, "procesor: [Process] job GetValue failed")
		}
		exe.processExternalKey(key)
	case model.StoreIDJob:
		v := new(model.StoreIDJobValue)
		if err := job.GetValue(v); err != nil {
			return errors.Wrapf(err, "procesor: [Process] job GetValue failed")
		}
		exe.processSingleStore(v.ExternalKey, v.ID)
	case model.BrandIDJob:
		var id int
		if err := job.GetValue(&id); err != nil {
			return errors.Wrapf(err, "procesor: [Process] job GetValue failed")
		}
		exe.workerGroup.fork(func() error {
			_, err := exe.processBrand(id)
			return err
		})
	default:
		return errors.Errorf("procesor: [Process] not supported job type:%v", job.Type)
	}

	syncStatus := model.SyncStatusSuccess
	err := errors.Wrapf(exe.workerGroup.waitAndClose(), "procesor: [Process] a series of processes failed")
	if err != nil {
		syncStatus = model.SyncStatusFailed
	}
	p.service.UpdateSyncStatus(&model.SyncStatus{
		LastSyncTime:   time.Now().UTC().String(),
		LastSyncStatus: syncStatus,
	})
	return err
}

type executor struct {
	service     model.Service
	seeker      seeker.Seeker
	workerGroup *errorWorkerGroup
}

func newExecutor(service model.Service, seeker seeker.Seeker, workerNum, poolSize int) *executor {
	return &executor{
		service:     service,
		seeker:      seeker,
		workerGroup: newErrorWorkerGroup(workerNum, poolSize),
	}
}

func (e *executor) processSingleStore(externalKey string, id int) {
	e.workerGroup.fork(func() error {
		stores, err := e.seeker.FetchStores(externalKey)
		if err != nil {
			return errors.Wrapf(err, "procesor: [processExternalKey] FetchStores failed")
		}

		var store *model.Store
		for _, s := range stores {
			if s.ID == id {
				store = s
				break
			}
		}

		if store != nil {
			e.processStores([]*model.Store{store}, externalKey)
		}
		return nil
	})
}

func (e *executor) processExternalKey(externalKey string) {
	e.workerGroup.fork(func() error {
		stores, err := e.seeker.FetchStores(externalKey)
		if err != nil {
			return errors.Wrapf(err, "procesor: [processExternalKey] FetchStores failed")
		}

		e.processStores(stores, externalKey)
		return nil
	})
}

func (e *executor) processStores(stores []*model.Store, externalKey string) {
	for _, store := range stores {

		store := store
		e.workerGroup.fork(func() error {
			brandSlug, err := e.processBrand(store.BrandID)
			if err != nil {
				return errors.Wrapf(err, "procesor: [processStores] processBrand failed")
			}
			if err := e.service.SyncCatalog(&model.Catalog{ID: store.CatalogID}); err != nil {
				return errors.Wrapf(err, "procesor: [processStores] SyncCatalog failed")
			}
			store.ExternalKey = externalKey
			if err := e.service.SyncStore(store); err != nil {
				return errors.Wrapf(err, "procesor: [processStores] SyncStore failed")
			}
			if err := e.processProducts(brandSlug, store.BrandID, store.ID, store.CatalogID); err != nil {
				return errors.Wrapf(err, "procesor: [processStores] processProducts failed")
			}
			return nil
		})
	}
}

func (e *executor) processBrand(brandID int) (string, error) {
	brand, err := e.seeker.FetchBrand(brandID)
	if err != nil {
		return "", errors.Wrapf(err, "procesor: [processBrand] FetchBrand failed")
	}
	if err = e.service.SyncBrand(brand); err != nil {
		return "", errors.Wrapf(err, "procesor: [processBrand] SyncBrand failed")
	}
	return brand.Slug, nil
}

func (e *executor) processProducts(brandSlug string, brandID, storeID, catalogID int) error {
	pages, err := e.manageProducts(brandSlug, brandID, storeID, catalogID, 1)
	if err != nil {
		return errors.Wrapf(err, "procesor: [processProducts] manageProducts failed")
	}

	for page := 2; page <= pages; page++ {

		page := page
		e.workerGroup.fork(func() error {
			_, err := e.manageProducts(brandSlug, brandID, storeID, catalogID, page)
			if err != nil {
				return errors.Wrapf(err, "procesor: [processProducts] manageProducts failed")
			}
			return nil
		})
	}
	return nil
}

func (e *executor) manageProducts(brandSlug string, brandID, storeID, catalogID, page int) (int, error) {
	products, pages, err := e.seeker.FetchProducts(storeID, page)
	if err != nil {
		return 0, errors.Wrapf(err, "procesor: [manageProducts] FetchProducts failed")
	}

	for _, product := range products {
		if len(product.Barcodes) == 0 {
			continue
		}
		product.CatalogID = catalogID
		product.BrandID = brandID
		product.BrandSlug = brandSlug
		if err = e.service.SyncProduct(product); err != nil {
			return 0, errors.Wrapf(err, "procesor: [manageProducts] SyncProduct failed")
		}

		bs := make([]*model.Barcode, 0, len(product.Barcodes))
		for _, barcode := range product.Barcodes {
			bs = append(bs, &model.Barcode{
				ProductID: product.ID,
				CatalogID: product.CatalogID,
				Barcode:   barcode,
				IsActive:  true,
			})
		}
		if err = e.service.SyncBarcode(bs); err != nil {
			return 0, errors.Wrapf(err, "procesor: [manageProducts] SyncBarcode failed")
		}
	}
	return pages, nil
}
