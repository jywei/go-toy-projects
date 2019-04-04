package handler_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/honestbee/habitat-catalog-cache/handler"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/honestbee/habitat-catalog-cache/model"
)

type cacheMockService struct {
	barcodes []string
}

func (c *cacheMockService) isBarcodesValid() bool {
	for _, barcode := range c.barcodes {
		switch barcode {
		case "":
			return false
		case " ":
			return false
		default:
			if strings.Contains(barcode, ",") {
				return false
			}
		}
	}
	return true
}

const (
	notFoundBrandID           = 45678910
	causeErrorBrandID         = 45678911
	causeStoreErrorBrandID    = 45678912
	causeGetStoreErrorBrandID = 45678913
)

func (*cacheMockService) GetExternalKeys() ([]string, error) { return []string{}, nil }
func (*cacheMockService) SyncCatalog(*model.Catalog) error   { return nil }
func (*cacheMockService) GetBrand(id int) (*model.Brand, error) {
	switch id {
	case notFoundBrandID:
		return nil, model.ErrNoRows
	case causeErrorBrandID:
		return nil, errors.New("error occured")
	case causeStoreErrorBrandID:
		return &model.Brand{
			ID: causeGetStoreErrorBrandID,
		}, nil
	}
	return &model.Brand{
		ID:      3345678,
		StoreID: 449,
	}, nil
}
func (*cacheMockService) SyncBrand(*model.Brand) error       { return nil }
func (*cacheMockService) PushJob(*model.Job) error           { return nil }
func (*cacheMockService) PopJob() (*model.Job, error)        { return nil, nil }
func (*cacheMockService) SyncProduct(*model.Product) error   { return nil }
func (*cacheMockService) SyncBarcode([]*model.Barcode) error { return nil }
func (c *cacheMockService) GetProductIDs(catalogID int, barcodes []string) ([]int, error) {
	switch barcodes[0] {
	case "cause-get-product-ids-failed-barcode":
		return nil, errors.New("error occured")
	case "cause-get-products-failed-barcode":
		return []int{450}, nil
	}

	c.barcodes = barcodes

	return []int{449}, nil
}
func (c *cacheMockService) GetProducts(ids []int) ([]*model.Product, error) {
	for _, id := range ids {
		switch id {
		case 450:
			return nil, errors.New("error occured")
		}
	}

	return []*model.Product{
		&model.Product{
			ID: 449,
		},
	}, nil
}
func (*cacheMockService) GetStoreByBrandID(id int) (*model.Store, error) {
	switch id {
	case causeGetStoreErrorBrandID:
		return nil, errors.New("error occured")
	}
	return &model.Store{
		ID: 449,
	}, nil
}
func (*cacheMockService) SyncStore(*model.Store) error { return nil }
func (*cacheMockService) SelectStores(exKey string, stype string) ([]*model.Store, error) {
	switch exKey {
	case "not-exist-external-key":
		return make([]*model.Store, 0), nil
	case "cause-error-external-key":
		return nil, errors.New("error occured")
	case "cause-select-stores-failed-external-key":
		return nil, errors.New("error occured")
	case "cause-get-brand-failed-external-key":
		return []*model.Store{
			&model.Store{ID: 3345678, BrandID: causeErrorBrandID},
		}, nil
	}
	return []*model.Store{
		&model.Store{ID: 3345678, CatalogID: 449},
		&model.Store{ID: 3345679, CatalogID: 449},
	}, nil
}
func (*cacheMockService) GetStoreByID(id int) (*model.Store, error) { return nil, nil }
func (*cacheMockService) Close() error                              { return nil }
func (*cacheMockService) UpdateSyncStatus(*model.SyncStatus) error  { return nil }
func (*cacheMockService) GetSyncStatus() (*model.SyncStatus, error) { return nil, nil }

var cacheEnv *handler.Env
var cacheService *cacheMockService

func init() {
	logger := zerolog.New(ioutil.Discard)
	cacheService = new(cacheMockService)
	cacheEnv = &handler.Env{
		Logger:  &logger,
		Service: cacheService,
	}
}

func TestGetStores(t *testing.T) {
	testCases := []struct {
		request *http.Request
		want    interface{}
		wantErr *handler.Error
	}{
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"cause-select-stores-failed-external-key"},
				"storeType":   []string{""},
				"fields":      []string{""},
			}},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.ServerInternalErrMsg},
		},
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"normal-external-key"},
				"storeType":   []string{"invalid-store-type"},
				"fields":      []string{""},
			}},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.InvalidAttributeErrMsg},
		},
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"normal-external-key"},
				"storeType":   []string{""},
				"fields":      []string{""},
			}},
			want: []*model.Store{
				&model.Store{ID: 3345678, CatalogID: 449},
				&model.Store{ID: 3345679, CatalogID: 449},
			},
			wantErr: nil,
		},
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"cause-get-brand-failed-external-key"},
				"storeType":   []string{""},
				"fields":      []string{"brand"},
			}},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.ServerInternalErrMsg},
		},
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"normal-external-key"},
				"storeType":   []string{""},
				"fields":      []string{"brand"},
			}},
			want: []*model.Store{
				&model.Store{ID: 3345678, CatalogID: 449, Brand: &model.Brand{
					ID:      3345678,
					StoreID: 449,
				}},
				&model.Store{ID: 3345679, CatalogID: 449, Brand: &model.Brand{
					ID:      3345678,
					StoreID: 449,
				}},
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with external key:%q, store type:%q, fields:%q",
			tc.request.Form["externalKey"], tc.request.Form["storeType"], tc.request.Form["fields"]), func(t *testing.T) {
			got, gotErr := handler.GetStores(cacheEnv, nil, tc.request)
			switch tc.wantErr {
			case nil:
				switch gotErr {
				case nil:
					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want:%v, got:%v", tc.want, got)
					}
				default:
					t.Errorf("error want: nil, got:%v", gotErr)
				}
			default:
				if err, ok := gotErr.(*handler.Error); ok {
					if tc.wantErr.OutputErr != err.OutputErr {
						t.Errorf("error want:%s, got:%s", tc.wantErr.OutputErr, err.OutputErr)
					}
				} else {
					t.Errorf("cast error:%v failed", gotErr)
				}
			}
		})
	}
}

func TestGetBrand(t *testing.T) {
	testCases := []struct {
		param   httprouter.Params
		want    interface{}
		wantErr *handler.Error
	}{
		{
			param: httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: "cast-failed-id",
				},
			},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.InvalidAttributeErrMsg},
		},
		{
			param: httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: strconv.Itoa(notFoundBrandID),
				},
			},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.RecordNotFoundErrMsg},
		},
		{
			param: httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: strconv.Itoa(causeErrorBrandID),
				},
			},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.ServerInternalErrMsg},
		},
		{
			param: httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: strconv.Itoa(causeStoreErrorBrandID),
				},
			},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.ServerInternalErrMsg},
		},
		{
			param: httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: strconv.Itoa(3345678),
				},
			},
			want: &model.Brand{
				ID:                       3345678,
				StoreID:                  449,
				MinimumOrderFreeDelivery: "0.0",
				DefaultDeliveryFee:       "0.0",
				MinimumSpend:             "0.0",
				MinimumSpendExtraFee:     "0.0",
				DefaultConciergeFee:      "0.0",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with id:%q", tc.param.ByName("id")), func(t *testing.T) {
			got, gotErr := handler.GetBrand(cacheEnv, tc.param, nil)
			switch tc.wantErr {
			case nil:
				switch gotErr {
				case nil:
					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want:%v, got:%v", tc.want, got)
					}
				default:
					t.Errorf("error want: nil, got:%v", gotErr)
				}
			default:
				if err, ok := gotErr.(*handler.Error); ok {
					if tc.wantErr.OutputErr != err.OutputErr {
						t.Errorf("error want:%s, got:%s", tc.wantErr.OutputErr, err.OutputErr)
					}
				} else {
					t.Errorf("cast error:%v failed", gotErr)
				}
			}
		})
	}
}

func TestGetProducts(t *testing.T) {
	testCases := []struct {
		request *http.Request
		want    interface{}
		wantErr *handler.Error
	}{
		{
			request: &http.Request{Form: url.Values{}},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.MissingParametersErrMsg},
		},
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"not-exist-external-key"},
				"storeType":   []string{""},
				"barcodes":    []string{"1,2,3,4"},
			}},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.RecordNotFoundErrMsg},
		},
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"cause-error-external-key"},
				"storeType":   []string{""},
				"barcodes":    []string{"1,2,3,4"},
			}},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.ServerInternalErrMsg},
		},
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"correct-external-key"},
				"storeType":   []string{""},
				"barcodes":    []string{"cause-get-product-ids-failed-barcode"},
			}},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.ServerInternalErrMsg},
		},
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"correct-external-key"},
				"storeType":   []string{""},
				"barcodes":    []string{"cause-get-products-failed-barcode"},
			}},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.ServerInternalErrMsg},
		},
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"correct-external-key"},
				"storeType":   []string{"not-exist-store-type"},
				"barcodes":    []string{"1,2,3,4"},
			}},
			want:    nil,
			wantErr: &handler.Error{OutputErr: handler.InvalidAttributeErrMsg},
		},
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"correct-external-key"},
				"storeType":   []string{""},
				"barcodes":    []string{"correct-barcode-1,correct-barcode-2"},
			}},
			want: []*model.Product{
				&model.Product{
					ID: 449,
				},
			},
			wantErr: nil,
		},
		{
			request: &http.Request{Form: url.Values{
				"externalKey": []string{"correct-external-key"},
				"storeType":   []string{""},
				"barcodes":    []string{"correct-barcode-1,correct-barcode-2,, ,"},
			}},
			want: []*model.Product{
				&model.Product{
					ID: 449,
				},
			},
			wantErr: nil,
		},
		{
			request: &http.Request{Form: url.Values{
				"external_key": []string{"correct-external-key"},
				"storeType":    []string{""},
				"barcodes":     []string{"correct-barcode-1,correct-barcode-2,, ,"},
			}},
			want: []*model.Product{
				&model.Product{
					ID: 449,
				},
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("with external key:%q, store type:%q, barcodes:%v",
			tc.request.Form["externalKey"], tc.request.Form["storeType"], tc.request.Form["barcodes"]), func(t *testing.T) {
			got, gotErr := handler.GetProducts(cacheEnv, nil, tc.request)
			switch tc.wantErr {
			case nil:
				switch gotErr {
				case nil:
					if !reflect.DeepEqual(tc.want, got) {
						t.Errorf("want:%v, got:%v", tc.want, got)
					}
				default:
					t.Errorf("error want: nil, got:%v", gotErr)
				}
			default:
				if err, ok := gotErr.(*handler.Error); ok {
					if tc.wantErr.OutputErr != err.OutputErr {
						t.Errorf("error want:%s, got:%s", tc.wantErr.OutputErr, err.OutputErr)
					}
				} else {
					t.Errorf("cast error:%v failed", gotErr)
				}
			}

			if !cacheService.isBarcodesValid() {
				t.Errorf("barcode:%v is not valid", cacheService.barcodes)
			}
		})
	}
}
