package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/honestbee/habitat-catalog-cache/model"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func rethrough(e *Env, job *model.Job, orgErr error) error {
	newErr := errors.Wrapf(
		e.Service.PushJob(job),
		"handler: [GetProducts] PushJob failed on job:%d, value:%v",
		job.Type, job.Value,
	)
	switch orgErr {
	case nil:
		return newErr
	default:
		switch newErr {
		case nil:
			return orgErr
		default:
			return errors.Wrapf(orgErr, ", %v", newErr)
		}
	}
}

// GetProducts handles get products request.
func GetProducts(e *Env, ps httprouter.Params, r *http.Request) (interface{}, error) {
	exkey := ValueFetcher("externalKey", r.FormValue)
	if exkey == "" {
		return nil, NewErr(MissingParametersErrCode, errors.New("handler: [GetProducts] external key is empty"))
	}
	storeType := ValueFetcher("storeType", r.FormValue)
	if !isStoreTypeValid(storeType) {
		return nil, NewErr(InvalidAttributeErrCode, errors.Errorf("unknown store type:%s", storeType))
	}

	barcodes := make([]string, 0)
	for _, barcode := range strings.Split(r.FormValue("barcodes"), ",") {
		barcode = strings.TrimSpace(barcode)
		if barcode != "" && barcode != " " {
			barcodes = append(barcodes, barcode)
		}
	}

	stores, err := e.Service.SelectStores(exkey, storeType)
	if err != nil {
		return nil, NewErr(ServerInternalErrCode, errors.Wrapf(err, "handler: [GetProducts] select stores failed"))
	}
	if len(stores) == 0 {
		return nil, NewErr(RecordNotFoundErrCode, errors.Wrapf(rethrough(e, &model.Job{
			Type:  model.ExternalKeyJob,
			Value: exkey,
		}, errors.New("handler: [GetProducts] len of select stores is 0")), "handler: [GetProducts] select stores failed"))
	}

	store := stores[0]
	pids, err := e.Service.GetProductIDs(store.CatalogID, barcodes)
	if err != nil {
		return nil, NewErr(ServerInternalErrCode, errors.Wrapf(err, "handler: [GetProducts] get product ids failed"))
	}

	if len(pids) == 0 {
		return []*model.Product{}, nil
	}

	products, err := e.Service.GetProducts(pids)
	if err != nil {
		return nil, NewErr(ServerInternalErrCode, errors.Wrapf(err, "handler: [GetProducts] get products failed"))
	}

	return products, nil
}

// GetBrand handles get brand request.
func GetBrand(e *Env, ps httprouter.Params, r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		return nil, NewErr(InvalidAttributeErrCode, errors.Wrapf(err, "handler: [GetBrand] convert brand id failed"))
	}

	brand, err := getBrand(e, id)
	if err != nil {
		return nil, err
	}

	store, err := e.Service.GetStoreByBrandID(brand.ID)
	if err != nil {
		return nil, NewErr(ServerInternalErrCode, errors.Wrapf(err, "handler: [GetBrand] get store failed"))
	}

	brand.Validate()
	brand.ValidateWithStore(store)

	return brand, nil
}

func getBrand(e *Env, id int) (*model.Brand, error) {
	brand, err := e.Service.GetBrand(id)
	if err != nil {
		code := 0
		switch err {
		case model.ErrNoRows:
			code = RecordNotFoundErrCode
		default:
			code = ServerInternalErrCode
		}
		return nil, NewErr(code, errors.Wrapf(err, "handler: [getBrand] get brand id:%d failed", id))
	}
	return brand, nil
}

func isStoreTypeValid(s string) bool {
	switch s {
	case "":
	case "habitat":
	case "value_store":
	default:
		return false
	}
	return true
}

// GetStores handles get stores request.
func GetStores(e *Env, ps httprouter.Params, r *http.Request) (interface{}, error) {
	exKey := ValueFetcher("externalKey", r.FormValue)
	storeType := ValueFetcher("storeType", r.FormValue)
	if !isStoreTypeValid(storeType) {
		return nil, NewErr(InvalidAttributeErrCode, errors.Errorf("unknown store type:%s", storeType))
	}

	stores, err := e.Service.SelectStores(exKey, storeType)
	if err != nil {
		return nil, NewErr(ServerInternalErrCode, errors.Wrapf(err, "handler: [GetStores] select stores failed"))
	}

	switch r.FormValue("fields") {
	case "brand":
		m := make(map[int]*model.Brand)
		for _, store := range stores {
			if _, exist := m[store.BrandID]; !exist {
				brand, err := getBrand(e, store.BrandID)
				if err != nil {
					return nil, err
				}
				m[store.BrandID] = brand
			}
			store.Brand = m[store.BrandID]
		}
	}
	return stores, nil
}
